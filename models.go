package main

type Config struct {
	Label      string     `json:"label"`
	Connection Connection `json:"connection"`
	Mode       string     `json:"mode"`
	Log        bool       `json:"log"`
	Tests      []Test     `json:"tests"`
}

type Connection struct {
	Provider         string `json:"provider"`
	ConnectionString string `json:"connectionString"`
}

type Test struct {
	Query      string `json:"query"`
	Iterations int    `json:"iterations"`
	Workers    int    `json:"workers"`
}
