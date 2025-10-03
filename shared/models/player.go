package models

type Player struct {
	PlayerHistory []PlayerHistory     `json:"history"`
	PlayerPast    []PlayerPastHistory `json:"history_past"`
}
