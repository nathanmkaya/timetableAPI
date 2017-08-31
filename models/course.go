package models

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"runtime"
	"time"

	"github.com/boltdb/bolt"
)

type (
	Course struct {
		Name    string
		Room    string
		Date    time.Time
		Shift   string
		Section string
	}
)

var db *bolt.DB
var open bool

func Open() error {
	var err error
	_, filename, _, _ := runtime.Caller(0)
	dbfile := path.Join(path.Dir(filename), "timetable.db")
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	open = true
	return nil
}

func Close() {
	open = false
	db.Close()
}

func (c *Course) save() error {
	if !open {
		return fmt.Errorf("db must be opened before saving!")
	}

	err := db.Update(func(tx *bolt.Tx) error {
		courses, err := tx.CreateBucketIfNotExists([]byte(c.Shift))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := c.encode()
		if err != nil {
			return fmt.Errorf("could not encode Cource %s", c.Name)
		}
		err = courses.Put([]byte(c.Name), enc)
		return err
	})
	return err
}

func GetCourses(shift string, names []string) ([]*Course, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}

	var courses []*Course
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(shift))
		c := b.Cursor()
		for _, name := range names {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if string(k) == name {
					course, err := decode(v)
					if err != nil {
						return err
					}
					return nil
					courses = append(courses, course)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Course %s in %s", names, shift)
		return nil, err
	}
	return courses, nil
}

func GetCourse(shift string, name string) (*Course, error) {
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}

	var course *Course
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		bucket := tx.Bucket([]byte(shift))
		key := []byte(name)
		course, err = decode(bucket.Get(key))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Course %s in %s", name, shift)
		return nil, err
	}
	return course, nil
}

func (c *Course) encode() ([]byte, error) {
	enc, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Course, error) {
	var c *Course
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
