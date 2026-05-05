package pgsp

import (
	"context"
	"fmt"
	"image/color"
	"text/template"

	_ "github.com/lib/pq"
)

// pg_stat_progress_cluster.
type Cluster struct {
	PID                 int    `db:"pid"`
	DATID               int    `db:"datid"`
	DATNAME             string `db:"datname"`
	RELID               int    `db:"relid"`
	Command             string `db:"command"`
	PHASE               string `db:"phase"`
	ClusterIndexRelid   int    `db:"cluster_index_relid"`
	HeapTuplesScanned   int64  `db:"heap_tuples_scanned"`
	HeapTuplesWritten   int64  `db:"heap_tuples_written"`
	HeapBlksTotal       int64  `db:"heap_blks_total"`
	HeapBlksScanned     int64  `db:"heap_blks_scanned"`
	IndexRebuildCount   int64  `db:"index_rebuild_count"`
	RELNAME             string
	ClusterIndexRelname string
}

var (
	ClusterTableName = "pg_stat_progress_cluster"
	ClusterQuery     string
	ClusterColumns   []string
	ClusterHeaders   []string
)

func GetCluster(ctx context.Context, pgsp *Pgsp) ([]Progress, error) {
	if len(ClusterColumns) == 0 {
		ClusterColumns = getColumns(Cluster{}, true)
	}
	if len(ClusterHeaders) == 0 {
		ClusterHeaders = getColumns(Cluster{}, false)
	}
	if ClusterQuery == "" {
		ClusterQuery = buildQuery(ClusterTableName, ClusterColumns)
	}
	return selectCluster(ctx, pgsp, ClusterQuery)
}

func selectCluster(ctx context.Context, pgsp *Pgsp, query string) ([]Progress, error) {
	db := pgsp.DB
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row Cluster
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
		row.ClusterIndexRelname, err = pgsp.GetItemName(OidAvailable{
			RELID:   row.ClusterIndexRelid,
			DATNAME: row.DATNAME,
		})
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v Cluster) Header() string {
	if v.ClusterIndexRelid != 0 {
		return fmt.Sprintf("%s: %s, %s, %s", ClusterTableName, v.DATNAME, v.RELNAME, v.ClusterIndexRelname)
	}
	return fmt.Sprintf("%s: %s, %s", ClusterTableName, v.DATNAME, v.RELNAME)
}

func (v Cluster) Pid() int {
	return v.PID
}

func (v Cluster) Color() (color.Color, color.Color) {
	return color.RGBA{R: 90, G: 86, B: 224}, color.RGBA{R: 238, G: 111, B: 248}
}

func (v Cluster) Template() *template.Template {

	return ClusterTemplate
}

func (v Cluster) Progress() float64 {
	return float64(v.HeapBlksScanned) / float64(v.HeapBlksTotal)
}
