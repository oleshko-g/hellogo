package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type (
	ogdb struct {
		queries
	}
	queries interface {
		QueryRow(q string, args ...any) *sql.Row
		Stats() sql.DBStats
		Close() error
	}
)

const (
	videos           = "videos.db"
	mostWatchedVideo = "SELECT title, channel_title, views from videos WHERE views = (SELECT MAX(views) from videos);"
)

func main() {
	db, err := newDB()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	result, err := db.selectMostWatchedVideo()
	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Println(result)
}

func newDB() (*ogdb, error) {
	db, err := openSQLiteDB(videos)
	if err != nil {
		return nil, err
	}

	return &ogdb{queries: db}, nil
}

func openSQLiteDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}

func (db *ogdb) selectMostWatchedVideo() (string, error) {
	row := db.QueryRow(mostWatchedVideo)

	type result struct {
		title, channelTitle string
		views               int
	}
	r := result{}

	err := row.Scan(&r.title, &r.channelTitle, &r.views)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%+v", r), nil
}
