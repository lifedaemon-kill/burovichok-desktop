package storage

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
)

// BlocksStorage определяет интерфейс для взаимодействия с хранилищем данных.
type BlocksStorage interface {
	// Методы для добавления данных (принимают срезы, как возвращает парсер)
	AddBlockOneData(data []models.BlockOne) error
	AddBlockTwoData(data []models.BlockTwo) error
	AddBlockThreeData(data []models.BlockThree) error

	// Методы для получения всех данных (возвращают копии для безопасности)
	GetAllBlockOneData() ([]models.BlockOne, error)
	GetAllBlockTwoData() ([]models.BlockTwo, error)
	GetAllBlockThreeData() ([]models.BlockThree, error)

	// Метод для очистки всего хранилища
	ClearAll() error

	// Можно добавить методы для получения количества записей (опционально)
	CountBlockOne() int
	CountBlockTwo() int
	CountBlockThree() int
}

type GuidebooksStorage interface {
	AddOilPlaces() ([]models.OilPlaces, error)
	AddInstrumentType() ([]models.InstrumentType, error)
	AddProductiveHorizon() ([]models.ProductiveHorizon, error)

	GetAllOilPlaces() ([]models.OilPlaces, error)
	GetAllInstrumentType() ([]models.InstrumentType, error)
	GetAllProductiveHorizon() ([]models.ProductiveHorizon, error)
}
