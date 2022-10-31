package bitwarden

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type HasString bool

// As far as we're concerned for this part, we only care
// about whether there's a "truthy" value here. This also
// prevents the login information being saved more in memory.
func (h *HasString) UnmarshalJSON(data []byte) error {
	if data == nil {
		*h = false
		return nil
	}

	if len(data) == 0 {
		*h = false
		return nil
	}

	if data[0] != '"' {
		// We know that this must either be null, or a string literal.
		// If it doesn't start with a quote, then it must be null.
		*h = false
		return nil
	}

	// We already know it starts with a ", so the only way
	// for it to have length 2 is for it to be an empty string.
	if len(data) == 2 {
		*h = false
		return nil
	}

	*h = true
	return nil
}

type LoginListItem struct {
	Name string `json:"name"`

	Login struct {
		Username HasString `json:"username"`
		Password HasString `json:"password"`
		TOTP     HasString `json:"totp"`
	}
}

func GetLoginList() ([]*LoginListItem, error) {
	cmd := exec.Command(CMD, "list", "items")

	cmd.Env = append(os.Environ(), fmt.Sprintf(`BW_SESSION="%s"`, config.SessionToken))
	buf := NewBuffer()
	cmd.Stdout = buf

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	logins := make([]*LoginListItem, 0)
	if err := json.Unmarshal(buf.Bytes(), &logins); err != nil {
		return nil, err
	}
	buf.Clear()

	return logins, nil
}
