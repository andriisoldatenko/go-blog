package main

import (
	//"fmt"

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

func main() {
	db := pg.Connect(&pg.Options{
		Database: "blog_db",
		User: "blog",
		Password: "blog_secret_password",
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := &Author{
		Name:   "admin",
	}
	err = db.Insert(user1)
	if err != nil {
		panic(err)
	}

	post1 := &Post{
		Title:    "Cool story",
		AuthorId: user1.Id,
	}
	err = db.Insert(post1)
	if err != nil {
		panic(err)
	}

	// Select user by primary key.
	user := &Author{Id: user1.Id}
	err = db.Select(user)
	if err != nil {
		panic(err)
	}

	// Select all users.
	var authors []Author
	err = db.Model(&authors).Select()
	if err != nil {
		panic(err)
	}

	// Select story and associated author in one query.
	post := new(Post)
	err = db.Model(post).
		Relation("Author").
		Where("post.id = ?", post1.Id).
		Select()
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	fmt.Println(authors)
	fmt.Println(post)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Author)(nil), (*Post)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
