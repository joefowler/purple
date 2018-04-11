package main

import (
	"fmt"
	"strings"
)

// Machine represents the settings and state of a PURPLE machine.
type Machine struct {
	sixes     *Switch
	twenties  [3]*Switch
	fast      *Switch
	middle    *Switch
	slow      *Switch
	alphabet  string
	plugboard map[byte]byte
}

// NewMachine creates a pointer to a new instance of a PURPLE machine, configured according to arguments.
func NewMachine(sixpos, tw1pos, tw2pos, tw3pos, fast, middle int, alphabet string) (*Machine, error) {
	m := new(Machine)
	m.sixes = sixesSwitch
	m.sixes.setPosition(sixpos)
	if fast == middle {
		return nil, fmt.Errorf("fast and middle (%d, %d) must be different", fast, middle)
	}
	if fast < 1 || fast > 3 {
		return nil, fmt.Errorf("fast = %d, must be in [1,3]", fast)
	}
	if middle < 1 || middle > 3 {
		return nil, fmt.Errorf("middle = %d, must be in [1,3]", middle)
	}
	twen := [3](*Switch){twenties1, twenties2, twenties3}
	m.twenties[0] = twen[fast-1]
	m.twenties[1] = twen[middle-1]
	for i := 1; i <= 3; i++ {
		if i != fast && i != middle {
			m.twenties[2] = twen[i-1]
			break
		}
	}
	m.twenties[0].setPosition(tw1pos)
	m.twenties[1].setPosition(tw2pos)
	m.twenties[2].setPosition(tw3pos)
	m.fast = m.twenties[0]
	m.middle = m.twenties[1]
	m.slow = m.twenties[2]

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

	m.plugboard = make(map[byte]byte)
	for i, c := range []byte(m.alphabet) {
		m.plugboard[c-'A'] = byte(i)
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
func (m *Machine) decipher(c byte) byte {
	n := m.plugboard[c]
	if n < 6 {
		return m.sixes.decipher(n)
	}
	return 6 + m.fast.decipher(m.middle.decipher(m.slow.decipher(n-6)))
}

// encipher converts p to ciphertext but does NOT step the machine
func (m *Machine) encipher(p byte) byte {
	n := m.plugboard[p]
	if n < 6 {
		return m.sixes.encipher(n)
	}
	return 6 + m.fast.encipher(m.middle.encipher(m.slow.encipher(n-6)))
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
		m.step()
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
		m.step()
	}
	return string(result)
}
