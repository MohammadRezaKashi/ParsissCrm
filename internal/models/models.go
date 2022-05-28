package models

import "time"

type Report struct {
	ID            int
	Date          time.Time
	Acccess_level int
	Created_at    time.Time
	Updated_at    time.Time
}
