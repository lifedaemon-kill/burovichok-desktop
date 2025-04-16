package importer

import (
	"math"
	"strconv"
	"time"

	excelize "github.com/360EntSecGroup-Skylar/excelize"
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
