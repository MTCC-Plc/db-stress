package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

var (
	config         Config
	err            error
	lastUpdateUnix int64
)

func main() {
	err = loadConfig()
	if err != nil {
		log.Println("Loading config: ", err.Error())
	}
	start := time.Now()
	if config.Mode == "parallel" {
		fmt.Println("Running", len(config.Tests), "tests in parallel")
	}
	var wg sync.WaitGroup
	for _, test := range config.Tests {
		if config.Mode == "parallel" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := runWorkers(test)
				if err != nil {
					log.Println("runner", err.Error())
				}
			}()
		} else {
			err := runWorkers(test)
			if err != nil {
				log.Println("runner", err.Error())
			}
		}
	}
	wg.Wait()
	fmt.Println("\nTotal Time: ", time.Since(start))
}

func loadConfig() error {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFile, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return err
	}
	configFile.Close()

	return nil
}

func connect() (*sql.DB, error) {
	conn, err := sql.Open(config.Connection.Provider, config.Connection.ConnectionString)
	if err != nil {
		return conn, nil
	}
	return conn, nil
}

func runWorkers(test Test) error {
	if test.Workers == 0 {
		test.Workers = 1
	}
	if config.Mode != "parallel" {
		fmt.Println("\nQUERY:\t\t", test.Query)
		fmt.Println("ITERATIONS:\t", test.Iterations)
		if test.Workers > 1 {
			fmt.Println("WORKERS:\t", test.Workers)
		}
	}
	start := time.Now()
	var wg sync.WaitGroup
	perWorker := test.Iterations / test.Workers
	remainder := test.Iterations % test.Workers
	fmt.Print("\033[s")
	for i := 0; i < test.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, nil := connect()
			if err != nil {
				log.Println("connect", err.Error())
				return
			}
			defer conn.Close()
			iterationStart := i * perWorker
			iterationEnd := iterationStart + perWorker
			if i != 0 && i == test.Workers-1 {
				iterationEnd += remainder
			}
			err := runner(conn, test, iterationStart, iterationEnd)
			if err != nil {
				log.Println("runner", err.Error())
			}
		}()
	}
	wg.Wait()
	fmt.Print("\033[u\033[K")
	if config.Mode == "parallel" {
		fmt.Println("\nQUERY:\t\t", test.Query)
		fmt.Println("ITERATIONS:\t", test.Iterations)
		if test.Workers > 1 {
			fmt.Println("WORKERS:\t", test.Workers)
		}
	}
	since := time.Since(start)
	fmt.Print("\tTime: ", since, "\n")
	fmt.Print("\tAverage Time: ", since/time.Duration(test.Iterations), "\n")
	return nil
}

func runner(conn *sql.DB, test Test, iterationStart int, iterationEnd int) error {
	runnerTotal := time.Duration(0)
	for i := iterationStart; i < iterationEnd; i++ {
		start := time.Now()
		_, err := conn.Exec(test.Query)
		if err != nil {
			return err
		}
		since := time.Since(start)
		runnerTotal += since
		progressUpdate(i, since, runnerTotal)
	}
	return nil
}

func progressUpdate(i int, iterationTime time.Duration, totalTime time.Duration) {
	if config.Mode == "parallel" {
		return
	}
	currentSecond := time.Now().Unix()
	if currentSecond <= lastUpdateUnix {
		return
	}
	lastUpdateUnix = currentSecond
	iteration := i + 1
	fmt.Print("\033[u\033[K")
	fmt.Print("\tIteration: ", iteration, ", Iteration Time: ", iterationTime, ", Average Time: ", totalTime/time.Duration(iteration))
}
