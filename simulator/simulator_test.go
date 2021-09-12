package simulator_test

import (
	"alien/simulator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewWorld(t *testing.T) {
	for msg, fn := range map[string]func(t *testing.T){
		"agent count invalid error": testInvalidAgentCount,
		"invalid city map error":    testInvalidCityMap,
		"valid city map":            testValidCityMap,
	} {
		t.Run(msg, fn)
	}
}

func testValidCityMap(t *testing.T) {
	mp := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	w, err := simulator.NewWorld(mp, 2, nil)
	require.NoError(t, err)
	require.Equal(t, len(w.AlienLocation), 2)
	require.Equal(t, w.RemSimulation, simulator.MaxSimulation)
}

func testInvalidAgentCount(t *testing.T) {
	_, err := simulator.NewWorld(nil, -1, nil)
	require.Contains(t, err.Error(), "alien count should be in range [0, 10,000]")
	_, err = simulator.NewWorld(nil, 111111, nil)
	require.Contains(t, err.Error(), "alien count should be in range [0, 10,000]")
}

func testInvalidCityMap(t *testing.T) {
	_, err := simulator.NewWorld(nil, 2, nil)
	require.Contains(t, err.Error(), "world must contain at least one city")
}
