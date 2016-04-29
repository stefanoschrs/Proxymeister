package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Proxy struct {
	Ip			string
	Port		int
	Status		int
	LastChecked string
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	checkErr(err)

	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Proxy(
			ip VARCHAR(15),
			port INTEGER,
			status INTEGER DEFAULT 0,
			lastChecked DATETIME DEFAULT 0,
			PRIMARY KEY ( ip, port )
		)
	`)
	checkErr(err)
}

func InsertProxy(db *sql.DB, ip string, port int) {
	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO Proxy(ip, port) VALUES (?, ?)`)
	checkErr(err)
	defer stmt.Close()

	_, err2 := stmt.Exec(ip, port)
	checkErr(err2)
}

func UpdateProxy(db *sql.DB, ip string, port, status int) {
	stmt, err := db.Prepare(`
		UPDATE Proxy
		SET status = ?,
			lastChecked = CURRENT_TIMESTAMP
		WHERE 	ip = ? AND
				port = ?
	`)
	checkErr(err)
	defer stmt.Close()

	_, err2 := stmt.Exec(status, ip, port)
	checkErr(err2)
}

func SelectProxies(db *sql.DB, status int) []Proxy {
	rows, err := db.Query(`
		SELECT *
		FROM Proxy
		WHERE STATUS = ?
		ORDER BY lastChecked DESC
	`, status)
	checkErr(err)
	defer rows.Close()

	var result []Proxy
	for rows.Next() {
		item := Proxy{}
		err2 := rows.Scan(&item.Ip, &item.Port, &item.Status, &item.LastChecked)
		checkErr(err2)

		result = append(result, item)
	}

	return result
}

func SelectRecentProxies(db *sql.DB) []Proxy {
	rows, err := db.Query(`
		SELECT *
		FROM Proxy
		WHERE lastChecked < date('now', '10 minutes')
		ORDER BY lastChecked ASC
	`)
	checkErr(err)
	defer rows.Close()

	var result []Proxy
	for rows.Next() {
		item := Proxy{}
		err2 := rows.Scan(&item.Ip, &item.Port, &item.Status, &item.LastChecked)
		checkErr(err2)

		result = append(result, item)
	}

	return result
}

func SelectAllProxies(db *sql.DB) []Proxy {
	rows, err := db.Query(`
		SELECT *
		FROM Proxy
		ORDER BY lastChecked DESC
	`)
	checkErr(err)
	defer rows.Close()

	var result []Proxy
	for rows.Next() {
		newProxy := Proxy{}
		err2 := rows.Scan(&newProxy.Ip, &newProxy.Port, &newProxy.Status, &newProxy.LastChecked)
		checkErr(err2)

		result = append(result, newProxy)
	}

	return result
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
