package models

type linkType string

const (
	LinkClassic linkType = "classic"
	LinkMusic   linkType = "music"
	LinkShows   linkType = "shows"
)

// Link is the base model for a link that can contain a list
// of sublinks associated with its type
type Link struct {
	ID       string
	UserID   string
	Type     linkType
	Title    *string
	URL      *string
	SubLinks []interface{}
}

type linkStatus string

const (
	StatusOnSale    linkStatus = "on-sale"
	StatusNotOnSale linkStatus = "not-on-sale"
	StatusSoldOut   linkStatus = "sold-out"
)

// Show is a sublink containing information about a single show
type Show struct {
	ID       string
	Date     string
	Name     string
	Venue    string
	Location string
	Status   linkStatus
}

// Platform is a sublink representing a song's streaming platform and its url
type Platform struct {
	ID   string
	Name string
	URL  string
}
