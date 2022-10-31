package yubikey

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unsafe"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"net/http"
	"net/url"

	"github.com/jack073/mobile_kb_api/logger"

	"github.com/google/uuid"
)

type configStruct struct {
	ClientID  string `json:"client_id"`
	SecretKey string `json:"secret_key"`
}

var config = &configStruct{}

func init() {
	f, err := os.Open("config/api-keys.json")
	if err != nil {
		panic(fmt.Errorf("unable to open API keys: %w", err))
	}

	if err = json.NewDecoder(f).Decode(config); err != nil {
		panic(fmt.Errorf("unable to load API keys: %w", err))
	}
}

func LookupOTP(otp string) (string, error) {
	uri, err := url.Parse("https://api.yubico.com/wsapi/2.0/verify")
	if err != nil {
		// How...????
		return "", err
	}

	for n := 0; n < 5; n++ {
		query := uri.Query()
		query.Set("id", config.ClientID)
		query.Set("otp", otp)
		query.Set("nonce", generateNonce())
		query.Set("sl", "50")
		query.Set("timeout", "10")

		uri.RawQuery = query.Encode()

		resp, err := http.Get(uri.String())
		if err != nil {
			continue
		}

		if 400 <= resp.StatusCode && resp.StatusCode < 600 {
			// From https://developers.yubico.com/OTP/Specifications/OTP_validation_protocol.html

			// If you get a 4xx or 5xx response you should retry your request a few times, as
			// intermediate proxies and gateways may cause transient errors. If you are using
			// a locally-hosted validation server on your own network, this may not be necessary.
			continue
		}

		// status 200
		if resp.StatusCode != http.StatusOK {
			return "", errors.New(fmt.Sprintf("unknown error: invalid http status returned: %d", resp.StatusCode))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// The body consists of a number of parameter=value pairs, separated by CR LF.
		// 13 = CR, 10 = LF
		lines := strings.Split(string(body), string([]byte{13, 10}))
		for _, line := range lines {
			if !strings.HasPrefix(line, "status=") {
				// Right now, the status is the only part of the response we care about.
				// This could change in the future? But right now there's no point in
				// implementing / finding a full parser for it/
				continue
			}

			// Strip the status= prefix
			return line[7:], nil
		}

		// How does this happen??
		return "", errors.New(fmt.Sprintf("unknown error: no status line in HTTP resp, full resp: \n%s", string(body)))
	}

	return "", errors.New("unable to successfully connect to yubico API")
}

func generateNonce() string {
	now := time.Now().UnixNano()
	timeBytes := *(*[8]byte)(unsafe.Pointer(&now))

	id, err := uuid.NewRandom()
	if err != nil {
		logger.Logger.Errorln(fmt.Errorf("nonce UUID gen error: %w", err))
		var out []byte
		src := make([]byte, 12)
		if n, err := rand.Read(src); err == nil {
			out = make([]byte, (n*2)+16)
			hex.Encode(out, timeBytes[:])
			hex.Encode(out[16:], src)
		} else {
			out = make([]byte, 16)
			hex.Encode(out, timeBytes[:])
		}
		return string(out)
	}

	out := &bytes.Buffer{}
	enc := base64.NewEncoder(
		base64.
			// This is based on base64.encodeStd / base64.encodeURL but with the
			// non-alphanumeric chars swapped for a and A. It appears to be that any
			// non-alphanumeric chars are rejected by the yubico API - causing an error
			// as it then doesn't recognise a nonce.
			// I can't see any documented behaviour confirming this rule,
			// it's just my observation but this seems to fix it.
			NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789aA").
			WithPadding(base64.NoPadding),
		out)
	_, _ = enc.Write(timeBytes[:])
	_, _ = enc.Write(id[:])

	_ = enc.Close()

	return out.String()
}
