package kstartup

import "time"

type Announcement struct {
	ID           int
	BoardID      string
	Group        string
	Title        string
	Name         string
	Organization string
	Due          time.Time
}
