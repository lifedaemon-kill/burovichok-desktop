package ui

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/browser"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/service/calc"
	chartService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/chart"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/service/database"
	inmemoryStorage "github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
)

type importer interface {
	ParseBlockOneFile(path string, cfg models.OperationConfig) ([]models.TableOne, error)
	ParseBlockTwoFile(path string) ([]models.TableTwo, error)
	ParseBlockThreeFile(path string) ([]models.TableThree, error)
	ParseBlockFourFile(path string) ([]models.TableFour, error)
}

type converterService interface {
	ParseFlexibleTime(raw string) (time.Time, error)
}

type Service struct {
	app              fyne.App
	window           fyne.Window
	zLog             logger.Logger
	importer         importer
	memBlocksStorage inmemoryStorage.InMemoryBlocksStorage
	db               *database.Service
	converter        converterService
	chart            chartService.Service
	calc             *calc.Service

	serverMutex      sync.Mutex
	serverListener   net.Listener
	serverPort       string
	chartHtmlToServe string

	loadingLabel *widget.Label
	progressBar  *widget.ProgressBarInfinite
}

func NewService(cfg config.UI, zLog logger.Logger, imp importer, converter converterService,
	memBlocksStorage inmemoryStorage.InMemoryBlocksStorage, db *database.Service, chart chartService.Service, calc *calc.Service) *Service {
	a := app.New()
	// chartSvc := chartService.NewService() // Сервис теперь передается снаружи
	win := a.NewWindow(cfg.Name)
	win.Resize(fyne.NewSize(float32(cfg.Width), float32(cfg.Height)))

	loadingLbl := widget.NewLabel("Идет обработка...")
	loadingLbl.Hide()
	progressBr := widget.NewProgressBarInfinite()
	progressBr.Hide()

	return &Service{
		app:              a,
		window:           win,
		zLog:             zLog,
		importer:         imp,
		memBlocksStorage: memBlocksStorage,
		db:               db,
		converter:        converter,
		chart:            chart,
		loadingLabel:     loadingLbl,
		progressBar:      progressBr,
		calc:             calc,
	}
}

// --- Методы для управления веб-сервером ---

func (s *Service) startLocalWebServer() error {
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()

	if s.serverListener != nil {
		return nil
	}

	s.zLog.Infow("Starting local chart web server...")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		s.zLog.Errorw("Failed to find free port for chart server", "error", err)
		return fmt.Errorf("не удалось найти порт для сервера графика: %w", err)
	}
	s.serverListener = listener
	s.serverPort = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	s.zLog.Infow("Chart server listening", "address", listener.Addr().String())

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serveChartHTML)

	// Запускаем сервер в отдельной горутине, чтобы не блокировать UI
	go func() {
		if err := http.Serve(s.serverListener, mux); err != nil && err != http.ErrServerClosed {
			s.zLog.Errorw("Chart server error", "error", err)
			s.serverMutex.Lock()
			s.serverListener = nil
			s.serverPort = ""
			s.serverMutex.Unlock()
		}
		s.zLog.Infow("Chart server stopped.")
	}()

	return nil
}

func (s *Service) serveChartHTML(w http.ResponseWriter, r *http.Request) {
	s.serverMutex.Lock()
	htmlPath := s.chartHtmlToServe
	s.serverMutex.Unlock()

	if htmlPath == "" {
		http.Error(w, "График еще не сгенерирован", http.StatusNotFound)
		s.zLog.Errorw("Request for chart before generation")
		return
	}

	s.zLog.Infow("Serving chart HTML", "path", htmlPath)
	http.ServeFile(w, r, htmlPath) // Отдаем файл
}

// --- Конец методов веб-сервера ---

type ratioLayout struct{ ratio float32 }

func (r *ratioLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	firstW := size.Width * r.ratio
	objects[0].Resize(fyne.NewSize(firstW, size.Height))
	objects[1].Resize(fyne.NewSize(size.Width-firstW, size.Height))
	objects[1].Move(fyne.NewPos(firstW, 0))
}

func (r *ratioLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var h float32
	for _, o := range objects {
		if hh := o.MinSize().Height; hh > h {
			h = hh
		}
	}
	return fyne.NewSize(0, h)
}

func (s *Service) Run(ctx context.Context) error {
	s.window.SetOnClosed(func() {
		s.serverMutex.Lock()
		if s.serverListener != nil {
			s.zLog.Infow("Closing chart server...")
			s.serverListener.Close()
			s.serverListener = nil
			s.serverPort = ""
		}
		s.serverMutex.Unlock()
	})

	chartBtn := widget.NewButton("Интерактивный График T/P (Блок 1)", func() {
		blockOneData, err := s.memBlocksStorage.GetAllBlockOneData()
		if err != nil {
			s.zLog.Errorw("Failed to get BlockOne data for chart", "error", err)
			dialog.ShowError(fmt.Errorf("не удалось получить данные Блока 1: %w", err), s.window)
			return
		}
		if len(blockOneData) == 0 {
			dialog.ShowInformation("Нет данных", "Недостаточно данных для построения графика", s.window)
			return
		}

		htmlPath, err := s.chart.GeneratePressureTempChart(blockOneData)
		if err != nil {
			s.zLog.Errorw("Failed to generate chart HTML", "error", err)
			dialog.ShowError(fmt.Errorf("ошибка генерации HTML графика: %w", err), s.window)
			return
		}

		s.serverMutex.Lock()
		s.chartHtmlToServe = htmlPath
		s.serverMutex.Unlock()

		if err := s.startLocalWebServer(); err != nil {
			dialog.ShowError(fmt.Errorf("ошибка запуска веб-сервера для графика: %w", err), s.window)
			return
		}

		// Даем серверу микро-паузу на старт (не самое элегантное решение, но простое)
		time.Sleep(100 * time.Millisecond)

		s.serverMutex.Lock()
		port := s.serverPort
		s.serverMutex.Unlock()

		if port == "" {
			dialog.ShowError(fmt.Errorf("не удалось получить порт веб-сервера"), s.window)
			return
		}
		chartURL := fmt.Sprintf("http://127.0.0.1:%s/", port)
		s.zLog.Infow("Opening chart in browser", "url", chartURL)

		err = browser.OpenURL(chartURL)
		if err != nil {
			s.zLog.Errorw("Failed to open browser for chart", "error", err)
			dialog.ShowError(fmt.Errorf("не удалось открыть браузер: %w", err), s.window)
			return
		}
	})

	// 1) поле выбора файла + кнопка
	pathEntry := widget.NewEntry()
	pathEntry.PlaceHolder = "Файл не выбран"

	chooseBtn := widget.NewButton("Выбрать файл", func() {
		d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			// 1) Ошибка при открытии диалога
			if err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			// 2) Пользователь нажал «Cancel» — r == nil и err == nil
			if r == nil {
				return
			}
			// 3) Успешно выбрали файл
			defer r.Close()
			pathEntry.SetText(r.URI().Path())
		}, s.window)

		d.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		d.Show()
	})

	// 2) выбор типа документа
	docTypes := []string{"TableOne", "TableTwo", "TableThree", "TableFour"}
	typeSelect := widget.NewSelect(docTypes, nil)
	typeSelect.PlaceHolder = "Выберите тип документа"

	// 3) кнопка Import
	importBtn := widget.NewButton("Import", func() {
		path := pathEntry.Text
		typ := typeSelect.Selected
		if path == "" || typ == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите файл и тип документа", s.window)
			return
		}

		// если TableOne — сначала собираем параметры
		if typ == "TableOne" {
			showTableOneForm(s, path)
		} else {
			s.doGenericImport(ctx, path, typ)
		}
	})

	// 4) кнопка очистки хранилища
	clearBtn := widget.NewButton("Очистить хранилище", func() {
		dialog.ShowConfirm("Подтверждение", "Удалить все данные?", func(ok bool) {
			if !ok {
				return
			}
			if err := s.memBlocksStorage.ClearAll(); err != nil {
				s.zLog.Errorw("Clear failed", "error", err)
				dialog.ShowError(fmt.Errorf("ошибка очистки: %w", err), s.window)
				return
			}
			dialog.ShowInformation("Очищено", "Данные удалены", s.window)
		}, s.window)
	})

	guidebookBtn := widget.NewButton("Заполнить Шапку Отчета (Блок 5)", func() {
		s.showBlockFiveForm(ctx)
	})

	header := container.New(&ratioLayout{ratio: 0.7}, pathEntry, chooseBtn)
	content := container.NewVBox(
		widget.NewLabel("1. Выберите файл и тип данных:"),
		header,
		typeSelect,
		widget.NewSeparator(),
		widget.NewLabel("2. Выполните действия:"),
		importBtn,
		guidebookBtn,
		chartBtn,
		widget.NewSeparator(),
		widget.NewLabel("3. Управление хранилищем:"),
		clearBtn,
		widget.NewSeparator(),
		s.loadingLabel,
		s.progressBar,
	)

	s.window.SetContent(content)
	s.window.ShowAndRun()
	return nil
}

// showTableOneForm показывает форму параметров гидростатики и по клику "Ок" запускает импорт.
func showTableOneForm(s *Service, path string) {
	ws := widget.NewEntry() // Работа: c
	we := widget.NewEntry() // Работа: по
	is := widget.NewEntry() // Простой: c (необязательно)
	ie := widget.NewEntry() // Простой: по (необязательно)
	wr := widget.NewEntry() // Плотность (работа)
	ir := widget.NewEntry() // Плотность (простои) (необязательно)
	dh := widget.NewEntry() // Δh (м)
	unit := widget.NewSelect([]string{"kgf/cm2", "bar", "atm"}, nil)

	ws.PlaceHolder = "YYYY-MM-DD"
	we.PlaceHolder = "YYYY-MM-DD"
	is.PlaceHolder = "YYYY-MM-DD"
	ie.PlaceHolder = "YYYY-MM-DD"

	items := []*widget.FormItem{
		{Text: "Работа: c", Widget: ws},
		{Text: "Работа: по", Widget: we},
		{Text: "Простой: c (необязательно)", Widget: is},
		{Text: "Простой: по (необязательно)", Widget: ie},
		{Text: "Плотность (работа)", Widget: wr},
		{Text: "Плотность (простои, по умолчанию = плотность работы)", Widget: ir},
		{Text: "Δh (м)", Widget: dh},
		{Text: "Единица давления", Widget: unit},
	}

	dlg := dialog.NewForm("Параметры гидростатики", "Ок", "Отмена", items,
		func(ok bool) {
			if !ok {
				return // пользователь отменил
			}

			var missing []string
			if strings.TrimSpace(ws.Text) == "" {
				missing = append(missing, "Работа: c")
			}
			if strings.TrimSpace(we.Text) == "" {
				missing = append(missing, "Работа: по")
			}
			if strings.TrimSpace(wr.Text) == "" {
				missing = append(missing, "Плотность (работа)")
			}
			if strings.TrimSpace(dh.Text) == "" {
				missing = append(missing, "Δh (м)")
			}
			if unit.Selected == "" {
				missing = append(missing, "Единица давления")
			}
			if len(missing) > 0 {
				dialog.ShowInformation(
					"Ошибка",
					"Пожалуйста, заполните обязательные поля:\n• "+strings.Join(missing, "\n• "),
					s.window,
				)
				return
			}

			var parseErrs []string

			workStart, err := s.converter.ParseFlexibleTime(ws.Text)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprintf("Работа: c — %v", err))
			}
			workEnd, err := s.converter.ParseFlexibleTime(we.Text)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprintf("Работа: по — %v", err))
			}

			var idleStart, idleEnd time.Time
			if strings.TrimSpace(is.Text) != "" {
				idleStart, err = s.converter.ParseFlexibleTime(is.Text)
				if err != nil {
					parseErrs = append(parseErrs, fmt.Sprintf("Простой: c — %v", err))
				}
			}
			if strings.TrimSpace(ie.Text) != "" {
				idleEnd, err = s.converter.ParseFlexibleTime(ie.Text)
				if err != nil {
					parseErrs = append(parseErrs, fmt.Sprintf("Простой: по — %v", err))
				}
			}

			if len(parseErrs) > 0 {
				dialog.ShowError(
					fmt.Errorf("Неверный формат даты/времени:\n%s", strings.Join(parseErrs, "\n")),
					s.window,
				)
				return
			}

			// дальше парсим плотности и формируем cfg как раньше
			workDens, _ := strconv.ParseFloat(wr.Text, 64)
			var idleDens float64
			if strings.TrimSpace(ir.Text) == "" {
				idleDens = workDens
			} else {
				idleDens, _ = strconv.ParseFloat(ir.Text, 64)
			}
			depthDiff, _ := strconv.ParseFloat(dh.Text, 64)

			cfg := models.OperationConfig{
				WorkStart:    workStart,
				WorkEnd:      workEnd,
				IdleStart:    idleStart,
				IdleEnd:      idleEnd,
				WorkDensity:  workDens,
				IdleDensity:  idleDens,
				DepthDiff:    depthDiff,
				PressureUnit: unit.Selected,
			}

			fileName := filepath.Base(path)
			s.showLoadingIndicator(fileName)

			go s.doTableOneImport(path, cfg)

		}, s.window)

	dlg.Show()
}

// --- НОВАЯ ФУНКЦИЯ: Показ формы для Блока 5 ---
func (s *Service) showBlockFiveForm(ctx context.Context) {
	s.zLog.Debugw("Opening Block 5 form")

	// 1. Загружаем справочники из БД
	oilFieldsModels, err := s.db.GetAllOilFields()
	if err != nil {
		s.zLog.Errorw("Failed to get oil fields", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки месторождений: %w", err), s.window)
		return
	}
	horizonsModels, err := s.db.GetAllProductiveHorizons()
	if err != nil {
		s.zLog.Errorw("Failed to get horizons", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки горизонтов: %w", err), s.window)
		return
	}
	instrumentTypesModels, err := s.db.GetAllInstrumentTypes()
	if err != nil {
		s.zLog.Errorw("Failed to get instrument types", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки типов приборов: %w", err), s.window)
		return
	}

	oilFields := make([]string, len(oilFieldsModels))
	for i, item := range oilFieldsModels {
		oilFields[i] = item.Name
	}
	horizons := make([]string, len(horizonsModels))
	for i, item := range horizonsModels {
		horizons[i] = item.Name
	}
	instrumentTypes := make([]string, len(instrumentTypesModels))
	for i, item := range instrumentTypesModels {
		instrumentTypes[i] = item.Name
	}

	// 2. Создаем виджеты для формы
	fieldNameSelect := widget.NewSelect(oilFields, nil) // Или Entry, если надо и вводить? Пока Select.
	fieldNumberEntry := widget.NewEntry()
	clusterNumberEntry := widget.NewEntry()
	horizonSelect := widget.NewSelect(horizons, nil)
	startTimeEntry := widget.NewEntry() // TODO: Может, Date Picker? Пока Entry.
	endTimeEntry := widget.NewEntry()   // TODO: Может, Date Picker? Пока Entry.
	instrumentTypeSelect := widget.NewSelect(instrumentTypes, nil)
	instrumentNumberEntry := widget.NewEntry()
	measuredDepthEntry := widget.NewEntry()

	// Заглушки для полей, которые не заполняем вручную или рассчитываются
	// Смотрим models.TableFive и create_table_five.sql
	// true_vertical_depth - REAL NOT NULL
	// true_vertical_depth_sub_sea - REAL NOT NULL
	// vdp_measured_depth - REAL NOT NULL
	// vdp_true_vertical_depth - REAL (nullable)
	// vdp_true_vertical_depth_sea - REAL (nullable)
	// diff_instrument_vdp - REAL (nullable)
	// density_oil - REAL NOT NULL
	// density_liquid_stopped - REAL NOT NULL
	// density_liquid_working - REAL NOT NULL
	// pressure_diff_stopped - REAL (nullable)
	// pressure_diff_working - REAL (nullable)
	// cluster_number - INTEGER (nullable)
	// instrument_number - INTEGER (nullable)

	// Добавляем валидаторы для обязательных полей (смотрим NOT NULL в SQL и не-указатели в struct)
	fieldNumberEntry.Validator = validation.NewRegexp(`^\d+$`, "Требуется число")                // Простой пример
	clusterNumberEntry.Validator = validation.NewRegexp(`^\d*$`, "Требуется число или пусто")    // Nullable
	instrumentNumberEntry.Validator = validation.NewRegexp(`^\d*$`, "Требуется число или пусто") // Nullable
	measuredDepthEntry.Validator = validation.NewRegexp(`^\d+(\.\d+)?$`, "Требуется число")      // С плавающей точкой
	startTimeEntry.Validator = validation.NewRegexp(`.+`, "Поле не может быть пустым")           // Проверка на непустоту
	endTimeEntry.Validator = validation.NewRegexp(`.+`, "Поле не может быть пустым")             // Проверка на непустоту

	// Создаем элементы формы
	formItems := []*widget.FormItem{
		widget.NewFormItem("Месторождение", fieldNameSelect), // field_name TEXT NOT NULL
		widget.NewFormItem("№ Скважины", fieldNumberEntry),   // field_number INTEGER NOT NULL
		widget.NewFormItem("№ Куста (опц.)", clusterNumberEntry),
		widget.NewFormItem("Горизонт", horizonSelect), // horizon TEXT NOT NULL
		widget.NewFormItem("Дата начала", startTimeEntry),
		widget.NewFormItem("Дата окончания", endTimeEntry),
		widget.NewFormItem("Тип прибора", instrumentTypeSelect), // instrument_type TEXT NOT NULL
		widget.NewFormItem("№ Прибора (опц.)", instrumentNumberEntry),
		widget.NewFormItem("Глубина замера (MD)", measuredDepthEntry), // measure_depth REAL NOT NULL
		widget.NewFormItem("", widget.NewLabel("Остальные параметры (TVD, Дебиты и т.д.)\nрассчитываются или пока не реализованы.")),
	}

	// 3. Создаем и показываем диалог формы
	formDialog := dialog.NewForm(
		"Заполнение Шапки Отчета (Блок 5)",
		"Сохранить",
		"Отмена",
		formItems,
		func(save bool) {
			if !save {
				s.zLog.Debugw("Block 5 form cancelled")
				return
			}
			s.zLog.Debugw("Block 5 form submitted")

			// 4. Валидация при сохранении (Fyne сделает это автоматически перед вызовом callback,
			if fieldNameSelect.Selected == "" || horizonSelect.Selected == "" || instrumentTypeSelect.Selected == "" {
				dialog.ShowInformation("Ошибка", "Пожалуйста, выберите значения из выпадающих списков.", s.window)
				return
			}

			// 5. Собираем данные из формы в models.TableFive
			var report models.TableFive

			report.FieldName = fieldNameSelect.Selected
			report.Horizon = horizonSelect.Selected
			report.InstrumentType = instrumentTypeSelect.Selected

			var convErrs []string
			var err error

			if report.FieldNumber, err = strconv.Atoi(fieldNumberEntry.Text); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("№ Скважины: %v", err))
			}
			if clusterNumberEntry.Text != "" {
				cn, err := strconv.Atoi(clusterNumberEntry.Text)
				if err != nil {
					convErrs = append(convErrs, fmt.Sprintf("№ Куста: %v", err))
				} else {
					report.ClusterNumber = &cn
				}
			}
			if instrumentNumberEntry.Text != "" {
				in, err := strconv.Atoi(instrumentNumberEntry.Text)
				if err != nil {
					convErrs = append(convErrs, fmt.Sprintf("№ Прибора: %v", err))
				} else {
					report.InstrumentNumber = &in
				}
			}
			if report.MeasuredDepth, err = strconv.ParseFloat(measuredDepthEntry.Text, 64); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Глубина замера (MD): %v", err))
			}
			if report.StartTime, err = s.converter.ParseFlexibleTime(startTimeEntry.Text); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Дата начала: %v", err))
			}
			if report.EndTime, err = s.converter.ParseFlexibleTime(endTimeEntry.Text); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Дата окончания: %v", err))
			}

			// --- Установка заглушек для NOT NULL полей, которые не заполняем/рассчитываем ---
			report.TrueVerticalDepth = 0.0       // REAL NOT NULL
			report.TrueVerticalDepthSubSea = 0.0 // REAL NOT NULL
			report.VDPMeasuredDepth = 0.0        // REAL NOT NULL
			report.DensityOil = 0.0              // REAL NOT NULL
			report.DensityLiquidStopped = 0.0    // REAL NOT NULL
			report.DensityLiquidWorking = 0.0    // REAL NOT NULL
			// Nullable поля (`*float64`, `*int`) останутся nil, если не были заданы.

			if len(convErrs) > 0 {
				dialog.ShowError(fmt.Errorf("Ошибки конвертации данных:\n%s", strings.Join(convErrs, "\n")), s.window)
				return
			}

			// 6. Делает авторасчет
			researchID := s.memBlocksStorage.GetResearchID()
			report = s.calc.CalcBlockFive(ctx, report, researchID)

			// 7. Вызываем сервис БД для сохранения
			id, err := s.db.SaveReport(report)
			if err != nil {
				s.zLog.Errorw("Failed to save report (Block 5)", "error", err)
				dialog.ShowError(fmt.Errorf("Ошибка сохранения отчета в БД: %w", err), s.window)
				return
			}

			s.zLog.Infow("Report (Block 5) saved successfully", "id", id)
			dialog.ShowInformation("Успех", fmt.Sprintf("Шапка отчета успешно сохранена (ID: %d)", id), s.window)

		},
		s.window,
	)
	formDialog.Resize(fyne.NewSize(450, 500))
	formDialog.Show()
}

// --- Методы для индикатора загрузки ---
func (s *Service) showLoadingIndicator(fileName string) {
	s.loadingLabel.SetText(fmt.Sprintf("Обработка файла: %s...", fileName))
	s.loadingLabel.Show()
	s.progressBar.Show()
	s.zLog.Debugw("Showing loading indicator", "file", fileName)
}

func (s *Service) hideLoadingIndicator(err error) {
	s.loadingLabel.Hide()
	s.progressBar.Hide()
	s.zLog.Debugw("Hiding loading indicator", "error", err)

	if err != nil {
		dialog.ShowError(fmt.Errorf("Ошибка: %w", err), s.window)
	}
}

// doTableOneImport делает парсинг TableOne, сохраняет, логирует и выводит результат.
func (s *Service) doTableOneImport(path string, cfg models.OperationConfig) {

	var finalErr error
	defer func() {
		s.hideLoadingIndicator(finalErr)
	}()

	start := time.Now()
	data, err := s.importer.ParseBlockOneFile(path, cfg)
	count := len(data)
	if err != nil {
		finalErr = err
		s.zLog.Errorw("TableOne import failed", "error", err, "duration", time.Since(start))
		return
	}
	if err2 := s.memBlocksStorage.AddBlockOneData(data); err2 != nil {
		finalErr = fmt.Errorf("ошибка сохранения: %w", err2)
		s.zLog.Errorw("TableOne save failed", "error", err2)
		return
	}

	// Если дошли сюда, ошибок не было (importErr == nil)
	elapsed := time.Since(start)
	s.zLog.Infow("TableOne import success", "count", count, "duration", elapsed)
	go func() {
		time.Sleep(100 * time.Millisecond)
		dialog.ShowInformation(
			"Готово",
			fmt.Sprintf("TableOne: %d записей импортировано за %s", count, elapsed.Round(time.Millisecond)),
			s.window,
		)
	}()
}

// doGenericImport обрабатывает TableTwo/TableThree.
func (s *Service) doGenericImport(ctx context.Context, path, typ string) {
	fileName := filepath.Base(path)
	s.showLoadingIndicator(fileName)

	// Запускаем сам импорт в горутине
	go func() {
		var finalErr error
		var count int
		start := time.Now()

		defer func() {
			s.hideLoadingIndicator(finalErr)
		}()

		switch typ {
		case "TableTwo":
			var data []models.TableTwo
			data, finalErr = s.importer.ParseBlockTwoFile(path)
			count = len(data)
			if finalErr == nil {
				finalErr = s.memBlocksStorage.AddBlockTwoData(data)
			}
		case "TableThree":
			var data []models.TableThree
			data, finalErr = s.importer.ParseBlockThreeFile(path)
			count = len(data)
			if finalErr == nil {
				finalErr = s.memBlocksStorage.AddBlockThreeData(data)
			}
		case "TableFour":
			var data []models.TableFour
			data, finalErr = s.importer.ParseBlockFourFile(path)
			count = len(data)
			if finalErr == nil {
				id, dbErr := s.db.SaveBlockFour(ctx, data)
				if dbErr == nil {
					s.memBlocksStorage.SetResearchID(id)
					finalErr = s.memBlocksStorage.AddBlockFourData(data)
				} else {
					finalErr = dbErr
				}
			}
		}

		if finalErr != nil {
			s.zLog.Errorw(typ+" import failed", "error", finalErr, "duration", time.Since(start))
			return
		}

		// Успех
		elapsed := time.Since(start)
		s.zLog.Infow(typ+" import success", "count", count, "duration", elapsed)

		// Показываем сообщение об успехе ПОСЛЕ скрытия индикатора
		// Запускаем это тоже в горутине, чтобы не блокировать выход из текущей
		go func() {
			time.Sleep(100 * time.Millisecond)
			dialog.ShowInformation(
				"Готово",
				fmt.Sprintf("%s: %d записей импортировано за %s", typ, count, elapsed.Round(time.Millisecond)),
				s.window,
			)
		}()
	}()
}
