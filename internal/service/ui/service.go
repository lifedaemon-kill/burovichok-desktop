package ui

import (
	"context"
	"errors"
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
	"github.com/google/uuid"
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
	memBlocksStorage inmemoryStorage.InMemoryBlocksStorage, db *database.Service, chart chartService.Service, calcSvc *calc.Service) *Service {

	a := app.New()
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
		calc:             calcSvc,
		loadingLabel:     loadingLbl,
		progressBar:      progressBr,
	}
}

// --- веб‑сервер для графиков ---

func (s *Service) startLocalWebServer() error {
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()
	if s.serverListener != nil {
		return nil
	}

	s.zLog.Infow("Starting local chart web server...")
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("не удалось найти порт для сервера графика: %w", err)
	}
	s.serverListener = ln
	s.serverPort = strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serveChartHTML)

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
	html := s.chartHtmlToServe
	s.serverMutex.Unlock()

	if html == "" {
		http.Error(w, "График еще не сгенерирован", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, html)
}

// --- Навигация между разделами ---

func (s *Service) showMainMenu(ctx context.Context) {
	cell := fyne.NewSize(400, 200)

	importsBtn := widget.NewButton("импорты", func() { s.showImportView(ctx) })
	reportsBtn := widget.NewButton("отчёты", func() { s.showReportsView(ctx) })
	chartsBtn := widget.NewButton("графики", func() { s.showChartsView(ctx) })
	guidebooksBtn := widget.NewButton("Справочники", func() { s.showGuidebookView(ctx) })

	// NewGridWrap принимает размер ячейки — и «упаковывает» каждый элемент в box этого размера
	grid := container.NewGridWrap(cell,
		importsBtn,
		reportsBtn,
		chartsBtn,
		guidebooksBtn,
	)

	s.window.SetContent(container.NewCenter(grid))
}

func (s *Service) showImportView(ctx context.Context) {
	back := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })
	content := s.buildImportContent(ctx)
	s.window.SetContent(container.NewBorder(back, nil, nil, nil, content))
}

func (s *Service) showReportsView(ctx context.Context) {
	back := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })
	// Отчёты: только блок 5
	reportBtn := widget.NewButton("Заполнить Шапку Отчета (Блок 5)", func() {
		s.showBlockFiveForm(ctx)
	})
	s.window.SetContent(container.NewBorder(back, nil, nil, nil,
		container.NewVBox(
			widget.NewLabel("Отчёты"),
			widget.NewSeparator(),
			reportBtn,
		),
	))
}

func (s *Service) showChartsView(ctx context.Context) {
	back := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })
	chartBtn := widget.NewButton("Интерактивный График Pзаб/Тзаб (Блок 1)", func() {
		// вставьте сюда тот же код из Run для chartBtn
		blockOneData, err := s.memBlocksStorage.GetAllBlockOneData()
		if err != nil {
			dialog.ShowError(fmt.Errorf("не удалось получить данные Блока 1: %w", err), s.window)
			return
		}
		if len(blockOneData) == 0 {
			dialog.ShowInformation("Нет данных", "Недостаточно данных для построения графика", s.window)
			return
		}
		htmlPath, err := s.chart.GenerateTableOneChart(blockOneData)
		if err != nil {
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
		time.Sleep(100 * time.Millisecond)
		s.serverMutex.Lock()
		port := s.serverPort
		s.serverMutex.Unlock()
		if port == "" {
			dialog.ShowError(fmt.Errorf("не удалось получить порт веб-сервера"), s.window)
			return
		}
		url := fmt.Sprintf("http://127.0.0.1:%s/", port)
		if err := browser.OpenURL(url); err != nil {
			dialog.ShowError(fmt.Errorf("не удалось открыть браузер: %w", err), s.window)
		}
	})

	s.window.SetContent(container.NewBorder(back, nil, nil, nil,
		container.NewVBox(
			widget.NewLabel("Графики"),
			widget.NewSeparator(),
			chartBtn,
		),
	))
}

// --- НОВАЯ ФУНКЦИЯ: Показ экрана управления справочниками ---
func (s *Service) showGuidebookView(ctx context.Context) {
	s.zLog.Debugw("Opening Guidebook Management view")

	backBtn := widget.NewButton("◀ Домой", func() { s.showMainMenu(ctx) })

	// Определяем типы справочников, которые можно редактировать
	guidebookTypes := []string{
		"Месторождение",         // -> oilfield
		"Продуктивный горизонт", // -> productive_horizon
		"Тип прибора",           // -> instrument_type
		"Вид исследования",      // -> research_type
	}
	guidebookTypeSelect := widget.NewSelect(guidebookTypes, nil)
	guidebookTypeSelect.PlaceHolder = "Выберите тип справочника"

	newValueEntry := widget.NewEntry()
	newValueEntry.PlaceHolder = "Введите новое значение"
	newValueEntry.Validator = validation.NewRegexp(`.+`, "Поле не может быть пустым")

	statusLabel := widget.NewLabel("")

	addBtn := widget.NewButton("Добавить значение", func() {
		statusLabel.SetText("")
		selectedType := guidebookTypeSelect.Selected
		newValue := strings.TrimSpace(newValueEntry.Text)

		if selectedType == "" || newValue == "" {
			dialog.ShowInformation("Ошибка", "Выберите тип справочника и введите непустое значение.", s.window)
			return
		}

		err := newValueEntry.Validate()
		if err != nil {
			statusLabel.SetText("Ошибка валидации: " + err.Error())
			return
		}

		err = s.addGuidebookEntry(ctx, selectedType, newValue)
		if err != nil {
			// Обрабатываем возможную ошибку уникальности (если уберем ON CONFLICT) или другую ошибку БД
			// В случае ON CONFLICT DO NOTHING ошибки не будет при дубликате
			s.zLog.Errorw("Failed to add guidebook entry", "type", selectedType, "value", newValue, "error", err)
			dialog.ShowError(fmt.Errorf("Не удалось добавить значение: %w", err), s.window)
			statusLabel.SetText("Ошибка при добавлении.")
		} else {
			s.zLog.Infow("Successfully added guidebook entry", "type", selectedType, "value", newValue)
			statusLabel.SetText(fmt.Sprintf("Значение '%s' добавлено в '%s'.", newValue, selectedType))
			newValueEntry.SetText("")
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Управление справочниками"),
		widget.NewSeparator(),
		guidebookTypeSelect,
		newValueEntry,
		addBtn,
		statusLabel,
	)

	s.window.SetContent(container.NewBorder(backBtn, nil, nil, nil, content))
}

// --- НОВАЯ ФУНКЦИЯ-ОБЕРТКА для вызова правильного метода БД ---
func (s *Service) addGuidebookEntry(ctx context.Context, guidebookType string, name string) error {
	switch guidebookType {
	case "Месторождение":
		return s.db.SaveOilFields(ctx, []models.OilField{{Name: name}})
	case "Продуктивный горизонт":
		return s.db.SaveProductiveHorizons(ctx, []models.ProductiveHorizon{{Name: name}})
	case "Тип прибора":
		return s.db.SaveInstrumentTypes(ctx, []models.InstrumentType{{Name: name}})
	case "Вид исследования":
		return s.db.SaveResearchTypes(ctx, []models.ResearchType{{Name: name}})
	default:
		return errors.New("неизвестный тип справочника")
	}
}

func (s *Service) Run(ctx context.Context) error {
	s.window.SetOnClosed(func() {
		s.serverMutex.Lock()
		if s.serverListener != nil {
			s.serverListener.Close()
			s.serverListener = nil
			s.serverPort = ""
		}
		s.serverMutex.Unlock()
	})

	s.showMainMenu(ctx)
	s.window.ShowAndRun()
	return nil
}

// --- Построение содержимого импортов ---

func (s *Service) buildImportContent(ctx context.Context) fyne.CanvasObject {
	// 1) путь
	pathEntry := widget.NewEntry()
	pathEntry.PlaceHolder = "Файл не выбран"
	chooseBtn := widget.NewButton("Выбрать файл", func() {
		d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			if r == nil {
				return
			}
			defer r.Close()
			pathEntry.SetText(r.URI().Path())
		}, s.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		d.Show()
	})

	// 2) тип документа
	docTypes := []string{"TableOne", "TableTwo", "TableThree", "TableFour"}
	typeSelect := widget.NewSelect(docTypes, nil)
	typeSelect.PlaceHolder = "Выберите тип документа"

	// 3) Import
	importBtn := widget.NewButton("Import", func() {
		path := pathEntry.Text
		typ := typeSelect.Selected
		if path == "" || typ == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите файл и тип документа", s.window)
			return
		}
		if typ == "TableOne" {
			showTableOneForm(s, path)
		} else {
			s.doGenericImport(ctx, path, typ)
		}
	})

	// 4) Очистка хранилища
	clearBtn := widget.NewButton("Очистить хранилище", func() {
		dialog.ShowConfirm("Подтверждение", "Удалить все данные?", func(ok bool) {
			if !ok {
				return
			}
			if err := s.memBlocksStorage.ClearAll(); err != nil {
				dialog.ShowError(fmt.Errorf("ошибка очистки: %w", err), s.window)
				return
			}
			dialog.ShowInformation("Очищено", "Данные удалены", s.window)
		}, s.window)
	})

	// Собираем всё в VBox
	return container.NewVBox(
		widget.NewLabel("Импорт данных"),
		widget.NewSeparator(),
		widget.NewLabel("1. Выберите файл и тип:"),
		container.New(&ratioLayout{ratio: 0.7}, pathEntry, chooseBtn),
		typeSelect,
		widget.NewSeparator(),
		widget.NewLabel("2. Действия:"),
		importBtn,
		clearBtn,
		widget.NewSeparator(),
		s.loadingLabel,
		s.progressBar,
	)
}

// --- Далее ваши существующие формы и методы ---

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
	unit.SetSelected("kgf/cm2")

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
	oilFieldsModels, err := s.db.GetAllOilFields(ctx)
	if err != nil {
		s.zLog.Errorw("Failed to get oil fields", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки месторождений: %w", err), s.window)
		return
	}
	horizonsModels, err := s.db.GetAllProductiveHorizons(ctx)
	if err != nil {
		s.zLog.Errorw("Failed to get horizons", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки горизонтов: %w", err), s.window)
		return
	}
	instrumentTypesModels, err := s.db.GetAllInstrumentTypes(ctx)
	if err != nil {
		s.zLog.Errorw("Failed to get instrument types", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки типов приборов: %w", err), s.window)
		return
	}

	researchTypesModels, err := s.db.GetAllResearchTypes(ctx)
	if err != nil {
		s.zLog.Errorw("Failed to get research types", "error", err)
		dialog.ShowError(fmt.Errorf("Ошибка загрузки видов исследований: %w", err), s.window)

		researchTypesModels = []models.ResearchType{}
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

	researchTypeNames := make([]string, len(researchTypesModels))
	for i, item := range researchTypesModels {
		researchTypeNames[i] = item.Name
	}

	// 2. Создаем виджеты для формы
	// Уже существующие:
	fieldNameSelect := widget.NewSelect(oilFields, nil)
	fieldNumberEntry := widget.NewEntry()
	clusterNumberEntry := widget.NewEntry()
	horizonSelect := widget.NewSelect(horizons, nil)
	startTimeEntry := widget.NewEntry()
	endTimeEntry := widget.NewEntry()
	instrumentTypeSelect := widget.NewSelect(instrumentTypes, nil)
	instrumentNumberEntry := widget.NewEntry()
	measuredDepthEntry := widget.NewEntry()
	vdpMeasuredDepthEntry := widget.NewEntry()
	densityOilEntry := widget.NewEntry()
	densityLiquidStoppedEntry := widget.NewEntry()
	densityLiquidWorkingEntry := widget.NewEntry()

	researchTypeSelect := widget.NewSelect(researchTypeNames, nil)
	researchTypeSelect.PlaceHolder = "Выберите вид исследования"

	// 3. Добавляем валидаторы
	fieldNumberEntry.Validator = validation.NewRegexp(`^\d+$`, "Требуется число")
	clusterNumberEntry.Validator = validation.NewRegexp(`^\d*$`, "Требуется число или пусто")
	instrumentNumberEntry.Validator = validation.NewRegexp(`^\d*$`, "Требуется число или пусто")
	measuredDepthEntry.Validator = validation.NewRegexp(`^\d+(\.\d+)?$`, "Требуется число")
	startTimeEntry.Validator = validation.NewRegexp(`.+`, "Поле не может быть пустым")
	endTimeEntry.Validator = validation.NewRegexp(`.+`, "Поле не может быть пустым")

	// Валидация для новых полей:
	vdpMeasuredDepthEntry.Validator = validation.NewRegexp(`^\d+(\.\d+)?$`, "Требуется число")
	densityOilEntry.Validator = validation.NewRegexp(`^\d+(\.\d+)?$`, "Требуется число")
	densityLiquidStoppedEntry.Validator = validation.NewRegexp(`^\d+(\.\d+)?$`, "Требуется число")
	densityLiquidWorkingEntry.Validator = validation.NewRegexp(`^\d+(\.\d+)?$`, "Требуется число")

	// 4. Формируем items
	formItems := []*widget.FormItem{
		widget.NewFormItem("Вид исследования", researchTypeSelect),
		widget.NewFormItem("Месторождение", fieldNameSelect),
		widget.NewFormItem("№ Скважины", fieldNumberEntry),
		widget.NewFormItem("№ Куста (опц.)", clusterNumberEntry),
		widget.NewFormItem("Горизонт", horizonSelect),
		widget.NewFormItem("Дата начала", startTimeEntry),
		widget.NewFormItem("Дата окончания", endTimeEntry),
		widget.NewFormItem("Тип прибора", instrumentTypeSelect),
		widget.NewFormItem("№ Прибора (опц.)", instrumentNumberEntry),
		widget.NewFormItem("Глубина замера (MD)", measuredDepthEntry),

		// Новые поля:
		widget.NewFormItem("MD ВДП (VDPMeasuredDepth)", vdpMeasuredDepthEntry),
		widget.NewFormItem("Плотность нефти (kg/m³)", densityOilEntry),
		widget.NewFormItem("Плотность жидкости в простое (kg/m³)", densityLiquidStoppedEntry),
		widget.NewFormItem("Плотность жидкости в работе (kg/m³)", densityLiquidWorkingEntry),

		widget.NewFormItem("", widget.NewLabel("Остальные параметры (TVD, ΔP и т.д.) рассчитываются автоматически.")),
	}

	// 5. Показ диалога и обработчик Save
	formDialog := dialog.NewForm(
		"Заполнение Шапки Отчета (Блок 5)",
		"Сохранить", "Отмена",
		formItems,
		func(save bool) {
			if !save {
				s.zLog.Debugw("Block 5 form cancelled")
				return
			}
			// Проверка обязательных селектов
			if fieldNameSelect.Selected == "" || horizonSelect.Selected == "" || instrumentTypeSelect.Selected == "" || researchTypeSelect.Selected == "" {
				dialog.ShowInformation("Ошибка", "Пожалуйста, выберите все необходимые значения.", s.window)
				return
			}

			var report models.TableFive
			var convErrs []string
			var err error

			report.FieldName = fieldNameSelect.Selected
			report.Horizon = horizonSelect.Selected
			report.InstrumentType = instrumentTypeSelect.Selected
			report.ResearchType = researchTypeSelect.Selected

			// Парсинг существующих
			if report.FieldNumber, err = strconv.Atoi(fieldNumberEntry.Text); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("№ Скважины: %v", err))
			}
			if cn := clusterNumberEntry.Text; cn != "" {
				if report.ClusterNumber, err = strconv.Atoi(cn); err != nil {
					convErrs = append(convErrs, fmt.Sprintf("№ Куста: %v", err))
				}
			}
			if report.MeasuredDepth, err = strconv.ParseFloat(measuredDepthEntry.Text, 64); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("MD: %v", err))
			}
			if report.StartTime, err = s.converter.ParseFlexibleTime(startTimeEntry.Text); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Дата начала: %v", err))
			}
			if report.EndTime, err = s.converter.ParseFlexibleTime(endTimeEntry.Text); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Дата окончания: %v", err))
			}
			if inTxt := instrumentNumberEntry.Text; inTxt != "" {
				if report.InstrumentNumber, err = strconv.Atoi(inTxt); err != nil {
					convErrs = append(convErrs, fmt.Sprintf("№ Прибора: %v", err))
				}
			}

			// Парсинг новых вручную:
			if report.VDPMeasuredDepth, err = strconv.ParseFloat(vdpMeasuredDepthEntry.Text, 64); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("MD ВДП: %v", err))
			}
			if report.DensityOil, err = strconv.ParseFloat(densityOilEntry.Text, 64); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Плотность нефти: %v", err))
			}
			if report.DensityLiquidStopped, err = strconv.ParseFloat(densityLiquidStoppedEntry.Text, 64); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Плотность жидкости в простое: %v", err))
			}
			if report.DensityLiquidWorking, err = strconv.ParseFloat(densityLiquidWorkingEntry.Text, 64); err != nil {
				convErrs = append(convErrs, fmt.Sprintf("Плотность жидкости в работе: %v", err))
			}

			if len(convErrs) > 0 {
				dialog.ShowError(fmt.Errorf("Ошибки конвертации:\n%s", strings.Join(convErrs, "\n")), s.window)
				return
			}

			// Расчет и сохранение остаётся без изменений
			researchID := s.memBlocksStorage.GetResearchID()
			if researchID == uuid.Nil {
				s.zLog.Errorw("Failed to Get ResearchID", "error", err)
				dialog.ShowError(fmt.Errorf("вы не импортировали блок 4 для отчетов"), s.window)
				return
			}

			report = s.calc.CalcBlockFive(ctx, report, researchID)

			id, err := s.db.SaveReport(ctx, report)
			if err != nil {
				s.zLog.Errorw("Failed to save report (Block 5)", "error", err)
				dialog.ShowError(fmt.Errorf("ошибка сохранения: %w", err), s.window)
				return
			}
			dialog.ShowInformation("Успех", fmt.Sprintf("ID отчёта: %d", id), s.window)
		},
		s.window,
	)
	formDialog.Resize(fyne.NewSize(500, 600))
	formDialog.Show()
}

// --- Методы для индикатора загрузки ---
func (s *Service) showLoadingIndicator(fileName string) {
	s.zLog.Debugw("Showing loading indicator", "file", fileName)
	s.loadingLabel.SetText(fmt.Sprintf("Обработка файла: %s...", fileName))
	s.loadingLabel.Show()
	s.progressBar.Show()
}

func (s *Service) hideLoadingIndicator(err error) {
	s.zLog.Debugw("Hiding loading indicator", "error", err)
	s.loadingLabel.Hide()
	s.progressBar.Hide()

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
