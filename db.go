package main

import . "github.com/ahmetb/go-linq"

var names = make(map[string]string)
var points []*Point

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

func addPoint(origin string, recipient string, amount int) {
	var possiblePoints []*Point
	From(points).WhereT(func(p *Point) bool {
		return p.origin == origin && p.recipient == recipient
	}).ToSlice(&possiblePoints)

	if len(possiblePoints) == 0 {
		point := &Point{
			origin:    origin,
			recipient: recipient,
			amount:    amount,
		}
		points = append(points, point)
	} else {
		From(points).ForEachT(func(p *Point) {
			p.amount += amount
		})
	}
}

type Point struct {
	origin    string
	recipient string
	amount    int
}
