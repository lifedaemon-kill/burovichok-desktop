// internal/service/archiver/service.go
package archiver

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"

	// "path/filepath"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
	"github.com/minio/minio-go/v7"
	"github.com/xuri/excelize/v2"
)

// Archiver определяет интерфейс для архивации и выгрузки
type Archiver interface {
	// ArchiveAndUpload собирает данные, создает XLSX, архивирует и выгружает в MinIO.
	// Возвращает имя созданного архива и ошибку.
	ArchiveAndUpload(ctx context.Context, baseName string) (string, error)
}

// service реализует интерфейс Archiver
type service struct {
	memStorage  inmemory.InMemoryBlocksStorage // Зависимость от интерфейса хранилища
	minioClient *minio.Client
	bucketName  string
	log         logger.Logger
}

// NewService создает новый экземпляр сервиса Archiver
func NewService(
	memStorage inmemory.InMemoryBlocksStorage,
	minioClient *minio.Client,
	bucketName string,
	log logger.Logger,
) Archiver {
	return &service{
		memStorage:  memStorage,
		minioClient: minioClient,
		bucketName:  bucketName,
		log:         log,
	}
}

// ArchiveAndUpload реализует основную логику
func (s *service) ArchiveAndUpload(ctx context.Context, baseName string) (string, error) {
	s.log.Infow("Starting archiving and upload process", "baseName", baseName)

	// 1. Получаем все данные из in-memory
	data1, err := s.memStorage.GetAllBlockOneData()
	if err != nil {
		return "", errors.Wrap(err, "failed to get block one data")
	}
	data2, err := s.memStorage.GetAllBlockTwoData()
	if err != nil {
		return "", errors.Wrap(err, "failed to get block two data")
	}
	data3, err := s.memStorage.GetAllBlockThreeData()
	if err != nil {
		return "", errors.Wrap(err, "failed to get block three data")
	}
	// data4, err := s.memStorage.GetAllBlockFourData() // Раскомментируй, когда реализуешь в inmemory
	// if err != nil {
	// 	 return "", errors.Wrap(err, "failed to get block four data")
	// }

	s.log.Debugw("Retrieved data from memory",
		"block1_count", len(data1),
		"block2_count", len(data2),
		"block3_count", len(data3),
		// "block4_count", len(data4),
	)

	// Используем буфер в памяти для создания ZIP-архива
	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	// 2. Создаем и добавляем XLSX файлы в архив
	// Блок 1
	if len(data1) > 0 {
		if err := addTableOneToZip(zipWriter, "Block_1_PressureTemp.xlsx", data1); err != nil {
			_ = zipWriter.Close() // Пытаемся закрыть writer при ошибке
			return "", errors.Wrap(err, "failed to add block 1 to zip")
		}
		s.log.Debugw("Added Block 1 to ZIP")
	}
	// Блок 2
	if len(data2) > 0 {
		if err := addTableTwoToZip(zipWriter, "Block_2_TubingAnnulus.xlsx", data2); err != nil {
			_ = zipWriter.Close()
			return "", errors.Wrap(err, "failed to add block 2 to zip")
		}
		s.log.Debugw("Added Block 2 to ZIP")
	}
	// Блок 3
	if len(data3) > 0 {
		if err := addTableThreeToZip(zipWriter, "Block_3_FlowRates.xlsx", data3); err != nil {
			_ = zipWriter.Close()
			return "", errors.Wrap(err, "failed to add block 3 to zip")
		}
		s.log.Debugw("Added Block 3 to ZIP")
	}
	// Блок 4 (когда будет готов)
	// if len(data4) > 0 {
	// 	 if err := addTableFourToZip(zipWriter, "Block_4_Inclinometry.xlsx", data4); err != nil {
	//      _ = zipWriter.Close()
	// 	 	   return "", errors.Wrap(err, "failed to add block 4 to zip")
	// 	 }
	// 	 s.log.Debugw("Added Block 4 to ZIP")
	// }

	// 3. Завершаем создание ZIP-архива
	if err := zipWriter.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close zip writer")
	}
	s.log.Infow("ZIP archive created in memory", "size_bytes", zipBuf.Len())

	// 4. Генерируем имя файла для MinIO
	timestamp := time.Now().Format("20060102_150405") // Формат YYYYMMDD_HHMMSS
	// Очищаем baseName от недопустимых символов (простой вариант)
	safeBaseName := strings.ReplaceAll(baseName, " ", "_")
	safeBaseName = strings.ReplaceAll(safeBaseName, "/", "-")
	// ... можно добавить еще замен или использовать regexp
	if safeBaseName == "" {
		safeBaseName = "archive"
	}
	objectName := fmt.Sprintf("%s_%s.zip", timestamp, safeBaseName)
	s.log.Infow("Generated object name for MinIO", "name", objectName)

	// 5. Загружаем архив в MinIO
	contentType := "application/zip"
	uploadInfo, err := s.minioClient.PutObject(
		ctx,
		s.bucketName,
		objectName,          // Имя объекта в MinIO
		zipBuf,              // Данные из буфера
		int64(zipBuf.Len()), // Размер данных
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		s.log.Errorw("Failed to upload archive to MinIO", "object", objectName, "error", err)
		return "", errors.Wrapf(err, "failed to upload '%s' to bucket '%s'", objectName, s.bucketName)
	}

	s.log.Infow("Successfully uploaded archive to MinIO",
		"object", objectName,
		"bucket", uploadInfo.Bucket,
		"etag", uploadInfo.ETag,
		"size", uploadInfo.Size,
	)

	return objectName, nil // Возвращаем имя загруженного файла
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
func addTableOneToZip(zipWriter *zip.Writer, filename string, data []models.TableOne) error {
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

// Аналогичные функции addTableTwoToZip, addTableThreeToZip, addTableFourToZip...
func addTableTwoToZip(zipWriter *zip.Writer, filename string, data []models.TableTwo) error {
	// ... реализация записи TableTwo в xlsx и добавления в zipWriter ...
	return nil // Заглушка
}

func addTableThreeToZip(zipWriter *zip.Writer, filename string, data []models.TableThree) error {
	// ... реализация записи TableThree в xlsx и добавления в zipWriter ...
	return nil // Заглушка
}

// func addTableFourToZip(zipWriter *zip.Writer, filename string, data []models.TableFour) error {
// 	 // ... реализация записи TableFour в xlsx и добавления в zipWriter ...
//	 return nil
// }
