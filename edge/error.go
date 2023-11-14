package edge

import "errors"

var (
	ErrEmptyFreq        = errors.New("frequency cannot be empty")
	ErrEmptySampleRate  = errors.New("samping rate cannot be empty")
	ErrEmptyReampleRate = errors.New("resamping rate cannot be empty")
)
