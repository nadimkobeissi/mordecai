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
	gg := [][]string{}
	turn := 0
	pegs := "0000"
	switch getDuplicatesAllowed() {
	case true:
		gg = genGuessesWithDuplicates(pieces, 4)
	case false:
		gg = genGuessesNoDuplicates(pieces, 0, [][]string{})
	}
	switch getGameMode() {
	case 0:
		for {
			gg, turn, pegs = makeGuess(gg, turn, pegs)
		}
	case 1:
		for {
			gg, turn, _ = makeGuess(gg, turn, pegs)
			gg, turn, pegs = getCombi(gg, turn)
		}
	case 2:
		for {
			gg, turn, pegs = getCombi(gg, turn)
			gg, turn, _ = makeGuess(gg, turn, pegs)
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

func genGuessesWithDuplicates(a []string, l int) [][]string {
	gg := [][]string{}
	gn := make([]int, l)
	for i := 0; i < l; i++ {
		gn[i] = 0
	}
	for {
		g := []string{}
		for i := range gn {
			g = append(g, a[gn[i]])
		}
		gg = append(gg, g)
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

func genGuessesNoDuplicates(a []string, i int, gg [][]string) [][]string {
	if i == 4 {
		aa := make([]string, 4)
		copy(aa, a)
		gg = append(gg, aa)
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

func makeGuess(gg [][]string, turn int, pegs string) ([][]string, int, string) {
	var guess []string
	lgg := [][]string{}
	switch turn {
	case 0:
		guess = []string{"B", "B", "R", "R"}
	default:
		guess = minMaxGuess(gg, pegs)
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
	return lgg, turn + 1, nPegs
}

func getCombi(gg [][]string, turn int) ([][]string, int, string) {
	lgg := [][]string{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your guess? ")
	combi, _ := reader.ReadString('\n')
	combi = strings.ToUpper(combi)
	guess := strings.Split(combi, "")
	nPegs := getPegs()
	for _, g := range gg {
		cPegs := calculatePegs(g, guess)
		if nPegs == cPegs {
			lgg = append(lgg, g)
		}
	}
	checkCombiFound(lgg)
	return lgg, turn + 1, nPegs
}

func getPegs() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pegs? ")
	pegs, _ := reader.ReadString('\n')
	p := strings.Split(pegs[:4], "")
	sort.Strings(p)
	return p[3] + p[2] + p[1] + p[0]
}

func minMaxGuess(gg [][]string, pegs string) []string {
	possiblePegs := []string{
		"0000", "1000", "2000",
		"1100", "2100", "2200",
		"1110", "2110", "2210",
		"1111", "2111", "2211",
		"2220", "2222",
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
	_, maxi := minMax(eliminationCount)
	return gg[maxi]
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

func checkCombiFound(gg [][]string) {
	if len(gg) == 1 {
		fmt.Printf("Combination found: %s.\n", colorPrint(gg[0]))
		os.Exit(0)
	}
}

func colorPrint(guess []string) string {
	colorGuess := ""
	for _, g := range guess {
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
