package bootstrap

import (
	"database/sql"
)

type Client struct {
	DB *sql.DB
}

func NewClient() *Client {
	return &Client{DB: &sql.DB{}}
}
