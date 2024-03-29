package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
)

type User struct {
	discordId string
	nickname  sql.NullString
}
type Points struct {
	giver    string
	receiver string
	amount   int64
}

/*

     users
__________________
discord_id: text PRIMARY KEY NOT NULL
nick_name: text








points
______________
id: text PRIMARY KEY NOT NULL
receiver_id: text NOT NULL
giver_id: text NOT NULL
amount: in64/bigint






prefixes
________________
guild_id: text PRIMARY KEY NOT NULL
prefix: text NOT NULL



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
	query := "INSERT INTO users (discord_id) VALUES "
	var values []interface{}
	for i, s := range users {
		values = append(values, s)
		numFields := 1
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1]
	query += "ON CONFLICT DO NOTHING"
	_, err := DB.Exec(query, values...)
	logErr(err)
}

func startupAddAllGuilds(guilds []string) {
	query := "INSERT INTO prefixes (guild_id, prefix) VALUES "
	var values []interface{}
	for i, s := range guilds {
		values = append(values, s, "?")
		numFields := 2
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1]
	query += "ON CONFLICT DO NOTHING"
	_, err := DB.Exec(query, values...)
	logErr(err)
}

func setPrefixForGuild(guildId string, prefix string) {
	stmt, err := DB.Prepare("UPDATE prefixes SET prefix = $2 WHERE guild_id = $1")
	logErr(err)
	_, err = stmt.Exec(guildId, prefix)
	logErr(err)
	defer stmt.Close()
}

func getGuildPrefix(guildId string) *string {
	rows, err := DB.Query("SELECT prefix FROM prefixes WHERE guild_id = $1", guildId)
	logErr(err)
	defer rows.Close()
	for rows.Next() {
		var prefix string
		err = rows.Scan(&prefix)
		logErr(err)
		return &prefix
	}
	return nil
}

func getUserFromNickname(nickname string) sql.NullString {
	rows, err := DB.Query("SELECT discord_id FROM users WHERE nick_name = $1", nickname)
	logErr(err)
	defer rows.Close()
	for rows.Next() {
		var discordId sql.NullString
		err = rows.Scan(&discordId)
		logErr(err)
		return discordId
	}
	return sql.NullString{
		String: "",
		Valid:  false,
	}
}

func setUsersNickname(user *User) {
	stmt, err := DB.Prepare("UPDATE users SET nick_name = $2 WHERE discord_id = $1")
	logErr(err)
	_, err = stmt.Exec(user.discordId, user.nickname)
	logErr(err)
	defer stmt.Close()
}

func getUser(discordId string) *User {
	rows, err := DB.Query("SELECT nick_name FROM users WHERE discord_id = $1", discordId)
	logErr(err)
	defer rows.Close()
	for rows.Next() {
		var nickname sql.NullString
		err = rows.Scan(&nickname)
		logErr(err)
		return &User{
			nickname:  nickname,
			discordId: discordId,
		}
	}
	return nil
}

func getUsersNicknameOr(discordId string, alternative sql.NullString) sql.NullString {
	user := getUser(discordId)
	if user == nil {
		return alternative
	}
	return user.nickname
}

func giveUserPoints(giver string, receiver string, amount int64) {
	id := giver + "_" + receiver
	rows, err := DB.Query("INSERT INTO points (id, receiver_id, giver_id, amount) VALUES ($4, $3, $2, $1) ON CONFLICT (id) DO UPDATE SET amount = (points.amount + $1) WHERE points.id = $4 RETURNING points.amount", amount, receiver, giver, id)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var amount int64
		err = rows.Scan(&amount)
		if amount == 0 {
			rows, err := DB.Query("DELETE FROM points WHERE id = $1", id)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
		}
	}
	defer rows.Close()
}

func getUsersPointsReceived(discordId string) ([]*Points, []sql.NullString) {
	rows, err := DB.Query("SELECT points.giver_id, users.nick_name, points.amount FROM points INNER JOIN users ON users.discord_id = points.giver_id WHERE points.receiver_id = $1 ORDER BY points.amount DESC", discordId)
	logErr(err)
	var points []*Points
	var nicknames []sql.NullString
	defer rows.Close()
	for rows.Next() {
		var giverId string
		var nickname sql.NullString
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

func getUsersPointsReceivedFromOtherUser(receiverId string, giverId string) sql.NullInt64 {
	rows, err := DB.Query("SELECT points.amount FROM points WHERE points.id = $1", receiverId+"_"+giverId)
	logErr(err)
	defer rows.Close()
	for rows.Next() {
		var amount sql.NullInt64
		err = rows.Scan(&amount)
		return amount
	}
	return sql.NullInt64{
		Int64: 0,
		Valid: false,
	}
}

func getUsersPointsGiven(discordId string) ([]*Points, []sql.NullString) {
	rows, err := DB.Query("SELECT points.receiver_id, users.nick_name, points.amount FROM points INNER JOIN users ON users.discord_id = points.receiver_id WHERE points.giver_id = $1 ORDER BY points.amount DESC", discordId)
	logErr(err)
	defer rows.Close()
	var points []*Points
	var nicknames []sql.NullString
	for rows.Next() {
		var receiverId string
		var nickname sql.NullString
		var amount int64
		err = rows.Scan(&receiverId, &nickname, &amount)
		logErr(err)
		nicknames = append(nicknames, nickname)
		points = append(points, &Points{
			giver:    discordId,
			receiver: receiverId,
			amount:   amount,
		})
	}
	return points, nicknames
}
