package main

type Config struct {
	Connection Connection `json:"connection"`
	Mode       string     `json:"mode"`
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
