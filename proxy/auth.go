package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	OsuApiUrl = "https://osu.ppy.sh/api/v2"
)

var clientid = os.Getenv("CLIENT_ID")
var clientsecret = os.Getenv("CLIENT_SECRET")
var redirect_uri = os.Getenv("REDIRECT_URI")
var scopes = []string{
	"public",
	"identify",
}

type OsuApiClient struct {
	User       User
	BaseURL    string
	HTTPclient *http.Client
}

func NewOsuApiClient(user User) (*OsuApiClient, error) {

	if user.Token == (Token{}) {
		return nil, errors.New("No Valid Credentials")
	}

	if time.Now().After(user.ExpireDate) {

	}

	return &OsuApiClient{
		User:    user,
		BaseURL: OsuApiUrl,
		HTTPclient: &http.Client{
			Timeout: time.Minute * 2},
	}, nil
}

func LoginRedirect(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Context())
	cookie, ok := r.Context().Value("cookie").(string)

	if !ok || cookie == "" {
		fmt.Println(cookie, ok)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r,
		fmt.Sprintf("https://osu.ppy.sh/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
			clientid,
			redirect_uri,
			strings.Join(scopes, " "),
			cookie),
		http.StatusPermanentRedirect)
}
