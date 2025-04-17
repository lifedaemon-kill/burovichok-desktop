// internal/pkg/storage/inmemory/store.go
package inmemory

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
	"sync"
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

//можешно импортировать несколько файлов одного или разных типов, данные будут накапливаться. Кнопка "Очистить хранилище" удалит всё из памяти. Хранятся в InMemoryStore
//в будущем, если понадобятся данные для графиков или экспорта, соответствующие сервисы (chart, export) также должны будут получить экземпляр BlocksStorage и вызывать методы GetAll...Data().

// InMemoryStore реализует интерфейс storage.BlocksStorage, храня данные в памяти.
type InMemoryStore struct {
	mu             sync.RWMutex // Mutex для безопасного доступа к данным
	blockOneData   []models.BlockOne
	blockTwoData   []models.BlockTwo
	blockThreeData []models.BlockThree
}

// NewStore создает новый экземпляр InMemoryStore.
func NewStore() BlocksStorage { // Возвращаем интерфейс!
	return &InMemoryStore{
		blockOneData:   make([]models.BlockOne, 0),
		blockTwoData:   make([]models.BlockTwo, 0),
		blockThreeData: make([]models.BlockThree, 0),
	}
}

// AddBlockOneData добавляет данные BlockOne в хранилище.
func (s *InMemoryStore) AddBlockOneData(data []models.BlockOne) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOneData = append(s.blockOneData, data...)
	return nil // В in-memory обычно нет ошибок добавления, кроме нехватки памяти (panic)
}

// AddBlockTwoData добавляет данные BlockTwo в хранилище.
func (s *InMemoryStore) AddBlockTwoData(data []models.BlockTwo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockTwoData = append(s.blockTwoData, data...)
	return nil
}

// AddBlockThreeData добавляет данные BlockThree в хранилище.
func (s *InMemoryStore) AddBlockThreeData(data []models.BlockThree) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockThreeData = append(s.blockThreeData, data...)
	return nil
}

// GetAllBlockOneData возвращает копию всех данных BlockOne.
func (s *InMemoryStore) GetAllBlockOneData() ([]models.BlockOne, error) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	// Возвращаем копию, чтобы внешние изменения не влияли на хранилище
	dataCopy := make([]models.BlockOne, len(s.blockOneData))
	copy(dataCopy, s.blockOneData)
	return dataCopy, nil
}

// GetAllBlockTwoData возвращает копию всех данных BlockTwo.
func (s *InMemoryStore) GetAllBlockTwoData() ([]models.BlockTwo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.BlockTwo, len(s.blockTwoData))
	copy(dataCopy, s.blockTwoData)
	return dataCopy, nil
}

// GetAllBlockThreeData возвращает копию всех данных BlockThree.
func (s *InMemoryStore) GetAllBlockThreeData() ([]models.BlockThree, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.BlockThree, len(s.blockThreeData))
	copy(dataCopy, s.blockThreeData)
	return dataCopy, nil
}

// ClearAll очищает все данные в хранилище.
func (s *InMemoryStore) ClearAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOneData = make([]models.BlockOne, 0)
	s.blockTwoData = make([]models.BlockTwo, 0)
	s.blockThreeData = make([]models.BlockThree, 0)
	return nil
}

// CountBlockOne возвращает количество записей BlockOne.
func (s *InMemoryStore) CountBlockOne() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockOneData)
}

// CountBlockTwo возвращает количество записей BlockTwo.
func (s *InMemoryStore) CountBlockTwo() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockTwoData)
}

// CountBlockThree возвращает количество записей BlockThree.
func (s *InMemoryStore) CountBlockThree() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockThreeData)
}
