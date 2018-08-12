// Package cron contains functions and helpers responsible for
// processing and caching top queried news.
package cron

import (
	"time"
)

// Key is the request parameter key for querying news.
type Key string

// TopQueried holds the mapping of top queried values to its request parameter key.
type TopQueried struct {
	Key
	Values []string
}

// Log is the response for querying top queried news.
type Log struct {
	ElapsedTime time.Duration
	Queried     []TopQueried
}
