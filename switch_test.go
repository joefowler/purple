package main

import (
	"testing"
)

func TestSixes(t *testing.T) {
	if len(sixesData) != 25 {
		t.Errorf("len(sixesData) = %d, want 25", len(sixesData))
	}
	for i := 0; i < len(sixesData); i++ {
		if len(sixesData[i]) != 6 {
			t.Errorf("len(sixesData[%d]) = %d, want 6", i, len(sixesData[i]))
		}
	}
}

func TestInvert(t *testing.T) {
	d := [][]int{{1, 2, 3, 4}, {3, 1, 2, 4}, {3, 2, 4, 1}, {1, 3, 2, 4}, {1, 3, 4, 2}}
	sd := datamaker(d)
	id := sd.invert()
	want := [][]int{{1, 2, 3, 4}, {2, 3, 1, 4}, {4, 2, 1, 3}, {1, 3, 2, 4}, {1, 4, 2, 3}}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if id[i][j] != want[i][j]-1 {
				t.Errorf("switch.inverse()[%d][%d] = %d, want %d", i, j, id[i][j], want[i][j]-1)
			}
		}
	}
}

func TestSwitch(t *testing.T) {
	p := [][]int{{1, 2}, {1, 2, 3}, {1, 2, 3, 4}}
	_, err := newSwitch(p)
	if err == nil {
		t.Errorf("newSwitch with unequal length lines should raise error but did not")
	}

	s := sixesSwitch
	s.step()
}
