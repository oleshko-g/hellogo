package main // revive:disable-line:package-comments

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

	list, err := db.queryTagVideos(limit)
	if err != nil {
		panic(err)
	}
	// для теста проверим, какие строки содержит v.Tags
	// выведем по 4 первых тега
	for _, v := range list {
		length := 4
		if len(v.Tags) < length {
			length = len(v.Tags)
		}
		fmt.Println(strings.Join(v.Tags[:length], " # "))
	}
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
	Tags  Tags
}

// limit — максимальное количество записей.
const limit = 20

func (db *ogdb) queryTagVideos(limit int) ([]Video, error) {
	videos := make([]Video, 0, limit)

	rows, err := db.Query("SELECT video_id, title, tags from videos "+
		"GROUP BY video_id ORDER BY views LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v Video
		// все теги должны автоматически преобразоваться в слайс v.Tags
		err = rows.Scan(&v.ID, &v.Title, &v.Tags)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return videos, nil
}

type Tags []string

func (tt Tags) Value() (driver.Value, error) {
	if len(tt) == 0 {
		return "", nil
	}
	return strings.Join(tt, "|"), nil
}

func (tt *Tags) Scan(src interface{}) error {
	if src == nil {
		tt = &Tags{}
		return nil
	}

	v, err := driver.String.ConvertValue(src)
	if err != nil {
		return fmt.Errorf("cannot scan value. %w", err)
	}

	sv, ok := v.(string)
	if !ok {
		return errors.New("cannot scan value. cannot convert value to string")
	}

	ss := strings.Split(sv, "|")
	for i, s := range ss {
		ss[i] = strings.Trim(s, `"`)
	}

	*tt = append(*tt, ss...)

	return nil
}
