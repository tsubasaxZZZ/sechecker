package main

import (
	"flag"
	"log"

	"github.com/gainings/pixela-go-client"
)

type PixelaConfig struct {
	UserID  string `json:"UserID"`
	GraphID string `json:"GraphID"`
	Secret  string `json:"Secret"`
}

func NewPixelaClient(userID, graphID, secret string) *PixelaConfig {
	c := &PixelaConfig{
		UserID:  userID,
		GraphID: graphID,
		Secret:  secret,
	}
	return c
}

func (api PixelaConfig) PostEvent() error {
	c := pixela.NewClient(api.UserID, api.Secret)
	err := c.IncrementPixelQuantity(api.GraphID)
	return err
}

func main() {
	c := PixelaConfig{}
	flag.StringVar(&c.UserID, "userId", "", "User ID")
	flag.StringVar(&c.GraphID, "graphId", "", "Graph ID")
	flag.StringVar(&c.Secret, "secret", "", "Secret")
	flag.Parse()

	if c.UserID == "" || c.GraphID == "" || c.Secret == "" {
		log.Fatalf("Need argument")
	}

	if err := c.PostEvent(); err != nil {
		log.Fatal(err)
	}
}
