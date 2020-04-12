# Linktree Challenge

## Author: Alessio Palumbo

## Problem

### Design a REST API to support multiple types of links:

#### Link types

* _Classic_
    * Title => no longer than 144 characters
    * Url   => may contain query parameters
* _Shows List_
    * Status => sold out, on sale, not yet on sale
* _Music Player_
    * Multiple links => requires a link to each platform

#### Objectives

* The client must be able to create a new link of each type.
* The client must be able to find all links matching a particular userId.
* The client must be able to find links matching a particular userId by dateCreated

## Solution

### DB schema (postgres)

* users: -- this is just to support links table so we just need an id
    * id UUID NOT NULL (PK)

* links:
    * id UUID NOT NULL (PK)
    * user_id UUID NOT NULL (FK)
    * type CHAR(10) NOT NULL default 'classic'
    * title VARCHAR(144) default NULL
    * url VARCHAR(500) default NULL -- TODO could use shortened urls
    * thumbnail VARCHAR(144) default NULL -- assuming is shortened and stored in an s3 bucket
    * created_at TIMESTAMPTZ DEFALT CURRENT_TIMESTAMP
    * TODO if needed an update_at timestamp could be added
    * TODO for the sake of ordering we could have an order_id field

* sublinks:
    * id UUID NOT NULL (PK)
    * link_id UUID NOT NULL (FK)
    * metadata JSONB NULL

### Models

#### Main Link model

* Link --
    * ID string (uuid)
    * UserID string (uuid)
    * Type string ('classic', 'music', 'shows')
    * Title *string
    * URL *string
    * SubLinks []interface{}

#### Shows List sublink model

* Show:
    * ID string (uuid)
    * Date string -- This could be a time.Time
    * Name string
    * Venue string
    * Location string -- TODO this could be a new struct Address
    * Status Status(string)

#### Music Player sublink model

* Platform:
    * ID uuid
    * Name string
    * URL string

#### Examples:

* Classic:
    Link{
        ID: "001",
        UserID: "43435423",
        Type: "classic",
        Title: "First Link",
        URL: "https://myfirstlink.com/1",
    }

* Show-Link:
    Link{
        ID: "002",
        UserID: "34232354",
        Type: "show",
        SubLinks: []Show{
            Show{
                ID: "s001",
                Date: "Apr 01 2019",
                Name: "Cats",
                Venue: "Princess Theatre",
                Location: "Melbourne",
                Status: "on-sale",
                URL: "https://www.ticketmaster.com.au/cats-the-musical-tickets/artist/843992",
            },
        },
    }

* Music player:
    Link{
        ID: "003",
        UserID: "43532454",
        Type: "music",
        Title: "All of me - John Legend"
        URL: "https://www.youtube.com/watch?v=450p7goxZqg",  -- This stores the url for the embed default player
        SubLinks: []Platform{
            Platform{
                ID: "p001",
                Name: "Spotify",
                URL: "https://open.spotify.com/album/1YdXQgntClL3BhIXB0xpgs",
            },
            Platform{
                ID: "p002",
                Name: "SoundCloud",
                URL: "https://soundcloud.com/johnlegend/all-of-me-3",
            },
        },
    }

### REST API

#### Authentication

The Api assumes it can fetch the userID of the request from the authentication middleware

#### Links Rest

* GET /api/links
    * Query params
        * order_by: created_at:asc,title:desc,type (optional, accepts multiple columns. TODO default to positionID)
    * Responses:
        * 200 OK
            ```
            {
                "links": [
                    {
                        "id": "001",
                        "type": "classic",
                        "url": "https://myfirstlink.com/1"
                    },
                    {
                        "id": "002",
                        "type": "show",
                        "url": "https://myfirstlink.com/1",
                        "sublinks": [
                            {
                                "id": "s001",
                                "date": "Apr 01 2019",
                                "venue": "Princess Theatre",
                                "location": "Melbourne",
                                "status": "sold-out"
                            }
                        ]
                    }
                ]
            }
            ```
        * 400 Bad Request
        * 404 Not found (user)

* GET /api/links/{link_id}
    * Responses:
        * 200 OK
            ```
            {
                "id": "001",
                "type": "classic",
                "url": "https://myfirstlink.com/1"
            }

* POST /api/links
    * Request:
        ```
        {
            "type": "music",
            "title": "All of me - John Legend",
            "url": "https://www.youtube.com/watch?v=450p7goxZqg",
            "sublinks": [
                {
                    "name": "Spotify",
                    "url": "https://open.spotify.com/album/1YdXQgntClL3BhIXB0xpgs"
                },
                {
                    "name": "SoundCloud",
                    "url": "https://soundcloud.com/johnlegend/all-of-me-3"
                }
            ]
        }
        ```
    * Responses:
        * 201 Created
        * 400 Bad Request

* PUT /api/links/{link_id}
    * Request:
        ```
        {
            "type": "classic",
            "title": "My second Link",
            "url": "https://www.mysecondlink.com/2"
        }
        ```
    * Responses:
        * 200 OK
            ```
            {
                "id": "004",
                "type": "classic",
                "title": "My second Link",
                "url": "https://www.mysecondlink.com/2"
            }
            ```
        * 400 Bad Request

* DELETE /api/links/{link_id} -- Remove link and any sublinks
    * Response:
        * 204 No Responses

#### SubLinks Rest (Only POST, PUT and DELETE)

* POST /api/links/{link_id}/sublinks
    * Request:
        ```
        {
            "name": "Spotify",
            "url": "https://open.spotify.com/album/1YdXQgntClL3BhIXB0xpgs"
        }
        ```
    * Responses:
        * 201 Created
        * 400 Bad Request

* PUT /api/links/{link_id}/sublinks/{sublink_id}
    * Request:
        ```
        {
            "type": "classic",
            "title": "My second Link",
            "url": "https://www.mysecondlink.com/2"
        }
        ```
    * Responses:
        * 200 OK
            ```
            {
                "id": "004",
                "type": "classic",
                "title": "My second Link",
                "url": "https://www.mysecondlink.com/2"
            }
            ```
        * 400 Bad Request

* DELETE /api/links/{link_id}/sublinks/{sublink_id}
    * Response:
        * 204 No Responses

## Language used: Go

### Setup

#### Installing Go (1.12 or higher)

Install Go following the official instructions: https://golang.org/doc/install

#### Run

* From the console `cd` into main folder `linktree-challenge`
* Run `go run server.go`

#### Build and Run

* Run `go build server.go`
* Run `./linktree-challenge`

#### Test

* Run all test `go test -v`

#### Notes about using GOPATH

* If running inside a GOPATH on Go 1.11/1.12 use the following command to build and test
  * GO111MODULE=on go run server.go
  * GO111MODULE=on go build .
  * GO111MOdULE=on go test -v