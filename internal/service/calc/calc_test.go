package calc

import (
	"fmt"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type plug struct {
	DataLoader
	logger.Logger
}

func getPlugService() *Service {
	return &Service{
		dataLoader: &plug{},
		logger:     &plug{},
	}
}
func TestTableOne(t *testing.T) {
	in := models.TableOne{
		PressureDepth:    28.0,
		TemperatureDepth: 7.0,
		Timestamp:        time.Now(),
	}

	s := getPlugService()
	from, _ := time.Parse(time.DateOnly, "2020-01-01")
	to, _ := time.Parse(time.DateOnly, "2025-01-01")
	idlefrom, _ := time.Parse(time.DateOnly, "2019-01-01")
	idleto, _ := time.Parse(time.DateOnly, "2019-01-02")

	conf := models.OperationConfig{
		PressureUnit: "kgf/cm2",
		DepthDiff:    4,
		WorkStart:    from,
		WorkEnd:      to,
		WorkDensity:  1.0,
		IdleStart:    idlefrom,
		IdleEnd:      idleto,
		IdleDensity:  1.0,
	}
	out := s.CalcTableOne(in, conf)

	fmt.Println("in: ", in, "\nout:", out)
	assert.NotNil(t, out.PressureAtVDP)
}
