package main

import (
	"database/sql"
	"fmt"
	"log"

	. "github.com/ahmetb/go-linq"
	_ "github.com/mattn/go-sqlite3"
)

var names = make(map[string]string)
var transactions []*Transaction

type Transaction struct {
	origin    string
	recipient string
	amount    int
}

type Points struct {
	origin string
	amount int
}

type UserInfo struct {
	id     int
	name   string
	points []Points
}

func openDBConnection(dbCon string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", dbCon)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

func noRows(err error) {
	if err == sql.ErrNoRows {
		fmt.Println("No rows to return!")
	}
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
	user_info						transaction_info

	id      |  name  |				giver  | recipient | amount
	int		| string |				string |  string   | int
*/

func getPointsList() []UserInfo {
	var (
		id   int
		name string
	)

	rows, err := DB.Query("SELECT user_id, name FROM user_info")
	logErr(err)

	var users []UserInfo

	for rows.Next() {
		err := rows.Scan(&id, &name)
		logErr(err)

		users = append(users, getUserPoints(id))
	}

	err = rows.Err()
	logErr(err)

	return users
}

func getUserPoints(userid int) UserInfo {
	var (
		username string
		giverid  int
	)
	rows, err := DB.Query(`
		SELECT giver.id, transaction.recipient
		FROM transaction_info transaction, user_info giver
			INNER JOIN user_info user 
				ON (recipient = user.name AND user.id = ?)
			WHERE giver.name = transaction.giver
		`, userid)
	logErr(err)

	var points []Points
	for rows.Next() {
		err := rows.Scan(&giverid, &username)
		logErr(err)

		points = append(points, getUserPointsFromGiver(userid, giverid))
	}

	err = rows.Err()
	logErr(err)

	return UserInfo{
		id:     userid,
		name:   username,
		points: points,
	}
}

func getUserPointsFromGiver(userid int, giverid int) Points {
	rows, err := DB.Query(`
		SELECT transaction_info.amount
		FROM transaction_info, user_info as user
			INNER JOIN user_info as giver 
				ON (transaction_info.giver = giver.name 
					AND giver.ID = ?)
			WHERE (transaction_info.recipient = user.name
					AND user.ID = ?)
	`, giverid, userid)
	logErr(err)

	var points Points
	err = DB.QueryRow(`
		SELECT name FROM user_info
		WHERE user_info.ID = ?
		`, giverid).Scan(&points.origin)

	noRows(err)
	logErr(err)

	currPoints := 0
	for rows.Next() {
		err := rows.Scan(&currPoints)
		logErr(err)

		points.amount += currPoints
	}

	return points
}

// IGNORE -- LOCAL TESTING

func setName(id string, name string) {
	names[id] = name
}

func getName(id string) (string, bool) {
	name, ok := names[id]
	return name, ok
}

func getNameOr(id string, otherwise string) string {
	name, ok := getName(id)
	if ok {
		return name
	}
	return otherwise
}

func addTransaction(origin string, recipient string, amount int) {
	var possiblePoints []*Transaction
	From(transactions).WhereT(func(p *Transaction) bool {
		return p.origin == origin && p.recipient == recipient
	}).ToSlice(&possiblePoints)

	if len(possiblePoints) == 0 {
		transaction := &Transaction{
			origin:    origin,
			recipient: recipient,
			amount:    amount,
		}
		transactions = append(transactions, transaction)
	} else {
		From(transactions).ForEachT(func(p *Transaction) {
			p.amount += amount
		})
	}
}
