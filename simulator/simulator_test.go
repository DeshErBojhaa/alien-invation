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
	w, err := simulator.NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	require.Equal(t, len(w.AlienLocation), 2)
	require.Equal(t, w.RemIteration, simulator.MaxIteration)
}

func testInvalidAgentCount(t *testing.T) {
	_, err := simulator.NewSimulator(nil, -1, nil)
	require.Contains(t, err.Error(), "alien count should be in range [0, 10,000]")
	_, err = simulator.NewSimulator(nil, 111111, nil)
	require.Contains(t, err.Error(), "alien count should be in range [0, 10,000]")
}

func testInvalidCityMap(t *testing.T) {
	_, err := simulator.NewSimulator(nil, 2, nil)
	require.Contains(t, err.Error(), "world must contain at least one city")
}

func TestWorldX_OneEpoch(t *testing.T) {
	for msg, fn := range map[string]func(t *testing.T){
		"all alien died":  testAllAlienDead,
		"simulation ends": testSimulationEnds,
	} {
		t.Run(msg, fn)
	}
}

func testSimulationEnds(t *testing.T) {
	mp := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	w, err := simulator.NewSimulator(mp, 1, nil)
	require.NoError(t, err)

	for i := simulator.MaxIteration; i > 0; i-- {
		msg, err := w.OneEpoch()
		require.NoError(t, err)
		require.Equal(t, msg, []string{})
	}
	_, err = w.OneEpoch()
	require.Equal(t, err, simulator.ErrMaxIterationReached)
}

func testAllAlienDead(t *testing.T) {
	mp := map[string][]string{
		"a": {},
	}
	w, err := simulator.NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	// City a should be destroyed
	msg, err := w.OneEpoch()
	require.NoError(t, err)
	require.Equal(t, msg, []string{"a destroyed by [1 2]"})
	// Error returned
	msg, err = w.OneEpoch()
	require.Equal(t, err, simulator.ErrNoAlienAlive)
	require.Equal(t, msg, []string(nil))
}
