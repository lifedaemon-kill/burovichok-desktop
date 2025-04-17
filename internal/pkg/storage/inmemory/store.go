// internal/pkg/storage/inmemory/store.go
package inmemory

import (
	"sync"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/storage" // Импортируем наш интерфейс
)

//можешно импортировать несколько файлов одного или разных типов, данные будут накапливаться. Кнопка "Очистить хранилище" удалит всё из памяти. Хранятся в InMemoryStore
//в будущем, если понадобятся данные для графиков или экспорта, соответствующие сервисы (chart, export) также должны будут получить экземпляр Storage и вызывать методы GetAll...Data().

// InMemoryStore реализует интерфейс storage.Storage, храня данные в памяти.
type InMemoryStore struct {
	mu             sync.RWMutex // Mutex для безопасного доступа к данным
	blockOneData   []models.TableOne
	blockTwoData   []models.TableTwo
	blockThreeData []models.TableThree
}

// NewStore создает новый экземпляр InMemoryStore.
func NewStore() storage.Storage { // Возвращаем интерфейс!
	return &InMemoryStore{
		blockOneData:   make([]models.TableOne, 0),
		blockTwoData:   make([]models.TableTwo, 0),
		blockThreeData: make([]models.TableThree, 0),
	}
}

// AddBlockOneData добавляет данные TableOne в хранилище.
func (s *InMemoryStore) AddBlockOneData(data []models.TableOne) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOneData = append(s.blockOneData, data...)
	return nil // В in-memory обычно нет ошибок добавления, кроме нехватки памяти (panic)
}

// AddBlockTwoData добавляет данные TableTwo в хранилище.
func (s *InMemoryStore) AddBlockTwoData(data []models.TableTwo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockTwoData = append(s.blockTwoData, data...)
	return nil
}

// AddBlockThreeData добавляет данные TableThree в хранилище.
func (s *InMemoryStore) AddBlockThreeData(data []models.TableThree) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockThreeData = append(s.blockThreeData, data...)
	return nil
}

// GetAllBlockOneData возвращает копию всех данных TableOne.
func (s *InMemoryStore) GetAllBlockOneData() ([]models.TableOne, error) {
	s.mu.RLock() // Блокировка на чтение
	defer s.mu.RUnlock()
	// Возвращаем копию, чтобы внешние изменения не влияли на хранилище
	dataCopy := make([]models.TableOne, len(s.blockOneData))
	copy(dataCopy, s.blockOneData)
	return dataCopy, nil
}

// GetAllBlockTwoData возвращает копию всех данных TableTwo.
func (s *InMemoryStore) GetAllBlockTwoData() ([]models.TableTwo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableTwo, len(s.blockTwoData))
	copy(dataCopy, s.blockTwoData)
	return dataCopy, nil
}

// GetAllBlockThreeData возвращает копию всех данных TableThree.
func (s *InMemoryStore) GetAllBlockThreeData() ([]models.TableThree, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dataCopy := make([]models.TableThree, len(s.blockThreeData))
	copy(dataCopy, s.blockThreeData)
	return dataCopy, nil
}

// ClearAll очищает все данные в хранилище.
func (s *InMemoryStore) ClearAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockOneData = make([]models.TableOne, 0)
	s.blockTwoData = make([]models.TableTwo, 0)
	s.blockThreeData = make([]models.TableThree, 0)
	return nil
}

// CountBlockOne возвращает количество записей TableOne.
func (s *InMemoryStore) CountBlockOne() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockOneData)
}

// CountBlockTwo возвращает количество записей TableTwo.
func (s *InMemoryStore) CountBlockTwo() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockTwoData)
}

// CountBlockThree возвращает количество записей TableThree.
func (s *InMemoryStore) CountBlockThree() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blockThreeData)
}
