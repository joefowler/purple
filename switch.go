package main

import "fmt"

// switchData represents a single encipher/deciphering series of permutations
type switchData [][]int

// invert makes a new switchData table, inverting each permutation
func (s switchData) invert() switchData {
	npositions := len(s)
	nperms := len(s[1])
	r := make(switchData, npositions)
	for i := 0; i < npositions; i++ {
		r[i] = make([]int, nperms)
		for j := 0; j < nperms; j++ {
			r[i][s[i][j]] = j
		}
	}
	return r
}

// Switch represents the wiring and state of a single Purple switch.
type Switch struct {
	nlevels        int // How many "levels" the switch permutes
	npositions     int // How many permutations the switch has
	position       int // The current position
	decipherWiring switchData
	encipherWiring switchData
}

func (s *Switch) step() int {
	s.position = (s.position + 1) % s.npositions
	return s.position
}

var sixesData, twenties1, twenties2, twenties3 switchData
var sixesSwitch *Switch

func init() {
	p := [][]int{
		{2, 1, 3, 5, 4, 6},
		{6, 3, 5, 2, 1, 4},
		{1, 5, 4, 6, 2, 3},
		{4, 3, 2, 1, 6, 5},
		{3, 6, 1, 4, 5, 2},
		{2, 1, 6, 5, 3, 4},
		{6, 5, 4, 2, 1, 3},
		{3, 6, 1, 4, 5, 2},
		{5, 4, 2, 6, 3, 1},
		{4, 5, 3, 2, 1, 6},
		{2, 1, 4, 5, 6, 3},
		{5, 4, 6, 3, 2, 1},
		{3, 1, 2, 6, 4, 5},
		{4, 2, 5, 1, 3, 6},
		{1, 6, 2, 3, 5, 4},
		{5, 4, 3, 6, 1, 2},
		{6, 2, 5, 3, 4, 1},
		{2, 3, 4, 1, 5, 6},
		{1, 2, 3, 5, 6, 4},
		{3, 1, 6, 4, 2, 5},
		{6, 5, 1, 2, 4, 3},
		{1, 3, 6, 4, 2, 5},
		{6, 4, 5, 1, 3, 2},
		{4, 6, 1, 2, 5, 3},
		{5, 2, 4, 3, 6, 1},
	}
	sixesData = datamaker(p)
	sixesSwitch, _ = newSwitch(p)
}

func datamaker(p [][]int) switchData {
	var data switchData
	for i := range p {
		data = append(data, make([]int, len(p[i])))
		for j := range data[i] {
			data[i][j] = p[i][j] - 1
		}
	}
	return data
}

func newSwitch(p [][]int) (*Switch, error) {
	sdata := datamaker(p)
	s := new(Switch)
	s.npositions = len(sdata)
	s.nlevels = len(sdata[0])
	for i, line := range sdata {
		if len(line) != s.nlevels {
			return nil, fmt.Errorf("newSwitch: all slices must be the same length, got %d in position %d, want %d", len(line), i, s.nlevels)
		}
	}
	s.decipherWiring = sdata
	s.encipherWiring = sdata.invert()
	return s, nil
}
