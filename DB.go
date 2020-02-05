package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Uplink to the Scrin mothership
func DBConnect(c *Config) (db *sql.DB, err error) {
	db, _ = sql.Open("mysql", c.Database.ConnectionString)
	err = db.Ping()
	return
}

func UserByName(db *sql.DB, uname string) (u UserProfile, err error) {
	selForm, err := db.Prepare(`SELECT
		Username, DisplayName, JoinedOn, Shadow, APISecret
		FROM Users WHERE Username=?`)
	if err != nil { return }
	err = selForm.QueryRow(uname).Scan(
		&u.Username,
		&u.DisplayName,
		&u.JoinedOn,
		&u.Shadow,
		&u.APISecret)

	/* if err == sql.ErrNoRows {
		return
	} else {
		return
	} */

	return
}


func (ws WebSession) Associated(db *sql.DB) (s Session, err error) {
	selForm, err := db.Prepare(`SELECT
		SessID, Username, Expires
		FROM Sessions WHERE SessID=?`)
	if err != nil { return }
	err = selForm.QueryRow(ws.SessID).Scan(
		&s.SessID,
		&s.Username,
		&s.Expires)

	return
}

func (ws WebSession) Associate(db *sql.DB, uname string) (error) {
	q := `INSERT INTO Sessions
		(SessID, Username) VALUES (?, ?)`
	insForm, err := db.Prepare(q)
	if err != nil { return err }
	_, err = insForm.Exec(ws.SessID, uname)
	return err
}

func (ws WebSession) Disassociate(db *sql.DB) (error) {
	q := `DELETE FROM Sessions
		WHERE SessID=?`
	delForm, err := db.Prepare(q)
	if err != nil { return err }
	_, err = delForm.Exec(ws.SessID)
	return err
}

func (u UserProfile) Bookmarks(db *sql.DB) (marks []Bookmark, err error) {
	q := `SELECT
		BId, Username, URL, Title, Unread, Archived, AddedOn
		FROM Bookmarks WHERE Username=?`
	selForm, err := db.Prepare(q)
	if err != nil { return }
	var m Bookmark
	rows, err := selForm.Query(u.Username)
	for rows.Next() { rows.Scan(
		&m.BId,
		&m.Username,
		&m.URL,
		&m.Title,
		&m.Unread,
		&m.Archived,
		&m.AddedOn)
		marks = append(marks, m)
	}
	return
}

func (u UserProfile) Create(db *sql.DB, pass string) (UserProfile, error) {
	shadow := DoShadow(pass)
	apisecret := APISecret()

	q := `INSERT INTO Users
		(Username, DisplayName, Shadow, APISecret) VALUES
		(?, ?, ?, ?)`

	insForm, err := db.Prepare(q)
	if err != nil { return UserProfile{}, err }
	_, err = insForm.Exec(
		u.Username,
		u.DisplayName,
		shadow,
		apisecret)
	if err != nil { return UserProfile{}, err }

	return UserByName(db, u.Username)
}

func LetMeIn(db *sql.DB, uname, pass string) (UserProfile, error) {
	u, err := UserByName(db, uname)
	if err != nil { return u, err }

	fail := CompareShadow(u.Shadow, pass)
	if fail != nil { return u, err }

	return u, nil
}