package simulator_test

import (
	"alien/simulator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSimulator(t *testing.T) {
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
	s, err := simulator.NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	require.Equal(t, len(s.AlienLocation), 2)
	require.Equal(t, s.RemIteration, simulator.MaxIteration)
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

func TestSimulator_OneEpoch(t *testing.T) {
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
	s, err := simulator.NewSimulator(mp, 1, nil)
	require.NoError(t, err)

	for i := simulator.MaxIteration; i > 0; i-- {
		msg, err := s.OneEpoch()
		require.NoError(t, err)
		require.Equal(t, msg, []string{})
	}
	_, err = s.OneEpoch()
	require.Equal(t, err, simulator.ErrMaxIterationReached)
}

func testAllAlienDead(t *testing.T) {
	mp := map[string][]string{
		"a": {},
	}
	s, err := simulator.NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	// City a should be destroyed
	msg, err := s.OneEpoch()
	require.NoError(t, err)
	require.Equal(t, msg, []string{"a destroyed by alien [1 2]!"})
	// Error returned
	msg, err = s.OneEpoch()
	require.Equal(t, err, simulator.ErrNoAlienAlive)
	require.Equal(t, msg, []string(nil))
}

func TestSimulator_Simulate(t *testing.T) {
	for msg, fn := range map[string]func(t *testing.T){
		"one alien":                           testLoneAlien,
		"two alien toggle cities":             testTwoAliensToggleCities,
		"three alien converge to center city": testThreeAlienConvergeToCenter,
	} {
		t.Run(msg, fn)
	}
}

func testThreeAlienConvergeToCenter(t *testing.T) {
	t.Skip()
	mp := map[string][]string{
		"center": {"a", "b", "c"},
		"a":      {"center"},
		"b":      {"center"},
		"c":      {"center"},
	}
	alienLocation := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}
	s, err := simulator.NewSimulator(mp, 2, alienLocation)
	require.NoError(t, err)
	err = s.Simulate()
	require.ErrorIs(t, err, simulator.ErrNoAlienAlive)
	require.Equal(t, len(s.AlienLocation), 0)
	require.Equal(t, len(s.CityDestroyed), 1)
	require.Equal(t, s.CityDestroyed["center"], true)
}

func testTwoAliensToggleCities(t *testing.T) {
	t.Skip()
	mp := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	s, err := simulator.NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	err = s.Simulate()
	require.ErrorIs(t, err, simulator.ErrMaxIterationReached)
	require.Equal(t, len(s.AlienLocation), 2)
}

func testLoneAlien(t *testing.T) {
	t.Skip()
	mp := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	s, err := simulator.NewSimulator(mp, 1, nil)
	require.NoError(t, err)
	err = s.Simulate()
	require.ErrorIs(t, err, simulator.ErrMaxIterationReached)
	require.Equal(t, len(s.AlienLocation), 1)
}
