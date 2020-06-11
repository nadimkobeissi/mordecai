package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	. "github.com/logrusorgru/aurora"
)

var pieces = []string{"B", "R", "G", "Y", "P", "W"}

func main() {
	gg := []string{}
	pgg := []string{}
	turn := 0
	pegs := "0000"
	switch getDuplicatesAllowed() {
	case true:
		gg = genGuessesWithDuplicates(pieces, 4)
		pgg = genGuessesWithDuplicates(pieces, 4)
	case false:
		gg = genGuessesNoDuplicates(pieces, 0, []string{})
		pgg = genGuessesNoDuplicates(pieces, 0, []string{})
	}
	switch getGameMode() {
	case 0:
		for {
			gg, pgg, turn, pegs = makeGuess(gg, pgg, turn, pegs)
		}
	case 1:
		for {
			gg, pgg, turn, _ = makeGuess(gg, pgg, turn, pegs)
			gg, pgg, turn, pegs = getCombi(gg, pgg, turn)
		}
	case 2:
		for {
			gg, pgg, turn, pegs = getCombi(gg, pgg, turn)
			gg, pgg, turn, _ = makeGuess(gg, pgg, turn, pegs)
		}
	}
}

func getDuplicatesAllowed() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Duplicates allowed? (y/n) ")
	allowed, _ := reader.ReadString('\n')
	return allowed[:1] == "y"
}

func getGameMode() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Single player? (y/n) ")
	single, _ := reader.ReadString('\n')
	if single[:1] == "y" {
		return 0
	}
	fmt.Print("Do I go first? (y/n) ")
	first, _ := reader.ReadString('\n')
	first = strings.ToLower(first)
	if first[:1] == "y" {
		return 1
	}
	return 2
}

func genGuessesWithDuplicates(a []string, l int) []string {
	gg := []string{}
	gn := make([]int, l)
	for i := 0; i < l; i++ {
		gn[i] = 0
	}
	for {
		g := []string{}
		for i := range gn {
			g = append(g, a[gn[i]])
		}
		gg = append(gg, strings.Join(g, ""))
		ll := l - 1
		gn[ll]++
		for ll > 0 {
			if gn[ll] == len(a) {
				gn[ll] = 0
				gn[ll-1]++
			}
			ll--
		}
		if gn[ll] == len(a) {
			return gg
		}
	}
}

func genGuessesNoDuplicates(a []string, i int, gg []string) []string {
	if i == 4 {
		aa := make([]string, 4)
		copy(aa, a)
		gg = append(gg, strings.Join(aa, ""))
		return gg
	}
	gg = genGuessesNoDuplicates(a, i+1, gg)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		gg = genGuessesNoDuplicates(a, i+1, gg)
		a[i], a[j] = a[j], a[i]
	}
	return gg
}

func makeGuess(gg []string, pgg []string, turn int, pegs string) ([]string, []string, int, string) {
	var guess string
	lgg := []string{}
	checkUnsolvable(gg)
	switch turn {
	case 0:
		guess = "BBRR"
	default:
		guess = minMaxGuess(gg, pgg, pegs)
	}
	for i := range pgg {
		if pgg[i] == guess {
			pgg = delFromSlice(i, pgg)
			break
		}
	}
	fmt.Printf("\nMy guess: %s\n", colorPrint(guess))
	nPegs := getPegs()
	for _, g := range gg {
		cPegs := calculatePegs(g, guess)
		if nPegs == cPegs {
			lgg = append(lgg, g)
		}
	}
	checkCombiFound(lgg)
	return lgg, pgg, turn + 1, nPegs
}

func getCombi(gg []string, pgg []string, turn int) ([]string, []string, int, string) {
	lgg := []string{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your guess? ")
	guess, _ := reader.ReadString('\n')
	guess = strings.ToUpper(guess[:4])
	for i := range pgg {
		if pgg[i] == guess {
			delFromSlice(i, pgg)
			break
		}
	}
	nPegs := getPegs()
	for _, g := range gg {
		cPegs := calculatePegs(g, guess)
		if nPegs == cPegs {
			lgg = append(lgg, g)
		}
	}
	checkCombiFound(lgg)
	return lgg, pgg, turn + 1, nPegs
}

func getPegs() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pegs? ")
	pegs, _ := reader.ReadString('\n')
	p := strings.Split(pegs[:4], "")
	sort.Strings(p)
	return p[3] + p[2] + p[1] + p[0]
}

func minMaxGuess(gg []string, pgg []string, pegs string) string {
	pPegs := []string{
		"0000", "1000", "1100", "1110", "1111", "2000", "2100",
		"2110", "2111", "2200", "2210", "2211", "2220", "2222",
	}
	ec := []int{}
	for _, g1 := range pgg {
		pegScores := []int{}
		for range pPegs {
			pegScores = append(pegScores, 0)
		}
		for _, g2 := range gg {
			p := calculatePegs(g1, g2)
			pegScores[strInSlice(p, pPegs)]++
		}
		_, maxi := minMax(pegScores)
		ec = append(ec, pegScores[maxi])
	}
	mini, _ := minMax(ec)
	for i, g := range pgg {
		if ec[i] != ec[mini] {
			continue
		}
		if strInSlice(g, gg) >= 0 {
			return g
		}
	}
	return pgg[mini]
}

func calculatePegs(guess string, code string) string {
	pegs := []string{}
	g := guess
	c := code
	for i := range g {
		if g[i] == c[i] {
			pegs = append(pegs, "2")
			g = g[:i] + "-" + g[i+1:]
			c = c[:i] + "+" + c[i+1:]
		}
	}
	for i := range g {
		ii := strings.Index(c, string(g[i]))
		if ii >= 0 {
			pegs = append(pegs, "1")
			g = g[:i] + "-" + g[i+1:]
			c = c[:ii] + "+" + c[ii+1:]
		}
	}
	for len(pegs) != 4 {
		pegs = append(pegs, "0")
	}
	return strings.Join(pegs, "")
}

func checkCombiFound(gg []string) {
	if len(gg) == 1 {
		fmt.Printf("Combination found: %s.\n", colorPrint(gg[0]))
		os.Exit(0)
	}
}

func checkUnsolvable(gg []string) {
	if len(gg) == 0 {
		fmt.Println("Combination is unsolvable.")
		os.Exit(1)
	}
}

func colorPrint(guess string) string {
	gs := strings.Split(guess, "")
	colorGuess := ""
	for _, g := range gs {
		switch g {
		case "B":
			colorGuess += Blue(g).BgBlack().String()
		case "R":
			colorGuess += Red(g).BgBlack().String()
		case "G":
			colorGuess += Green(g).BgBlack().String()
		case "Y":
			colorGuess += Yellow(g).BgBlack().String()
		case "P":
			colorGuess += Magenta(g).BgBlack().String()
		default:
			colorGuess += White(g).BgBlack().String()
		}
	}
	return colorGuess
}

func minMax(array []int) (int, int) {
	var min int = array[0]
	var max int = array[0]
	var mini = 0
	var maxi = 0
	for i, value := range array {
		if max < value {
			max = value
			maxi = i
		}
		if min > value {
			min = value
			mini = i
		}
	}
	return mini, maxi
}

func strInSlice(x string, a []string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func delFromSlice(i int, a []string) []string {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = ""
	a = a[:len(a)-1]
	return a
}
