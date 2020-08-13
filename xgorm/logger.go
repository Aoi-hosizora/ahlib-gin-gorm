package xgorm

import (
	"database/sql/driver"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// logrus.Logger

type GormLogrus struct {
	logger *logrus.Logger
}

func NewGormLogrus(logger *logrus.Logger) *GormLogrus {
	return &GormLogrus{logger: logger}
}

// See gorm.LogFormatter for details.
func (g *GormLogrus) Print(v ...interface{}) {
	if len(v) <= 1 {
		return
	}

	// info
	if len(v) == 2 {
		g.logger.WithFields(logrus.Fields{
			"module": "gorm",
			"type":   v[0],
			"info":   v[1],
		}).Infof(fmt.Sprintf("[Gorm] %v", v[1]))
		return
	}

	// sql
	if v[0] == "sql" {
		source := v[1]
		duration := v[2]
		sql := render(v[3].(string), v[4])
		rows := v[5]
		g.logger.WithFields(logrus.Fields{
			"module":   "gorm",
			"type":     "sql",
			"source":   source,
			"duration": duration,
			"sql":      sql,
			"rows":     rows,
		}).Info(fmt.Sprintf("[Gorm] #: %3d | %10s | %s | %s", rows, duration, sql, source))
		return
	}

	// other
	g.logger.WithFields(logrus.Fields{
		"module": "gorm",
		"type":   v[0],
	}).Info(fmt.Sprintf("[Gorm] %s", fmt.Sprint(v[2:]...)))
}

// log.Logger

type GormLogger struct {
	logger *log.Logger
}

func NewGormLogger(logger *log.Logger) *GormLogger {
	return &GormLogger{logger: logger}
}

func (g *GormLogger) Print(v ...interface{}) {
	if len(v) <= 1 {
		return
	}

	if len(v) == 2 {
		g.logger.Printf("[Gorm] %v", v[1])
		return
	}

	if v[0] == "sql" {
		source := v[1]
		duration := v[2]
		sql := render(v[3].(string), v[4])
		rows := v[5]
		g.logger.Printf("[Gorm] #: %3d | %10s | %s | %s", rows, duration, sql, source)
		return
	}

	g.logger.Printf("[Gorm] %s", fmt.Sprint(v[2:]...))
}

// render

func render(sql string, param interface{}) string {
	values := make([]interface{}, 0)
	for _, value := range param.([]interface{}) {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() { // valid
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok { // time
				values = append(values, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
			} else if b, ok := value.([]byte); ok { // bytes
				values = append(values, fmt.Sprintf("'%v'", string(b)))
			} else if r, ok := value.(driver.Valuer); ok { // driver
				if value, err := r.Value(); err == nil && value != nil {
					values = append(values, fmt.Sprintf("'%v'", value))
				} else {
					values = append(values, "NULL")
				}
			} else { // other value
				values = append(values, fmt.Sprintf("'%v'", value))
			}
		} else { // invalid
			values = append(values, fmt.Sprintf("'%v'", value))
		}
	}

	result := fmt.Sprintf(sqlRegexp.ReplaceAllString(sql, "%v"), values...)
	return strings.TrimSpace(result)
}
