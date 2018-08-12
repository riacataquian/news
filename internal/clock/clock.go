// Package clock abstracts interaction with the built-in package time.
package clock

import "time"

// Time describes a clock interface.
type Time interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// Clock abstracts dealing with time.
//
// Clock implements a clock interface.
type Clock struct{}

// New returns an interface to Time.
func New() Time {
	return Clock{}
}

// Now abstracts time.Now().
func (c Clock) Now() time.Time {
	return time.Now()
}

// Since abstracts time.Since(start).
func (c Clock) Since(start time.Time) time.Duration {
	return time.Since(start)
}
