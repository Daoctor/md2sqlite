// Package savedb use to save content to sqlite3
package savedb

import (
	"database/sql"
	"time"
	"fmt"
	"log"
)

//TendingData hold all link for one language in one day
type TendingData struct {
	Date     string
	Language string
	Links    []string
}

func createDB(path string) error {
	db, err := sql.Open("sqlite3", path)
	createSQL := `
	create table  if not exists trendings (
		id integer not null primary key autoincrement,
		date date,
		link char(255),
		stars integer default NULL,
		language char(50),
		update_time datetime
		);
		create unique index if not exists avoid_dup on trendings (date, link, language);
	`
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(createSQL)
	if err != nil {
		log.Printf("%q: %s\n", err, createSQL)
		return err
	}
	return nil
}

func insertTo(stmt *sql.Stmt, lang string, date string, link string) {
	t := time.Now()
	updateTime := fmt.Sprintf("%d-%d-%d %d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	_, err := stmt.Exec(lang, date, link, updateTime)
	if err != nil {
		fmt.Print(err, "\n")
		fmt.Printf("insert %s %s %s failed\n", lang, date, link)
	}
}

// Save : save all link to sqlite
func Save(dbPath string, data TendingData) {
	createDB(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err, "\n")
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into trendings(language, date, link,update_time) values(?,?,?,?)")
	lang := data.Language
	links := data.Links
	date := data.Date
	for _, link := range links {
		insertTo(stmt, lang, date, link)
	}
	tx.Commit()
}

// Test for function
func Test() {
	var d TendingData
	d.Language = "shell"
	d.Date = "2012-10-1"
	d.Links = []string{"hello", "world"}
	Save("./test.db", d)
}
