package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	UserID    int    `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	EndPoint  string `json:"endpoint"`
	Token     `json:"-"`
}

type Token struct {
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpireDate   time.Time `json:"-"`
}

var db *sql.DB

const (
	layout = "2006-01-02 15:04:05.999999999-07:00"
)

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	//users
	u := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		endpoint TEXT,
		avatar_url TEXT,
		access_token TEXT,
		refresh_token TEXT,
		expire_date TEXT
	);`

	_, err = db.Exec(u)
	if err != nil {
		log.Fatal(err)
	}

	//cookiejar
	a := `
	CREATE TABLE IF NOT EXISTS cookiejar (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		cookie TEXT NOT NULL UNIQUE,
		FOREIGN KEY(user_id) REFERENCES users(id) On DELETE CASCADE
	);`

	_, err = db.Exec(a)
	if err != nil {
		log.Fatal(err)
	}

}

func GetUserByCookie(cookie string) (*User, error) {

	query := `
	SELECT users.id, users.name, users.endpoint, users.avatar_url, users.access_token, users.refresh_token, users.expire_date
	FROM users users
	JOIN cookiejar co ON users.id = co.user_id
	WHERE co.cookie = ?`
	row := db.QueryRow(query, cookie)

	var user User
	var ExpireStr string
	err := row.Scan(&user.UserID, &user.Name, &user.EndPoint, &user.AvatarUrl, &user.AccessToken, &user.RefreshToken, &ExpireStr)
	if err != nil {
		fmt.Println(err)
		return &User{}, err
	}

	user.ExpireDate, err = time.Parse(layout, ExpireStr)
	if err != nil {
		fmt.Println(err)
		return &User{}, err
	}

	return &user, nil
}

func SaveUser(user User) error {
	query := `INSERT INTO users (id, name, endpoint, avatar_url, access_token, refresh_token, expire_date) VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET 
			name = excluded.name,
			endpoint = excluded.endpoint,
			avatar_url = excluded.avatar_url,
			access_token = excluded.access_token,
			refresh_token = excluded.refresh_token,
			expire_date = excluded.expire_date;
	`

	_, err := db.Exec(query, user.UserID, user.Name, user.EndPoint, user.AvatarUrl, user.AccessToken, user.RefreshToken, user.ExpireDate)
	return err
}

func SaveCookie(userID int, cookie string) error {
	query := "INSERT INTO cookiejar (user_id, cookie) VALUES (?, ?)"
	_, err := db.Exec(query, userID, cookie)
	return err
}

func UpdateUserTokens(userID int, auth Token) error {
	query := "UPDATE users SET auth_token = ?, refresh_token = ?, expire_date = ? WHERE id = ?"
	_, err := db.Exec(query, auth.AccessToken, auth.RefreshToken, auth.ExpireDate, userID)
	return err
}

func UpdateUserEndPoint(userID int, endPoint string) error {
	query := "UPDATE users SET endpoint = ? WHERE id = ?"
	_, err := db.Exec(query, endPoint, userID)
	return err
}
