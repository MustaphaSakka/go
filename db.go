package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type article struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func GetArticle(db *sql.DB, id int) (*article, error) {
	query := "SELECT * FROM article WHERE id = ?"
	row := db.QueryRow(query, id)

	a := &article{}
	err := row.Scan(&a.Id, &a.Title)
	if err != nil {
		return nil, err
	}
	return a, nil
}

var articles []article
var a article
var db *sql.DB
var err error

func init() {
	db, err := sql.Open("mysql", "root:@/api")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	fmt.Println("Connexion de database is OK")

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	fmt.Println("select query")
	rows, err := db.Query("SELECT id, title from article")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var id, title string
		rows.Scan(&id, &title)
		fmt.Printf("ID = %s\nTITLE = %s\n\n", id, title)
		a := article{Id: id, Title: title}
		articles = append(articles, a)
	}
	defer rows.Close()
}
