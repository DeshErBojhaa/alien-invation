package world

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type World struct {
	connections map[string][]city
	inputOrder  []string
}

type city struct {
	name        string
	orientation string
}

func CreateFromFile(f io.ReadCloser) *World {
	defer func() { _ = f.Close() }()
	world := &World{
		connections: make(map[string][]city),
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		curCity := ""
		for i, part := range parts {
			// i == 0 means the current city name
			if i == 0 {
				world.inputOrder = append(world.inputOrder, part)
				world.connections[part] = make([]city, 0)
				curCity = part
				continue
			}
			orientation, neighbour := func(s string) (string, string) {
				x := strings.Split(part, "=")
				return x[0], x[1]
			}(part)
			world.connections[curCity] = append(world.connections[curCity],
				city{name: neighbour, orientation: orientation})
		}
	}
	return world
}

func (w *World) GetConnections() map[string][]string {
	conn := make(map[string][]string)
	for cityName, neighbours := range w.connections {
		if _, ok := conn[cityName]; !ok {
			conn[cityName] = make([]string, 0)
		}
		for _, neighbour := range neighbours {
			neighbourName := neighbour.name
			conn[cityName] = append(conn[cityName], neighbourName)
			// Connections are bidirectional
			if _, ok := conn[neighbourName]; !ok {
				conn[neighbourName] = make([]string, 0)
			}
			conn[neighbourName] = append(conn[neighbourName], cityName)
		}
	}
	// De-duplicate neighbours
	for cityName, neighbours := range conn {
		uniqueNeighbours := make(map[string]bool)
		for _, name := range neighbours {
			uniqueNeighbours[name] = true
		}
		list := make([]string, 0)
		for name := range uniqueNeighbours {
			list = append(list, name)
		}
		sort.Strings(list)
		conn[cityName] = list
	}
	return conn
}

func (w *World) Report(log *log.Logger, cityDestroyed map[string]bool) {
	for _, curName := range w.inputOrder {
		if cityDestroyed[curName] {
			continue
		}
		neighbours := w.connections[curName]
		var sb strings.Builder
		for _, city := range neighbours {
			if cityDestroyed[city.name] {
				continue
			}
			sb.WriteString(fmt.Sprintf(" %s=%s", city.orientation, city.name))
		}
		if sb.Len() == 0 {
			sb.WriteString(" ---")
		}
		log.Println(curName, sb.String())
	}
}
