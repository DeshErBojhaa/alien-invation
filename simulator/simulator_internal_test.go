package simulator

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"log"
	"os"
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
	s, err := NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	require.Equal(t, len(s.alienLocation), 2)
	require.Equal(t, s.remIteration, MaxIteration)
}

func testInvalidAgentCount(t *testing.T) {
	_, err := NewSimulator(nil, -1, nil)
	require.Contains(t, err.Error(), "alien count should be in range [0, 10,000]")
	_, err = NewSimulator(nil, 111111, nil)
	require.Contains(t, err.Error(), "alien count should be in range [0, 10,000]")
}

func testInvalidCityMap(t *testing.T) {
	_, err := NewSimulator(nil, 2, nil)
	require.Contains(t, err.Error(), "world must contain at least one city")
}

func TestSimulator_oneEpoch(t *testing.T) {
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
	s, err := NewSimulator(mp, 1, nil)
	require.NoError(t, err)

	for i := MaxIteration; i > 0; i-- {
		msg, err := s.oneEpoch()
		require.NoError(t, err)
		require.Equal(t, msg, []string{})
	}
	_, err = s.oneEpoch()
	require.Equal(t, err, ErrMaxIterationReached)
}

func testAllAlienDead(t *testing.T) {
	mp := map[string][]string{
		"a": {},
	}
	s, err := NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	// City a should be destroyed
	msg, err := s.oneEpoch()
	require.NoError(t, err)
	require.Equal(t, msg, []string{"a destroyed by alien [1 2]!"})
	// Error returned
	msg, err = s.oneEpoch()
	require.Equal(t, err, ErrNoAlienAlive)
	require.Equal(t, msg, []string(nil))
}

func TestSimulator_Simulate(t *testing.T) {
	for msg, fn := range map[string]func(t *testing.T){
		"one alien":                           testLoneAlien,
		"two alien toggle cities":             testTwoAliensToggleCities,
		"three alien converge to center city": testThreeAlienConvergeToCenter,
		"four alien converge trapping one":    testFourAlienConvergeTrappingOne,
	} {
		t.Run(msg, fn)
	}
}

func testFourAlienConvergeTrappingOne(t *testing.T) {
	mp := map[string][]string{
		"left_edge1":  {"center_a"},
		"left_edge2":  {"center_a"},
		"center_a":    {"left_edge1", "left_edge2", "center_b"},
		"center_b":    {"center_a", "center_c"},
		"center_c":    {"center_b", "center_d"},
		"center_d":    {"center_c", "center_e"},
		"center_e":    {"center_d", "right_edge1", "right_edge2"},
		"right_edge1": {"center_e"},
		"right_edge2": {"center_e"},
	}
	alienLocation := map[int]string{
		1: "left_edge1",
		2: "left_edge2",
		3: "center_c",
		4: "right_edge1",
		5: "right_edge2",
	}
	expectedMsg := []string{
		"center_e destroyed by alien [4 5]!",
		"center_a destroyed by alien [1 2]!",
	}

	s, err := NewSimulator(mp, 2, alienLocation)
	require.NoError(t, err)

	reader, writer, err := os.Pipe()
	defer func() { _ = reader.Close() }()
	require.NoError(t, err)

	logger := log.New(writer, "", 0)
	s.setLogger(logger)

	err = s.Simulate()
	require.ErrorIs(t, err, ErrMaxIterationReached)

	err = writer.Close()
	require.NoError(t, err)

	receivedMsg := make([]string, 0)

	sc := bufio.NewScanner(reader)
	for sc.Scan() {
		receivedMsg = append(receivedMsg, sc.Text())
	}

	require.ElementsMatch(t, expectedMsg, receivedMsg)
	require.Equal(t, len(s.alienLocation), 1)
	require.Equal(t, len(s.cityDestroyed), 2)
	require.Equal(t, s.cityDestroyed["center_a"], true)
	require.Equal(t, s.cityDestroyed["center_e"], true)
	require.Equal(t, s.remIteration, 0)
}

func testThreeAlienConvergeToCenter(t *testing.T) {
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
	s, err := NewSimulator(mp, 2, alienLocation)
	require.NoError(t, err)
	err = s.Simulate()
	require.ErrorIs(t, err, ErrNoAlienAlive)
	require.Equal(t, len(s.alienLocation), 0)
	require.Equal(t, len(s.cityDestroyed), 1)
	require.Equal(t, s.cityDestroyed["center"], true)
}

func testTwoAliensToggleCities(t *testing.T) {
	mp := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	s, err := NewSimulator(mp, 2, nil)
	require.NoError(t, err)
	err = s.Simulate()
	require.ErrorIs(t, err, ErrMaxIterationReached)
	require.Equal(t, len(s.alienLocation), 2)
}

func testLoneAlien(t *testing.T) {
	mp := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	s, err := NewSimulator(mp, 1, nil)
	require.NoError(t, err)
	err = s.Simulate()
	require.ErrorIs(t, err, ErrMaxIterationReached)
	require.Equal(t, len(s.alienLocation), 1)
}
