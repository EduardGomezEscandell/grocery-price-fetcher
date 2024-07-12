package google

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/database/dbtypes"
	"github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/httputils"
)

type Application struct {
	clientID     string
	redirectURI  string
	clientSecret string
}

type Settings struct {
	ClientID         string `yaml:"client-id"`
	ClientSecretFile string `yaml:"client-secret-file"`
	RedirectURI      string `yaml:"redirect-uri"`
}

func NewApplication(s Settings) (Application, error) {
	out, err := os.ReadFile(s.ClientSecretFile)
	if err != nil {
		return Application{}, fmt.Errorf("could not read client secret file: %w", err)
	} else if len(out) == 0 {
		return Application{}, errors.New("empty client secret file")
	}

	return Application{
		clientID:     s.ClientID,
		clientSecret: strings.TrimSpace(string(out)),
		redirectURI:  s.RedirectURI,
	}, nil
}

type accessTokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	JWT          string `json:"id_token"`
}

type jwtData struct {
	NotAfter int64  `json:"exp"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Subject  string `json:"sub"`
}

func (a Application) Login(ctx context.Context, code string) (dbtypes.Session, error) {
	if code == "" {
		return dbtypes.Session{}, errors.New("empty code")
	}

	out, err := postGoogleAuth(ctx, a, map[string]any{
		"code":       code,
		"grant_type": "authorization_code",
	})
	if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not get token: POST: %v", err)
	}

	hash := sha512.Sum512([]byte(code))
	cred := base64.StdEncoding.EncodeToString(hash[:])

	return a.newSession(cred, out)
}

func (a Application) Refresh(ctx context.Context, session dbtypes.Session) (dbtypes.Session, error) {
	out, err := postGoogleAuth(ctx, a, map[string]any{
		"refresh_token": session.RefreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not refresh token: POST: %v", err)
	}

	return a.newSession(session.ID, out)
}

func (a Application) newSession(ID string, accessData []byte) (dbtypes.Session, error) {
	var access accessTokenData
	if err := json.Unmarshal(accessData, &access); err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not get token: could not decode response body: %v", err)
	}

	if access.TokenType != "Bearer" {
		return dbtypes.Session{}, fmt.Errorf("unexpected token type: %s", access.TokenType)
	}

	jwtData, err := decodeJWT(access.JWT)
	if err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not get token: %v", err)
	}

	if err := a.validate(jwtData); err != nil {
		return dbtypes.Session{}, fmt.Errorf("could not validate token: %v", err)
	}

	return dbtypes.Session{
		ID:           ID,
		AccessToken:  access.AccessToken,
		RefreshToken: access.RefreshToken,
		NotAfter:     time.Unix(jwtData.NotAfter, 0),
		User:         jwtData.Subject,
	}, nil
}

func postGoogleAuth(ctx context.Context, app Application, data map[string]any) ([]byte, error) {
	data["client_id"] = app.clientID
	data["client_secret"] = app.clientSecret
	data["redirect_uri"] = app.redirectURI

	var w bytes.Buffer
	if err := json.NewEncoder(&w).Encode(data); err != nil {
		return nil, fmt.Errorf("could not encode request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", &w)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	req.Header.Set("Content-Type", fmt.Sprint(httputils.MediaTypeJSON))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %v", resp.Status)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	return body, nil
}

func decodeJWT(s string) (jwtData, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return jwtData{}, errors.New("invalid JWT format")
	}

	b, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return jwtData{}, fmt.Errorf("could not decode JWT: %v", err)
	}

	var info jwtData
	if err := json.Unmarshal(b, &info); err != nil {
		return jwtData{}, fmt.Errorf("could not decode JWT: %v", err)
	}

	return info, nil
}

func (a Application) validate(jwt jwtData) error {
	switch jwt.Issuer {
	case "https://accounts.google.com", "accounts.google.com":
	default:
		return fmt.Errorf("invalid issuer: %s", jwt.Issuer)
	}

	if jwt.Audience != a.clientID {
		return fmt.Errorf("invalid audience: %s", jwt.Audience)
	}

	if time.Now().Unix() > jwt.NotAfter {
		return fmt.Errorf("token expired")
	}

	return nil
}
