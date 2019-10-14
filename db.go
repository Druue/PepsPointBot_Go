package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
)

type User struct {
	discordId string
	nickname  string
}
type Points struct {
	giver    string
	receiver string
	amount   int64
}

/*

     users
__________________
discord_id: text
nick_name: text








points
______________
id: text
receiver_id: text
giver_id: text
amount: in64/bigint



*/

func openDBConnection() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		SECRET.DB_HOST, SECRET.DB_PORT, SECRET.DB_USER, SECRET.DB_PASSWORD, SECRET.DB_NAME)
	db, err := sql.Open("postgres", psqlInfo)
	logErr(err)
	err = db.Ping()
	logErr(err)
	DB = db
}

func startupAddAllUsers(users []string) {
	q := ""
	usersInterface := make([]interface{}, len(users))
	for i := 0; i < len(users); i++ {
		usersInterface[i] = users[i]
		q += "INSERT INTO users (discord_id) VALUES ($" + strconv.Itoa(i+1) + ")"
	}
	_, err := DB.Query(q, usersInterface...)
	logErr(err)
}

func setUsersNickname(user User) {
	stmt, err := DB.Prepare("UPDATE users SET nick_name = $2 WHERE discord_id = $1")
	logErr(err)
	_, err = stmt.Exec(user.discordId, user.nickname)
	logErr(err)
}

func getUser(discordId string) *User {
	rows, err := DB.Query("SELECT nick_name FROM users WHERE discord_id = $1", discordId)
	logErr(err)
	for rows.Next() {
		var nickname string
		err = rows.Scan(&nickname)
		logErr(err)
		return &User{
			nickname:  nickname,
			discordId: discordId,
		}
	}
	return nil
}

func getUsersNicknameOr(discordId string, alternative string) string {
	user := getUser(discordId)
	if user == nil {
		return alternative
	}
	return user.nickname
}

func giveUserPoints(giver string, receiver string, amount int64) {
	DB.QueryRow("INSERT INTO points (id, receiver_id, giver_id, amount) VALUES ($4, $3, $2, $1) ON CONFLICT (id) UPDATE points SET amount = amount + $1 WHERE id = $4", amount, giver, receiver, giver+"_"+receiver)
}

func getUsersPointsReceived(discordId string) ([]*Points, []string) {
	rows, err := DB.Query("SELECT points.giver_id, users.nick_name, points.amount FROM points INNER JOIN users ON users.discord_id = points.giver_id WHERE points.receiver_id = $1", discordId)
	logErr(err)
	var points []*Points
	var nicknames []string
	for rows.Next() {
		var giverId string
		var nickname string
		var amount int64
		err = rows.Scan(&giverId, &nickname, &amount)
		logErr(err)
		nicknames = append(nicknames, nickname)
		points = append(points, &Points{
			giver:    giverId,
			receiver: discordId,
			amount:   amount,
		})
	}
	return points, nicknames
}

func getUsersPointsGiven(discordId string) ([]*Points, []string) {
	rows, err := DB.Query("SELECT points.receiver_id, users.nick_name, points.amount FROM points INNER JOIN users ON users.discord_id = points.receiver_id WHERE points.giver_id = $1", discordId)
	logErr(err)
	var points []*Points
	var nicknames []string
	for rows.Next() {
		var giverId string
		var nickname string
		var amount int64
		err = rows.Scan(&giverId, &nickname, &amount)
		logErr(err)
		nicknames = append(nicknames, nickname)
		points = append(points, &Points{
			giver:    giverId,
			receiver: discordId,
			amount:   amount,
		})
	}
	return points, nicknames
}
