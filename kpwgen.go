package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

// Default password length
// https://blog.codinghorror.com/your-password-is-too-damn-short/
const defaultLength = 15

// Default number of passwords to generate
// https://en.wikipedia.org/wiki/The_Magical_Number_Seven,_Plus_or_Minus_Two
const defaultPasswords = 7

// Filtered koremutake syllables
// http://shorl.com/koremutake.php
var syllables = []string{"ba", "be", "bi", "bo", "bu", "by", "da", "de", "di",
	"do", "du", "dy", "fe", "fi", "fo", "fu", "fy", "ga", "ge", "gi", "go", "gu",
	"gy", "ha", "he", "hi", "ho", "hu", "hy", "ja", "je", "ji", "jo", "ju", "jy",
	"ka", "ke", "ko", "ku", "ky", "la", "le", "li", "lo", "", "lu", "ly", "ma",
	"me", "mi", "mo", "mu", "my", "na", "ne", "ni", "no", "nu", "ny", "pa", "pe",
	"pi", "po", "pu", "py", "ra", "re", "ri", "ro", "", "ru", "ry", "sa", "se",
	"si", "so", "su", "sy", "ta", "te", "ti", "to", "tu", "ty", "va", "ve", "vi",
	"vo", "vu", "vy", "bra", "bre", "bri", "", "bro", "bru", "bry", "dra", "dre",
	"dri", "dro", "dru", "dry", "fra", "fre", "fri", "fro", "fru", "fry", "gra",
	"gre", "gri", "", "gro", "gru", "gry", "pra", "pre", "pri", "pro", "pru",
	"pry", "sta", "ste", "sti", "sto", "stu", "sty", "tra", "tre", "er", "", "ed",
	"in", "ex", "al", "en", "an", "ad", "or", "at", "ca", "ap", "el", "ci", "an",
	"et", "it", "ob", "of", "af", "au", "cy", "im", "op", "co", "", "up", "ing",
	"con", "ter", "com", "per", "ble", "der", "cal", "man", "est", "for", "mer",
	"col", "ful", "get", "low", "son", "", "tle", "day", "pen", "pre", "ten",
	"tor", "ver", "ber", "can", "ple", "fer", "gen", "den", "mag", "sub", "sur",
	"men", "min", "", "out", "tal", "but", "cit", "cle", "cov", "dif", "ern",
	"eve", "hap", "ket", "nal", "sup", "ted", "tem", "tin", "tro", "tro"}

// Punctuation characters to use. Exclude ' because it's often disallowed
// in an amateurish attempt at preventing SQL injection, exclude \ because
// it can cause problems due to quoting rules.
var punctuation = []rune{'!', '"', '#', '$', '%', '&', '(', ')', '*', '+',
	',', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@', '[', ']', '^',
	'_', '`', '{', '|', '}', '~'}

var digit = flag.Bool("d", false, "put a number in each password")
var upper = flag.Bool("u", false, "put an uppercase letter in each password")
var punct = flag.Bool("p", false, "put a punctuation character in each password")
var numpass = flag.Int("n", defaultPasswords, "number of passwords to generate")

var numSyllables = len(syllables)

// Number of random bytes to ask for at a time
const poolSize = 128

type RNG struct {
	pool []byte
	i    int
}

func NewRNG() *RNG {
	return &RNG{
		pool: make([]byte, poolSize),
		i:    0,
	}
}

func (r *RNG) GetByte() int {
	if r.i == 0 {
		_, err := rand.Read(r.pool)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't generate random numbers: %s", err)
			os.Exit(1)
		}
	}
	b := r.pool[r.i]
	r.i = (r.i + 1) % poolSize
	return int(b)
}

func (r *RNG) GetPassword(pwlen int) (string, error) {
	// bytes.Buffer is the best way to do string append.
	// Because we built the password from syllables of up to 3 characters,
	// the password might be up to 2 characters longer than asked for.
	var pw bytes.Buffer
	pw.Grow(pwlen + 2)
	for pw.Len() < pwlen {
		n := r.GetByte() & (len(syllables) - 1)
		syl := syllables[n]
		pw.Write([]byte(syl))
	}
	// From this point it's easier to deal with a byte slice.
	pwb := pw.Bytes()
	di := -1
	if *digit {
		di = r.GetByte() % (len(pwb) - 1)
		d := '0' + r.GetByte()&9
		pwb[di] = byte(d)
	}
	ui := -1
	if *upper {
		// Get an index which isn't the one we just picked for putting a digit in,
		// if we did that
		for {
			ui = r.GetByte() % (len(pwb) - 1)
			if ui != di {
				break
			}
		}
		// This would be wrong if we didn't know the byte was an ASCII character
		pwb[ui] = byte(unicode.ToUpper(rune(pwb[ui])))
	}
	if *punct {
		pi := -1
		// Get an index which isn't the uppercase or the digit one
		for {
			pi = r.GetByte() % (len(pwb) - 1)
			if pi != di && pi != ui {
				break
			}
		}
		pn := r.GetByte() & (len(punctuation) - 1)
		pc := punctuation[pn]
		pwb[pi] = byte(pc)
	}
	return string(pwb), nil
}

func main() {
	flag.Parse()
	pwlen := defaultLength
	var err error
	if len(flag.Args()) > 0 {
		sx := flag.Args()[0]
		pwlen, err = strconv.Atoi(sx)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: argument must be an integer number of characters")
			os.Exit(1)
		}
	}
	r := NewRNG()
	for i := 0; i < *numpass; i++ {
		pw, err := r.GetPassword(pwlen)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating password: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(pw)
	}
}
