// package client implements a client for the Rick and Morty API hosted at https://rickandmortyapi.com/
// It uses the GraphQL endpoint at https://rickandmortyapi.com/graphql

package client

import "time"

//go:generate generate-gpt -p client -i spec.go -o client.go -s Client

type Client interface {
	GetCharacter(id int) (Character, error)
}

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

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
