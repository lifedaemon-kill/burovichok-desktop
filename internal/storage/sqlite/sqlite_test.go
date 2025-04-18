package sqlite

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Выполняется первым
func TestOpen(t *testing.T) {
	DSN := "tests/database/example.db"

	db, err := NewDB(config.DBConf{DSN: DSN})
	assert.Nil(t, err)

	_, err = NewGuidebookStorage(db)
	assert.Nil(t, err)
	_, err = NewBlockStorage(db)
	assert.Nil(t, err)
}
