package postgres

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func dollarQuote(s string) string {
	return "$escape$" + s + "$escape$"
}

func QueryBuilder(table string, pkey string, data map[string]interface{}) (string, error) {
	if data["after"] == nil {
		return "", errors.New("received message with empty data field. Returning empty query")
	}

	after := data["after"].(map[string]interface{})
	var fields []string
	var values []string

	for k, v := range after {
		fields = append(fields, k)

		if v != nil {
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				values = append(values, dollarQuote(v.(string)))
			case reflect.Int:
				values = append(values, strconv.Itoa(v.(int)))
			case reflect.Float64:
				values = append(values, strconv.FormatFloat(v.(float64), 'f', -1, 64))
			default:
				values = append(values, "NULL")
			}
		} else {
			values = append(values, "NULL")
		}
	}

	builder := strings.Builder{}

	for index, field := range fields {
		value := values[index]
		if index != len(fields)-1 {
			builder.WriteString(fmt.Sprintf(" %s = %s, ", field, value))
		} else {
			builder.WriteString(fmt.Sprintf(" %s = %s;", field, value))
		}
	}

	updateFields := builder.String()

	upsertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET",
		table,
		strings.Join(fields, ", "),
		strings.Join(values, ", "),
		pkey,
	)

	return upsertStmt + updateFields, nil
}
