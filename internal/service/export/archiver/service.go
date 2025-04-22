package archiver

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/cockroachdb/errors"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
	"github.com/xuri/excelize/v2"
)

// Service реализует интерфейс Archiver
type Service struct {
	memStorage inmemory.InMemoryBlocksStorage // Зависимость от интерфейса хранилища

	bucketName string
	log        logger.Logger
}

// NewService создает новый экземпляр сервиса Archiver
func NewService(
	log logger.Logger,
) *Service {
	return &Service{
		log: log,
	}
}

// Archive возвращает t1 t2 t3 t4 t5 в виде зип архиве в буфере
func (s *Service) Archive(
	t1 []models.TableOne,
	t2 []models.TableTwo,
	t3 []models.TableThree,
	t4 []models.TableFour,
	t5 models.TableFive,
) (*bytes.Buffer, error) {
	s.log.Infow("Starting archiving and upload process")

	// Используем буфер в памяти для создания ZIP-архива
	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	if (len(t1)*len(t2)*len(t3)*len(t4)) == 0 || t5 == (models.TableFive{}) {
		return nil, errors.Wrap(errors.New("Один из блоков не заполнен"), "")
	}
	// 2. Создаем и добавляем XLSX файлы в архив
	// Блок 1
	if err := tableOneToXLSXBuffer(zipWriter, "Block_1_PressureTemp.xlsx", t1); err != nil {
		_ = zipWriter.Close()
		return nil, errors.Wrap(err, "failed to add block 1 to zip")
	}

	// Блок 2
	if err := tableTwoToXLSXBuffer(zipWriter, "Block_2_TubingAnnulus.xlsx", t2); err != nil {
		_ = zipWriter.Close()
		return nil, errors.Wrap(err, "failed to add block 2 to zip")
	}

	// Блок 3

	if err := tableThreeToXLSXBuffer(zipWriter, "Block_3_FlowRates.xlsx", t3); err != nil {
		_ = zipWriter.Close()
		return nil, errors.Wrap(err, "failed to add block 3 to zip")
	}
	// Блок 4
	if err := tableFourToXLSXBuffer(zipWriter, "Block_4_Inclinometry.xlsx", t4); err != nil {
		_ = zipWriter.Close()
		return nil, errors.Wrap(err, "failed to add block 4 to zip")
	}
	//Блок 5
	if err := tableFiveToXLSXBuffer(zipWriter, "Block_5_Tech_card.xlsx", t5); err != nil {
		_ = zipWriter.Close()
		return nil, errors.Wrap(err, "failed to add block 4 to zip")
	}

	// 3. Завершаем создание ZIP-архива
	if err := zipWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to close zip writer")
	}
	s.log.Infow("ZIP archive created in memory", "size_bytes", zipBuf.Len())

	return zipBuf, nil
}

// --- Хелперы для записи данных в XLSX и добавления в ZIP ---

// Вспомогательная функция для записи данных в io.Writer внутри ZIP
func addFileToZip(zipWriter *zip.Writer, filename string, generatorFunc func(w io.Writer) error) error {
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to create '%s' in zip", filename)
	}
	// Используем буфер, чтобы сначала сгенерировать XLSX, потом записать
	excelBuf := new(bytes.Buffer)
	xlsxFile := excelize.NewFile()

	// Генерация контента XLSX файла (пример для TableOne)
	sheetName := "Data" // Имя листа
	// Заголовки
	headers := []string{"timestamp", "pressure_depth", "temperature_depth", "pressure_at_vdp"} // Используй models.TableOne{}.Columns() если они есть
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = xlsxFile.SetCellValue(sheetName, cell, h)
	}
	// Данные (нужно реализовать запись data1 в xlsxFile)
	// Примерно так:
	// for rowIdx, rowData := range data1 {
	//    _ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx+2), rowData.Timestamp.Format(time.RFC3339))
	//    _ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx+2), rowData.PressureDepth)
	//    // ... и так далее для всех полей ...
	//    // Не забыть про *float64 - проверять на nil
	//    if rowData.PressureAtVDP != nil {
	//        _ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx+2), *rowData.PressureAtVDP)
	//    } else {
	//		  _ = xlsxFile.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx+2), nil) // или пустую строку
	//    }
	// }

	if err := xlsxFile.Write(excelBuf); err != nil {
		return errors.Wrapf(err, "failed to write excel data to buffer for '%s'", filename)
	}

	// Копируем сгенерированный Excel из буфера в ZIP
	if _, err := io.Copy(fileWriter, excelBuf); err != nil {
		return errors.Wrapf(err, "failed to copy excel buffer to zip writer for '%s'", filename)
	}

	return nil
}

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

// Аналогичные функции tableTwoToXLSXBuffer, tableThreeToXLSXBuffer, tableFourToXLSXBuffer...
func tableTwoToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableTwo) error {
	// ... реализация записи TableTwo в xlsx и добавления в zipWriter ...
	return nil // Заглушка
}

func tableThreeToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableThree) error {
	// ... реализация записи TableThree в xlsx и добавления в zipWriter ...
	return nil // Заглушка
}

func tableFourToXLSXBuffer(zipWriter *zip.Writer, filename string, data []models.TableFour) error {
	// ... реализация записи TableFour в xlsx и добавления в zipWriter ...
	return nil
}

func tableFiveToXLSXBuffer(zipWriter *zip.Writer, filename string, data models.TableFive) error {
	// ... реализация записи TableFour в xlsx и добавления в zipWriter ...
	return nil
}
