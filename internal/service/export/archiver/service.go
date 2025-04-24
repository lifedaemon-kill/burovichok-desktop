package archiver

import (
	"archive/zip"
	"bytes"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
	"github.com/xuri/excelize/v2"
)

type Archiver interface {
	// Archive собирает данные t1-t5 и возвращает ZIP-архив в буфере
	Archive(
		t1 []models.TableOne,
		t2 []models.TableTwo,
		t3 []models.TableThree,
		t4 []models.TableFour,
		t5 models.TableFive,
	) (*bytes.Buffer, error)
}

// Service реализует интерфейс Archiver
// service реализует интерфейс Archiver
type service struct {
	memStorage inmemory.InMemoryBlocksStorage // Зависимость от интерфейса хранилища (если нужна в будущем)
	log        logger.Logger
}

// NewService создает новый экземпляр сервиса Archiver
func NewService(
	log logger.Logger,
	// memStorage inmemory.InMemoryBlocksStorage, // Убираем memStorage из аргументов, т.к. данные приходят снаружи
) Archiver { // Возвращаем интерфейс
	return &service{
		// memStorage:  memStorage,
		log: log,
	}
}

// Archive возвращает t1 t2 t3 t4 t5 в виде зип архива в буфере
func (s *service) Archive(
	t1 []models.TableOne,
	t2 []models.TableTwo,
	t3 []models.TableThree,
	t4 []models.TableFour,
	t5 models.TableFive,
) (*bytes.Buffer, error) {
	s.log.Infow("Starting archiving process")

	//Проверяем, что все блоки импортированы
	if (len(t1) * len(t2) * len(t3) * len(t4)) == 0 {
		return nil, errors.New("сначала ипортируйте все 4 блока")
	}
	//Проверяем, что техническая карта составлена
	if t5 == (models.TableFive{}) {
		return nil, errors.New("сначала заполните Тех. карту")
	}

	// 1. Используем буфер в памяти для создания ZIP-архива
	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	var finalErr error // Для сбора ошибок из хелперов

	// 2. Создаем и добавляем XLSX файлы в архив
	// Блок 1
	if err := tableOneToXLSXBuffer(zipWriter, "Block_1_PressureTemp.xlsx", t1); err != nil {
		finalErr = errors.Wrap(err, "failed to add block 1 to zip") // Собираем ошибки
		s.log.Errorw("Archiver error", "error", finalErr)           // Логируем
		// Не выходим сразу, пытаемся добавить другие файлы
	} else {
		s.log.Debugw("Added Block 1 to ZIP")
	}

	// Блок 2
	if err := tableTwoToXLSXBuffer(zipWriter, "Block_2_TubingAnnulus.xlsx", t2, s.log); err != nil {
		finalErr = errors.Wrap(err, "failed to add block 2 to zip")
		s.log.Errorw("Archiver error", "error", finalErr)
	} else {
		s.log.Debugw("Added Block 2 to ZIP")
	}

	// Блок 3
	if err := tableThreeToXLSXBuffer(zipWriter, "Block_3_FlowRates.xlsx", t3, s.log); err != nil {
		finalErr = errors.Wrap(err, "failed to add block 3 to zip")
		s.log.Errorw("Archiver error", "error", finalErr)
	} else {
		s.log.Debugw("Added Block 3 to ZIP")
	}

	// Блок 4
	if err := tableFourToXLSXBuffer(zipWriter, "Block_4_Inclinometry.xlsx", t4, s.log); err != nil {
		finalErr = errors.Wrap(err, "failed to add block 4 to zip")
		s.log.Errorw("Archiver error", "error", finalErr)
	} else {
		s.log.Debugw("Added Block 4 to ZIP")
	}

	//Блок 5
	if err := tableFiveToXLSXBuffer(zipWriter, "Block_5_Tech_card.xlsx", t5, s.log); err != nil {
		finalErr = errors.Wrap(err, "failed to add block 5 to zip")
		s.log.Errorw("Archiver error", "error", finalErr)
	} else {
		s.log.Debugw("Added Block 5 to ZIP")
	}

	// 3. Завершаем создание ZIP-архива
	if err := zipWriter.Close(); err != nil {
		// Если уже была ошибка при добавлении файла, возвращаем её, иначе ошибку закрытия
		if finalErr != nil {
			return nil, finalErr // Возвращаем первую ошибку добавления файла
		}
		return nil, errors.Wrap(err, "failed to close zip writer")
	}

	// Если была ошибка при добавлении файла, но закрытие прошло успешно, все равно возвращаем ошибку добавления
	if finalErr != nil {
		return nil, finalErr
	}

	s.log.Infow("ZIP archive created in memory", "size_bytes", zipBuf.Len())
	return zipBuf, nil // Ошибки нет, возвращаем буфер
}

// --- Хелперы для записи данных в XLSX и добавления в ZIP ---

// Создадим отдельные функции для каждого типа таблицы
func tableOneToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableOne) error {
	xlsxFile := excelize.NewFile()
	sheetName := "Block1_PressureTemp"
	_ = xlsxFile.SetSheetName("Sheet1", sheetName) // Переименуем лист

	// Заголовки (лучше брать из модели, если есть Columns())
	headers := []string{"timestamp", "pressure_depth", "temperature_depth", "pressure_at_vdp"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = xlsxFile.SetCellValue(sheetName, cell, h)
	}

	// Данные
	for rowIdx, rowData := range data {
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx+2), rowData.Timestamp) // Excelize сам может форматнуть время
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx+2), rowData.PressureDepth)
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx+2), rowData.TemperatureDepth)
		//надо подумоть
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx+2), rowData.PressureAtVDP)
	}

	// Запись файла в ZIP
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "zip.Create failed for %s", filename)
	}
	if err := xlsxFile.Write(fileWriter); err != nil {
		return errors.Wrapf(err, "xlsxFile.Write failed for %s", filename)
	}
	return nil
}

func tableTwoToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableTwo, log logger.Logger) error { // Добавлен аргумент log
	xlsxFile := excelize.NewFile()
	sheetName := "Block2_TubingAnnulus"
	_ = xlsxFile.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"timestamp_tubing", "pressure_tubing",
		"timestamp_annulus", "pressure_annulus",
		"timestamp_linear", "pressure_linear",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = xlsxFile.SetCellValue(sheetName, cell, h)
	}

	// Данные TableTwo
	for rowIdx, rowData := range data {
		col := 1
		// Используем новый способ получения имени колонки
		colName, err := excelize.ColumnNumberToName(col)
		if err != nil {
			log.Errorw("Failed to get col name", "col", col, "error", err) // Используем переданный логгер
			return errors.Wrapf(err, "col num %d", col)
		}
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx+2), rowData.TimestampTubing)
		col++

		colName, err = excelize.ColumnNumberToName(col)
		if err != nil {
			log.Errorw("Failed to get col name", "col", col, "error", err)
			return errors.Wrapf(err, "col num %d", col)
		}
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx+2), rowData.PressureTubing)
		col++

		colName, err = excelize.ColumnNumberToName(col)
		if err != nil {
			log.Errorw("Failed to get col name", "col", col, "error", err)
			return errors.Wrapf(err, "col num %d", col)
		}
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx+2), rowData.TimestampAnnulus)
		col++

		colName, err = excelize.ColumnNumberToName(col)
		if err != nil {
			log.Errorw("Failed to get col name", "col", col, "error", err)
			return errors.Wrapf(err, "col num %d", col)
		}
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx+2), rowData.PressureAnnulus)
		col++

		colName, err = excelize.ColumnNumberToName(col)
		if err != nil {
			log.Errorw("Failed to get col name", "col", col, "error", err)
			return errors.Wrapf(err, "col num %d", col)
		}
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx+2), rowData.TimestampLinear)
		col++

		colName, err = excelize.ColumnNumberToName(col)
		if err != nil {
			log.Errorw("Failed to get col name", "col", col, "error", err)
			return errors.Wrapf(err, "col num %d", col)
		}
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("%s%d", colName, rowIdx+2), rowData.PressureLinear)
	}

	// Запись файла в ZIP
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "zip.Create failed for %s", filename)
	}
	if err := xlsxFile.Write(fileWriter); err != nil {
		return errors.Wrapf(err, "xlsxFile.Write failed for %s", filename)
	}
	return nil
}

func tableThreeToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableThree, log logger.Logger) error { // Добавлен логгер
	xlsxFile := excelize.NewFile()
	sheetName := "Block3_FlowRates"
	_ = xlsxFile.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"timestamp", "flow_liquid", "water_cut", "flow_gas",
		"oil_flow_rate", "water_flow_rate", "gas_oil_ratio",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = xlsxFile.SetCellValue(sheetName, cell, h)
	}

	for rowIdx, rowData := range data {
		col := 1
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.Timestamp, log) // Передаем логгер в setCellValue
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.LiquidFlowRate, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.WaterCut, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.GasFlowRate, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.OilFlowRate, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.WaterFlowRate, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.GasFactor, log) // GasFactor поле называется
	}

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "zip.Create failed for %s", filename)
	}
	if err := xlsxFile.Write(fileWriter); err != nil {
		return errors.Wrapf(err, "xlsxFile.Write failed for %s", filename)
	}
	return nil
}

func tableFourToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableFour, log logger.Logger) error { // Добавлен логгер
	xlsxFile := excelize.NewFile()
	sheetName := "Block4_Inclinometry"
	_ = xlsxFile.SetSheetName("Sheet1", sheetName)

	headers := []string{"research_id", "measure_depth", "true_vertical_depth", "true_vertical_depth_sub_sea"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = xlsxFile.SetCellValue(sheetName, cell, h)
	}

	for rowIdx, rowData := range data {
		col := 1
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.ResearchID.String(), log) // Передаем логгер
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.MeasuredDepth, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.TrueVerticalDepth, log)
		col++
		setCellValue(xlsxFile, sheetName, col, rowIdx+2, rowData.TrueVerticalDepthSubSea, log)
	}

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "zip.Create failed for %s", filename)
	}
	if err := xlsxFile.Write(fileWriter); err != nil {
		return errors.Wrapf(err, "xlsxFile.Write failed for %s", filename)
	}
	return nil
}

func tableFiveToXLSXBuffer(zipWriter *zip.Writer, filename string, data models.TableFive, log logger.Logger) error { // Добавлен логгер
	xlsxFile := excelize.NewFile()
	sheetName := "Block5_TechCard"
	_ = xlsxFile.SetSheetName("Sheet1", sheetName)

	headers := data.Columns()
	for i, h := range headers {
		_ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("A%d", i+1), h)
	}

	dataMap := data.Map()
	for i, h := range headers {
		if val, ok := dataMap[h]; ok {
			setCellValue(xlsxFile, sheetName, 2, i+1, val, log) // Передаем логгер
		} else {
			// Можно добавить обработку, если ключ не найден в Map (хотя не должно быть)
			log.Errorw("Key not found in TableFive map", "key", h)
			setCellValue(xlsxFile, sheetName, 2, i+1, nil, log) // Пишем nil на всякий случай
		}
	}

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "zip.Create failed for %s", filename)
	}
	if err := xlsxFile.Write(fileWriter); err != nil {
		return errors.Wrapf(err, "xlsxFile.Write failed for %s", filename)
	}
	return nil
}

// --- НАЧАЛО: Хелпер для записи значения с учетом типа (особенно указателей) ---
func setCellValue(f *excelize.File, sheet string, col, row int, value interface{}, log logger.Logger) {
	// Используем ColumnNumberToName для получения имени колонки
	colName, err := excelize.ColumnNumberToName(col)
	if err != nil {
		log.Errorw("Failed to get col name in setCellValue", "col", col, "error", err)
		// В хелпере просто не будем записывать значение при ошибке имени колонки
		return
	}
	cellName := fmt.Sprintf("%s%d", colName, row) // Собираем имя ячейки A1, B2 и т.д.

	switch v := value.(type) {
	case *float64:
		if v != nil {
			_ = f.SetCellValue(sheet, cellName, *v)
		} else {
			_ = f.SetCellValue(sheet, cellName, nil)
		}
	case *string:
		if v != nil {
			_ = f.SetCellValue(sheet, cellName, *v)
		} else {
			_ = f.SetCellValue(sheet, cellName, nil)
		}
	case *int:
		if v != nil {
			_ = f.SetCellValue(sheet, cellName, *v)
		} else {
			_ = f.SetCellValue(sheet, cellName, nil)
		}
	case time.Time:
		// Можно указать формат, если нужно
		// Или excelize сам попробует
		_ = f.SetCellValue(sheet, cellName, v)
		// Можно задать стиль для даты, если надо
		// style, _ := f.NewStyle(`{"number_format": 14}`) // Пример для формата dd-mm-yy
		// _ = f.SetCellStyle(sheet, cellName, cellName, style)
	default:
		_ = f.SetCellValue(sheet, cellName, v)
	}
}
