package postgres

import (
	"context"

	"github.com/cockroachdb/errors"
	_ "github.com/lib/pq"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// GetAllOilField возвращает все записи из таблицы oilfield
func (p *Postgres) GetAllOilField(ctx context.Context) ([]models.OilField, error) {
	var items []models.OilField
	qb := psql().
		Select(models.OilField{}.Columns()...).
		From(models.OilField{}.TableName())

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building GetAllOilField query")
	}
	if err := p.DB.SelectContext(ctx, &items, sqlStr, args...); err != nil {
		return nil, errors.Wrap(err, "executing GetAllOilField query")
	}
	return items, nil
}

// GetAllInstrumentType возвращает все записи из таблицы instrument_type
func (p *Postgres) GetAllInstrumentType(ctx context.Context) ([]models.InstrumentType, error) {
	var items []models.InstrumentType
	qb := psql().
		Select(models.InstrumentType{}.Columns()...).
		From(models.InstrumentType{}.TableName())

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building GetAllInstrumentType query")
	}
	if err := p.DB.SelectContext(ctx, &items, sqlStr, args...); err != nil {
		return nil, errors.Wrap(err, "executing GetAllInstrumentType query")
	}
	return items, nil
}

// GetAllProductiveHorizon возвращает все записи из таблицы productive_horizon
func (p *Postgres) GetAllProductiveHorizon(ctx context.Context) ([]models.ProductiveHorizon, error) {
	var items []models.ProductiveHorizon
	qb := psql().
		Select(models.ProductiveHorizon{}.Columns()...).
		From(models.ProductiveHorizon{}.TableName())

	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building GetAllProductiveHorizon query")
	}
	if err := p.DB.SelectContext(ctx, &items, sqlStr, args...); err != nil {
		return nil, errors.Wrap(err, "executing GetAllProductiveHorizon query")
	}
	return items, nil
}

// GetAllResearchType возвращает все записи из таблицы research_type
func (p *Postgres) GetAllResearchType(ctx context.Context) ([]models.ResearchType, error) {
	var items []models.ResearchType
	qb := psql().
		Select(models.ResearchType{}.Columns()...).
		From(models.ResearchType{}.TableName())
	sqlStr, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building GetAllResearchType query")
	}
	if err := p.DB.SelectContext(ctx, &items, sqlStr, args...); err != nil {
		return nil, errors.Wrap(err, "executing GetAllResearchType query")
	}
	return items, nil
}

// AddOilField вставляет записи в таблицу oilfield
func (p *Postgres) AddOilField(ctx context.Context, fields []models.OilField) error {
	for _, f := range fields {
		qb := psql().
			Insert(models.OilField{}.TableName()).
			SetMap(f.Map())

		sqlStr, args, err := qb.ToSql()
		if err != nil {
			return errors.Wrap(err, "building AddOilField query")
		}
		if _, err := p.DB.ExecContext(ctx, sqlStr, args...); err != nil {
			return errors.Wrap(err, "executing AddOilField query")
		}
	}
	return nil
}

// AddInstrumentType вставляет записи в таблицу instrument_type
func (p *Postgres) AddInstrumentType(ctx context.Context, items []models.InstrumentType) error {
	for _, it := range items {
		qb := psql().
			Insert(models.InstrumentType{}.TableName()).
			SetMap(it.Map())

		sqlStr, args, err := qb.ToSql()
		if err != nil {
			return errors.Wrap(err, "building AddInstrumentType query")
		}
		if _, err := p.DB.ExecContext(ctx, sqlStr, args...); err != nil {
			return errors.Wrap(err, "executing AddInstrumentType query")
		}
	}
	return nil
}

// AddProductiveHorizon вставляет записи в таблицу productive_horizon
func (p *Postgres) AddProductiveHorizon(ctx context.Context, items []models.ProductiveHorizon) error {
	for _, ph := range items {
		qb := psql().
			Insert(models.ProductiveHorizon{}.TableName()).
			SetMap(ph.Map())

		sqlStr, args, err := qb.ToSql()
		if err != nil {
			return errors.Wrap(err, "building AddProductiveHorizon query")
		}
		if _, err := p.DB.ExecContext(ctx, sqlStr, args...); err != nil {
			return errors.Wrap(err, "executing AddProductiveHorizon query")
		}
	}
	return nil
}

// AddResearchType вставляет записи в таблицу research_type
func (p *Postgres) AddResearchType(ctx context.Context, items []models.ResearchType) error {
	for _, it := range items {
		qb := psql().
			Insert(models.ResearchType{}.TableName()).
			SetMap(it.Map())
		sqlStr, args, err := qb.ToSql()
		if err != nil {
			return errors.Wrap(err, "building AddResearchType query")
		}
		if _, err := p.DB.ExecContext(ctx, sqlStr, args...); err != nil {
			return errors.Wrap(err, "executing AddResearchType query")
		}
	}
	return nil
}
