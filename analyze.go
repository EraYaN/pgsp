package pgsp

import (
	"context"
	"fmt"
	"image/color"
	"text/template"

	_ "github.com/lib/pq"
)

// pg_stat_progress_analyze.
type Analyze struct {
	PID                      int     `db:"pid"`
	DATID                    int     `db:"datid"`
	DATNAME                  string  `db:"datname"`
	RELID                    int     `db:"relid"`
	PHASE                    string  `db:"phase"`
	SampleBLKSTotal          int64   `db:"sample_blks_total"`
	SampleBLKSScanned        int64   `db:"sample_blks_scanned"`
	ExtStatsTotal            int64   `db:"ext_stats_total"`
	ExtStatsComputed         int64   `db:"ext_stats_computed"`
	ChildTablesTotal         int64   `db:"child_tables_total"`
	ChildTablesDone          int64   `db:"child_tables_done"`
	CurrentChildTableRelid   int     `db:"current_child_table_relid"`
	DelayTime                float64 `db:"delay_time"`
	RELNAME                  string
	CurrentChildTableRelName string
}

var (
	AnalyzeTableName = "pg_stat_progress_analyze"
	AnalyzeQuery     string
	AnalyzeColumns   []string
	AnalyzeHeaders   []string
)

func GetAnalyze(ctx context.Context, pgsp *Pgsp) ([]Progress, error) {
	if len(AnalyzeColumns) == 0 {
		AnalyzeColumns = getColumns(Analyze{}, true)
	}
	if len(AnalyzeHeaders) == 0 {
		AnalyzeHeaders = getColumns(Analyze{}, false)
	}
	if AnalyzeQuery == "" {
		AnalyzeQuery = buildQuery(AnalyzeTableName, AnalyzeColumns)
	}
	return selectAnalyze(ctx, pgsp, AnalyzeQuery)
}

func selectAnalyze(ctx context.Context, pgsp *Pgsp, query string) ([]Progress, error) {
	db := pgsp.DB
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row Analyze
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		row.RELNAME, err = pgsp.GetItemName(OidAvailable{
			RELID:   row.RELID,
			DATNAME: row.DATNAME,
		})
		if err != nil {
			return nil, err
		}
		row.CurrentChildTableRelName, err = pgsp.GetItemName(OidAvailable{
			RELID:   row.CurrentChildTableRelid,
			DATNAME: row.DATNAME,
		})
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v Analyze) Header() string {
	if v.CurrentChildTableRelid != 0 {
		return fmt.Sprintf("%s: %s, %s, %s", AnalyzeTableName, v.DATNAME, v.RELNAME, v.CurrentChildTableRelName)
	}
	return fmt.Sprintf("%s: %s, %s", AnalyzeTableName, v.DATNAME, v.RELNAME)
}

func (v Analyze) Pid() int {
	return v.PID
}

func (v Analyze) Color() (color.Color, color.Color) {
	return color.RGBA{R: 255, G: 124, B: 203}, color.RGBA{R: 253, G: 255, B: 140}
}

func (v Analyze) Template() *template.Template {
	return AnalyzeTemplate
}

func (v Analyze) Progress() float64 {
	if v.ChildTablesTotal != 0 {
		return float64(v.ChildTablesDone) / float64(v.ChildTablesTotal)
	}
	if v.ExtStatsTotal != 0 {
		return float64(v.ExtStatsComputed) / float64(v.ExtStatsTotal)
	}
	return float64(v.SampleBLKSScanned) / float64(v.SampleBLKSTotal)
}
