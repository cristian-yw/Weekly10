package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Cinema struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
