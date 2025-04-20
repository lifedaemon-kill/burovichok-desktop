package postgres

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// GetAllTableFive возвращает все записи из таблицы reports
func (p *Postgres) GetAllTableFive(ctx context.Context) ([]models.TableFive, error) {
	var reports []models.TableFive
	qb := psql().
		Select(models.TableFive{}.Columns()...).
		From(models.TableFive{}.TableName())

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building GetAllTableFive query")
	}
	if err := p.DB.SelectContext(ctx, &reports, sqlStr, args...); err != nil {
		return nil, errors.Wrap(err, "executing GetAllTableFive query")
	}
	return reports, nil
}

// AddBlockFive вставляет запись в таблицу reports и возвращает сгенерированный ID
func (p *Postgres) AddBlockFive(ctx context.Context, data models.TableFive) (int64, error) {
	qb := psql().
		Insert(models.TableFive{}.TableName()).
		Columns(models.TableFive{}.Columns()...).
		SetMap(data.Map()).
		Suffix("RETURNING id")

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "building AddBlockFive query")
	}
	var id int64
	if err = p.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "executing AddBlockFive query")
	}
	return id, nil
}
