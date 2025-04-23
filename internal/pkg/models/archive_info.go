// internal/pkg/models/archive_info.go
package models

import "time"

// ArchiveInfo хранит метаданные о загруженном архиве
type ArchiveInfo struct {
	ObjectName string    `db:"object_name"` // Имя файла в MinIO (PK)
	BucketName string    `db:"bucket_name"`
	Size       int64     `db:"size"`
	ETag       string    `db:"etag"`
	UploadedAt time.Time `db:"uploaded_at"`
}

// TableName возвращает имя таблицы
func (ArchiveInfo) TableName() string {
	return "archive_info"
}

// Columns возвращает список колонок
func (ArchiveInfo) Columns() []string {
	return []string{"object_name", "bucket_name", "size", "etag", "uploaded_at"}
}

// Map для удобства вставки (опционально, если используется SetMap)
func (ai ArchiveInfo) Map() map[string]interface{} {
	return map[string]interface{}{
		"object_name": ai.ObjectName,
		"bucket_name": ai.BucketName,
		"size":        ai.Size,
		"etag":        ai.ETag,
		"uploaded_at": ai.UploadedAt,
	}
}