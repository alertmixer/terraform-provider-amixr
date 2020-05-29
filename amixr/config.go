package amixr

import (
	"github.com/alertmixer/amixr-go-client"
)

type Config struct {
	Token string
}

func (c *Config) Client() (interface{}, error) {
	client, err := amixr.NewClient(c.Token)
	return client, err
}
