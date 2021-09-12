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
	Graph         map[string][]string
	CityDestroyed map[string]bool
	AlienLocation map[int]string
	RemSimulation int
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewWorld(graph map[string][]string, n int, alienLocation map[int]string) (*worldX, error) {
	if n < 0 || n > MaxAliens {
		return nil, fmt.Errorf("alien count should be in range [0, 10,000], found %d", n)
	}
	if graph == nil || (len(graph) == 0 && n > 0) {
		return nil, fmt.Errorf("world must contain at least one city")
	}

	world := &worldX{
		Graph:         graph,
		RemSimulation: MaxSimulation,
		AlienLocation: make(map[int]string),
	}
	// Number of aliens can be greater than the number of cities.
	// Iterate until all aliens are assigned a city each.
	//
	// Use map's random traversal order to get a city.
	if alienLocation != nil {
		world.AlienLocation = alienLocation
	} else {
		for i := n; i > 0; i-- {
			for k := range graph {
				world.AlienLocation[i] = k
			}
		}
	}

	return world, nil
}

// OneEpoch simulates moves made by all aliens.
func (w *worldX) OneEpoch() ([]string, error) {
	if w.RemSimulation <= 0 {
		return nil, ErrSimulationEnds
	}
	w.RemSimulation--

	if len(w.AlienLocation) == 0 {
		return nil, ErrNoAlienAlive
	}

	// Alien count per city
	cityInvaderCnt := make(map[string]int)
	// Cities that end up with more than one alien after this epoch.
	fightZones := make(map[string][]int)

	for alien, curCity := range w.AlienLocation {
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
		if w.CityDestroyed[city] {
			return nil, fmt.Errorf("city %s already destroyed", city)
		}
		w.CityDestroyed[city] = true
		// Make aliens unavailable for subsequent epoch.
		for _, a := range aliens {
			delete(w.AlienLocation, a)
		}
		msg = append(msg, fmt.Sprintf("%s destroyed by %v", city, aliens))
	}
	return msg, nil
}

// AlienMove return the next city for an alien from a list of
// valid cities. If no valid city is found, i.e. alien is trapped
// it return the current city of the alien.
func (w *worldX) AlienMove(alien int) (string, error) {
	curCity := w.AlienLocation[alien]
	if curCity == "" {
		return "", fmt.Errorf("current location for alien %d not found", alien)
	}
	if w.CityDestroyed[curCity] {
		return "", fmt.Errorf("city %s is already destroyed", curCity)
	}

	validNeighbours := make([]string, 0)
	for _, neighbour := range w.Graph[curCity] {
		if w.CityDestroyed[neighbour] {
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
