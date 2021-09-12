package main

import (
	"alien/simulator"
	"alien/world"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	fileName := fmt.Sprintf("%s/data.txt", dir)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = file.Close() }()

	var numAlien int
	if len(os.Args) >= 2 {
		numAlien, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}
	}

	w := world.CreateFromFile(file)
	conn := w.GetConnections()
	sim, err := simulator.NewSimulator(conn, numAlien, nil)
	if err != nil {
		log.Fatalln(err)
	}
	err = sim.Simulate()
	if !errors.Is(err, simulator.ErrMaxIterationReached) && !errors.Is(err, simulator.ErrNoAlienAlive) {
		log.Fatalln(err)
	}
	fmt.Println("Report: ------------------------")
	logger := log.New(os.Stdout, "", 0)
	w.Report(logger, sim.GetDestroyedCities())
}
