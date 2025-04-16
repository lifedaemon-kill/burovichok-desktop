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

// Importer умеет парсить файлы блоков.
type Importer interface {
	ParseBlockOneFile(path string) ([]models.BlockOne, error)
}

// Service отвечает за инициализацию и запуск UI приложения.
type Service struct {
	app      fyne.App
	window   fyne.Window
	zLog     logger.Logger
	importer Importer
}

// NewService создаёт новый UI‑сервис с указанным заголовком и размерами окна.
func NewService(title string, width, height int, zLog logger.Logger, importer Importer) *Service {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	return &Service{app: a, window: w, zLog: zLog, importer: importer}
}

// Run создаёт окно с яркой кнопкой и меткой и запускает приложение.
func (s *Service) Run() error {
	// Простая метка и кнопка, как в примере Fyne Test
	label := widget.NewLabel("Нажмите кнопку, чтобы загрузить .xlsx файл")
	button := widget.NewButton("Загрузить файл", func() {
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
				dialog.ShowError(err, s.window)
				return
			}

			// Пока просто печатаем данные в консоль
			fmt.Println(rows)
		}, s.window)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlg.Show()
	})

	// Компоновка: вертикальный бокс с меткой и кнопкой
	content := container.NewVBox(
		label,
		button,
	)

	s.window.SetContent(content)
	s.window.ShowAndRun()
	return nil
}
