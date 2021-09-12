package simulator

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	graph         map[string][]string
	cityDestroyed map[string]bool
	alienLocation map[int]string
	remIteration  int
	log           *log.Logger
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
		graph:         graph,
		remIteration:  MaxIteration,
		alienLocation: make(map[int]string),
		cityDestroyed: make(map[string]bool),
		log:           log.New(os.Stdout, "", 0),
	}
	// Number of aliens can be greater than the number of cities.
	// Iterate until all aliens are assigned a city each.
	//
	// Use map's random traversal order to get a city.
	if alienLocation != nil {
		world.alienLocation = alienLocation
	} else {
		alien, i := n, 1
		for alien > 0 {
			for k := range graph {
				world.alienLocation[i] = k
				i++
				alien--
				if alien <= 0 {
					break
				}
			}
		}
	}

	return world, nil
}

// Simulate runs until one of the condition met
//   1. All aliens are destroyed
//   2. All cities are destroyed, which is also covered by case 1
//   3. Max iteration reached
//   4. Fatal error
func (s *Simulator) Simulate() error {
	for {
		msg, err := s.oneEpoch()
		if errors.Is(err, ErrMaxIterationReached) || errors.Is(err, ErrNoAlienAlive) {
			return fmt.Errorf("simulation ends: %w", err)
		}
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}
		if msg != nil && len(msg) > 0 {
			for _, m := range msg {
				s.log.Println(m)
			}
		}
	}
}

// oneEpoch simulates moves made by all aliens.
func (s *Simulator) oneEpoch() ([]string, error) {
	if s.remIteration <= 0 {
		return nil, ErrMaxIterationReached
	}
	s.remIteration--

	if len(s.alienLocation) == 0 {
		return nil, ErrNoAlienAlive
	}

	// Alien count per city
	cityInvaderCnt := make(map[string]int)
	// Cities that end up with more than one alien after this epoch.
	fightZones := make(map[string][]int)

	for alien, curCity := range s.alienLocation {
		nxtCity, err := s.alienMove(alien)
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
		if len(aliens) < 2 {
			continue
		}
		if s.cityDestroyed[city] {
			return nil, fmt.Errorf("city %s already destroyed", city)
		}
		s.cityDestroyed[city] = true
		// Make aliens unavailable for subsequent epoch.
		for _, a := range aliens {
			delete(s.alienLocation, a)
		}
		sort.Ints(aliens)
		msg = append(msg, fmt.Sprintf("%s destroyed by alien %v!", city, aliens))
	}
	return msg, nil
}

// alienMove return the next city for an alien from a list of
// valid cities. If no valid city is found, i.e. alien is trapped
// it return the current city of the alien.
func (s *Simulator) alienMove(alien int) (string, error) {
	curCity := s.alienLocation[alien]
	if curCity == "" {
		return "", fmt.Errorf("current location for alien %d not found", alien)
	}
	if s.cityDestroyed[curCity] {
		return "", fmt.Errorf("city %s is already destroyed", curCity)
	}

	validNeighbours := make([]string, 0)
	for _, neighbour := range s.graph[curCity] {
		if s.cityDestroyed[neighbour] {
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

func (s *Simulator) setLogger(log *log.Logger) {
	s.log = log
}

func randomSelect(cities []string) string {
	return cities[rand.Intn(len(cities))]
}
