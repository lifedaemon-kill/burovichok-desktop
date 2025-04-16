package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/models"
)

// Importer умеет парсить файлы блоков.
type Importer interface {
	ParseBlockOneFile(path string) ([]models.BlockOne, error)
}

// Service отвечает за UI.
type Service struct {
	app      fyne.App
	window   fyne.Window
	zLog     logger.Logger
	importer Importer
}

// NewService создаёт сервис UI.
func NewService(title string, width, height int, zLog logger.Logger, importer Importer) *Service {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	return &Service{app: a, window: w, zLog: zLog, importer: importer}
}

// Run строит окно, рисует фон и кликабельный квадратик, запускает UI.
func (s *Service) Run() error {
	// Белый фон всего окна
	bg := canvas.NewRectangle(color.White)
	bg.Move(fyne.NewPos(0, 0))
	bg.Resize(s.window.Canvas().Size())

	// Кликабельный прямоугольник
	rect := canvas.NewRectangle(color.NRGBA{R: 0, G: 122, B: 204, A: 255})
	rect.StrokeWidth = 2
	rect.StrokeColor = color.Black
	rect.Move(fyne.NewPos(31, 86))
	rect.Resize(fyne.NewSize(721, 677))

	// Невидимая кнопка для обработки клика
	btn := widget.NewButton("", func() {
		dlg := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			if r == nil {
				return
			}
			defer r.Close()
			path := r.URI().Path()
			rows, err := s.importer.ParseBlockOneFile(path)
			if err != nil {
				s.zLog.Errorw("failed to parse file", "error", err)
				return
			}
			fmt.Println(rows)
		}, s.window)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlg.Show()
	})
	// Задаём размер и положение кнопки вместо SetMinSize
	btn.Resize(fyne.NewSize(721, 677))
	btn.Move(fyne.NewPos(31, 86))

	// Контейнер для абсолютного позиционирования
	content := container.NewWithoutLayout(bg, rect, btn)

	s.window.SetContent(content)
	s.window.ShowAndRun()
	return nil
}
