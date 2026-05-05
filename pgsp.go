package pgsp

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"log"
	"net/netip"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type SPTaget string

const (
	SPAnalyze     SPTaget = "Analyze"
	SPCreateIndex SPTaget = "CreateIndex"
	SPVacuum      SPTaget = "Vacuum"
	SPCluster     SPTaget = "Cluster"
	SPBaseBackup  SPTaget = "BaseBackup"
	SPCopy        SPTaget = "Copy"
)

type SPTable struct {
	Enable bool
	Get    func(ctx context.Context, pgsp *Pgsp) ([]Progress, error)
}

type StatProgress map[SPTaget]*SPTable

type Pgsp struct {
	DB           *sqlx.DB
	MetaDBs      map[string]*sqlx.DB
	MetaNames    map[string]string
	BaseConfig   *pq.Config
	StatProgress StatProgress
}

type OidAvailable struct {
	RELID   int
	DATNAME string
}

func (o *OidAvailable) Key() string {
	return fmt.Sprintf("%s.%d", o.DATNAME, o.RELID)
}

type Progress interface {
	Header() string
	Pid() int
	Color() (color.Color, color.Color)
	Display() string
	Progress() float64
}

func New(dsn string, resolveNames bool) (*Pgsp, error) {
	db, err := Connect(dsn)
	if err != nil {
		return nil, err
	}
	var baseConfig *pq.Config
	if resolveNames {
		config, err := pq.NewConfig(dsn)
		if err != nil {
			return nil, err
		}
		config.Database = "postgres"
		baseConfig = &config
	}

	monitor := NewMonitor()
	return &Pgsp{
		DB:           db,
		BaseConfig:   baseConfig,
		MetaDBs:      make(map[string]*sqlx.DB),
		MetaNames:    make(map[string]string),
		StatProgress: monitor,
	}, nil
}

func NewMonitor() StatProgress {
	return StatProgress{
		SPAnalyze: {
			Get: GetAnalyze,
		},
		SPCreateIndex: {
			Get: GetCreateIndex,
		},
		SPVacuum: {
			Get: GetVacuum,
		},
		SPCluster: {
			Get: GetCluster,
		},
		SPBaseBackup: {
			Get: GetBaseBackup,
		},
		SPCopy: {
			Get: GetCopy,
		},
	}
}

func getDSN(config *pq.Config) string {
	dsn := ""
	t := reflect.TypeOf(*config)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		j := field.Tag.Get("postgres")
		if j == "" || j == "-" {
			continue
		}
		value := reflect.ValueOf(*config).Field(i).Interface()
		if value == nil {
			continue
		}
		valueAddr, ok := value.(netip.Addr)
		if ok && !valueAddr.IsValid() {
			continue
		}

		valueStr := ""
		log.Printf("%s Value %v -> '%s'\n", j, value, valueStr)
		valueTime, ok := value.(time.Duration)
		if ok && valueTime.Seconds() < 0.0001 {
			valueStr = "0"
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}
		if valueStr == "" {
			continue
		}
		dsn += fmt.Sprintf("%s='%s' ", j, valueStr)
	}
	log.Printf("DSN: %s\n", dsn)
	return dsn
}

func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (p *Pgsp) ConnectForDatabase(database string) (*sqlx.DB, error) {
	if db, ok := p.MetaDBs[database]; ok {
		return db, nil
	}

	if p.BaseConfig == nil {
		return nil, fmt.Errorf("base config is not set")
	}
	config := (*p.BaseConfig).Clone()
	config.Database = database
	db, err := sqlx.Connect("postgres", getDSN(&config))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	p.MetaDBs[database] = db
	return db, nil
}

func (p *Pgsp) GetItemName(oidAvailable OidAvailable) (string, error) {
	if p.BaseConfig == nil {
		return "", fmt.Errorf("base config is not set")
	}
	key := oidAvailable.Key()
	if name, ok := p.MetaNames[key]; ok {
		return name, nil
	}
	db, err := p.ConnectForDatabase(oidAvailable.DATNAME)
	if err != nil {
		return "", err
	}
	var name string
	err = db.Get(&name, "SELECT relname FROM pg_catalog.pg_class WHERE oid = $1", oidAvailable.RELID)
	if err != nil {
		name = fmt.Sprintf("not-found-%d", oidAvailable.RELID)
	}
	p.MetaNames[key] = name
	return name, nil
}

func (p *Pgsp) Disconnect() []error {
	errors := make([]error, 0)
	for _, db := range p.MetaDBs {
		if err := db.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	if err := p.DB.Close(); err != nil {
		errors = append(errors, err)
	}
	return errors
}

func (p *Pgsp) Targets(target []string) {
	if len(target) != 0 {
		enableF := false
		for _, t := range target {
			if v, ok := p.StatProgress[SPTaget(t)]; ok {
				enableF = true
				v.Enable = true
			}
		}
		// Return if there is even one target.
		if enableF {
			return
		}
	}

	// All targets.
	for _, v := range p.StatProgress {
		v.Enable = true
	}
}

func (p *Pgsp) TargetString() string {
	var ms []string
	for n, v := range p.StatProgress {
		if v.Enable {
			ms = append(ms, string(n))
		}
	}
	sort.Strings(ms)
	return strings.Join(ms, " ")
}

func (p *Pgsp) ConnectionCount() int {
	return len(p.MetaDBs) + 1
}

func buildQuery(tableName string, columns []string) string {
	buff := new(bytes.Buffer)
	buff.WriteString("SELECT ")
	buff.WriteString(strings.Join(columns, ", "))
	buff.WriteString(" FROM ")
	buff.WriteString(tableName)
	return buff.String()
}

func getColumns(s interface{}, db_only bool) []string {
	t := reflect.TypeOf(s)
	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		j := field.Tag.Get("db")
		if j == "" {
			if db_only {
				continue
			}
			j = field.Name
		}
		columns = append(columns, j)
	}
	return columns
}
