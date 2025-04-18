package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// TODO достать все, добавить строку
func NewBlockStorage(db *sqlx.DB) (BlocksStorage, error) {
	return &sqlite{
		DB: db,
	}, nil
}

func (s sqlite) GetTableOneByReportID() ([]models.TableOne, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetTableTwoByReportID() ([]models.TableTwo, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetTableThreeByReportID() ([]models.TableThree, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetTableFourByReportID() ([]models.TableFour, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetTableFive() ([]models.TableFive, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockOne(data []models.TableOne, reportID int32) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockTwo(data []models.TableTwo, reportID int32) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockThree(data []models.TableThree, reportID int32) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFour(data []models.TableFour, reportID int32) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFive(data []models.TableFive) (reportID int32, err error) {
	//TODO implement me
	panic("implement me")
}
