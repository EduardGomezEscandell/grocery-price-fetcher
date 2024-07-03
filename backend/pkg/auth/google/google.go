package google

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
)

type JWTDecoded struct {
	NotAfter int64  `json:"exp"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Subject  string `json:"sub"`
}

func DecodeJWT(s string) (JWTDecoded, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return JWTDecoded{}, errors.New("invalid JWT format")
	}

	b, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return JWTDecoded{}, fmt.Errorf("could not decode JWT: %v", err)
	}

	var info JWTDecoded
	if err := json.Unmarshal(b, &info); err != nil {
		return JWTDecoded{}, fmt.Errorf("could not decode JWT: %v", err)
	}

	return info, nil
}

func (a Application) Validate(info JWTDecoded) error {
	if err := info.validateIssuer(); err != nil {
		return err
	}

	if err := info.validateAudience(a); err != nil {
		return err
	}

	if info.Expired() {
		return fmt.Errorf("token expired")
	}

	return nil
}

func (info JWTDecoded) validateIssuer() error {
	switch info.Issuer {
	case "https://accounts.google.com", "accounts.google.com":
		return nil
	default:
		return fmt.Errorf("invalid issuer: %s", info.Issuer)
	}
}

func (info JWTDecoded) validateAudience(a Application) error {
	if info.Audience != a.ClientID {
		return fmt.Errorf("invalid audience: %s", info.Audience)
	}

	return nil
}

func (info JWTDecoded) Expired() bool {
	var (
		now      = time.Now()
		notAfter = time.Unix(info.NotAfter, 0)
	)

	return now.After(notAfter)
}

type Application struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type Access struct {
	AccessToken  string `json:"access_token"`
	ID           string `json:"id_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (a Application) Login(code string) (_ Access, err error) {
	var w bytes.Buffer
	if err := json.NewEncoder(&w).Encode(map[string]string{
		"code":          code,
		"client_id":     a.ClientID,
		"client_secret": a.ClientSecret,
		"redirect_uri":  a.RedirectURI,
		"grant_type":    "authorization_code",
	}); err != nil {
		return Access{}, fmt.Errorf("could not encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/token", &w)
	if err != nil {
		return Access{}, fmt.Errorf("could not create request: %v", err)
	}

	req.Header.Set("Content-Type", fmt.Sprint(httputils.MediaTypeJSON))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return Access{}, fmt.Errorf("could not fetch token: %v", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil || len(body) == 0 {
			return Access{}, fmt.Errorf("could not fetch token: %v", resp.Status)
		}
		return Access{}, fmt.Errorf("could not fetch token: %v: %s", resp.Status, body)
	}

	var tk Access
	if err := json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		return Access{}, fmt.Errorf("could not decode response body: %v", err)
	}

	return tk, nil
}
