package importer

import (
	"math"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/models"
)

// Service отвечает за логику импорта данных из Excel.
type Service struct{}

// NewService создает новый экземпляр сервис импорта.
func NewService() *Service {
	return &Service{}
}

// ParseBlockOneFile читает XLSX-файл и возвращает срез записей.
func (s *Service) ParseBlockOneFile(path string) ([]models.BlockOne, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	const sheet = "Инструментальный замер"
	rows := f.GetRows(sheet)

	var result []models.BlockOne
	for i, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}

		// Обрабатываем Timestamp: может быть числом или строкой
		tsRaw := row[0]
		var ts time.Time
		if num, errPars := strconv.ParseFloat(tsRaw, 64); errPars == nil {
			// Преобразуем serial date в time.Time
			ts, errPars = excelDateToTime(num, false)
			if errPars != nil {
				return nil, errors.Wrapf(errPars, "convert excel date on row %d", i+2)
			}
		} else {
			// Ожидаем формат DD/MM/YYYY HH:MM:SS
			//TODO должно быть 5 различных форматирований времени
			ts, err = time.Parse("02/01/2006 15:04:05", tsRaw)
			if err != nil {
				return nil, errors.Wrapf(err, "parse timestamp on row %d", i+2)
			}
		}

		// Парсим давление
		pres, errParse := strconv.ParseFloat(row[1], 64)
		if errParse != nil {
			return nil, errors.Wrapf(errParse, "parse pressure on row %d", i+2)
		}

		// Парсим температуру
		temp, errParse := strconv.ParseFloat(row[2], 64)
		if errParse != nil {
			return nil, errors.Wrapf(errParse, "parse temperature on row %d", i+2)
		}

		result = append(result, models.BlockOne{
			Timestamp:   ts,
			Pressure:    pres,
			Temperature: temp,
		})
	}
	return result, nil
}

// ParseBlockTwoFile читает XLSX с листом «Инструментальный замер» (блок 2)
// и возвращает срез BlockTwo.
func (s *Service) ParseBlockTwoFile(path string) ([]models.BlockTwo, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	const sheet = "Инструментальный замер"
	rows := f.GetRows(sheet)

	var out []models.BlockTwo
	// Пропускаем первые две строки заголовков
	for i, row := range rows[2:] {
		if len(row) < 6 {
			continue
		}
		// A,B: трубное давление
		tsTub, err := parseExcelDateOrString(row[0])
		if err != nil {
			return nil, errors.Wrapf(err, "parse tubing timestamp on row %d", i+3)
		}
		presTub, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse tubing pressure on row %d", i+3)
		}

		// C,D: затрубное давление
		tsAnn, err := parseExcelDateOrString(row[2])
		if err != nil {
			return nil, errors.Wrapf(err, "parse annulus timestamp on row %d", i+3)
		}
		presAnn, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse annulus pressure on row %d", i+3)
		}

		// E,F: линейное давление
		tsLin, err := parseExcelDateOrString(row[4])
		if err != nil {
			return nil, errors.Wrapf(err, "parse linear timestamp on row %d", i+3)
		}
		presLin, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse linear pressure on row %d", i+3)
		}

		out = append(out, models.BlockTwo{
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

// ParseBlockThreeFile читает дебиты: дата/время, Qж, W, Qг.
func (s *Service) ParseBlockThreeFile(path string) ([]models.BlockThree, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	const sheet = "Инструментальный замер"
	rows := f.GetRows(sheet)

	var out []models.BlockThree
	for i, row := range rows[1:] {
		if len(row) < 4 {
			continue
		}
		// Время
		ts, err := parseExcelDateOrString(row[0])
		if err != nil {
			return nil, errors.Wrapf(err, "parse timestamp on row %d", i+2)
		}
		// Qж
		flowL, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse flow liquid on row %d", i+2)
		}
		// W
		wc, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse water cut on row %d", i+2)
		}
		// Qг
		flowG, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse flow gas on row %d", i+2)
		}

		out = append(out, models.BlockThree{
			Timestamp:  ts,
			FlowLiquid: flowL,
			WaterCut:   wc,
			FlowGas:    flowG,
		})
	}
	return out, nil
}

// parseExcelDateOrString умеет разобрать и serial‑дату Excel (число) и строковый формат DD/MM/YYYY HH:MM:SS.
func parseExcelDateOrString(raw string) (time.Time, error) {
	// попытаемся сначала как число
	if num, err := strconv.ParseFloat(raw, 64); err == nil {
		return excelSerialToTime(num, false)
	}
	// иначе ожидаем строку вида “09/11/2024 17:21:21”
	return time.Parse("02/01/2006 15:04:05", raw)
}

// excelSerialToTime конвертирует serial date Excel в time.Time.
func excelSerialToTime(serial float64, date1904 bool) (time.Time, error) {
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

// excelDateToTime конвертирует Excel serial date в time.Time.
// date1904=false означает систему 1900 (Windows Excel).
func excelDateToTime(serial float64, date1904 bool) (time.Time, error) {
	// Опорная дата
	var epoch time.Time
	if date1904 {
		epoch = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		epoch = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	}

	// Коррекция для ложного 29 февраля 1900
	days := math.Floor(serial)
	if !date1904 && days >= 61 {
		days -= 1
	}
	// Целая и дробная часть
	frac := serial - math.Floor(serial)

	// Добавляем дни и долю дня
	d := epoch.Add(time.Duration(days) * 24 * time.Hour)
	t := d.Add(time.Duration(frac * 24 * float64(time.Hour)))
	return t, nil
}
