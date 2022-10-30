package yubikey

import (
	"fmt"
	"github.com/jack073/mobile_kb_api/logger"
	"io"
	"os"
	"strings"
)

type yubikey struct {
	name string
	id   string
}

var allowedYubikeys = make([]*yubikey, 0)

func init() {
	f, err := os.Open("config/yubikeys.txt")
	if err != nil {
		panic(fmt.Errorf("error when opening yubikey file: %w", err))
	}

	rawData, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("error when reading yubikey file: %w", err))
	}

	for _, line := range strings.Split(string(rawData), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		split := strings.Split(line, ":")
		if len(split) != 2 {
			logger.Logger.Warn("unable to parse yubikey config: ", line)
			continue
		}

		ID := strings.TrimSpace(split[1])
		if len(ID) != 12 {
			logger.Logger.Warn("invalid yubikey ID: ", ID)
			continue
		}

		allowedYubikeys = append(allowedYubikeys, &yubikey{
			name: strings.TrimSpace(split[0]),
			id:   ID,
		})
	}
}

func CheckKey(key string) (string, bool) {
	for _, k := range allowedYubikeys {
		if k.id == key {
			return k.name, true
		}
	}

	return "", false
}
