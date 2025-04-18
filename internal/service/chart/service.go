package chart

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/wcharczuk/go-chart/v2" // Импортируем go-chart
)

type Service interface {
	GeneratePressureTempChart(data []models.TableOne) (image.Image, error)
}

type chartService struct{}

func NewService() Service {
	return &chartService{}
}

func (s *chartService) GeneratePressureTempChart(data []models.TableOne) (image.Image, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("недостаточно данных для построения графика (нужно минимум 2 точки)")
	}

	timeValues := make([]time.Time, len(data))
	pressureValues := make([]float64, len(data))
	tempValues := make([]float64, len(data))

	for i, point := range data {
		timeValues[i] = point.Timestamp
		pressureValues[i] = point.PressureDepth
		tempValues[i] = point.TemperatureDepth
	}

	pressureSeries := chart.TimeSeries{
		Name:    "Давление (кгс/см2)",
		XValues: timeValues,
		YValues: pressureValues,
		Style: chart.Style{
			StrokeColor: chart.ColorBlue,
			FillColor:   chart.ColorBlue.WithAlpha(50),
		},
		YAxis: chart.YAxisPrimary,
	}

	tempSeries := chart.TimeSeries{
		Name:    "Температура (°C)",
		XValues: timeValues,
		YValues: tempValues,
		Style: chart.Style{
			StrokeColor: chart.ColorRed,
		},
		YAxis: chart.YAxisSecondary,
	}

	graph := chart.Chart{
		Title: "Давление и Температура от Времени (Блок 1)",
		XAxis: chart.XAxis{
			Name: "Время",
			// Можно настроить формат тиков, если нужно
			// ValueFormatter: chart.TimeValueFormatterWithFormat("02.01 15:04"),
		},
		YAxis: chart.YAxis{
			Name: "Давление (кгс/см2)",
			// Можно настроить сетку, диапазон и т.д.
			// GridMajorStyle: chart.Style{ StrokeColor: chart.ColorAlternateGray, StrokeWidth: 1.0 },
		},
		YAxisSecondary: chart.YAxis{
			Name: "Температура (°C)",
		},
		Series: []chart.Series{
			pressureSeries,
			tempSeries,
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer) // Рендерим в PNG
	if err != nil {
		return nil, fmt.Errorf("ошибка рендеринга графика: %w", err)
	}

	img, _, err := image.Decode(buffer)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения графика: %w", err)
	}

	return img, nil
}
