package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
)

// Postgres инкапсулирует подключение к PostgreSQL через sqlx.DB
type Postgres struct {
	DB   *sqlx.DB
	zLog logger.Logger
}

// New подключается к базе по DSN, настраивает пул и возвращает обертку Postgres
func New(ctx context.Context, cfg config.DBConf, zLog logger.Logger) (*Postgres, error) {
	var (
		db  *sqlx.DB
		err error
	)

	// Попытки переподключения с экспоненциальным бэкоффом
	backoff := time.Second
	for i := 0; i < cfg.MaxRetries; i++ {
		// Контекст с таймаутом на соединение
		ctxConn, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		db, err = sqlx.ConnectContext(ctxConn, "postgres", cfg.DSN)
		if err == nil {
			break
		}

		zLog.Infow(fmt.Sprintf("Postgres connect attempt %d failed: %v; retrying in %s", i+1, err, backoff))
		time.Sleep(backoff)
		backoff *= 2
	}
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.ConnectContext")
	}

	// Проверяем пинг
	if err = db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "db.PingContext")
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return &Postgres{DB: db, zLog: zLog}, nil
}

func (p *Postgres) GetSqlDB() *sql.DB {
	return p.DB.DB
}

// Postgres specific squirrel builder
func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
