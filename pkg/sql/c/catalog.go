// Code generated by "sqlc-gen-zz". DO NOT EDIT.

package c

import (
	"time"
)

type Contact struct {
	Blob      []byte    `json:"blob"`
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	MetaJson  []byte    `json:"meta_json"`
	Name      string    `json:"name"`
}
