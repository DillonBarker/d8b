package db

import (
	"database/sql"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

func loadConfig() (Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return config, err
	}
	return config, nil
}

func ExecuteQuery(db *sql.DB, query string) ([][]string, []string) {
	var result [][]string

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(columns))
	for i := range values {
		valuePointers[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePointers...)
		if err != nil {
			panic(err)
		}

		var row []string
		for _, v := range values {
			switch val := v.(type) {
			case nil:
				row = append(row, "NULL")
			case []byte:
				row = append(row, string(val))
			default:
				row = append(row, fmt.Sprintf("%v", val))
			}
		}
		result = append(result, row)
	}

	return result, columns
}

func ConnectToDb() (*sql.DB, error) {
	config, err := loadConfig()

	if err != nil {
		panic(err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
