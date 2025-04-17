package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	appStorage "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/storage"
)

// ratioLayout располагает два объекта в контейнере в пропорции ratio к (1-ratio).
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

// Importer умеет парсить три типа блоков.
type Importer interface {
	ParseBlockOneFile(path string, cfg models.OperationConfig) ([]models.TableOne, error)
	ParseBlockTwoFile(path string) ([]models.TableTwo, error)
	ParseBlockThreeFile(path string) ([]models.TableThree, error)
}

// Service отвечает за инициализацию и запуск UI приложения.
type Service struct {
	app      fyne.App
	window   fyne.Window
	zLog     logger.Logger
	importer Importer
	store    appStorage.Storage
}

// NewService создаёт новый UI‑сервис.
func NewService(title string, w, h int, zLog logger.Logger, imp Importer, store appStorage.Storage) *Service {
	a := app.New()
	win := a.NewWindow(title)
	win.Resize(fyne.NewSize(float32(w), float32(h)))
	return &Service{app: a, window: win, zLog: zLog, importer: imp, store: store}
}

// Run строит интерфейс и запускает приложение.
func (s *Service) Run() error {
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
	docTypes := []string{"TableOne", "TableTwo", "TableThree"}
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
			if err := s.store.ClearAll(); err != nil {
				s.zLog.Errorw("Clear failed", "error", err)
				dialog.ShowError(fmt.Errorf("ошибка очистки: %w", err), s.window)
				return
			}
			dialog.ShowInformation("Очищено", "Данные удалены", s.window)
		}, s.window)
	})

	// компоновка
	header := container.New(&ratioLayout{ratio: 0.7}, pathEntry, chooseBtn)
	content := container.NewVBox(
		header,
		typeSelect,
		importBtn,
		widget.NewSeparator(),
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

			workStart, err := parseFlexibleTime(ws.Text)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprintf("Работа: c — %v", err))
			}
			workEnd, err := parseFlexibleTime(we.Text)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprintf("Работа: по — %v", err))
			}

			// для необязательных полей простоя парсим только если не пусто
			var idleStart, idleEnd time.Time
			if strings.TrimSpace(is.Text) != "" {
				idleStart, err = parseFlexibleTime(is.Text)
				if err != nil {
					parseErrs = append(parseErrs, fmt.Sprintf("Простой: c — %v", err))
				}
			}
			if strings.TrimSpace(ie.Text) != "" {
				idleEnd, err = parseFlexibleTime(ie.Text)
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
	if err2 := s.store.AddBlockOneData(data); err2 != nil {
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
			err = s.store.AddBlockTwoData(data)
		}
	case "TableThree":
		var data []models.TableThree
		data, err = s.importer.ParseBlockThreeFile(path)
		count = len(data)
		if err == nil {
			err = s.store.AddBlockThreeData(data)
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

// parseFlexibleTime пытается разобрать время как Excel‑serial или ISO/RFC строки.
func parseFlexibleTime(raw string) (time.Time, error) {
	if num, err := strconv.ParseFloat(raw, 64); err == nil {
		return excelDateToTime(num, false)
	}

	layouts := []string{
		time.RFC3339,          // 2024-11-09T17:21:21Z
		"2006-01-02T15:04:05", // 2024-11-09T17:21:21
		"2006-01-02 15:04:05", // 2024-11-09 17:21:21
		"2006-01-02",          // 2024-11-10
		"02/01/2006 15:04:05", // 09/11/2024 17:21:21
		"02/01/2006",          // 09/11/2024
		"02.01.2006 15:04:05", // 09.11.2024 17:21:21
		"02.01.2006",          // 09.11.2024
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.Errorf("unsupported time format %q", raw)
}

// excelDateToTime конвертирует serial date Excel в time.Time.
func excelDateToTime(serial float64, date1904 bool) (time.Time, error) {
	var epoch time.Time
	if date1904 {
		epoch = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		epoch = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	}
	days := math.Floor(serial)
	if !date1904 && days >= 61 {
		days--
	}
	frac := serial - math.Floor(serial)
	d := epoch.Add(time.Duration(days) * 24 * time.Hour)
	return d.Add(time.Duration(frac * 24 * float64(time.Hour))), nil
}
