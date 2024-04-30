package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"github.com/m-colson/dystopia/shared/graph"
)

func GenerateMap(w io.Writer, nodes int, edges int, avgcost float64) {
	median := avgcost / math.Ln2

	for i := 0; i < edges; i++ {
		from := rand.IntN(nodes)
		to := rand.IntN(nodes)
		cost := math.Ceil(-median * math.Log2(1-rand.Float64()))

		fmt.Fprintf(w, "%d->%d:%d\n", from, to, int(cost))
	}
}

func ParseMap(file string) graph.Graph {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	links := make([]graph.Link, 0)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if len(line) == 0 {
			break
		}

		links = append(links, ParseLink(line))
	}

	return graph.New(links...)
}

func ParseLink(line []byte) graph.Link {
	off := 0

	fromStr := strings.Builder{}
	for _, c := range line {
		if c < '0' || c > '9' {
			break
		}
		fromStr.WriteByte(c)
		off += 1
	}
	from, err := graph.ParseID(fromStr.String())
	if err != nil {
		panic(err)
	}

	if line[off] != '-' {
		panic("expected '-'")
	}
	off += 1
	if line[off] != '>' {
		panic("expected '>'")
	}
	off += 1

	toStr := strings.Builder{}
	for _, c := range line[off:] {
		if c < '0' || c > '9' {
			break
		}
		toStr.WriteByte(c)
		off += 1

	}
	to, err := graph.ParseID(toStr.String())
	if err != nil {
		panic(err)
	}

	if line[off] != ':' {
		panic("expected ':'")
	}
	off += 1

	costStr := strings.Builder{}
	for _, c := range line[off:] {
		if c < '0' || c > '9' {
			break
		}
		costStr.WriteByte(c)
		off += 1
	}
	cost, err := strconv.ParseUint(costStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	if off != len(line) {
		panic(fmt.Errorf("unexpected characters at end of line %q", line[off:]))
	}

	return graph.Link{From: from, To: to, Cost: graph.Cost(cost)}
}
