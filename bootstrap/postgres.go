package bootstrap

import (
	"database/sql"
)

type Client struct {
	Db *sql.DB
}

func NewClient() *Client {
	return &Client{Db: &sql.DB{}}
}
