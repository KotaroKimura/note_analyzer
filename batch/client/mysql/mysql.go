package mysql

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	loc      = "Asia%2FTokyo"
	timeZone = "%27Asia%2FTokyo%27"
)

var (
	dbURL = os.Getenv("DB_URL")
)

type Conn struct {
	M *sql.DB
}

func NewConn() (*Conn, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true&loc=%s&time_zone=%s", dbURL, loc, timeZone))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(0)

	return &Conn{
		M: db,
	}, nil
}

func (c *Conn) Close() {
	c.M.Close()
}
