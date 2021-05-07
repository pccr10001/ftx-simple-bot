package model

type Signal struct {
	Direction string `json:"direction"`
	Time      int64  `json:"time"`
	Symbol    string `json:"symbol"`
	Interval  int    `json:"interval"`
	Latest    bool   `json:"latest"`
}
