package bitwarden

import (
	"encoding/json"
	"fmt"
	"os"
)

type configStruct struct {
	// ClientID     string `json:"client_id"`
	// ClientSecret string `json:"client_secret"`

	authenticated bool

	// This is probably going to get changed to
	// something where its created and managed locally by the API,
	// this just makes dev work easier while it's still being worked on.
	SessionToken string `json:"session_token"`
}

var config = &configStruct{}

func init() {
	f, err := os.Open("config/bw_cli.json")
	if err != nil {
		panic(fmt.Errorf("unable to open bitwarden CLI keys: %w", err))
	}

	if err = json.NewDecoder(f).Decode(config); err != nil {
		panic(fmt.Errorf("unable to load bitwarden CLI keys: %w", err))
	}
}
