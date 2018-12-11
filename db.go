package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Author struct {
	Name        string
	Email       string
}

func (u Author) String() string {
	return fmt.Sprintf("Author<%s %s>", u.Name, u.Email)
}

type Post struct {
	Id          int64
	Title       string
	Content     string
	AuthorEmail string
}

func (s Post) String() string {
	return fmt.Sprintf("Post<%d %s %s>", s.Id, s.Title, s.AuthorEmail)
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

func createAuthorProfile(db *pg.DB, profile map[string]string) {
	fmt.Println(profile)
	author := &Author{
		Name: profile["name"],
		Email: profile["email"],
	}
	db.Model(author).Where("email = ?", author.Email).SelectOrInsert()
}