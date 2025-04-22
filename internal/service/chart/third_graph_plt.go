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

func generateEchartsTableThreeData(data []models.TableThree) ([][]opts.ScatterData, []string) {
	yLabels := make([][]opts.ScatterData, 6)

	xLabels := make([]string, 0, len(data))

	for _, point := range data {
		yLabels[0] = append(yLabels[0], opts.ScatterData{Value: point.LiquidFlowRate, Name: point.Timestamp.Format(time.RFC3339)})
		yLabels[1] = append(yLabels[1], opts.ScatterData{Value: point.OilFlowRate, Name: point.Timestamp.Format(time.RFC3339)})
		yLabels[2] = append(yLabels[2], opts.ScatterData{Value: point.WaterFlowRate, Name: point.Timestamp.Format(time.RFC3339)})

		yLabels[3] = append(yLabels[3], opts.ScatterData{Value: point.WaterCut, Name: point.Timestamp.Format(time.RFC3339)})
		yLabels[4] = append(yLabels[4], opts.ScatterData{Value: point.GasFlowRate, Name: point.Timestamp.Format(time.RFC3339)})
		yLabels[5] = append(yLabels[5], opts.ScatterData{Value: point.GasFactor, Name: point.Timestamp.Format(time.RFC3339)})

		xLabels = append(xLabels, point.Timestamp.Format("02.01.02 15:04")) // Только время для краткости оси X
	}

	return yLabels, xLabels
}
func (s *chartService) GenerateTableThreeChart(data []models.TableThree) (string, error) {
	if len(data) == 0 {
		return "", errors.Wrap(errors.New("Нет данных, для построения графика"), "GenerateTableOneChart")
	}

	line := charts.NewScatter()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "",
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
			Name: "Дебит (м²/сут)",
			Type: "value",
		}),

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
					Name:  "third_block_chart",
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

	tableOneData, xLabels := generateEchartsTableThreeData(data)

	line.SetXAxis(xLabels).
		AddSeries("Дебит жидкости", tableOneData[0], charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"})).
		AddSeries("Дебит нефти", tableOneData[1], charts.WithLineStyleOpts(opts.LineStyle{Color: "green"})).
		AddSeries("Дебит воды", tableOneData[2], charts.WithLineStyleOpts(opts.LineStyle{Color: "red"})).
		AddSeries("Обводненность", tableOneData[3], charts.WithLineStyleOpts(opts.LineStyle{Color: "cyan"})).
		AddSeries("Дебит газа", tableOneData[4], charts.WithLineStyleOpts(opts.LineStyle{Color: "grey"})).
		AddSeries("Газовый фактор", tableOneData[5], charts.WithLineStyleOpts(opts.LineStyle{Color: "brown"})).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(false)}),
		)

	f, err := os.Create(HTMLFileNameThree)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл %s: %w", HTMLFileNameThree, err)
	}
	defer f.Close()

	err = line.Render(f)
	if err != nil {
		return "", fmt.Errorf("не удалось отрендерить график в файл: %w", err)
	}

	return HTMLFileNameThree, nil
}
