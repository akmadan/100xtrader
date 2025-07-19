package data

import (
	"database/sql"
	"io/ioutil"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func InitDB(dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		logrus.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func RunMigrations(db *sql.DB, migrationFile string) {
	content, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		logrus.Fatalf("failed to read migration file: %v", err)
	}
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			logrus.Fatalf("failed to execute migration statement: %v\nSQL: %s", err, stmt)
		}
	}
}
