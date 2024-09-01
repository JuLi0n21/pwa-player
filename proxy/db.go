package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	UserID    int
	Name      string
	AvatarUrl string
	EndPoint  string
	Token
}

type Token struct {
	AuthToken    string
	RefreshToken string
	ExpireDate   time.Time
}

var db *sql.DB

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
		auth_token TEXT,
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
		cookie TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) On DELETE CASCADE
	);`

	_, err = db.Exec(a)
	if err != nil {
		log.Fatal(err)
	}

}

func GetUserByCookie(cookie string) (*User, error) {
	query := `
	SELECT users.id, users.name, users.endpoint, users.avatar_url, users.auth_token, users.refresh_token, users.expire_date
	FROM users users
	JOIN cookiejar cookie ON users.id = cookie.user_id
	WHERE cookie.cookie = ?`
	row := db.QueryRow(query, cookie)

	var user User
	err := row.Scan(&user.UserID, &user.Name, &user.EndPoint, &user.AvatarUrl, &user.AuthToken, &user.RefreshToken, &user.ExpireDate)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func SaveCookie(userID int, cookie string) error {
	query := "INSERT INTO cookiejar (user_id, cookie) VALUES (?, ?)"
	_, err := db.Exec(query, userID, cookie)
	return err
}

func UpdateUserTokens(userID int, auth Token) error {
	query := "UPDATE users SET auth_token = ?, refresh_token = ?, expire_date = ? WHERE id = ?"
	_, err := db.Exec(query, auth.AuthToken, auth.RefreshToken, auth.ExpireDate, userID)
	return err
}

func UpdateUserEndPoint(userID int, endPoint string) error {
	query := "UPDATE users SET endpoint = ? WHERE id = ?"
	_, err := db.Exec(query, endPoint, userID)
	return err
}
