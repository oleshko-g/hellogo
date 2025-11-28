package main // revive:disable-line:package-comments

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
		Query(q string, args ...any) (*sql.Rows, error)
		QueryRow(q string, args ...any) *sql.Row
		Stats() sql.DBStats
		Close() error
	}
)

const (
	videos           = "videos.db"
	id               = "0EbFotkXOiA"
	mostWatchedVideo = "SELECT title, channel_title, views from videos WHERE views = (SELECT MAX(views) from videos);"
)

func main() {
	db, err := newDB()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	result, err := db.queryVideos(limit)
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

func (db *ogdb) getDesc(id string) (string, error) {
	row := db.QueryRow("SELECT description FROM videos WHERE video_id = ?", id)

	var desc sql.NullString

	err := row.Scan(&desc)
	if err != nil {
		return "", err
	}
	if desc.Valid {
		return desc.String, nil
	}
	return "-----", nil
}

func (db *ogdb) queryVideos(limit int) ([]Video, error) {
	videos := make([]Video, 0, limit)

	rows, err := db.Query("SELECT video_id, title, views from videos ORDER BY views LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var v Video
		err = rows.Scan(&v.ID, &v.Title, &v.Views)
		if err != nil {
			return nil, err
		}

		videos = append(videos, v)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// Video структура видео.
type Video struct {
	ID    string
	Title string
	Views int64
}

// limit — максимальное количество записей.
const limit = 20
