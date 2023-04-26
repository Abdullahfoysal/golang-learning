package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

// var tableName = "user_info"
var DBName = "user_info_db"

type UserInfo struct {
	ID       int64
	Name     string
	BirthDay string
}

func main() {
	// Capture connection properties.
	// Get a database handle.
	sqlDatbaseSetup()

	userId, err := addUserInfo(UserInfo{
		ID:       2,
		Name:     "foysal",
		BirthDay: dateToString(stringToDate("2023-04-24")),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added user: %v\n", userId)

	user, err := userById(2)
	fmt.Println(user)
	message, err := birthdayMessage(user.ID)
	//fmt.Print("BirthDayMessage" + message + "**" + err.Error())
	fmt.Println(message)

}

func dateToString(date time.Time) string {
	return date.Format("2006-01-02")
}
func stringToDate(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

func birthdayMessage(userId int64) (string, error) {
	// 1. get user info
	user, err := userById(userId)

	userBirthday := user.BirthDay
	//fmt.Println(userBirthday)

	message := daysToBirthday(stringToDate(userBirthday))

	return message, err

}

func daysToBirthday(birthday time.Time) string {
	now := time.Now()
	today := stringToDate(dateToString(now))
	fmt.Println(today)

	if birthday.Month() == today.Month() && birthday.Day() == today.Day() {
		return "Happy Birthday"
	} else if birthday.After(today) {
		day := birthday.Sub(today).Hours() / 24
		return strconv.Itoa(int(day)) + " days left for your birthday"
	} else {
		day := birthday.Sub(today).Hours() / 24
		return strconv.Itoa(int(-day)) + " days ago was your birthday"

	}

}

func sqlDatbaseSetup() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: DBName,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	createUserInfoTable()

}

// userByID queries for the user with the specified ID.
func userById(id int64) (UserInfo, error) {
	// An album to hold data from the returned row.
	var alb UserInfo

	row := db.QueryRow("SELECT * FROM user_info WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Name, &alb.BirthDay); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("userById %d: no such user", id)
		}
		return alb, fmt.Errorf("userById %d: %v", id, err)
	}

	return alb, nil
}

// addUser adds the specified user to the database,
// returning the user ID of the new entry
func addUserInfo(alb UserInfo) (int64, error) {

	result, err := db.Exec("REPLACE INTO user_info (id,name, birthday) VALUES (?, ?, ?)", alb.ID, alb.Name, alb.BirthDay)
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	return id, nil
}

func createUserInfoTable() (int64, error) {

	result, err := db.Exec("CREATE TABLE IF NOT EXISTS user_info (id INT NOT NULL,name  VARCHAR(128) NOT NULL, birthday VARCHAR(128) NOT NULL,PRIMARY KEY (`id`));")

	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}

	return id, nil
}
