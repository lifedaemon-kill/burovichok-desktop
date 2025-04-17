package importer

import (
	"os"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/thedatashed/xlsxreader"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

type calcService interface {
	CalcTableOne(rec models.TableOne, cfg models.OperationConfig) models.TableOne
}

type converterService interface {
	ParseFlexibleTime(raw string) (time.Time, error)
}

// Service отвечает за логику импорта данных из Excel.
type Service struct {
	calc      calcService
	converter converterService
}

// NewService создает новый экземпляр сервис импорта.
func NewService(calc calcService, converter converterService) *Service {
	return &Service{
		calc:      calc,
		converter: converter,
	}
}

// ParseBlockOneFile читает XLSX‑файл через xlsxreader и возвращает []TableOne.
func (s *Service) ParseBlockOneFile(path string, cfg models.OperationConfig) ([]models.TableOne, error) {
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
		ts, err := s.converter.ParseFlexibleTime(cells[0].Value)
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

		rec := models.TableOne{
			Timestamp:     ts,
			PressureDepth: pres,
			Temperature:   temp,
		}

		// 4) автоматический расчёт ВДП
		rec = s.calc.CalcTableOne(rec, cfg)
		out = append(out, rec)
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
		tsTub, err := s.converter.ParseFlexibleTime(cells[0].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse tubing timestamp row %d", row.Index)
		}
		presTub, err := strconv.ParseFloat(cells[1].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse tubing pressure row %d", row.Index)
		}
		// Annulus pressure
		tsAnn, err := s.converter.ParseFlexibleTime(cells[2].Value)
		if err != nil {
			return nil, errors.Wrapf(err, "parse annulus timestamp row %d", row.Index)
		}
		presAnn, err := strconv.ParseFloat(cells[3].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse annulus pressure row %d", row.Index)
		}
		// Linear pressure
		tsLin, err := s.converter.ParseFlexibleTime(cells[4].Value)
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
		ts, err := s.converter.ParseFlexibleTime(cells[0].Value)
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

// ParseBlockFourFile читает XLSX‑файл и возвращает []Inclinometry.
func (s *Service) ParseBlockFourFile(path string) ([]models.Inclinometry, error) {
	// 1) читаем файл
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %s", path)
	}

	// 2) создаём ридер xlsxreader
	xl, err := xlsxreader.NewReader(data)
	if err != nil {
		return nil, errors.Wrap(err, "xlsxreader.NewReader")
	}

	var out []models.Inclinometry

	// 3) перебираем все строки на первом листе
	for row := range xl.ReadRows(xl.Sheets[0]) {
		// пропускаем заголовки и строку с единицами (Excel: строки 1–4)
		if row.Index <= 4 {
			continue
		}

		cells := row.Cells
		if len(cells) < 3 {
			continue
		}

		// 4) парсим MeasuredDepth (колонка A / cells[0])
		md, err := strconv.ParseFloat(cells[0].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse MeasuredDepth block4 row %d", row.Index)
		}

		// 5) парсим TrueVerticalDepth (колонка I / cells[1])
		tvd, err := strconv.ParseFloat(cells[1].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse TrueVerticalDepth block4 row %d", row.Index)
		}

		// 6) парсим TrueVerticalDepthSubSea (колонка J / cells[2])
		tvdss, err := strconv.ParseFloat(cells[2].Value, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse TrueVerticalDepthSubSea block4 row %d", row.Index)
		}

		// 7) собираем запись и добавляем в выходной срез
		rec := models.Inclinometry{
			MeasuredDepth:           md,
			TrueVerticalDepth:       tvd,
			TrueVerticalDepthSubSea: tvdss,
		}
		out = append(out, rec)
	}

	return out, nil
}
