package model

import "time"

type Event struct {
	Title string    `json:"title"`
	Time  time.Time `json:"time"`
}
