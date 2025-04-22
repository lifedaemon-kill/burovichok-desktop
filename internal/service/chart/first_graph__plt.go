package chart

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"os"
	"time"
)

func generateEchartsTableOneData(data []models.TableOne) ([][]opts.LineData, []string) {
	yLabels := make([][]opts.LineData, 3)

	xLabels := make([]string, 0, len(data))

	for _, point := range data {
		yLabels[0] = append(yLabels[0], opts.LineData{Value: point.PressureDepth, Name: point.Timestamp.Format(time.RFC3339)})
		yLabels[1] = append(yLabels[1], opts.LineData{Value: point.PressureAtVDP, Name: point.Timestamp.Format(time.RFC3339)})
		yLabels[2] = append(yLabels[2], opts.LineData{Value: point.TemperatureDepth, Name: point.Timestamp.Format(time.RFC3339)})

		xLabels = append(xLabels, point.Timestamp.Format("02.01.02 15:04")) // Только время для краткости оси X
	}

	return yLabels, xLabels
}
func (s *chartService) GenerateTableOneChart(data []models.TableOne) (string, error) {
	if len(data) == 0 {
		return "", errors.Wrap(errors.New("Нет данных, для построения графика"), "GenerateTableOneChart")
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Забойное давление и температура",
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
			Type:       "inside",
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
					Name:  "first_block_chart",
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
		//charts.WithTheme(types.ThemeInfographic),
	)

	tableOneData, xLabels := generateEchartsTableOneData(data)

	line.SetXAxis(xLabels).
		AddSeries("Рзаб на глубине", tableOneData[0], charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"})).
		AddSeries("Рзаб на ВДП", tableOneData[1], charts.WithLineStyleOpts(opts.LineStyle{Color: "green"})).
		AddSeries("Tзаб на глубине", tableOneData[2], charts.WithLineStyleOpts(opts.LineStyle{Color: "red"})).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}),
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(false)}),
		)

	f, err := os.Create(chartHTMLFilenameOne)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл %s: %w", chartHTMLFilenameOne, err)
	}
	defer f.Close()

	err = line.Render(f)
	if err != nil {
		return "", fmt.Errorf("не удалось отрендерить график в файл: %w", err)
	}

	return chartHTMLFilenameOne, nil
}
