package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func NewBlockStorage(db *sqlx.DB) (BlocksStorage, error) {
	return &sqlite{
		db: db,
	}, nil
}

func (s sqlite) GetAllTableFive() ([]models.TableFive, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFive(data models.TableFive) (reportID int32, err error) {
	//TODO implement me
	panic("implement me")
}
