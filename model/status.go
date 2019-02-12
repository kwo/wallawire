package model

import (
	"time"
)

type Status struct {
	ServiceName string    `json:"service"`
	Version     string    `json:"version"`
	Runtime     string    `json:"runtime"`
	BuildTime   time.Time `json:"buildTime"`
	StartTime   time.Time `json:"start"`
	SystemTime  time.Time `json:"time"`
	Uptime      string    `json:"uptime"`
}

func (z *Status) PopulateNow() {
	z.Populate(time.Now().Truncate(time.Second).UTC())
}

func (z *Status) Populate(now time.Time) {
	z.SystemTime = now
	z.Uptime = z.SystemTime.Sub(z.StartTime).String()
}
