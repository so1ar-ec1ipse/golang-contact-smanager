package models

import "errors"

type Contact struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	Email  string `json:"email"`
}

type Response struct {
	Status  string    `json:"status"`
	Data    []Contact `json:"data"`
	Message string    `json:"message"`
}

var list = []Contact{
	{Name: "one", Number: "8908901234", Email: "onea@one.com"},
	{Name: "two", Number: "8908902345", Email: "two@two.com"},
	{Name: "three", Number: "8908903456", Email: "three@three.com"},
	{Name: "four", Number: "8908904567", Email: "four@four.com"},
	{Name: "five", Number: "8908905678", Email: "five@five.com"},
	{Name: "six", Number: "8908906789", Email: "six@six.com"},
}

func Getlist() ([]Contact, error) {
	out := []Contact{}
	out = append(out, list...)

	if out != nil {
		return out, nil
	} else {
		return out, errors.New("no contacts found")
	}
}

func Setlist(in []Contact) {
	list = in
}
