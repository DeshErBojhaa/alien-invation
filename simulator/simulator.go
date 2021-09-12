package simulator

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	MaxAliens     = 10000
	MaxSimulation = 10000
)

var (
	// ErrSimulationEnds is returned when we have run total MaxSimulation number of epoch.
	ErrSimulationEnds = errors.New("simulation ended")

	// ErrNoAlienAlive is returned when all alien have died fighting each other. Or we had
	// no aliens to begin with.
	ErrNoAlienAlive = errors.New("no alien alive")
)

type worldX struct {
	graph         map[string][]string
	cityDestroyed map[string]bool
	alienLocation map[int]string
	remSimulation int
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewWorld(graph map[string][]string, n int, alienLocation map[int]string) (*worldX, error) {
	if n < 0 || n > MaxAliens {
		return nil, fmt.Errorf("alien count sould be in range [0, 10,000], found %d", n)
	}
	world := &worldX{
		graph: graph,
		remSimulation: MaxSimulation,
	}
	// Number of aliens can be greater than the number of cities.
	// Iterate until all aliens are assigned a city each.
	//
	// Use map's random traversal order to get a city.
	if alienLocation == nil {
		for i := n; i > 0; i-- {
			for k := range graph {
				world.alienLocation[i] = k
			}
		}
	} else {
		world.alienLocation = alienLocation
	}

	return world, nil
}

// OneEpoch simulates moves made by all aliens.
func (w *worldX) OneEpoch() ([]string, error) {
	if w.remSimulation <= 0 {
		return nil, ErrSimulationEnds
	}
	w.remSimulation--

	if len(w.alienLocation) == 0 {
		return nil, ErrNoAlienAlive
	}

	// Alien count per city
	cityInvaderCnt := make(map[string]int)
	// Cities that end up with more than one alien after this epoch.
	fightZones := make(map[string][]int)

	for alien, curCity := range w.alienLocation {
		nxtCity, err := w.AlienMove(alien)
		if err != nil {
			return nil, err
		}
		cityInvaderCnt[curCity]--
		cityInvaderCnt[nxtCity]++
		if cityInvaderCnt[nxtCity] > 1 {
			fightZones[nxtCity] = append(fightZones[nxtCity], alien)
		}
	}

	msg := make([]string, 0)
	// Destroy cities and aliens where more than one alien collide.
	for city, aliens := range fightZones {
		if w.cityDestroyed[city] {
			return nil, fmt.Errorf("city %s already destroyed", city)
		}
		w.cityDestroyed[city] = true
		// Make aliens unavailable for subsequent epoch.
		for _, a := range aliens {
			delete(w.alienLocation, a)
		}
		msg = append(msg, fmt.Sprintf("%s destroyed by %v", city, aliens))
	}
	return msg, nil
}

// AlienMove return the next city for an alien from a list of
// valid cities. If no valid city is found, i.e. alien is trapped
// it return the current city of the alien.
func (w *worldX) AlienMove(alien int) (string, error) {
	curCity := w.alienLocation[alien]
	if curCity == "" {
		return "", fmt.Errorf("current location for alien %d not found", alien)
	}
	if w.cityDestroyed[curCity] {
		return "", fmt.Errorf("city %s is already destroyed", curCity)
	}

	validNeighbours := make([]string, 0)
	for _, neighbour := range w.graph[curCity] {
		if w.cityDestroyed[neighbour] {
			continue
		}
		validNeighbours = append(validNeighbours, neighbour)
	}
	// Alien trapped!
	if len(validNeighbours) == 0 {
		return curCity, nil
	}
	return randomSelect(validNeighbours), nil
}

func randomSelect(cities []string) string {
	return cities[rand.Intn(len(cities))]
}
