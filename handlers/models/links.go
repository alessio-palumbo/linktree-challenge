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
	Type      linkType      `json:"type"`
	Title     *string       `json:"title"`
	URL       *string       `json:"url"`
	Thumbnail *string       `json:"thumbnail,omitempty"`
	CreatedAt time.Time     `json:"-"`
	SubLinks  []interface{} `json:"sublinks,omitempty"`
}

// LinkPayload validates a request to create a new Link
type LinkPayload struct {
	Type      linkType          `json:"type" validate:"required,oneof=classic music shows"`
	Title     *string           `json:"title" validate:"omitempty,max=144"`
	URL       *string           `json:"url" validate:"omitempty,max=500"`
	Thumbnail *string           `json:"thumbnail,omitempty" validate:"omitempty,max=144"`
	SubLinks  []json.RawMessage `json:"sublinks,omitempty"`
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
// TODO Date could be a more general format so that FE can use the local to display the date
type Show struct {
	ID       string     `json:"id"`
	Date     string     `json:"date" validate:"required,lkDate"`
	Name     string     `json:"name"`
	Venue    string     `json:"venue" validate:"required_without=Location"`
	Location string     `json:"location" validate:"required_without=Venue"`
	Status   showStatus `json:"status" validate:"required,oneof=on-sale sold-out not-on-sale"`
	URL      string     `json:"url" validate:"required"`
}

// Platform is a sublink representing a song's streaming platform and its url
type Platform struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required"`
}
