package storage

import "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"

// Storage определяет интерфейс для взаимодействия с хранилищем данных.
type Storage interface {
	// Методы для добавления данных (принимают срезы, как возвращает парсер)
	AddBlockOneData(data []models.TableOne) error
	AddBlockTwoData(data []models.TableTwo) error
	AddBlockThreeData(data []models.TableThree) error
	AddBlockFourData(data []models.TableFour) error

	// Методы для получения всех данных (возвращают копии для безопасности)
	GetAllBlockOneData() ([]models.TableOne, error)
	GetAllBlockTwoData() ([]models.TableTwo, error)
	GetAllBlockThreeData() ([]models.TableThree, error)

	// Метод для очистки всего хранилища
	ClearAll() error

	// Можно добавить методы для получения количества записей (опционально)
	CountBlockOne() int
	CountBlockTwo() int
	CountBlockThree() int
}
