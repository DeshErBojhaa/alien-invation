package world_test

import (
	"alien/world"
	"bufio"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestWorld(t *testing.T) {
	file := `
		A north=B south=C east=D west=E
		B south=A west=X east=Y
		X south=E
		E east=A south=M
		C west=M
	`
	readCloser := ioutil.NopCloser(strings.NewReader(file))
	w := world.CreateFromFile(readCloser)
	gotConn := w.GetConnections()
	expectedConn := map[string][]string{
		"A": {"B", "C", "D", "E"},
		"B": {"A", "X", "Y"},
		"C": {"A", "M"},
		"D": {"A"},
		"E": {"A", "M", "X"},
		"M": {"C", "E"},
		"X": {"B", "E"},
		"Y": {"B"},
	}
	require.Equal(t, expectedConn, gotConn)

	reader, writer, err := os.Pipe()
	defer func() { _ = reader.Close() }()
	require.NoError(t, err)

	logger := log.New(writer, "", 0)
	w.Report(logger, map[string]bool{"A": true, "X": true})

	err = writer.Close()
	require.NoError(t, err)

	receivedMsg := make([]string, 0)

	sc := bufio.NewScanner(reader)
	for sc.Scan() {
		receivedMsg = append(receivedMsg, sc.Text())
	}
	expectedMsg := []string{
		"B  east=Y",
		"E  south=M",
		"C  west=M",
	}
	require.Equal(t, expectedMsg, receivedMsg)
}
