package ui

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/service/database"

	chartService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/chart"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	inmemoryStorage "github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
	"github.com/pkg/browser"
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
	db               database.Service
	converter        converterService
	chart            chartService.Service

	serverMutex      sync.Mutex
	serverListener   net.Listener
	serverPort       string
	chartHtmlToServe string
}

func NewService(title string, w, h int, zLog logger.Logger, imp importer, converter converterService,
	memBlocksStorage inmemoryStorage.InMemoryBlocksStorage, db database.Service, chart chartService.Service) *Service {
	a := app.New()
	// chartSvc := chartService.NewService() // Сервис теперь передается снаружи
	win := a.NewWindow(title)
	win.Resize(fyne.NewSize(float32(w), float32(h)))
	return &Service{
		app:              a,
		window:           win,
		zLog:             zLog,
		importer:         imp,
		memBlocksStorage: memBlocksStorage,
		db:               db,
		converter:        converter,
		chart:            chart,
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

func (s *Service) Run() error {
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
			if r == nil || err != nil {
				dialog.ShowError(err, s.window)
				return
			}
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
			// сразу импортим во всех остальных случаях
			go s.doGenericImport(path, typ)
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

	header := container.New(&ratioLayout{ratio: 0.7}, pathEntry, chooseBtn)
	content := container.NewVBox(
		widget.NewLabel("1. Выберите файл и тип данных:"),
		header,
		typeSelect,
		widget.NewSeparator(),
		widget.NewLabel("2. Выполните действия:"),
		importBtn,
		chartBtn,
		widget.NewSeparator(),
		widget.NewLabel("3. Управление хранилищем:"),
		clearBtn,
	)

	s.window.SetContent(content)
	s.window.ShowAndRun()
	return nil
}

// showTableOneForm показывает форму параметров гидростатики и по клику "Ок" запускает импорт.
func showTableOneForm(s *Service, path string) {
	// поля формы
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

			// Проверка обязательных полей
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

			// Теперь парсим время и собираем ошибки
			var parseErrs []string

			workStart, err := s.converter.ParseFlexibleTime(ws.Text)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprintf("Работа: c — %v", err))
			}
			workEnd, err := s.converter.ParseFlexibleTime(we.Text)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprintf("Работа: по — %v", err))
			}

			// для необязательных полей простоя парсим только если не пусто
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

			// Если были ошибки парсинга — показываем и выходим
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

			go s.doTableOneImport(path, cfg)
		}, s.window)

	dlg.Show()
}

// doTableOneImport делает парсинг TableOne, сохраняет, логирует и выводит результат.
func (s *Service) doTableOneImport(path string, cfg models.OperationConfig) {
	start := time.Now()
	data, err := s.importer.ParseBlockOneFile(path, cfg)
	count := len(data)
	if err != nil {
		s.zLog.Errorw("TableOne import failed", "error", err, "duration", time.Since(start))
		dialog.ShowError(err, s.window)
		return
	}
	if err2 := s.memBlocksStorage.AddBlockOneData(data); err2 != nil {
		s.zLog.Errorw("TableOne save failed", "error", err2)
		dialog.ShowError(fmt.Errorf("ошибка сохранения: %w", err2), s.window)
		return
	}
	elapsed := time.Since(start)
	s.zLog.Infow("TableOne import success", "count", count, "duration", elapsed)
	dialog.ShowInformation(
		"Готово",
		fmt.Sprintf("TableOne: %d записей импортировано за %s", count, elapsed.Round(time.Millisecond)),
		s.window,
	)
}

// doGenericImport обрабатывает TableTwo/TableThree.
func (s *Service) doGenericImport(path, typ string) {
	start := time.Now()
	var count int
	var err error

	switch typ {
	case "TableTwo":
		var data []models.TableTwo
		data, err = s.importer.ParseBlockTwoFile(path)
		count = len(data)
		if err == nil {
			err = s.memBlocksStorage.AddBlockTwoData(data)
		}
	case "TableThree":
		var data []models.TableThree
		data, err = s.importer.ParseBlockThreeFile(path)
		count = len(data)
		if err == nil {
			err = s.memBlocksStorage.AddBlockThreeData(data)
		}

	case "TableFour":
		// Инклинометрия
		var data []models.TableFour
		data, err = s.importer.ParseBlockFourFile(path)
		count = len(data)
		if err == nil {
			// Реализуйте в BlocksStorage метод AddBlockFourData
			err = s.memBlocksStorage.AddBlockFourData(data)
		}
	}

	if err != nil {
		s.zLog.Errorw(typ+" import failed", "error", err, "duration", time.Since(start))
		dialog.ShowError(err, s.window)
		return
	}
	elapsed := time.Since(start)
	s.zLog.Infow(typ+" import success", "count", count, "duration", elapsed)
	dialog.ShowInformation(
		"Готово",
		fmt.Sprintf("%s: %d записей импортировано за %s", typ, count, elapsed.Round(time.Millisecond)),
		s.window,
	)
}
