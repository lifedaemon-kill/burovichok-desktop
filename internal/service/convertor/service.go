package convertor

import (
	"math"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ParseFlexibleTime пытается разобрать время как Excel‑serial или ISO/RFC строки.
func (s *Service) ParseFlexibleTime(raw string) (time.Time, error) {
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
