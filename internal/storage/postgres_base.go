package storage

import (
	"reflect"

	"github.com/jackc/pgx/v5"
)

func StructToNamedArgs(entity interface{}) pgx.NamedArgs {
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)

	args := pgx.NamedArgs{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			args[dbTag] = value
		}
	}
	return args
}
