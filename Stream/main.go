package Stream

import (
	getstream "github.com/GetStream/stream-go"
	"errors"
)

var Client *getstream.Client

func Connect(api_key string, api_secret string, api_region string) error {
	var err error
	if api_key == "" || api_secret == "" || api_region == "" {
		return errors.New("Config not complete")
	}

	Client, err = getstream.New(&getstream.Config{
		APIKey: api_key,
		APISecret: api_secret,
		Location: api_region,
	})
	return err
}
