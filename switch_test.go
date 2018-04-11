package main

import (
	"testing"
)

func TestInvert(t *testing.T) {
	data := [][]byte{{1, 2, 3, 4}, {3, 1, 2, 4}, {3, 2, 4, 1}, {1, 3, 2, 4}, {1, 3, 4, 2}}
	want := [][]byte{{1, 2, 3, 4}, {2, 3, 1, 4}, {4, 2, 1, 3}, {1, 3, 2, 4}, {1, 4, 2, 3}}
	sd := datamaker(data)
	id := sd.invert()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if id[i][j] != want[i][j]-1 {
				t.Errorf("switch.inverse()[%d][%d] = %d, want %d", i, j, id[i][j], want[i][j]-1)
			}
		}
	}
}

func TestSwitch(t *testing.T) {
	p := [][]byte{{1, 2}, {1, 2, 3}, {1, 2, 3, 4}}
	_, err := newSwitch(p)
	if err == nil {
		t.Errorf("newSwitch with unequal length lines should raise error but did not")
	}

	// Test a few encipher/decipher positions from the sixes switch
	s := sixesSwitch
	var tests = []struct {
		position    int
		permutation []byte
	}{
		{0, []byte{2, 1, 3, 5, 4, 6}},
		{1, []byte{6, 3, 5, 2, 1, 4}},
		{9, []byte{4, 5, 3, 2, 1, 6}},
	}

	for _, test := range tests {
		s.setPosition(test.position)

		for i, want := range test.permutation {
			p := s.decipher(byte(i))
			if p != want-1 {
				t.Errorf("sixesSwitch at pos=%d level %d deciphers to %d, want %d",
					s.position, i, p, want-1)
			}
			c := s.encipher(byte(p))
			if c != byte(i) {
				t.Errorf("sixesSwitch at pos=%d level %d enciphers to %d, want %d",
					s.position, p, c, i)
			}
		}
	}

	// Make sure stepping works
	for i := 0; i < 25; i += 5 {
		s.setPosition(i)
		s.step()
		if s.position != i+1 {
			t.Errorf("sixesSwitch is at position %d, want %d", s.position, i+1)
		}
	}

	// Test a few encipher/decipher positions from the twenties switch #1.
	s = twenties1
	var tests20_1 = []struct {
		position    int
		permutation []byte
	}{
		{0, []byte{6, 19, 14, 1, 10, 4, 2, 7, 13, 9, 8, 16, 3, 18, 15, 11, 5, 12, 20, 17}},
		{1, []byte{4, 5, 16, 17, 14, 1, 20, 15, 3, 8, 18, 11, 12, 13, 10, 19, 2, 6, 9, 7}},
		{9, []byte{12, 8, 17, 9, 3, 20, 4, 10, 14, 5, 7, 18, 2, 16, 13, 6, 1, 19, 15, 11}},
	}

	for _, test := range tests20_1 {
		s.setPosition(test.position)

		for i, want := range test.permutation {
			p := s.decipher(byte(i))
			if p != want-1 {
				t.Errorf("twenties1 at pos=%d level %d deciphers to %d, want %d",
					s.position, i, p, want-1)
			}
			c := s.encipher(p)
			if c != byte(i) {
				t.Errorf("twenties1 at pos=%d level %d enciphers to %d, want %d",
					s.position, p, c, i)
			}
		}
	}
}
