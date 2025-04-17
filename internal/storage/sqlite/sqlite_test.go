package sqlite

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpen(t *testing.T) {
	DSN := "tests/database/example.db"

	_, err := New(config.DBConf{DSN: DSN})
	assert.Nil(t, err)
}
