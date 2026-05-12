package utils

import "errors"

var (
	PollNotFound    = errors.New("poll not found")
	OptionsNotFound = errors.New("option not found")
	PollIDMissing   = errors.New("poll id is missing")
	GenericError    = errors.New("something went wrong")
)
