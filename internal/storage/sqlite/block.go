package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func NewBlockStorage(db *sqlx.DB) (BlocksStorage, error) {
	return &sqlite{
		db: db,
	}, nil
}

func (s sqlite) GetAllTableFive() ([]models.TableFive, error) {
	var reports []models.TableFive

	err := s.db.Select(&reports, "SELECT * FROM reports")
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s sqlite) AddBlockFive(data models.TableFive) (reportID int64, err error) {
	query := `
		INSERT INTO reports (
			field_name,
			field_number,
			cluster_number,
			horizon,
			start_time,
			end_time,
			instrument_type,
			instrument_number,
			measure_depth,
			true_vertical_depth,
			true_vertical_depth_sub_sea,
			vdp_measured_depth,
			vdp_true_vertical_depth,
			vdp_true_vertical_depth_sea,
			diff_instrument_vdp,
			density_oil,
			density_liquid_stopped,
			density_liquid_working,
			pressure_diff_stopped,
			pressure_diff_working
		) VALUES (
			:field_name,
			:field_number,
			:cluster_number,
			:horizon,
			:start_time,
			:end_time,
			:instrument_type,
			:instrument_number,
			:measure_depth,
			:true_vertical_depth,
			:true_vertical_depth_sub_sea,
			:vdp_measured_depth,
			:vdp_true_vertical_depth,
			:vdp_true_vertical_depth_sea,
			:diff_instrument_vdp,
			:density_oil,
			:density_liquid_stopped,
			:density_liquid_working,
			:pressure_diff_stopped,
			:pressure_diff_working
		)
	`

	result, err := s.db.NamedExec(query, data)
	if err != nil {
		return 0, err
	}
	// Получаем ID вставленной записи
	reportID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return reportID64, nil
}
