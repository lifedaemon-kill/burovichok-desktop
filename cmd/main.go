package main

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/logger/z"
	"image/color"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"
)

const configPath = "config/config.yaml"

func main() {
	ctx := context.Background()
	err := bootstrap(ctx)
	if err != nil {
		log.Fatalf("main bootstrap failed: %v", err)
	}
}

// Функция для генерации случайных данных
func generateData(size int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(100) // Генерируем случайные числа от 0 до 99
	}
	return data
}

// Функция для создания графика (прямоугольники)
func createPlot(data []int, index int) *canvas.Rectangle {
	// Создаем прямоугольник, представляющий значение
	height := float32(data[index%len(data)])                // Высота прямоугольника
	rect := canvas.NewRectangle(color.RGBA{0, 0, 255, 255}) // Синий цвет
	rect.SetMinSize(fyne.NewSize(50, height))               // Ширина 50, высота по значению
	return rect
}

func bootstrap(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	//SET UP
	conf := config.Load(configPath)

	if err := z.InitLogger(conf.Logger); err != nil {
		log.Fatalf("init logger error: %s", err)
	}
	defer z.Log.Sync()
	z.Log.Info("init logger and config success")

	//APP
	burApp := app.New()
	myWindow := burApp.NewWindow("burovichok")
	// Генерируем случайные данные
	data := generateData(20) // Изменил на 20 для примера

	// Создаем контейнер для графиков
	graphContainer := container.NewVBox()

	// Добавляем графики в контейнер
	for i := 0; i < 3; i++ {
		graph := createPlot(data, i)
		graphContainer.Add(graph)
	}

	// Добавляем кнопку для обновления данных
	updateButton := widget.NewButton("Обновить данные", func() {
		data = generateData(20)    // Генерируем новые данные
		graphContainer.RemoveAll() // Очищаем контейнер
		for i := 0; i < 3; i++ {
			graph := createPlot(data, i)
			graphContainer.Add(graph) // Добавляем новые графики
		}
		graphContainer.Refresh() // Обновляем контейнер
	})

	// Создаем основной контейнер
	content := container.NewVBox(graphContainer, updateButton)
	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()

	//SHUTTING DOWN
	//<-ctx.Done()
	z.Log.Info("context done, shutting down...")

	//	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	defer shutdownCancel()

	//TODO завершить sqlite
	z.Log.Info("shutting down completed")
	return nil
}
