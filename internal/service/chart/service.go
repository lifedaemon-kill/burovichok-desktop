package chart

import (
	"fmt"
	"os" 
	"time"

	"github.com/go-echarts/go-echarts/v2/charts" 
	"github.com/go-echarts/go-echarts/v2/opts"   

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

const (

	chartHTMLFilename = "burovichok_chart.html"
)

type Service interface {
	// GeneratePressureTempChart генерирует HTML файл графика и возвращает путь к нему
	GeneratePressureTempChart(data []models.TableOne) (string, error)
}

type chartService struct{}

func NewService() Service {
	return &chartService{}
}

func generateEchartsData(data []models.TableOne) ([]opts.LineData, []string) {
	items := make([]opts.LineData, 0, len(data))
	xLabels := make([]string, 0, len(data))
	for _, point := range data {
		items = append(items, opts.LineData{Value: point.PressureDepth, Name: point.Timestamp.Format(time.RFC3339)})
		xLabels = append(xLabels, point.Timestamp.Format("15:04:05")) // Только время для краткости оси X
	}
	return items, xLabels
}

func generateTempEchartsData(data []models.TableOne) []opts.LineData {
	items := make([]opts.LineData, 0, len(data))
	for _, point := range data {
		items = append(items, opts.LineData{Value: point.TemperatureDepth, Name: point.Timestamp.Format(time.RFC3339)})
	}
	return items
}

func (s *chartService) GeneratePressureTempChart(data []models.TableOne) (string, error) {
	if len(data) == 0 { 
		return "", fmt.Errorf("нет данных для построения графика")
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Давление и Температура (Блок 1)",
			Subtitle: "Интерактивный график",
		}),
		charts.WithTooltipOpts(opts.Tooltip{ // Всплывающие подсказки при наведении
			Show:      opts.Bool(true),
			Trigger:   "axis", // Показывает данные для всех линий на оси X
			TriggerOn: "mousemove|click",
		}),
		charts.WithXAxisOpts(opts.XAxis{ // Настройка оси X
			Name: "Время",
			Type: "category", // Ось категорий (наши метки времени)
		}),
		charts.WithYAxisOpts(opts.YAxis{ // Основная ось Y (Давление)
			Name: "Давление (кгс/см2)",
			Type: "value",
		}),
		// Важно: go-echarts напрямую не поддерживает вторую Y-ось так просто, как go-chart/v2.
		// Обычно вторую ось добавляют через opts.Grid и позиционирование, либо
		// создают второй график и совмещают. Пока сделаем только давление для простоты.
		// TODO: Разобраться с добавлением второй оси Y для температуры в go-echarts.
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}), // Показываем легенду
		charts.WithDataZoomOpts(opts.DataZoom{ 
			Type:       "slider", 
			Start:      0,        
			End:        100,      
			XAxisIndex: []int{0}, 
		}),
		charts.WithToolboxOpts(opts.Toolbox{ 
			Show: opts.Bool(true),
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{ 
					Show:  opts.Bool(true),
					Type:  "png",
					Name:  "pressure_chart",
					Title: "Сохранить PNG",
				},
				DataZoom: &opts.ToolBoxFeatureDataZoom{ 
					Show:  opts.Bool(true),
					Title: map[string]string{"zoom": "Зум", "back": "Сброс"},
				},
				Restore: &opts.ToolBoxFeatureRestore{ 
					Show:  opts.Bool(true),
					Title: "Сброс",
				},
			},
		}),
		// Можно выбрать тему оформления
		// charts.WithTheme(types.ThemeInfographic),
	)

	pressureData, xLabels := generateEchartsData(data)

	line.SetXAxis(xLabels).
		AddSeries("Давление", pressureData, charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"})). 
		// AddSeries("Температура", tempData, charts.WithLineStyleOpts(opts.LineStyle{Color:"red"})). // Пока без температуры
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}), 
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(false)}),           
		)

	f, err := os.Create(chartHTMLFilename)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл %s: %w", chartHTMLFilename, err)
	}
	defer f.Close()

	err = line.Render(f)
	if err != nil {
		return "", fmt.Errorf("не удалось отрендерить график в файл: %w", err)
	}

	return chartHTMLFilename, nil
}
