package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Machine represents the settings and state of a PURPLE machine.
type Machine struct {
	sixes        *Switch
	twenties     [3]*Switch
	fast         *Switch
	middle       *Switch
	slow         *Switch
	alphabet     string
	plugboardIn  [26]byte
	plugboardOut [26]byte
}

// NewMachineFromKey creates a pointer to a new instance of a PURPLE machine, configured according to arguments.
//
//     switches: must be a string of the form 'a-b,c,d-ef' where
//     a - starting position of the sixes switch (1-25)
//     b - starting position of the twenties switch #1 (1-25)
//     c - starting position of the twenties switch #2 (1-25)
//     d - starting position of the twenties switch #3 (1-25)
//     e - which switch is the fast switch (1-3)
//     f - which switch is the middle switch (1-3)
//
// Example: '9-1,24,6-23'
func NewMachineFromKey(key, alphabet string) (*Machine, error) {
	parts := strings.Split(key, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Key was not of the form 9-1,24,6-23")
	}
	sixes, twenties, permutation := parts[0], parts[1], parts[2]

	tparts := strings.Split(twenties, ",")
	if len(tparts) != 3 {
		return nil, fmt.Errorf("Key was not of the form 9-1,24,6-23")
	}

	sixpos, err := strconv.Atoi(sixes)
	if err != nil {
		return nil, err
	}
	tw1pos, err := strconv.Atoi(tparts[0])
	if err != nil {
		return nil, err
	}
	tw2pos, err := strconv.Atoi(tparts[1])
	if err != nil {
		return nil, err
	}
	tw3pos, err := strconv.Atoi(tparts[2])
	if err != nil {
		return nil, err
	}
	permnum, err := strconv.Atoi(permutation)
	if err != nil {
		return nil, err
	}
	fast := permnum / 10
	middle := permnum % 10
	return NewMachine(sixpos, tw1pos, tw2pos, tw3pos, fast, middle, alphabet)
}

// NewMachine creates a pointer to a new instance of a PURPLE machine, configured according to arguments.
func NewMachine(sixpos, tw1pos, tw2pos, tw3pos, fast, middle int, alphabet string) (*Machine, error) {
	if sixpos < 1 || sixpos > 25 ||
		tw1pos < 1 || tw1pos > 25 ||
		tw2pos < 1 || tw2pos > 25 ||
		tw3pos < 1 || tw3pos > 25 {
		return nil, fmt.Errorf("switch positions [%d, %d, %d, %d] should all be in range 1-25",
			sixpos, tw1pos, tw2pos, tw3pos)
	}
	m := new(Machine)
	m.sixes = sixesSwitch
	m.sixes.setPosition(sixpos - 1)
	if fast == middle {
		return nil, fmt.Errorf("fast and middle (%d, %d) must be different", fast, middle)
	}
	if fast < 1 || fast > 3 {
		return nil, fmt.Errorf("fast = %d, must be in [1,3]", fast)
	}
	if middle < 1 || middle > 3 {
		return nil, fmt.Errorf("middle = %d, must be in [1,3]", middle)
	}
	m.twenties[0] = twenties1
	m.twenties[1] = twenties2
	m.twenties[2] = twenties3
	m.fast = m.twenties[fast-1]
	m.middle = m.twenties[middle-1]
	for i := 1; i <= 3; i++ {
		if i != fast && i != middle {
			m.slow = m.twenties[i-1]
			break
		}
	}
	m.twenties[0].setPosition(tw1pos - 1)
	m.twenties[1].setPosition(tw2pos - 1)
	m.twenties[2].setPosition(tw3pos - 1)

	// Validate the alphabet
	if len(alphabet) != 26 {
		return nil, fmt.Errorf("alphabet length=%d, should be 26", len(alphabet))
	}
	m.alphabet = strings.ToUpper(alphabet)
	count := make(map[rune]int)
	for _, c := range m.alphabet {
		count[c]++
	}
	if len(count) != 26 {
		return nil, fmt.Errorf("alphabet has %d unique values, should be 26", len(count))
	}
	for _, c := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if count[c] != 1 {
			return nil, fmt.Errorf("alphabet contains %d '%c', should be 1", count[c], c)
		}
	}

	for i, c := range []byte(m.alphabet) {
		m.plugboardIn[c-'A'] = byte(i)
		m.plugboardOut[i] = c - 'A'
	}

	return m, nil
}

// step advances the sixes switch and exactly one twenties switch, according
// the the rules of the machine.
func (m *Machine) step() {
	if m.middle.position == 24 && m.sixes.position == 23 {
		m.slow.step()
	} else if m.sixes.position == 24 {
		m.middle.step()
	} else {
		m.fast.step()
	}
	m.sixes.step()
}

// decipher converts c from cipher to plain text but does NOT step the machine
func (m *Machine) decipher(c byte) (p byte) {
	n := m.plugboardIn[c]
	if n < 6 {
		p = m.sixes.decipher(n)
	} else {
		p = 6 + m.twenties[0].decipher(m.twenties[1].decipher(m.twenties[2].decipher(n-6)))
	}
	return m.plugboardOut[p]
}

// encipher converts p to ciphertext but does NOT step the machine
func (m *Machine) encipher(p byte) (c byte) {
	n := m.plugboardIn[p]
	if n < 6 {
		c = m.sixes.encipher(n)
	} else {
		c = 6 + m.twenties[2].encipher(m.twenties[1].encipher(m.twenties[0].encipher(n-6)))
	}
	return m.plugboardOut[c]
}

func (m *Machine) decipherMessage(cipher string) string {
	result := make([]byte, len(cipher))
	for i, c := range []byte(cipher) {
		if c >= 'A' && c <= 'Z' {
			result[i] = m.decipher(c-'A') + 'A'
		} else if c >= 'a' && c <= 'z' {
			result[i] = m.decipher(c-'a') + 'a'
		} else {
			result[i] = c
		}
		if c != ' ' && c != '\n' {
			m.step()
		}
	}
	return string(result)
}

func (m *Machine) encipherMessage(plain string) string {
	result := make([]byte, len(plain))
	for i, p := range []byte(plain) {
		if p >= 'A' && p <= 'Z' {
			result[i] = m.encipher(p-'A') + 'A'
		} else if p >= 'a' && p <= 'z' {
			result[i] = m.encipher(p-'a') + 'a'
		} else {
			result[i] = p
		}
		if p != ' ' && p != '\n' {
			m.step()
		}
	}
	return string(result)
}
