// internal/storage/postgres/archive_info.go
package postgres

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// SaveArchiveInfo сохраняет метаданные архива в БД
func (p *Postgres) SaveArchiveInfo(ctx context.Context, info models.ArchiveInfo) error {
	qb := psql().
		Insert(info.TableName()).
		Columns(info.Columns()...).
		Values(
			info.ObjectName,
			info.BucketName,
			info.Size,
			info.ETag,
			info.UploadedAt, // Или можно положиться на DEFAULT NOW() в БД
		).
        // Добавляем ON CONFLICT на случай повторной загрузки с тем же именем (маловероятно из-за timestamp)
        Suffix("ON CONFLICT (object_name) DO UPDATE SET bucket_name = EXCLUDED.bucket_name, size = EXCLUDED.size, etag = EXCLUDED.etag, uploaded_at = EXCLUDED.uploaded_at")


	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "building SaveArchiveInfo query")
	}

	if _, err := p.DB.ExecContext(ctx, sqlStr, args...); err != nil {
		return errors.Wrap(err, "executing SaveArchiveInfo query")
	}
	return nil
}