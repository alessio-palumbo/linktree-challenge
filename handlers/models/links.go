package models

import (
	"encoding/json"
	"time"
)

type linkType string

const (
	// LinkClassic is a standard link with title and url
	LinkClassic linkType = "classic"
	// LinkMusic is a link to a song record. It allows to specify
	// multiple sublinks with a link to each platform
	LinkMusic linkType = "music"
	// LinkShows represent a list with links to multiple shows
	LinkShows linkType = "shows"
)

// Link is the base model for a link that can contain a list
// of sublinks associated with its type
type Link struct {
	ID        string        `json:"id"`
	UserID    string        `json:"-"`
	Type      linkType      `json:"type" validate:"required,oneof=classic music shows"`
	Title     *string       `json:"title" validate:"omitempty,max=144"`
	URL       *string       `json:"url" validate:"omitempty,max=500"`
	Thumbnail *string       `json:"thumbnail,omitempty" validate:"omitempty,max=144"`
	CreatedAt time.Time     `json:"-"`
	SubLinks  []interface{} `json:"-" validate:"dive"`
}

// Sublink contains the metadata of a sublink
type Sublink struct {
	ID       string
	UserID   string
	Metadata json.RawMessage
}

type showStatus string

// The status a given show can be set to
const (
	StatusOnSale    showStatus = "on-sale"
	StatusNotOnSale showStatus = "not-on-sale"
	StatusSoldOut   showStatus = "sold-out"
)

// Show is a sublink containing information about a single show
type Show struct {
	ID       string
	Date     string
	Name     string
	Venue    string
	Location string
	Status   showStatus
}

// Platform is a sublink representing a song's streaming platform and its url
type Platform struct {
	ID   string
	Name string
	URL  string
}
