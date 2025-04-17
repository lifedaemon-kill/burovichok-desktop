package storage

import "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"

// Storage определяет интерфейс для взаимодействия с хранилищем данных.
type Storage interface {
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
