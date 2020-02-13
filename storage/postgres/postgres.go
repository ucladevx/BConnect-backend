package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Table interface for creating new tables
type Table interface {
	create()
}

//CreatePostgresTables creates postgresql tables
func CreatePostgresTables(tables ...Table) {
	for _, table := range tables {
		table.create()
	}
}

// Connect connects to postgresql
func Connect(databaseHost string, username string, databaseName string, password string) *gorm.DB {
	db, err := gorm.Open("postgres", dbURI(databaseHost, username, databaseName, password))
	if err != nil {
		print(err.Error())
	}
	//defer db.Close()
	fmt.Println("Successfully connected!", db)
	return db
}

func dbURI(databaseHost string, username string, databaseName string, password string) string {
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
}