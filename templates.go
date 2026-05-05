package pgsp

import (
	"fmt"
	"text/template"
)

var (
	FuncMap = template.FuncMap{
		"BoldStyle": BoldStyle.Render,
		"percent": func(a, b int64) string {
			if b == 0 {
				return "N/A"
			}
			return fmt.Sprintf("%.1f%%", float64(a)/float64(b)*100)
		},
	}

	AnalyzeTemplate = template.Must(template.New("analyze").Funcs(FuncMap).Parse(`Phase: {{.PHASE | BoldStyle}}; PID: {{.PID}}; RELID: {{.RELID}}; {{if .CurrentChildTableRelid }}{{.CurrentChildTableRelid}}{{end}}
{{ if .SampleBLKSTotal -}}
Sample blocks: {{.SampleBLKSScanned}}/{{.SampleBLKSTotal}} ({{percent .SampleBLKSScanned .SampleBLKSTotal}})
{{end -}}
{{- if .ExtStatsTotal -}}
External stats: {{.ExtStatsComputed}}/{{.ExtStatsTotal}} ({{percent .ExtStatsComputed .ExtStatsTotal}})
{{end -}}
{{- if .ChildTablesTotal -}}
Child tables: {{.ChildTablesDone}}/{{.ChildTablesTotal}} ({{percent .ChildTablesDone .ChildTablesTotal}})
{{end -}}
{{- if .DelayTime -}}
Delay time: {{printf "%.2f" .DelayTime}}s
{{end -}}`))

	BaseBackupTemplate = template.Must(template.New("baseBackup").Funcs(FuncMap).Parse(`Phase: {{.PHASE | BoldStyle}}; PID: {{.PID}}
{{ if .BackupTotal -}}
Backup: {{.BackupStreamed}}/{{.BackupTotal}}
{{end -}}
{{- if .TablespacesTotal -}}
Tablespaces: {{.TablespacesStreamed}}/{{.TablespacesTotal}} ({{percent .TablespacesStreamed .TablespacesTotal}})
{{end -}}`))

	ClusterTemplate = template.Must(template.New("cluster").Funcs(FuncMap).Parse(`Command: {{.Command | BoldStyle}}; Phase: {{.PHASE | BoldStyle}}; PID: {{.PID}}; RELID: {{.RELID}}; {{if .ClusterIndexRelid }}{{.ClusterIndexRelid}}{{end}}
{{ if .HeapTuplesScanned -}}
Heap Tuples Scanned: {{.HeapTuplesScanned}}; Written: {{.HeapTuplesWritten}}
{{end -}}
{{- if .HeapBlksTotal -}}
Heap Blocks: {{.HeapBlksScanned}}/{{.HeapBlksTotal}} ({{percent .HeapBlksScanned .HeapBlksTotal}})
{{end -}}
{{- if .IndexRebuildCount -}}
Index Rebuild Count: {{.IndexRebuildCount}}
{{end -}}`))

	CopyTemplate = template.Must(template.New("copy").Funcs(FuncMap).Parse(`Command: {{.COMMAND | BoldStyle}}; Type: {{.CTYPE | BoldStyle}}; PID: {{.PID}}; RELID: {{.RELID}};
{{ if .BYTESTotal -}}
Bytes: {{.BYTESProcessed}}/{{.BYTESTotal}} ({{percent .BYTESProcessed .BYTESTotal}})
{{end -}}
{{- if .TUPLESProcessed -}}
Tuples Processed: {{.TUPLESProcessed}}; Excluded: {{.TUPLESExcluded}}; Skipped: {{.TUPLESSkipped}}
{{end -}}`))

	CreateIndexTemplate = template.Must(template.New("createIndex").Funcs(FuncMap).Parse(`Command: {{.Command | BoldStyle}}; Phase: {{.PHASE | BoldStyle}}; PID: {{.PID}}; RELID: {{.RELID}}; {{if .IndexRelid }}{{.IndexRelid}}{{end}}
{{ if .BlocksTotal -}}
Blocks: {{.BlocksDone}}/{{.BlocksTotal}} ({{percent .BlocksDone .BlocksTotal}})
{{end -}}
{{- if .TuplesTotal -}}
Tuples: {{.TuplesDone}}/{{.TuplesTotal}} ({{percent .TuplesDone .TuplesTotal}})
{{end -}}
{{- if .PartitionsTotal -}}
Partitions: {{.PartitionsDone}}/{{.PartitionsTotal}} ({{percent .PartitionsDone .PartitionsTotal}})
{{end -}}`))

	VacuumTemplate = template.Must(template.New("vacuum").Funcs(FuncMap).Parse(`Phase: {{.PHASE | BoldStyle}}; PID: {{.PID}}; RELID: {{.RELID}}
{{ if .HeapBLKSTotal -}}
Heap Blocks Scanned: {{.HeapBLKSScanned}}/{{.HeapBLKSTotal}} ({{percent .HeapBLKSScanned .HeapBLKSTotal}}); Vacuumed: {{.HeapBLKSVacuumed}} ({{percent .HeapBLKSVacuumed .HeapBLKSTotal}})
{{end -}}
{{- if .IndexesTotal -}}
Indexes Vacuumed: {{.IndexVacuumCount}} ({{percent .IndexVacuumCount .IndexesTotal}})
Indexes: {{.IndexesProcessed}}/{{.IndexesTotal}} ({{percent .IndexesProcessed .IndexesTotal}})
{{end -}}
{{- if .DeadTupleBytes -}}
Dead Tuples: {{.DeadTupleBytes}} (max: {{.MaxDeadTupleBytes}})
{{end -}}
{{- if .DelayTime -}}
Delay time: {{printf "%.2f" .DelayTime}}s
{{end -}}`))
)
