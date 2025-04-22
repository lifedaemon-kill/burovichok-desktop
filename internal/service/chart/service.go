package chart

import (
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

const (
	chartHTMLFilenameOne   = "burovichok_first_chart.html"
	chartHTMLFilenameTwo   = "burovichok_second_chart.html"
	chartHTMLFilenameThree = "burovichok_third_chart.html"
)

type Service interface {
	// GenerateTableOneChart генерирует HTML файл графика и возвращает путь к нему
	GenerateTableOneChart(data []models.TableOne) (string, error)
	GenerateTableTwoChart(data []models.TableTwo, units string) (string, error)
	GenerateTableThreeChart(data []models.TableThree) (string, error)
}

type chartService struct{}

func NewService() Service {
	return &chartService{}
}

func generateTempEchartsData(data []models.TableOne) []opts.LineData {
	items := make([]opts.LineData, 0, len(data))
	for _, point := range data {
		items = append(items, opts.LineData{Value: point.TemperatureDepth, Name: point.Timestamp.Format(time.RFC3339)})
	}
	return items
}
