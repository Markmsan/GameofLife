package models

type SaveGameRequest struct {
	Name   string  `json:"name"`
	Grid   [][]int `json:"grid"`
	Height int     `json:"height"`
	Length int     `json:"length"`
}

type LoadGameResponse struct {
	Name   string  `json:"name"`
	Grid   [][]int `json:"grid"`
	Height int     `json:"height"`
	Length int     `json:"length"`
}
