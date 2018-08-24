package client

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// GenericSlackAPIClient is the slack client interface
type GenericSlackAPIClient interface {
	// Notify sends messages to a channel
	Notify(channel, text string, params slack.PostMessageParameters) error
}

// SlackAPIClient is an implementation of the interface
type SlackAPIClient struct {
	client *slack.Client
}

// NewSlackAPIClient creates a new slack api client
func NewSlackAPIClient(client *slack.Client) GenericSlackAPIClient {
	return &SlackAPIClient{
		client: client,
	}
}

// Notify sends messages to a channel
func (c *SlackAPIClient) Notify(channel, text string, params slack.PostMessageParameters) error {
	log.Infof("sending message to slack channel %s", channel)
	_, _, err := c.client.PostMessage(channel, text, params)
	return err
}
