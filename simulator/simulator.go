package simulator

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const (
	MaxAliens    = 10000
	MaxIteration = 10000
)

var (
	// ErrMaxIterationReached is returned when we have run total MaxIteration number of epoch.
	ErrMaxIterationReached = errors.New("max iteration reached")

	// ErrNoAlienAlive is returned when all alien have died fighting each other. Or we had
	// no aliens to begin with.
	//
	// If all cities are destroyed, all aliens would have died too. So this error is sufficient
	// for that case also.
	ErrNoAlienAlive = errors.New("no alien alive")
)

type Simulator struct {
	Graph         map[string][]string
	CityDestroyed map[string]bool
	AlienLocation map[int]string
	RemIteration  int
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewSimulator(graph map[string][]string, n int, alienLocation map[int]string) (*Simulator, error) {
	if n < 0 || n > MaxAliens {
		return nil, fmt.Errorf("alien count should be in range [0, 10,000], found %d", n)
	}
	if graph == nil || (len(graph) == 0 && n > 0) {
		return nil, fmt.Errorf("world must contain at least one city")
	}

	world := &Simulator{
		Graph:         graph,
		RemIteration:  MaxIteration,
		AlienLocation: make(map[int]string),
		CityDestroyed: make(map[string]bool),
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
func (s *Simulator) OneEpoch() ([]string, error) {
	if s.RemIteration <= 0 {
		return nil, ErrMaxIterationReached
	}
	s.RemIteration--

	if len(s.AlienLocation) == 0 {
		return nil, ErrNoAlienAlive
	}

	// Alien count per city
	cityInvaderCnt := make(map[string]int)
	// Cities that end up with more than one alien after this epoch.
	fightZones := make(map[string][]int)

	for alien, curCity := range s.AlienLocation {
		nxtCity, err := s.AlienMove(alien)
		if err != nil {
			return nil, err
		}
		cityInvaderCnt[curCity]--
		cityInvaderCnt[nxtCity]++
		// If curCity == nxtCity increase the count by 1
		if curCity == nxtCity {
			cityInvaderCnt[nxtCity]++
		}

		fightZones[nxtCity] = append(fightZones[nxtCity], alien)
	}

	msg := make([]string, 0)
	// Destroy cities and aliens where more than one alien collide.
	for city, aliens := range fightZones {
		if len(aliens) == 1 {
			continue
		}
		if s.CityDestroyed[city] {
			return nil, fmt.Errorf("city %s already destroyed", city)
		}
		s.CityDestroyed[city] = true
		// Make aliens unavailable for subsequent epoch.
		for _, a := range aliens {
			delete(s.AlienLocation, a)
		}
		sort.Ints(aliens)
		msg = append(msg, fmt.Sprintf("%s destroyed by %v", city, aliens))
	}
	return msg, nil
}

// AlienMove return the next city for an alien from a list of
// valid cities. If no valid city is found, i.e. alien is trapped
// it return the current city of the alien.
func (s *Simulator) AlienMove(alien int) (string, error) {
	curCity := s.AlienLocation[alien]
	if curCity == "" {
		return "", fmt.Errorf("current location for alien %d not found", alien)
	}
	if s.CityDestroyed[curCity] {
		return "", fmt.Errorf("city %s is already destroyed", curCity)
	}

	validNeighbours := make([]string, 0)
	for _, neighbour := range s.Graph[curCity] {
		if s.CityDestroyed[neighbour] {
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
