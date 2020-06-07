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

var pieces = []string{"B", "R", "G", "Y", "P", "W"}
var guesses = [][]string{}

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
	default:
		rand.Seed(time.Now().Unix())
		guess = gg[rand.Intn(len(gg))]
	}
	fmt.Printf("\nMy guess: %s\n", colorPrint(guess))
	fmt.Printf("Guess certainty: %d%%\n", 100/len(gg))
	r := getResponse()
	fmt.Println("")
	for _, g := range gg {
		rr := calculateResponse(g, guess)
		if r == rr {
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
	r := getResponse()
	for _, g := range gg {
		rr := calculateResponse(g, guess)
		if r == rr {
			leftGuesses = append(leftGuesses, g)
		}
	}
	if len(leftGuesses) == 1 {
		fmt.Printf("Combination found: %s.\n", colorPrint(leftGuesses[0]))
		os.Exit(0)
	}
	return leftGuesses
}

func getResponse() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Response? ")
	response, _ := reader.ReadString('\n')
	return response[:4]
}

func calculateResponse(guess []string, code []string) string {
	response := []string{}
	g := make([]string, 4)
	c := make([]string, 4)
	copy(g, guess)
	copy(c, code)
	for i := range g {
		if g[i] == c[i] {
			response = append(response, "2")
			g[i] = "-"
			c[i] = "+"
		}
	}
	for i := range g {
		ii := strInSlice(g[i], c)
		if ii >= 0 {
			response = append(response, "1")
			g[i] = "-"
			c[ii] = "+"
		}
	}
	for len(response) != 4 {
		response = append(response, "0")
	}
	return strings.Join(response, "")
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
