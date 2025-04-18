// internal/pkg/storage/inmemory/store.go
package inmemory

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"sync"
)

//можешно импортировать несколько файлов одного или разных типов, данные будут накапливаться. Кнопка "Очистить хранилище" удалит всё из памяти. Хранятся в Storage
//в будущем, если понадобятся данные для графиков или экспорта, соответствующие сервисы (chart, export) также должны будут получить экземпляр InMemoryBlocksStorage и вызывать методы GetAll...Data().

// InMemoryBlocksStorage определяет интерфейс для взаимодействия с хранилищем данных.
type InMemoryBlocksStorage interface {
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

// Storage реализует интерфейс storage.InMemoryBlocksStorage, храня данные в памяти.
type Storage struct {
	mu             sync.RWMutex // Mutex для безопасного доступа к данным
	blockOneData   []models.TableOne
	blockTwoData   []models.TableTwo
	blockThreeData []models.TableThree
	inclinometry   []models.TableFour
}

// NewInMemoryBlocksStorage создает новый экземпляр Storage.
func NewInMemoryBlocksStorage() InMemoryBlocksStorage { // Возвращаем интерфейс!
	return &Storage{
		blockOneData:   make([]models.TableOne, 0),
		blockTwoData:   make([]models.TableTwo, 0),
		blockThreeData: make([]models.TableThree, 0),
		inclinometry:   make([]models.TableFour, 0),
	}
}

// AddBlockOneData добавляет данные TableOne в хранилище.
func (s *Storage) AddBlockOneData(data []models.TableOne) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOneData = append(s.blockOneData, data...)
	return nil // В in-memory обычно нет ошибок добавления, кроме нехватки памяти (panic)
}

// AddBlockTwoData добавляет данные TableTwo в хранилище.
func (s *Storage) AddBlockTwoData(data []models.TableTwo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockTwoData = append(s.blockTwoData, data...)
	return nil
}

// AddBlockThreeData добавляет данные TableThree в хранилище.
func (s *Storage) AddBlockThreeData(data []models.TableThree) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockThreeData = append(s.blockThreeData, data...)
	return nil
}

// GetAllBlockOneData возвращает копию всех данных TableOne.
func (s *Storage) GetAllBlockOneData() ([]models.TableOne, error) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	// Возвращаем копию, чтобы внешние изменения не влияли на хранилище
	dataCopy := make([]models.TableOne, len(s.blockOneData))
	copy(dataCopy, s.blockOneData)
	return dataCopy, nil
}

// GetAllBlockTwoData возвращает копию всех данных TableTwo.
func (s *Storage) GetAllBlockTwoData() ([]models.TableTwo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableTwo, len(s.blockTwoData))
	copy(dataCopy, s.blockTwoData)
	return dataCopy, nil
}

// GetAllBlockThreeData возвращает копию всех данных TableThree.
func (s *Storage) GetAllBlockThreeData() ([]models.TableThree, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableThree, len(s.blockThreeData))
	copy(dataCopy, s.blockThreeData)
	return dataCopy, nil
}

// ClearAll очищает все данные в хранилище.
func (s *Storage) ClearAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOneData = make([]models.TableOne, 0)
	s.blockTwoData = make([]models.TableTwo, 0)
	s.blockThreeData = make([]models.TableThree, 0)
	return nil
}

// CountBlockOne возвращает количество записей TableOne.
func (s *Storage) CountBlockOne() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockOneData)
}

// CountBlockTwo возвращает количество записей TableTwo.
func (s *Storage) CountBlockTwo() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockTwoData)
}

// CountBlockThree возвращает количество записей TableThree.
func (s *Storage) CountBlockThree() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockThreeData)
}

func (s *Storage) AddBlockFourData(data []models.TableFour) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inclinometry = append(s.inclinometry, data...)
	return nil // В in-memory обычно нет ошибок добавления, кроме нехватки памяти (panic)
}
