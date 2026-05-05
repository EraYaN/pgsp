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

// pg_stat_progress_copy
type Copy struct {
	PID             int    `db:"pid"`
	DATID           int    `db:"datid"`
	DATNAME         string `db:"datname"`
	RELID           int    `db:"relid"`
	COMMAND         string `db:"command"`
	CTYPE           string `db:"type"`
	BYTESProcessed  int64  `db:"bytes_processed"`
	BYTESTotal      int64  `db:"bytes_total"`
	TUPLESProcessed int64  `db:"tuples_processed"`
	TUPLESExcluded  int64  `db:"tuples_excluded"`
	TUPLESSkipped   int64  `db:"tuples_skipped"`
	RELNAME         string
}

var (
	CopyTableName = "pg_stat_progress_copy"
	CopyQuery     string
	CopyColumns   []string
	CopyHeaders   []string
)

func GetCopy(ctx context.Context, pgsp *Pgsp) ([]Progress, error) {
	if len(CopyColumns) == 0 {
		CopyColumns = getColumns(Copy{}, true)
	}
	if len(CopyHeaders) == 0 {
		CopyHeaders = getColumns(Copy{}, false)
	}
	if CopyQuery == "" {
		CopyQuery = buildQuery(CopyTableName, CopyColumns)
	}
	return selectCopy(ctx, pgsp, CopyQuery)
}

func selectCopy(ctx context.Context, pgsp *Pgsp, query string) ([]Progress, error) {
	db := pgsp.DB
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []Progress
	for rows.Next() {
		var row Copy
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
		as = append(as, row)
	}
	return as, nil
}

func (v Copy) Header() string {
	return fmt.Sprintf("%s: %s, %s", CopyTableName, v.DATNAME, v.RELNAME)
}

func (v Copy) Pid() int {
	return v.PID
}

func (v Copy) Color() (color.Color, color.Color) {
	return color.RGBA{R: 90, G: 246, B: 255}, color.RGBA{R: 124, G: 255, B: 203}
}

func (v Copy) Display() string {
	value := str.ToStrStruct(v)
	buff := new(bytes.Buffer)
	t := tablewriter.NewWriter(buff)
	t.Header(CopyHeaders)
	t.Append(value)
	t.Render()
	return buff.String()
}

func (v Copy) Progress() float64 {
	if v.BYTESTotal == 0 {
		return float64(0.5)
	}
	return float64(v.BYTESProcessed) / float64(v.BYTESTotal)
}
