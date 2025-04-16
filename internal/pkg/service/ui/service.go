package ui

import (
	"bytes"
	"image/png"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/cockroachdb/errors"
	"github.com/wcharczuk/go-chart/v2"
)

// Service отвечает за инициализацию и запуск UI приложения.
type Service struct {
	App    fyne.App
	Window fyne.Window
}

// NewService создает новый UI‑сервис с указанным заголовком и размерами окна.
func NewService(title string, width, height int) *Service {
	a := app.New()
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	return &Service{
		App:    a,
		Window: w,
	}
}

// GenerateData возвращает пример данных (синусоиду) для отображения графика.
func GenerateData(points int) ([]float64, []float64) {
	x := make([]float64, points)
	y := make([]float64, points)
	for i := 0; i < points; i++ {
		val := float64(i) * 2 * math.Pi / float64(points-1)
		x[i] = val
		y[i] = math.Sin(val)
	}
	return x, y
}

// RenderChart создаёт изображение графика в виде canvas.Image по переданным данным.
func RenderChart(xvals, yvals []float64, width, height int) (*canvas.Image, error) {
	graph := chart.Chart{
		Width:  width,
		Height: height,
		Series: []chart.Series{
			chart.ContinuousSeries{XValues: xvals, YValues: yvals},
		},
	}

	var buf bytes.Buffer
	if err := graph.Render(chart.PNG, &buf); err != nil {
		return nil, errors.Wrap(err, "graph.Render")
	}

	imgData, err := png.Decode(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "png.Decode")
	}

	img := canvas.NewImageFromImage(imgData)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(float32(width), float32(height)))

	return img, nil
}

// Run инициализирует данные, рендерит график и запускает UI. Возвращает ошибку при неудаче.
func (s *Service) Run() error {
	// Генерация данных
	x, y := GenerateData(100)
	// Рендеринг графика
	img, err := RenderChart(x, y, 600, 300)
	if err != nil {
		return errors.Wrap(err, "RenderChart")
	}
	// Установка содержимого окна
	s.Window.SetContent(container.NewCenter(img))
	// Запуск UI (блокирующий)
	s.Window.ShowAndRun()

	return nil
}
