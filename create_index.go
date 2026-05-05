package pgsp

import (
	"bytes"
	"context"
	"fmt"
	"image/color"

	"github.com/EraYaN/pgsp/str"
	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

// pg_stat_progress_create_index
type CreateIndex struct {
	PID             int    `db:"pid"`
	DATID           int    `db:"datid"`
	DATNAME         string `db:"datname"`
	RELID           int    `db:"relid"`
	IndexRelid      int    `db:"index_relid"`
	Command         string `db:"command"`
	PHASE           string `db:"phase"`
	LockersTotal    int64  `db:"lockers_total"`
	LockersDone     int64  `db:"lockers_done"`
	LockersPid      int64  `db:"current_locker_pid"`
	BlocksTotal     int64  `db:"blocks_total"`
	BlocksDone      int64  `db:"blocks_done"`
	TuplesTotal     int64  `db:"tuples_total"`
	TuplesDone      int64  `db:"tuples_done"`
	PartitionsTotal int64  `db:"partitions_total"`
	PartitionsDone  int64  `db:"partitions_done"`
	RELNAME         string
	IndexRelname    string
}

var CreateIndexTableName = "pg_stat_progress_create_index"

var (
	CreateIndexQuery   string
	CreateIndexColumns []string
	CreateIndexHeaders []string
)

func GetCreateIndex(ctx context.Context, pgsp *Pgsp) ([]Progress, error) {
	if len(CreateIndexColumns) == 0 {
		CreateIndexColumns = getColumns(CreateIndex{}, true)
	}
	if len(CreateIndexHeaders) == 0 {
		CreateIndexHeaders = getColumns(CreateIndex{}, false)
	}
	if CreateIndexQuery == "" {
		CreateIndexQuery = buildQuery(CreateIndexTableName, CreateIndexColumns)
	}
	return selectCreateIndex(ctx, pgsp, CreateIndexQuery)
}

func selectCreateIndex(ctx context.Context, pgsp *Pgsp, query string) ([]Progress, error) {
	db := pgsp.DB
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row CreateIndex
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
		row.IndexRelname, err = pgsp.GetItemName(OidAvailable{
			RELID:   row.IndexRelid,
			DATNAME: row.DATNAME,
		})
		if err != nil {
			return nil, err
		}
		as = append(as, row)
	}
	return as, rows.Err()
}

func (v CreateIndex) Header() string {
	if v.IndexRelid != 0 {
		return fmt.Sprintf("%s: %s, %s, %s", CreateIndexTableName, v.DATNAME, v.RELNAME, v.IndexRelname)
	}
	return fmt.Sprintf("%s: %s, %s", CreateIndexTableName, v.DATNAME, v.RELNAME)
}

func (v CreateIndex) Pid() int {
	return v.PID
}

func (v CreateIndex) Color() (color.Color, color.Color) {
	return color.RGBA{R: 238, G: 111, B: 248}, color.RGBA{R: 90, G: 86, B: 224}
}

func (v CreateIndex) Display() string {
	value := str.ToStrStruct(v)
	buff := new(bytes.Buffer)

	t := tablewriter.NewWriter(buff)
	t.Header(CreateIndexHeaders)
	t.Append(value)
	t.Render()

	return buff.String()
}

func (v CreateIndex) Progress() float64 {
	if v.BlocksTotal != 0 {
		return float64(v.BlocksDone) / float64(v.BlocksTotal)
	}
	if v.PartitionsTotal != 0 {
		return float64(v.PartitionsDone) / float64(v.PartitionsTotal)
	}
	return float64(v.TuplesDone) / float64(v.TuplesTotal)
}
