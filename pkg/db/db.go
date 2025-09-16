package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

const schema = `CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(128),
				comment TEXT,
				repeat VARCHAR(128)
				);`

var db *sql.DB

func Init(dbFile string) error {
	var install bool
	var err error
	if _, err = os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return fmt.Errorf("stat %q failed: %w", dbFile, err)
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed create a database: %w", err)
		}
	}

	return nil
}
