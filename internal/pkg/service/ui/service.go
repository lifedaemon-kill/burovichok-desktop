// internal/pkg/service/ui/ui.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Service отвечает за инициализацию и запуск UI приложения.
type Service struct {
	App    fyne.App
	Window fyne.Window
}

// NewService создает новый UI-сервис с указанным заголовком и размерами окна.
func NewService(title string, width, height int) *Service {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	// Здесь можно задать начальное содержимое.
	content := container.NewVBox(widget.NewLabel("hello, burovichok!"))
	w.SetContent(content)
	return &Service{
		App:    a,
		Window: w,
	}
}

// Run запускает окно приложения.
// Этот метод блокирующий, т.е. выполнение продолжится после закрытия окна.
func (s *Service) Run() {
	s.Window.ShowAndRun()
}
