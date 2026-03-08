package main // revive:disable-line:package-comments

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type (
	ogdb struct {
		queries
	}
	queries interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		Exec(q string, args ...interface{}) (sql.Result, error)
		Query(q string, args ...interface{}) (*sql.Rows, error)
		QueryRow(q string, args ...interface{}) *sql.Row
		Stats() sql.DBStats
		Close() error
	}
)

const (
	videos                     = "videos.db"
	id                         = "0EbFotkXOiA"
	mostWatchedVideo           = "SELECT title, channel_title, views from videos WHERE views = (SELECT MAX(views) from videos);"
	_CSVFile                   = "USvideos.csv"
	USVideos_videos_short_full = "USvideos short full.csv"
)

func main() {
	db, err := newDB()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	vv, err := readVideoCSV(USVideos_videos_short_full, 2)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	start := time.Now()
	err = db.insertVideosShortTx(ctx, vv)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	finish := time.Since(start)
	slog.Info("Insreted videos_short in %v", "duration", finish.String())
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
	ID          string
	Title       string
	PublishTime time.Time // publish_time
	Tags        Tags
	Views       int64
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

func (tt Tags) String() string {
	if len(tt) == 0 {
		return ""
	}
	return strings.Join(tt, "|")
}

func (tt Tags) Value() (driver.Value, error) {
	return tt.String(), nil
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

type trend struct {
	trendDate
	Count int
}

type trendDate struct{ T time.Time }

func (t *trendDate) Scan(src interface{}) error {
	if t == nil {
		return errors.New("can't scan into a nil pointer")
	}

	if src == nil {
		t.T = time.Time{}
		return nil
	}

	v, _ := driver.String.ConvertValue(src)

	sv, ok := v.(string)
	if !ok {
		return errors.New("sql src value is not a string")
	}

	tv, err := time.Parse(trendDateLayout, sv)
	if err != nil {
		return err
	}

	t.T = tv

	return nil
}

func (db *ogdb) trendingCount() ([]trend, error) {
	q := "SELECT trending_date, COUNT(DISTINCT video_id) as trending_video_count FROM videos GROUP BY trending_date;"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	var tt []trend
	for rows.Next() {
		var t trend
		err = rows.Scan(&t.trendDate, &t.Count)
		if err != nil {
			return nil, err
		}
		tt = append(tt, t)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tt, nil
}

// YY.DD.MM
const trendDateLayout = "06.02.01"

func readVideoCSV(csvFilePath string, readFromLine int) ([]Video, error) {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(csvFile)
	r.FieldsPerRecord = 5

	var vv []Video
	for i := 1; ; i++ {
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			return nil, err
		}

		if i < readFromLine {
			continue
		}

		var v Video
		err = parseCSVRecord(record,
			&v.ID,
			&v.Title,
			&v.PublishTime,
			&v.Tags,
			&v.Views,
		)
		if err != nil {
			return nil, err
		}
		vv = append(vv, v)
	}

	return vv, nil
}

func parseCSVRecord(record []string, dsts ...interface{}) error {
	if record == nil {
		return errors.New("error record is nil")
	}

	if len(record) == 0 {
		return errors.New("error record has zero length")
	}

	if dsts == nil {
		return errors.New("error dsts is nil")
	}

	if len(dsts) == 0 {
		return errors.New("error dsts has zero length")
	}

	if len(record) != len(dsts) {
		return errors.New("error parsing csv record. record's field number doesn't match the number of dsts")
	}

	for i, r := range record {
		switch v := dsts[i].(type) {
		case *string:
			*v = r
		case *time.Time:
			tv, err := time.Parse(time.RFC3339, r)
			if err != nil {
				return err
			}
			*v = tv
		case *Tags:
			*v = strings.Split(r, "|")
		case *int64:
			iv, err := strconv.Atoi(r)
			if err != nil {
				return err
			}
			*v = int64(iv)
		default:
			return errors.New("error parsing csv record. unsupported dst")
		}
	}

	return nil
}

func (db *ogdb) insertVideosShort(vv []Video) error {
	if db == nil {
		return errors.New("error db is nil")
	}

	if vv == nil {
		return errors.New("error vv is nil")
	}

	q := `INSERT INTO
					videos_short (video_id, title, publish_time, tags, views)
				VALUES
					($1, $2, $3, $4, $5);
				`

	for _, v := range vv {
		_, err := db.Exec(q, v.ID, v.Title, v.PublishTime, v.Tags, v.Views)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *ogdb) insertVideosShortTx(ctx context.Context, vv []Video) error {
	if db == nil {
		return errors.New("error db is nil")
	}

	if ctx == nil {
		return errors.New("ctx is nil")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if vv == nil {
		return errors.New("error vv is nil")
	}

	q := `INSERT INTO
					videos_short (video_id, title, publish_time, tags, views)
				VALUES
					($1, $2, $3, $4, $5);
				`

	for i, v := range vv {
		_, err := tx.ExecContext(ctx, q, v.ID, v.Title, v.PublishTime, v.Tags, v.Views)
		if err != nil {
			tx.Rollback()
			return err
		}
		a := i
		_ = a
	}

	return nil
}
