package main

import (
	"testing"
)

func TestMachine(t *testing.T) {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if _, err := NewMachine(0, 0, 0, 0, 0, 3, alphabet); err == nil {
		t.Errorf("Fast switch = 0 should raise error")
	}
	if _, err := NewMachine(0, 0, 0, 0, 1, 0, alphabet); err == nil {
		t.Errorf("Middle switch = 0 should raise error")
	}
	if _, err := NewMachine(0, 0, 0, 0, 2, 2, alphabet); err == nil {
		t.Errorf("Fast switch = middle switch should raise error")
	}

	// Test bad alphabets
	if _, err := NewMachine(0, 0, 0, 0, 1, 2, "AB"); err == nil {
		t.Errorf("Short alphabet should raise error")
	}
	if _, err := NewMachine(0, 0, 0, 0, 1, 2, "AABCDEFGHIJKLMNOPQRSTUVWXYZ"); err == nil {
		t.Errorf("Long alphabet should raise error")
	}
	if _, err := NewMachine(0, 0, 0, 0, 1, 2, "AABCDEFGHIJKLMNOPQRSTUVWXY"); err == nil {
		t.Errorf("Alphabet with repeated characters should raise error")
	}
	if _, err := NewMachine(0, 0, 0, 0, 1, 2, "ABCDEFGHIJKLMNOPQRS1234567"); err == nil {
		t.Errorf("Alphabet with non-letter characters should raise error")
	}
}

// Test that the switches move according to Fig 10 in FSW paper.
func TestSwitchMotion(t *testing.T) {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	m, err := NewMachine(20, 0, 24, 4, 1, 2, alphabet)
	if err != nil {
		t.Fatalf("Could not complete NewMachine: %s", err.Error())
	}
	var tests = [][]int{
		{20, 0, 24, 4},
		{21, 1, 24, 4},
		{22, 2, 24, 4},
		{23, 3, 24, 4},
		{24, 3, 24, 5},
		{0, 3, 0, 5},
		{1, 4, 0, 5},
		{2, 5, 0, 5},
	}
	for i, test := range tests {
		if m.sixes.position != test[0] ||
			m.fast.position != test[1] ||
			m.middle.position != test[2] ||
			m.slow.position != test[3] {
			t.Errorf("TestSwitchMotion fails after %d steps: [%d,%d,%d,%d], want %v", i,
				m.sixes.position, m.fast.position, m.middle.position, m.slow.position, test)
		}
		m.step()
	}
}
