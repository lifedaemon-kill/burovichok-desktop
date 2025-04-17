package ui

import (
	"fmt"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	inmemory "github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
)

// ratioLayout располагает два объекта в контейнере в пропорции ratio к (1-ratio).
type ratioLayout struct{ ratio float32 }

// Layout вычисляет размеры и позиции дочерних элементов.
func (r *ratioLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	firstWidth := size.Width * r.ratio
	objects[0].Resize(fyne.NewSize(firstWidth, size.Height))
	objects[1].Resize(fyne.NewSize(size.Width-firstWidth, size.Height))
	objects[1].Move(fyne.NewPos(firstWidth, 0))
}

// MinSize возвращает минимальный размер контейнера по высоте самого "высокого" элемента.
func (r *ratioLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var height float32
	for _, o := range objects {
		h := o.MinSize().Height
		if h > height {
			height = h
		}
	}
	return fyne.NewSize(0, height)
}

// Importer умеет парсить три типа блоков.
type Importer interface {
	ParseBlockOneFile(path string) ([]models.BlockOne, error)
	ParseBlockTwoFile(path string) ([]models.BlockTwo, error)
	ParseBlockThreeFile(path string) ([]models.BlockThree, error)
}

// Service отвечает за инициализацию и запуск UI приложения.
type Service struct {
	app                 fyne.App
	window              fyne.Window
	zLog                logger.Logger
	importer            Importer
	memBlocksStorage    inmemory.BlocksStorage
	blocksRepository    sqlite.BlocksStorage
	guidebookRepository sqlite.GuidebooksStorage
}

// NewService создаёт новый UI‑сервис с заголовком и размерами окна.
func NewService(
	ui config.UI,
	zLog logger.Logger,
	importer Importer,
	memBlock inmemory.BlocksStorage,
	blocksRepository sqlite.BlocksStorage,
	guidebooks sqlite.GuidebooksStorage) *Service {
	a := app.New()
	w := a.NewWindow(ui.Name)
	w.Resize(fyne.NewSize(float32(ui.Width), float32(ui.Height)))
	return &Service{
		app:                 a,
		window:              w,
		zLog:                zLog,
		importer:            importer,
		memBlocksStorage:    memBlock,
		blocksRepository:    blocksRepository,
		guidebookRepository: guidebooks}
}

// Run строит интерфейс с тремя контролами и запускает приложение.
func (s *Service) Run() error {
	// Поле и кнопка выбора файла
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Файл не выбран")
	chooseBtn := widget.NewButton("Выбрать файл", func() {
		dlg := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r == nil || err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			defer r.Close()
			pathEntry.SetText(r.URI().Path())
		}, s.window)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlg.Show()
	})

	// Выбор типа документа
	docTypes := []string{"BlockOne", "BlockTwo", "BlockThree"}
	typeSelect := widget.NewSelect(docTypes, func(string) {})
	typeSelect.PlaceHolder = "Выберите тип документа"

	// Import кнопка
	importBtn := widget.NewButton("Import", func() {
		path := pathEntry.Text
		docType := typeSelect.Selected
		if path == "" || docType == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите файл и тип документа", s.window)
			return
		}

		// Засекаем время
		start := time.Now()

		s.zLog.Infow("Start import", "path", path, "type", docType)
		var count int
		var err error
		var storeErr error

		switch docType {
		case "BlockOne":
			data, parseErr := s.importer.ParseBlockOneFile(path)
			if parseErr != nil {
				err = parseErr // Приоритет у ошибки парсинга
			} else {
				count = len(data)
				storeErr = s.memBlocksStorage.AddBlockOneData(data) // <-- Сохраняем данные в хранилище
			}
		case "BlockTwo":
			data, parseErr := s.importer.ParseBlockTwoFile(path)
			if parseErr != nil {
				err = parseErr
			} else {
				count = len(data)
				storeErr = s.memBlocksStorage.AddBlockTwoData(data) // <-- Сохраняем данные в хранилище
			}
		case "BlockThree":
			data, parseErr := s.importer.ParseBlockThreeFile(path)
			if parseErr != nil {
				err = parseErr
			} else {
				count = len(data)
				storeErr = s.memBlocksStorage.AddBlockThreeData(data) // <-- Сохраняем данные в хранилище
			}
		}
		// Вычисление затраченного времени
		elapsed := time.Since(start)

		if err != nil {
			s.zLog.Errorw("Import failed (parsing)", "error", err, "duration", elapsed)
			dialog.ShowError(err, s.window)
			return
		}

		if storeErr != nil {
			// Ошибка сохранения - это скорее внутренняя проблема
			s.zLog.Errorw("Import failed (storing)", "error", storeErr, "duration", elapsed)
			dialog.ShowError(fmt.Errorf("внутренняя ошибка при сохранении данных: %w", storeErr), s.window)
			return
		}

		// Получаем текущее общее количество записей в хранилище (опционально, для информации)
		totalCount := 0
		switch docType {
		case "BlockOne":
			totalCount = s.memBlocksStorage.CountBlockOne()
		case "BlockTwo":
			totalCount = s.memBlocksStorage.CountBlockTwo()
		case "BlockThree":
			totalCount = s.memBlocksStorage.CountBlockThree()
		}

		// Информируем пользователя
		s.zLog.Infow("Import successful", "type", docType, "count", count, "total_in_store", totalCount, "duration", elapsed)
		dialog.ShowInformation(
			"Готово",
			fmt.Sprintf("Импортировано %d записей типа '%s'.\nВсего в памяти: %d.\nВремя: %s", count, docType, totalCount, elapsed.Round(time.Millisecond)),
			s.window,
		)
	})

	// <-- НОВАЯ КНОПКА: Очистка хранилища -->
	clearBtn := widget.NewButton("Очистить хранилище", func() {
		// Запрос подтверждения
		dialog.ShowConfirm("Подтверждение", "Вы уверены, что хотите удалить все загруженные данные из памяти?", func(confirm bool) {
			if !confirm {
				return
			}
			err := s.memBlocksStorage.ClearAll()
			if err != nil {
				// Маловероятно для in-memory, но на всякий случай
				s.zLog.Errorw("Failed to clear memBlocksStorage", "error", err)
				dialog.ShowError(fmt.Errorf("ошибка при очистке хранилища: %w", err), s.window)
				return
			}
			s.zLog.Infow("In-memory memBlocksStorage cleared by user")
			dialog.ShowInformation("Хранилище очищено", "Все данные в памяти удалены.", s.window)

		}, s.window)
	})

	// Компоновка: 70% для поля и 30% для кнопки
	head := container.New(&ratioLayout{ratio: 0.7}, pathEntry, chooseBtn)

	// Собираем всё вместе
	content := container.NewVBox(
		head,
		typeSelect,
		importBtn,
		widget.NewSeparator(), // <-- Разделитель для красоты
		clearBtn,              // <-- Добавляем кнопку очистки
	)

	s.window.SetContent(content)
	s.window.ShowAndRun()
	s.zLog.Infow("UI service stopped") // Это сообщение появится после закрытия окна
	return nil
}
