package main

import (
	"fmt"
	"strings"
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

func TestMachineFromKey(t *testing.T) {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	validKeys := []string{
		"9-1,2,3-23",
		"1-1,1,1-13",
		"25-25,25,25-31",
		"5-20,7,18-21",
	}
	for _, key := range validKeys {
		if m, err := NewMachineFromKey(key, alphabet); m == nil || err != nil {
			t.Errorf("NewMachineFromKey(\"%s\", \"%s\") failed", key, alphabet)
		}
	}

	badKeys := []string{
		"0-1,2,3-13",
		"26-1,2,3-13",
		"1-1,0,3-13",
		"1-1,2,26-13",
		"1-1,2,26-03",
		"1-1,2,26-00",
		"1-1,2,26-14",
		"bad string",
		"1-2-1,2,26-14",
		"a-9,2,20-13",
		"1-a,2,20-13",
		"1-9,a,20-13",
		"1-9,2,a-13",
		"1-9,2,20-a3",
		"1-9,2,20-1a",
		"1-9,2,20-123",
		"1-9,2,20-123-4",
		"1-9,2,20,5-123",
	}
	for _, key := range badKeys {
		if _, err := NewMachineFromKey(key, alphabet); err == nil {
			t.Errorf("NewMachineFromKey(\"%s\", \"%s\") should have failed", key, alphabet)
		}
	}
}

// Test that the switches move according to Fig 10 in FSW paper.
func TestSwitchMotion(t *testing.T) {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	m, err := NewMachine(21, 1, 25, 5, 1, 2, alphabet)
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

func Test14PartMessage(t *testing.T) {
	// This is part 1 of the famous 14 part message. The cipherlines and plainlines are
	// laid on top of each other, every other line, with blank lines inserted for
	// clarity. Garbles are indicated by '-' characters.
	//
	part1 := `
    ZTXODNWKCCMAVNZXYWEETUQTCIMNVEUVIWBLUAXRRTLVA
    FOVTATAKIDASINIMUIMINOMOXIWOIRUBESIFYXXFCKZZR

    RGNTPCNOIUPJLCIVRTPJKAUHVMUDTHKTXYZELQTVWGBUHFAWSH
    DXOOVBTNFYXFAEMEMORANDUMFIOFOVOOMOJIBAKARIFYXRAICC

    ULBFBHEXMYHFLOWD-KWHKKNXEBVPYHHGHEKXIOHQHUHWIKYJYH
    YLFCBBCFCTHEGOVE-NMENTOFJAPANLFLPROMPTEDBYAGENUINE

    PPFEALNNAKIBOOZNFRLQCFLJTTSSDDOIOCVT-ZCKQTSHXTIJCN
    DESIRETOCOMETOANAMICABLEUNDERSTANDIN-WITHTHEGOVERN

    WXOKUFNQR-TAOIHWTATWVHOTGCGAKVANKZANMUIN
    MENTOFTHE-NITEDSTATESINORDERTHATTHETWOCO

    YOYJFSRDKKSEQBWKIOORJAUWKXQGUWPDUDZNDRMDHVHYPNIZXB
    UNTRIESBYTHEIRJOINTEFFORTSMAYSECURETHEPEACEOFTHEPA

    GICXRMAWMFTIUDBXIENLONOQVQKYCOTVSHVNZZQPDLMXVNRUUN
    CIFICAREAANDTHEREBYCONTRIBUTETOWARDTHEREALIZATIONO

    QFTCDFECZDFGMXEHHWYONHYNJDOVJUNCSUVKKEIWOLKRBUUSOZ
    FWORLDPEACELFLHASCONTINUEDNEGOTIATIONSWITHTHEUTMOS

    UIGNISMWUOSBOBLJXERZJEQYQMTFTXBJNCMJKVRKOTSOPBOYMK
    TSINCERITYSINCEAPRILLASTWITHTHEGOVERNMENTOFTHEUNIT

    IRETINCPSQJAWVHUFKRMAMXNZUIFNOPUEMHGLOEJHZOOKHHEED
    EDSTATESREGARDINGTHEADJUSTMENTANDADVANCEMENTOFJAPA

    NIHXFXFXGPDZBSKAZABYEKYEPNIYSHVKFRFPVCJTPTOYCNEIQB
    NESEVVFAMERICANRELATIONSANDTHESTABILIZATIONOFTHEPA

    FEXMERMIZLGDRXZORLZFSQYPZFATZCHUGRNHWDDTAIHYOOCOOD
    CIFICAREACFCCCFTHEJAPANESEQOVERNMENXHASTHEHONORTOS

    UZYIWJROOJUMUIHRBEJFONAXGNCKAOARDIHCDZKIXPR--DIMUW
    TATEFRANKLYITSVIEWSCONCERNINGTHECLAIMSTHEAM--VCANG

    OMHLTJSOUXPFKGEPWJOMTUVKMWRKTACUPIGAFEDFVRKXFXLFGU
    OVERNMENTHASUERSISTENTLYMAINTAINEDASWELLASTHEMEASU

    RDETJIYOLKBHZKXOJDDOVRHMMUQBFOWRODMRMUWNAYKYPISDLH
    RESTHEUNITEDSTATESANDGREATBRITAINHAVETAKENTOWARDJA

    ECKINLJORKWNWXADAJOLONOEVMUQDFIDSPEBBPWROFBOPAZJEU
    PANDURINGTHKSEEIGHTMONTHSCYCCCFLFCDDCFCITISTHEIMMU

    USBHGIORCSUUQKIIEHPCTJRWSOGLETZLOUKKEOJOSMKJBWUCDD
    TABLXPOLWCYOFTHEJAPANESEGOVERNMENTTOINSURETHESTABI

    CPYUUWCSSKWWVLIUPKYXGKQOKAZTEZFHGVPJFEWEUBKLIZLWKK
    LITYOFEASTASIAANDTOPROMOTEWORLZPEACELFLANDTHEREBYT

    OBXLEPQPDATWUSUUPKYRHNWDZXXGTWDDNSHDCBCJXAOOEEPUBP
    OEIABLEALLNATIONSTOFINDEACHITSPROPERPLACEINTHEWORL

    WFRBQSFXSEZJJYAANMG-WLYMGWAQDGIVNOHKOUTIXYFOKNGGBF
    DCFCCCFEVERSINCETHE-HINAAFFAIRBROKEOUTOWINGTOTHEFA

    GANPWTUYLBEFFKUFLEXOIUUANVMMJEQUSFHFDOHQLAKWTBYYYL
    ILUREONTHEPARTOFCHINATOCOMPREHENLJAPANVCFSTRUEYNTE

    NTLYTSXCGKCEEWQRYAVGRKXIANPXNOFVXGKJFAVKLTHOCXCIVK
    NTIONSLFLTWEJAPANESEGOVERNMENTHASSTRIVENFORTHEREST

    OLXTJTUNCLQCICRUIIWQDDMOTPRVTJKKSKFHXFKMDIKIZWROGZ
    ORATIONOFPEACEANDIMHASCONSISTENTLYEXERTEDITSBESTEF

    JYMTMNOVMFJ-OKTEIVMYANOHNNYPDLEXCFRRNEBLMNYEBGNHCZ
    FORTSTOPREV-NTTHEEXTENTIONOFWARVVFLIKEVISTURBANCES

    ZCFNWGGRHRIUUTTILKLODUYZKQOZMMNHASXHLPVTNGHQDAJIUG
    CFCNSIASALSOTOTHATENDTNATINSEPTEMBERLASTYEARJAPANC

    OOSZ-----ZRTGWFBLKI--------YBDABJ-----WYOEANV---OM
    ONCL-----HETRIPAITI--------THGERM-----DYTALYC---OV
    `

	pt1lines := strings.Split(part1, "\n")
	var cipherlines []string
	var plainlines []string
	for i := 1; i < len(pt1lines); i += 3 {
		cipherlines = append(cipherlines, pt1lines[i])
	}
	for i := 2; i < len(pt1lines); i += 3 {
		plainlines = append(plainlines, pt1lines[i])
	}

	ciphertext := strings.Join(cipherlines, "\n")
	plaintext := strings.Join(plainlines, "\n")

	key := "9-1,24,6-23"
	alphabet := "NOKTYUXEQLHBRMPDICJASVWGZF"
	machine, err := NewMachineFromKey(key, alphabet)
	if err != nil {
		t.Fatalf("Could not make machine from key, alphabet: %s, %s", key, alphabet)
	}

	fmt.Printf("%s\n", machine.decipherMessage(ciphertext))
	fmt.Printf("%s\n", machine.encipherMessage(plaintext))
}
