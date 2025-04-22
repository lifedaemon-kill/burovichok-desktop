// internal/pkg/storage/inmemory/store.go
package inmemory

import (
	"sync"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

//можешно импортировать несколько файлов одного или разных типов, данные будут накапливаться. Кнопка "Очистить хранилище" удалит всё из памяти. Хранятся в Storage
//в будущем, если понадобятся данные для графиков или экспорта, соответствующие сервисы (chart, export) также должны будут получить экземпляр InMemoryBlocksStorage и вызывать методы GetAll...Data().

// InMemoryBlocksStorage определяет интерфейс для взаимодействия с хранилищем данных.
type InMemoryBlocksStorage interface {
	// Методы для добавления данных (принимают срезы, как возвращает парсер)
	PutTableOneData(data []models.TableOne) error
	PutTableTwoData(data []models.TableTwo) error
	PutTableThreeData(data []models.TableThree) error
	PutTableFourData(data []models.TableFour) error
	PutTableFiveData(data models.TableFive) error

	// Методы для получения всех данных (возвращают копии для безопасности)
	GetTableOneData() ([]models.TableOne, error)
	GetTableTwoData() ([]models.TableTwo, error)
	GetTableThreeData() ([]models.TableThree, error)
	GetTableFourData() ([]models.TableFour, error)
	GetTableFiveData() (models.TableFive, error)

	// Метод для очистки всего хранилища
	ClearAll() error

	// Можно добавить методы для получения количества записей (опционально)
	CountBlockOne() int
	CountBlockTwo() int
	CountBlockThree() int
}

// Storage реализует интерфейс storage.InMemoryBlocksStorage, храня данные в памяти.
type Storage struct {
	mu         sync.RWMutex // Mutex для безопасного доступа к данным
	blockOne   []models.TableOne
	blockTwo   []models.TableTwo
	blockThree []models.TableThree
	blockFour  []models.TableFour
	blockFive  models.TableFive
}

// NewInMemoryBlocksStorage создает новый экземпляр Storage.
func NewInMemoryBlocksStorage() InMemoryBlocksStorage { // Возвращаем интерфейс!
	return &Storage{
		blockOne:   make([]models.TableOne, 0),
		blockTwo:   make([]models.TableTwo, 0),
		blockThree: make([]models.TableThree, 0),
		blockFour:  make([]models.TableFour, 0),
	}
}

// PutTableOneData добавляет данные TableOne в хранилище.
func (s *Storage) PutTableOneData(data []models.TableOne) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOne = append(s.blockOne, data...)
	return nil // В in-memory обычно нет ошибок добавления, кроме нехватки памяти (panic)
}

// PutTableTwoData добавляет данные TableTwo в хранилище.
func (s *Storage) PutTableTwoData(data []models.TableTwo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockTwo = append(s.blockTwo, data...)
	return nil
}

// PutTableThreeData добавляет данные TableThree в хранилище.
func (s *Storage) PutTableThreeData(data []models.TableThree) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockThree = append(s.blockThree, data...)
	return nil
}

func (s *Storage) PutTableFourData(data []models.TableFour) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockFour = append(s.blockFour, data...)
	return nil // В in-memory обычно нет ошибок добавления, кроме нехватки памяти (panic)
}

func (s *Storage) PutTableFiveData(data models.TableFive) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockFive = data
	return nil
}

// GetAllBlockOneData возвращает копию всех данных TableOne.
func (s *Storage) GetTableOneData() ([]models.TableOne, error) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	// Возвращаем копию, чтобы внешние изменения не влияли на хранилище
	dataCopy := make([]models.TableOne, len(s.blockOne))
	copy(dataCopy, s.blockOne)
	return dataCopy, nil
}

// GetAllBlockTwoData возвращает копию всех данных TableTwo.
func (s *Storage) GetTableTwoData() ([]models.TableTwo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableTwo, len(s.blockTwo))
	copy(dataCopy, s.blockTwo)
	return dataCopy, nil
}

// GetAllBlockThreeData возвращает копию всех данных TableThree.
func (s *Storage) GetTableThreeData() ([]models.TableThree, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableThree, len(s.blockThree))
	copy(dataCopy, s.blockThree)
	return dataCopy, nil
}

// GetAllBlockFourData возвращает копию всех данных TableFour
func (s *Storage) GetTableFourData() ([]models.TableFour, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableFour, len(s.blockFour))
	copy(dataCopy, s.blockFour)
	return dataCopy, nil
}

func (s *Storage) GetTableFiveData() (models.TableFive, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := s.blockFive
	return dataCopy, nil
}

// ClearAll очищает все данные в хранилище.
func (s *Storage) ClearAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOne = make([]models.TableOne, 0)
	s.blockTwo = make([]models.TableTwo, 0)
	s.blockThree = make([]models.TableThree, 0)
	return nil
}

// CountBlockOne возвращает количество записей TableOne.
func (s *Storage) CountBlockOne() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockOne)
}

// CountBlockTwo возвращает количество записей TableTwo.
func (s *Storage) CountBlockTwo() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockTwo)
}

// CountBlockThree возвращает количество записей TableThree.
func (s *Storage) CountBlockThree() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockThree)
}
