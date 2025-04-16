package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/models"
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

// MinSize возвращает минимальный размер контейнера, равный высоте самого "высокого" элемента.
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

// Importer умеет парсить два типа блоков.
type Importer interface {
	ParseBlockOneFile(path string) ([]models.BlockOne, error)
	ParseBlockTwoFile(path string) ([]models.BlockTwo, error)
}

// Service отвечает за инициализацию и запуск UI приложения.
type Service struct {
	app      fyne.App
	window   fyne.Window
	zLog     logger.Logger
	importer Importer
}

// NewService создаёт новый UI‑сервис с заголовком и размерами окна.
func NewService(title string, width, height int, zLog logger.Logger, importer Importer) *Service {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	return &Service{app: a, window: w, zLog: zLog, importer: importer}
}

// Run строит интерфейс с тремя контролами и запускает приложение.
func (s *Service) Run() error {
	// 1) Поле для пути и кнопка выбора
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Файл не выбран")

	chooseBtn := widget.NewButton("Выбрать файл", func() {
		dlg := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
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
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlg.Show()
	})

	// Компоновка поля и кнопки с соотношением 70/30
	head := container.New(&ratioLayout{ratio: 0.7}, pathEntry, chooseBtn)

	// 2) Выпадающий список типов документов
	docTypes := []string{"BlockOne", "BlockTwo"}
	typeSelect := widget.NewSelect(docTypes, func(string) {})
	typeSelect.PlaceHolder = "Выберите тип документа"

	// 3) Кнопка Import
	importBtn := widget.NewButton("Import", func() {
		path := pathEntry.Text
		docType := typeSelect.Selected
		if path == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите файл", s.window)
			return
		}
		if docType == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите тип документа", s.window)
			return
		}

		s.zLog.Infow("Start import", "path", path, "type", docType)
		var count int
		var err error
		switch docType {
		case "BlockOne":
			data, e := s.importer.ParseBlockOneFile(path)
			err = e
			count = len(data)
		case "BlockTwo":
			data, e := s.importer.ParseBlockTwoFile(path)
			err = e
			count = len(data)
		}
		if err != nil {
			s.zLog.Errorw("Import failed", "error", err)
			dialog.ShowError(err, s.window)
			return
		}
		dialog.ShowInformation("Завершено", fmt.Sprintf("Импортировано %d записей", count), s.window)
	})

	// Сборка всего контента
	content := container.NewVBox(
		head,
		typeSelect,
		importBtn,
	)

	s.window.SetContent(content)
	s.window.ShowAndRun()
	return nil
}
