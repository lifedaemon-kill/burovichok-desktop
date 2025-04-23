package chart

import (
	"github.com/cockroachdb/errors"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"os"
	"slices"
	"time"
)

func generateEchartsTableTwoData(data []models.TableTwo) ([][]opts.LineData, []string) {
	yLabels := make([][]opts.LineData, 3)

	uniqueXLabels := make(map[string]struct{})

	for _, point := range data {
		yLabels[0] = append(yLabels[0], opts.LineData{Value: point.PressureAnnulus, Name: point.TimestampAnnulus.Format(time.RFC3339)})
		yLabels[1] = append(yLabels[1], opts.LineData{Value: point.PressureTubing, Name: point.TimestampTubing.Format(time.RFC3339)})
		yLabels[2] = append(yLabels[2], opts.LineData{Value: point.PressureLinear, Name: point.TimestampLinear.Format(time.RFC3339)})

		uniqueXLabels[point.TimestampAnnulus.Format("02.01.02 15:04")] = struct{}{}
		uniqueXLabels[point.TimestampTubing.Format("02.01.02 15:04")] = struct{}{}
		uniqueXLabels[point.TimestampLinear.Format("02.01.02 15:04")] = struct{}{}
	}
	uniqueSlice := make([]string, 0, len(uniqueXLabels))
	for key := range uniqueXLabels {
		uniqueSlice = append(uniqueSlice, key)
	}
	slices.Sort(uniqueSlice)

	return yLabels, uniqueSlice
}

func (s *chartService) GenerateTableTwoChart(data []models.TableTwo, units string) (string, error) {
	if len(data) == 0 {
		return "", errors.Wrap(errors.New("Нет данных, для построения графика"), "GenerateTableTwoChart")
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Устьевое давление и температура",
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
			Name: "Давление " + units,
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
					Name:  "second_block_chart",
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

	tableOneData, xLabels := generateEchartsTableTwoData(data)

	line.SetXAxis(xLabels).
		AddSeries("Ртр", tableOneData[0], charts.WithLineStyleOpts(opts.LineStyle{Color: "yellow"})).
		AddSeries("Pзтр", tableOneData[1], charts.WithLineStyleOpts(opts.LineStyle{Color: "black"})).
		AddSeries("Pлин", tableOneData[2], charts.WithLineStyleOpts(opts.LineStyle{Color: "brown"})).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}),
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(false)}),
		)

	f, err := os.Create(HTMLFileNameTwo)
	if err != nil {
		return "", errors.Wrap(err, "не удалось создать файл "+HTMLFileNameOne+"_second_block.html")
	}
	defer f.Close()

	err = line.Render(f)
	if err != nil {
		return "", errors.Wrap(err, "не удалось отрендерить график в файл "+
		HTMLFileNameTwo+"_second_block.html")
	}

	return HTMLFileNameTwo, nil
}
