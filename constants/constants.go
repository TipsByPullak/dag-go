package constants

import "time"

var (
	Port            string        = "8080"
	NumberOfWorkers int           = 4
	TickerTime      time.Duration = time.Second * 2
)
