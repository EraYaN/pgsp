package pgsp

import (
	"bytes"
	"log"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func CreateInput(v interface{}) interface{} {
	modelType := reflect.TypeOf(v)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		t := field.Type

		switch t.Kind() {
		case reflect.Int, reflect.Int64:
			reflect.ValueOf(&v).Field(i).SetInt(1)
		case reflect.String:
			reflect.ValueOf(&v).Field(i).SetString("string")
		case reflect.Float32, reflect.Float64:
			reflect.ValueOf(&v).Field(i).SetFloat(1.0)
		case reflect.Bool:
			reflect.ValueOf(&v).Field(i).SetBool(true)
		}
	}
	log.Printf("Model %v", v)
	return v
}

func TestTemplates(t *testing.T) {
	tests := []struct {
		name  string
		model Progress
	}{
		{
			name:  "Analyze",
			model: Analyze{},
		},
		{
			name:  "BaseBackup",
			model: BaseBackup{},
		},
		{
			name:  "Cluster",
			model: Cluster{},
		},
		{
			name:  "Copy",
			model: Copy{},
		},
		{
			name:  "CreateIndex",
			model: CreateIndex{},
		},
		{
			name:  "Vacuum",
			model: Vacuum{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var doc bytes.Buffer
			err := tt.model.Template().Execute(&doc, tt.model)
			if err != nil {
				t.Errorf("%s.Template() = \n%v\n", tt.name, err)
			}
		})
	}
}
