package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
)

// TODO достать все, добавить строку
func NewBlockRepository(db *sqlx.DB) (BlocksStorage, error) {
	return &sqlite{
		DB: db,
	}, nil
}

func (s sqlite) GetBlockOne() ([]models.BlockOne, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockTwo() ([]models.BlockTwo, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockThree() ([]models.BlockThree, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockFour() ([]models.Inclinometry, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockFive() ([]models.GeneralInformation, error) {
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
