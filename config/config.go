package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func LoadEnv() (err error) {
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	return
}

func InitDB() (*sql.DB, error) {
	DBUsername := viper.GetString("DB_USERNAME")
	DBPassword := viper.GetString("DB_PASSWORD")
	DBHost := viper.GetString("DB_HOST")
	DBPort := viper.GetString("DB_PORT")
	DBName := viper.GetString("DB_NAME")

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName)

	// Open a database connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")

	return db, nil
}
