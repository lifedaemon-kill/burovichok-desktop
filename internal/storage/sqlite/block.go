package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// TODO достать все, добавить строку
func NewBlockRepository(db *sqlx.DB) (BlocksStorage, error) {
	return &sqlite{
		DB: db,
	}, nil
}

func (s sqlite) GetBlockOne() ([]models.TableOne, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockTwo() ([]models.TableTwo, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockThree() ([]models.TableThree, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockFour() ([]models.TableFour, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockFive() ([]models.TableFive, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockOne() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockTwo() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockThree() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFour() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFive() (int32, error) {
	//TODO implement me
	panic("implement me")
}
