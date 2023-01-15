package domain

type Category struct {
	Id       string `json:"id"`
	SquareId string `json:"squareId"`
	Name     string `json:"name"`
	Version  int64  `json:"version"`
}
