package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

var guesses = [][]string{}
var pieces = []string{"B", "R", "G", "Y", "P", "W"}
var lastPegs = ""

func main() {
	genGuesses(pieces, 0)
	turn := 0
	if getDoIGoFirst() {
		leftGuesses := makeGuess(guesses, turn)
		turn++
		leftGuesses = getCombi(leftGuesses)
		turn++
		leftGuesses = makeGuess(leftGuesses, turn)
		turn++
		leftGuesses = getCombi(leftGuesses)
		turn++
		leftGuesses = makeGuess(leftGuesses, turn)
		turn++
		leftGuesses = getCombi(leftGuesses)
		turn++
		leftGuesses = makeGuess(leftGuesses, turn)
		getCombi(leftGuesses)
	} else {
		leftGuesses := getCombi(guesses)
		turn++
		leftGuesses = makeGuess(leftGuesses, turn)
		turn++
		leftGuesses = getCombi(leftGuesses)
		turn++
		leftGuesses = makeGuess(leftGuesses, turn)
		turn++
		leftGuesses = getCombi(leftGuesses)
		turn++
		leftGuesses = makeGuess(leftGuesses, turn)
		turn++
		leftGuesses = getCombi(leftGuesses)
		turn++
		makeGuess(leftGuesses, turn)
	}
}

func makeGuess(gg [][]string, turn int) [][]string {
	var guess []string
	leftGuesses := [][]string{}
	switch turn {
	case 0:
		rand.Seed(time.Now().UTC().UnixNano())
		a := 0
		b := 0
		for a == b {
			a = rand.Intn(len(pieces))
			b = rand.Intn(len(pieces))
		}
		guess = []string{pieces[a], pieces[a], pieces[b], pieces[b]}
	default:
		guess = minMaxGuess(gg, lastPegs)
	}
	fmt.Printf("\nMy guess: %s\n", colorPrint(guess))
	fmt.Printf("Guess certainty: %d%%\n", 100/len(gg))
	pegs := getPegs()
	fmt.Println("")
	for _, g := range gg {
		cPegs := calculatePegs(g, guess)
		lastPegs = cPegs
		if pegs == cPegs {
			leftGuesses = append(leftGuesses, g)
		}
	}
	if len(leftGuesses) == 1 {
		fmt.Printf("Combination found: %s.\n", colorPrint(leftGuesses[0]))
		os.Exit(0)
	}
	return leftGuesses
}

func colorPrint(guess []string) string {
	colorGuess := ""
	for _, g := range guess {
		switch g {
		case "B":
			colorGuess += fmt.Sprintf("%s", Blue(g))
		case "R":
			colorGuess += fmt.Sprintf("%s", Red(g))
		case "G":
			colorGuess += fmt.Sprintf("%s", Green(g))
		case "Y":
			colorGuess += fmt.Sprintf("%s", Yellow(g))
		case "P":
			colorGuess += fmt.Sprintf("%s", Magenta(g))
		case "W":
			colorGuess += fmt.Sprintf("%s", White(g))
		}
	}
	return colorGuess
}

func getDoIGoFirst() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do I go first? (y/n): ")
	first, _ := reader.ReadString('\n')
	first = strings.ToLower(first)
	return first[:1] == "y"
}

func getCombi(gg [][]string) [][]string {
	leftGuesses := [][]string{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your guess? ")
	combi, _ := reader.ReadString('\n')
	combi = strings.ToUpper(combi)
	guess := strings.Split(combi, "")
	pegs := getPegs()
	for _, g := range gg {
		cPegs := calculatePegs(g, guess)
		lastPegs = cPegs
		if pegs == cPegs {
			leftGuesses = append(leftGuesses, g)
		}
	}
	if len(leftGuesses) == 1 {
		fmt.Printf("Combination found: %s.\n", colorPrint(leftGuesses[0]))
		os.Exit(0)
	}
	return leftGuesses
}

func minMaxGuess(gg [][]string, pegs string) []string {
	possiblePegs := []string{
		"0000", "1000", "2000",
		"1100", "2100", "2200",
		"1110", "2110", "2210",
		"1111", "2111", "2211",
	}
	targetPegs := possiblePegs[strInSlice(pegs, possiblePegs):]
	eliminationCount := make([]int, len(gg))
	for i := range eliminationCount {
		eliminationCount[i] = 0
	}
	for i1, g1 := range gg {
		for i2, g2 := range gg {
			if i1 == i2 {
				continue
			}
			cPegs := calculatePegs(g2, g1)
			if strInSlice(cPegs, targetPegs) < 0 {
				eliminationCount[i1]++
			}
		}
	}
	_, max := minMax(eliminationCount)
	ei := []int{}
	for i := range eliminationCount {
		if eliminationCount[i] == max {
			ei = append(ei, i)
		}
	}
	return gg[ei[rand.Intn(len(ei))]]
}

func getPegs() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pegs? ")
	pegs, _ := reader.ReadString('\n')
	return pegs[:4]
}

func calculatePegs(guess []string, code []string) string {
	pegs := []string{}
	g := make([]string, 4)
	c := make([]string, 4)
	copy(g, guess)
	copy(c, code)
	for i := range g {
		if g[i] == c[i] {
			pegs = append(pegs, "2")
			g[i] = "-"
			c[i] = "+"
		}
	}
	for i := range g {
		ii := strInSlice(g[i], c)
		if ii >= 0 {
			pegs = append(pegs, "1")
			g[i] = "-"
			c[ii] = "+"
		}
	}
	for len(pegs) != 4 {
		pegs = append(pegs, "0")
	}
	return strings.Join(pegs, "")
}

func genGuesses(a []string, i int) {
	if i == 4 {
		aa := make([]string, 4)
		copy(aa, a)
		guesses = append(guesses, aa)
		return
	}
	genGuesses(a, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		genGuesses(a, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func strInSlice(x string, a []string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func minMax(array []int) (int, int) {
	var min int = array[0]
	var max int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
