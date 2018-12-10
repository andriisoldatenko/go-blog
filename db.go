package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Author struct {
	Id          int64
	Name        string
	Email       string
	Password    string
}

func (u Author) String() string {
	return fmt.Sprintf("Author<%d %s %v>", u.Id, u.Name, u.Email)
}

type Post struct {
	Id       int64
	Title    string
	AuthorId int64
	Author   *Author
}

func (s Post) String() string {
	return fmt.Sprintf("Post<%d %s %s>", s.Id, s.Title, s.Author)
}


func DBConn() (db *pg.DB) {
	db = pg.Connect(&pg.Options{
		Database: "blog_db",
		User: "blog",
		Password: "blog_secret_password",
	})
	return db
}

func init() {
	db := DBConn()
	createSchema(db)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Author)(nil), (*Post)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
