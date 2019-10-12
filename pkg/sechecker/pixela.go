package sechecker

import (
	"github.com/gainings/pixela-go-client"
)

type PixelaConfig struct {
	UserID  string `json:"UserID"`
	GraphID string `json:"GraphID"`
	Secret  string `json:"Secret"`
}

func NewPixelaClient(userID, graphID, secret string) (*PixelaConfig, error) {
	return &PixelaConfig{
		UserID:  userID,
		GraphID: graphID,
		Secret:  secret,
	}, nil
}

func (api *PixelaConfig) PostEvent(metadata MetaData) error {
	c := pixela.NewClient(api.UserID, api.Secret)
	err := c.IncrementPixelQuantity(api.GraphID)
	return err
}
