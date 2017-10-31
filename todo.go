package main

import (
	"log"

	"github.com/covrom/go-echo-vue/handlers"

	bolt "github.com/coreos/bbolt"
	"github.com/labstack/echo"
)

func main() {

	db := initDB("todo.db")
	defer db.Close()

	migrate(db)

	e := echo.New()

	e.File("/", "public/index.html")
	e.GET("/tasks", handlers.GetTasks(db))
	e.PUT("/tasks", handlers.PutTask(db))
	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	e.Start(":8000")
}

func initDB(filepath string) *bolt.DB {
	db, err := bolt.Open(filepath, 0600, nil)

	// Here we check for any db errors then exit
	if err != nil {
		log.Fatal(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		log.Fatal("db nil")
	}
	return db
}

func migrate(db *bolt.DB) {
	// sql := `
	// CREATE TABLE IF NOT EXISTS tasks(
	// 	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	name VARCHAR NOT NULL
	// );
	// `
	err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("tasks"))
		return nil
	})
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		log.Fatal(err)
	}
}
