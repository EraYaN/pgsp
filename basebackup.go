package pgsp

import (
	"context"
	"database/sql"
	"image/color"
	"text/template"

	_ "github.com/lib/pq"
)

// pg_stat_progress_basebackup.
type BaseBackup struct {
	PID                 int           `db:"pid"`
	PHASE               string        `db:"phase"`
	BackupTotal         sql.NullInt64 `db:"backup_total"`
	BackupStreamed      int64         `db:"backup_streamed"`
	TablespacesTotal    int64         `db:"tablespaces_total"`
	TablespacesStreamed int64         `db:"tablespaces_streamed"`
}

var (
	BaseBackupTableName = "pg_stat_progress_basebackup"
	BaseBackupQuery     string
	BaseBackupColumns   []string
	BaseBackupHeaders   []string
)

func GetBaseBackup(ctx context.Context, pgsp *Pgsp) ([]Progress, error) {
	if len(BaseBackupColumns) == 0 {
		BaseBackupColumns = getColumns(BaseBackup{}, true)
	}
	if len(BaseBackupHeaders) == 0 {
		BaseBackupHeaders = getColumns(BaseBackup{}, false)
	}
	if BaseBackupQuery == "" {
		BaseBackupQuery = buildQuery(BaseBackupTableName, BaseBackupColumns)
	}
	return selectBaseBackup(ctx, pgsp, BaseBackupQuery)
}

func selectBaseBackup(ctx context.Context, pgsp *Pgsp, query string) ([]Progress, error) {
	db := pgsp.DB
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row BaseBackup
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v BaseBackup) Header() string {
	return BaseBackupTableName
}

func (v BaseBackup) Pid() int {
	return v.PID
}

func (v BaseBackup) Color() (color.Color, color.Color) {
	return color.RGBA{R: 253, G: 255, B: 140}, color.RGBA{R: 255, G: 124, B: 203}
}

func (v BaseBackup) Template() *template.Template {
	return BaseBackupTemplate
}

func (v BaseBackup) Progress() float64 {
	total := v.BackupTotal.Int64
	if total != 0 {
		return float64(v.BackupStreamed) / float64(total)
	}
	return float64(v.TablespacesStreamed) / float64(v.TablespacesTotal)
}
