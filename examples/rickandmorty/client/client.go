package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Rick and Morty API GraphQL endpoint
const apiUrl = "https://rickandmortyapi.com/graphql"

// Character represents a character from the show
type Character struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Species  string    `json:"species"`
	Type     string    `json:"type"`
	Gender   string    `json:"gender"`
	Origin   Location  `json:"origin"`
	Location Location  `json:"location"`
	Image    string    `json:"image"`
	Episode  []string  `json:"episode"`
	URL      string    `json:"url"`
	Created  time.Time `json:"created"`
}

// Location represents a location
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// RequestPayload represents a GraphQL request payload
type RequestPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// ResponsePayload represents a GraphQL response payload
type ResponsePayload struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// NewClient returns a new client that implements the Client interface
func NewClient() Client {
	return &client{}
}

type client struct{}

func (c *client) GetCharacter(id int) (Character, error) {
	var result Character

	payload := RequestPayload{
		Query: `
			query getCharacter($id: ID!) {
				character(id: $id) {
					id
					name
					status
					species
					type
					gender
					origin {
						name
						url
					}
					location {
						name
						url
					}
					image
					episode
					url
					created
				}
			}
		`,
		Variables: map[string]interface{}{
			"id": id,
		},
	}

	var response ResponsePayload

	if err := doPost(apiUrl, &payload, &response); err != nil {
		return result, err
	}

	if len(response.Errors) > 0 {
		return result, errors.New(response.Errors[0].Message)
	}

	data, err := json.Marshal(response.Data.(map[string]interface{})["character"])
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

func doPost(url string, payload, response interface{}) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payloadJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}

	return nil
}
