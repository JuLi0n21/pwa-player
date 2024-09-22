package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	OsuApiUrl = "https://osu.ppy.sh/api/v2"
)

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
		//request new token?
	}

	return &OsuApiClient{
		User:    user,
		BaseURL: OsuApiUrl,
		HTTPclient: &http.Client{
			Timeout: time.Minute * 2},
	}, nil
}

func (c *OsuApiClient) Me() (*ApiUser, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/me", OsuApiUrl), nil)
	if err != nil {
		return nil, err
	}

	res := ApiUser{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil

}

func (c *OsuApiClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.User.AccessToken))

	res, err := c.HTTPclient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func LoginMiddlePage(w http.ResponseWriter, r *http.Request) {
	cookie, ok := r.Context().Value("cookie").(string)

	if !ok || cookie == "" {
		fmt.Println(cookie, ok)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var clientid = os.Getenv("CLIENT_ID")
	var redirect_uri = os.Getenv("REDIRECT_URI") + "/oauth/code"

	html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Login Required</title>
		</head>
		<body>
			<p>Redirecting...</p>
			<a href="%s">Click here if ur not being Redirected!</a>
		<script>
			window.location.href = "%s";
		</script>
		</body>
		</html>
		`

	loginURL := fmt.Sprintf("https://osu.ppy.sh/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		clientid,
		redirect_uri,
		strings.Join(scopes, " "),
		cookie)
	fmt.Fprintf(w, html, loginURL, loginURL)
	return
}

func Oauth(w http.ResponseWriter, r *http.Request) {
	cookie, ok := r.Context().Value("cookie").(string)

	if !ok || cookie == "" {
		fmt.Println(cookie, ok)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	q := r.URL.Query()

	code := q.Get("code")
	state := q.Get("state")

	if state != cookie {

		fmt.Println(state, cookie)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var clientid = os.Getenv("CLIENT_ID")
	var client_secret = os.Getenv("CLIENT_SECRET")
	var redirect_uri = os.Getenv("REDIRECT_URI") + "/oauth/code"
	//request accesstoken
	body := url.Values{
		"client_id":     {clientid},
		"client_secret": {client_secret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirect_uri},
	}

	req, err := http.NewRequest("POST", "https://osu.ppy.sh/oauth/token", bytes.NewBufferString(body.Encode()))
	if err != nil {
		fmt.Println("Error: ", err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpclient := &http.Client{
		Timeout: time.Minute,
	}

	res, err := httpclient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		fmt.Println("Error: ", err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var authToken AuthToken

	err = json.Unmarshal(data, &authToken)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: ", err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user := User{
		Token: Token{
			ExpireDate:   time.Now().Add(time.Second * time.Duration(authToken.ExpiresIn)),
			RefreshToken: authToken.RefreshToken,
			AccessToken:  authToken.AccessToken,
		},
	}

	c, err := NewOsuApiClient(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	apiuser, err := c.Me()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user.UserID = apiuser.ID
	user.Name = apiuser.Username
	user.AvatarUrl = apiuser.AvatarURL
	user.Share = false

	SaveCookie(user.UserID, cookie)
	if err = SaveUser(user); err != nil {
		fmt.Println(err)
	}

	var html = fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Login Success</title>
		<body>
		
		<input type="password" value="%s" id="myInput" disabled>
		<button onclick="copyToClipboard()">Copy Text</button>

			<script>
			function copyToClipboard() {
				var copyText = document.getElementById("myInput");
				copyText.select();
				copyText.setSelectionRange(0, 99999); // For mobile devices
				  navigator.clipboard.writeText(copyText.value);
			}

			window.close(); // Close the window after copy
			</script>
		</head>
		</html>
	`, cookie)

	fmt.Fprint(w, html)
	return

}

type AuthToken struct {
	Tokentype    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ApiUser struct {
	AvatarURL     string    `json:"avatar_url,omitempty"`
	CountryCode   string    `json:"country_code,omitempty"`
	DefaultGroup  string    `json:"default_group,omitempty"`
	ID            int       `json:"id,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	IsBot         bool      `json:"is_bot,omitempty"`
	IsDeleted     bool      `json:"is_deleted,omitempty"`
	IsOnline      bool      `json:"is_online,omitempty"`
	IsSupporter   bool      `json:"is_supporter,omitempty"`
	LastVisit     time.Time `json:"last_visit,omitempty"`
	PmFriendsOnly bool      `json:"pm_friends_only,omitempty"`
	ProfileColour string    `json:"profile_colour,omitempty"`
	Username      string    `json:"username,omitempty"`
	CoverURL      string    `json:"cover_url,omitempty"`
	Discord       string    `json:"discord,omitempty"`
	HasSupported  bool      `json:"has_supported,omitempty"`
	Interests     any       `json:"interests,omitempty"`
	JoinDate      time.Time `json:"join_date,omitempty"`
	Kudosu        struct {
		Total     int `json:"total,omitempty"`
		Available int `json:"available,omitempty"`
	} `json:"kudosu,omitempty"`
	Location     any      `json:"location,omitempty"`
	MaxBlocks    int      `json:"max_blocks,omitempty"`
	MaxFriends   int      `json:"max_friends,omitempty"`
	Occupation   any      `json:"occupation,omitempty"`
	Playmode     string   `json:"playmode,omitempty"`
	Playstyle    []string `json:"playstyle,omitempty"`
	PostCount    int      `json:"post_count,omitempty"`
	ProfileOrder []string `json:"profile_order,omitempty"`
	Title        any      `json:"title,omitempty"`
	Twitter      string   `json:"twitter,omitempty"`
	Website      string   `json:"website,omitempty"`
	Country      Country  `json:"country,omitempty"`
	Cover        struct {
		CustomURL string `json:"custom_url,omitempty"`
		URL       string `json:"url,omitempty"`
		ID        any    `json:"id,omitempty"`
	} `json:"cover,omitempty"`
	IsRestricted           bool  `json:"is_restricted,omitempty"`
	AccountHistory         []any `json:"account_history,omitempty"`
	ActiveTournamentBanner any   `json:"active_tournament_banner,omitempty"`
	Badges                 []struct {
		AwardedAt   time.Time `json:"awarded_at,omitempty"`
		Description string    `json:"description,omitempty"`
		Image2XURL  string    `json:"image@2x_url,omitempty"`
		ImageURL    string    `json:"image_url,omitempty"`
		URL         string    `json:"url,omitempty"`
	} `json:"badges,omitempty"`
	FavouriteBeatmapsetCount int `json:"favourite_beatmapset_count,omitempty"`
	FollowerCount            int `json:"follower_count,omitempty"`
	GraveyardBeatmapsetCount int `json:"graveyard_beatmapset_count,omitempty"`
	Groups                   []struct {
		ID          int    `json:"id,omitempty"`
		Identifier  string `json:"identifier,omitempty"`
		Name        string `json:"name,omitempty"`
		ShortName   string `json:"short_name,omitempty"`
		Description string `json:"description,omitempty"`
		Colour      string `json:"colour,omitempty"`
	} `json:"groups,omitempty"`
	LovedBeatmapsetCount int `json:"loved_beatmapset_count,omitempty"`
	MonthlyPlaycounts    []struct {
		StartDate string `json:"start_date,omitempty"`
		Count     int    `json:"count,omitempty"`
	} `json:"monthly_playcounts,omitempty"`
	Page struct {
		HTML string `json:"html,omitempty"`
		Raw  string `json:"raw,omitempty"`
	} `json:"page,omitempty"`
	PendingBeatmapsetCount int   `json:"pending_beatmapset_count,omitempty"`
	PreviousUsernames      []any `json:"previous_usernames,omitempty"`
	RankedBeatmapsetCount  int   `json:"ranked_beatmapset_count,omitempty"`
	ReplaysWatchedCounts   []struct {
		StartDate string `json:"start_date,omitempty"`
		Count     int    `json:"count,omitempty"`
	} `json:"replays_watched_counts,omitempty"`
	ScoresFirstCount int `json:"scores_first_count,omitempty"`
	Statistics       struct {
		Level struct {
			Current  int `json:"current,omitempty"`
			Progress int `json:"progress,omitempty"`
		} `json:"level,omitempty"`
		Pp                     float64 `json:"pp,omitempty"`
		GlobalRank             int     `json:"global_rank,omitempty"`
		RankedScore            int     `json:"ranked_score,omitempty"`
		HitAccuracy            float64 `json:"hit_accuracy,omitempty"`
		PlayCount              int     `json:"play_count,omitempty"`
		PlayTime               int     `json:"play_time,omitempty"`
		TotalScore             int     `json:"total_score,omitempty"`
		TotalHits              int     `json:"total_hits,omitempty"`
		MaximumCombo           int     `json:"maximum_combo,omitempty"`
		ReplaysWatchedByOthers int     `json:"replays_watched_by_others,omitempty"`
		IsRanked               bool    `json:"is_ranked,omitempty"`
		GradeCounts            struct {
			Ss  int `json:"ss,omitempty"`
			SSH int `json:"ssh,omitempty"`
			S   int `json:"s,omitempty"`
			Sh  int `json:"sh,omitempty"`
			A   int `json:"a,omitempty"`
		} `json:"grade_counts,omitempty"`
		Rank struct {
			Global  int `json:"global,omitempty"`
			Country int `json:"country,omitempty"`
		} `json:"rank,omitempty"`
	} `json:"statistics,omitempty"`
	SupportLevel     int `json:"support_level,omitempty"`
	UserAchievements []struct {
		AchievedAt    time.Time `json:"achieved_at,omitempty"`
		AchievementID int       `json:"achievement_id,omitempty"`
	} `json:"user_achievements,omitempty"`
	RankHistory struct {
		Mode string `json:"mode,omitempty"`
		Data []int  `json:"data,omitempty"`
	} `json:"rank_history,omitempty"`
}

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
