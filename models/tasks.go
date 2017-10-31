package models

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
)

// Task is a struct containing Task data
type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TaskCollection is collection of Tasks
type TaskCollection struct {
	Tasks []Task `json:"items"`
}

// GetTasks from the DB
func GetTasks(db *bolt.DB) TaskCollection {
	result := TaskCollection{}
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("tasks"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := Task{ID: string(k), Name: string(v)}
			result.Tasks = append(result.Tasks, task)
		}
		return nil
	})
	return result
}

// PutTask into DB
func PutTask(db *bolt.DB, name string) (string, error) {
	var bts string
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		id, _ := b.NextSequence()
		bts = fmt.Sprintf("%08d", id)
		return b.Put([]byte(bts), []byte(name))
	})
	return bts, err
}

// DeleteTask from DB
func DeleteTask(db *bolt.DB, id string) (string, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		return b.Delete([]byte(id))
	})
	return id, err
}
