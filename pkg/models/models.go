package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type MaterialGuidelineResults struct {
	MID          int32
	CommunityID  string
	Category     string
	YesNo        string
	CategoryType string
	Material     string
}
