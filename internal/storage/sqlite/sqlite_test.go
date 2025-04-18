package sqlite

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Выполняется первым
func TestOpen(t *testing.T) {
	DSN := "tests/database/example.db"

	db, err := NewDB(config.DBConf{DSN: DSN})
	assert.Nil(t, err)

	_, err = NewGuidebookRepository(db)
	assert.Nil(t, err)
	_, err = NewBlockRepository(db)
	assert.Nil(t, err)
}
