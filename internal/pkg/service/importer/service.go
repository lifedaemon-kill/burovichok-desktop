package importer

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/thedatashed/xlsxreader"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// Service отвечает за логику импорта данных из Excel.
type Service struct{}

// NewService создает новый экземпляр сервис импорта.
func NewService() *Service {
	return &Service{}
}

// ParseBlockOneFile читает XLSX‑файл через xlsxreader и возвращает []TableOne.
func (s *Service) ParseBlockOneFile(path string) ([]models.TableOne, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", path)
	}
	xl, err := xlsxreader.NewReader(data)
	if err != nil {
		return nil, errors.Wrap(err, "xlsxreader.NewReader")
	}

	var out []models.TableOne
	for row := range xl.ReadRows(xl.Sheets[0]) {
		if row.Index == 1 {
			continue
		}
		cells := row.Cells
		if len(cells) < 3 {
			continue
		}
		ts, err := parseFlexibleTime(cells[0].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse timestamp block1 row %d", row.Index)
		}
		pres, err := strconv.ParseFloat(cells[1].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse pressure block1 row %d", row.Index)
		}
		temp, err := strconv.ParseFloat(cells[2].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse temperature block1 row %d", row.Index)
		}
		out = append(out, models.TableOne{Timestamp: ts, PressureDepth: pres, Temperature: temp})
	}
	return out, nil
}

// ParseBlockTwoFile читает XLSX‑файл и возвращает []TableTwo.
func (s *Service) ParseBlockTwoFile(path string) ([]models.TableTwo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", path)
	}
	xl, err := xlsxreader.NewReader(data)
	if err != nil {
		return nil, errors.Wrap(err, "xlsxreader.NewReader")
	}

	var out []models.TableTwo
	for row := range xl.ReadRows(xl.Sheets[0]) {
		if row.Index <= 2 {
			continue
		}
		cells := row.Cells
		if len(cells) < 6 {
			continue
		}
		// Tubing pressure
		tsTub, err := parseFlexibleTime(cells[0].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse tubing timestamp row %d", row.Index)
		}
		presTub, err := strconv.ParseFloat(cells[1].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse tubing pressure row %d", row.Index)
		}
		// Annulus pressure
		tsAnn, err := parseFlexibleTime(cells[2].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse annulus timestamp row %d", row.Index)
		}
		presAnn, err := strconv.ParseFloat(cells[3].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse annulus pressure row %d", row.Index)
		}
		// Linear pressure
		tsLin, err := parseFlexibleTime(cells[4].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse linear timestamp row %d", row.Index)
		}
		presLin, err := strconv.ParseFloat(cells[5].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse linear pressure row %d", row.Index)
		}

		out = append(out, models.TableTwo{
			TimestampTubing:  tsTub,
			PressureTubing:   presTub,
			TimestampAnnulus: tsAnn,
			PressureAnnulus:  presAnn,
			TimestampLinear:  tsLin,
			PressureLinear:   presLin,
		})
	}
	return out, nil
}

// ParseBlockThreeFile читает XLSX‑файл и возвращает []TableThree.
func (s *Service) ParseBlockThreeFile(path string) ([]models.TableThree, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", path)
	}
	xl, err := xlsxreader.NewReader(data)
	if err != nil {
		return nil, errors.Wrap(err, "xlsxreader.NewReader")
	}

	var out []models.TableThree
	for row := range xl.ReadRows(xl.Sheets[0]) {
		// пропускаем заголовок
		if row.Index == 1 {
			continue
		}
		cells := row.Cells
		if len(cells) < 4 {
			continue
		}
		ts, err := parseFlexibleTime(cells[0].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse timestamp block3 row %d", row.Index)
		}
		flowL, err := strconv.ParseFloat(cells[1].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse flow liquid row %d", row.Index)
		}
		wc, err := strconv.ParseFloat(cells[2].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse water cut row %d", row.Index)
		}
		flowG, err := strconv.ParseFloat(cells[3].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse flow gas row %d", row.Index)
		}

		out = append(out, models.TableThree{
			Timestamp:  ts,
			FlowLiquid: flowL,
			WaterCut:   wc,
			FlowGas:    flowG,
		})
	}
	return out, nil
}

// parseFlexibleTime пытается разобрать время как Excel‑serial или ISO/RFC строки.
func parseFlexibleTime(raw string) (time.Time, error) {
	if num, err := strconv.ParseFloat(raw, 64); err == nil {
		return excelDateToTime(num, false)
	}

	layouts := []string{
		time.RFC3339,          // 2024-11-09T17:21:21Z
		"2006-01-02T15:04:05", // 2024-11-09T17:21:21
		"2006-01-02 15:04:05", // 2024-11-09 17:21:21
		"2006-01-02",          // 2024-11-10
		"02/01/2006 15:04:05", // 09/11/2024 17:21:21
		"02/01/2006",          // 09/11/2024
		"02.01.2006 15:04:05", // 09.11.2024 17:21:21
		"02.01.2006",          // 09.11.2024
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.Errorf("unsupported time format %q", raw)
}

// excelDateToTime конвертирует serial date Excel в time.Time.
func excelDateToTime(serial float64, date1904 bool) (time.Time, error) {
	var epoch time.Time
	if date1904 {
		epoch = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		epoch = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	}
	days := math.Floor(serial)
	if !date1904 && days >= 61 {
		days--
	}
	frac := serial - math.Floor(serial)
	d := epoch.Add(time.Duration(days) * 24 * time.Hour)
	return d.Add(time.Duration(frac * 24 * float64(time.Hour))), nil
}
