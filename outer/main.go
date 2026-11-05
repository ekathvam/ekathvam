// Ekathvam welcome banner.
//
// A minimal, centered terminal greeting. Standard library only,
// so it compiles to one self-contained binary.
//
// Build:  CGO_ENABLED=0 go build -ldflags="-s -w" -o ekathvam-welcome
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
	"ekathvam/outer/bootstrap"
)

const (
	version = "v0.0.1"
	stage   = "alpha"
)

// ---------------------------------------------------------------
// Colors
// ---------------------------------------------------------------

const (
	colorNone  = 0 // no colors (NO_COLOR set, or output is not a terminal)
	colorBasic = 1 // 8 colors (the Linux text console)
	color256   = 2 // 256 colors (almost every modern terminal)
)

func detectColor() int {
	// https://no-color.org: any value means "do not use color".
	if _, set := os.LookupEnv("NO_COLOR"); set {
		return colorNone
	}

	// If output goes to a file or another program, skip escape codes.
	info, err := os.Stdout.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		return colorNone
	}

	switch term := os.Getenv("TERM"); {
	case term == "" || term == "dumb":
		return colorNone
	case term == "linux":
		return colorBasic
	default:
		return color256
	}
}

type palette struct {
	level int
}

// paint wraps text in color codes. A color is a 256-color number plus
// a basic fallback (0 black, 1 red, 2 green, 3 yellow, 4 blue,
// 5 magenta, 6 cyan, 7 white).
func (p palette) paint(text string, c256, basic int, bold bool) string {
	switch p.level {
	case colorNone:
		return text
	case colorBasic:
		weight := "0"
		if bold {
			weight = "1"
		}
		return "\x1b[" + weight + ";" + strconv.Itoa(30+basic) + "m" + text + "\x1b[0m"
	default:
		weight := ""
		if bold {
			weight = "1;"
		}
		return "\x1b[" + weight + "38;5;" + strconv.Itoa(c256) + "m" + text + "\x1b[0m"
	}
}

func (p palette) dim(text string) string {
	switch p.level {
	case colorNone:
		return text
	case colorBasic:
		return "\x1b[1;30m" + text + "\x1b[0m" // bright black = grey
	default:
		return "\x1b[38;5;245m" + text + "\x1b[0m"
	}
}

// ---------------------------------------------------------------
// The big "EKATHVAM" lettering (5 rows per letter)
// ---------------------------------------------------------------

var font = map[rune][5]string{
	'E': {"█████", "█    ", "████ ", "█    ", "█████"},
	'K': {"█   █", "█  █ ", "███  ", "█  █ ", "█   █"},
	'A': {" ███ ", "█   █", "█████", "█   █", "█   █"},
	'T': {"█████", "  █  ", "  █  ", "  █  ", "  █  "},
	'H': {"█   █", "█   █", "█████", "█   █", "█   █"},
	'V': {"█   █", "█   █", "█   █", " █ █ ", "  █  "},
	'M': {"█   █", "██ ██", "█ █ █", "█   █", "█   █"},
}

// One color per letter: a gradient from cyan to violet.
var gradient256 = []int{51, 45, 39, 33, 63, 99, 135, 171}
var gradientBasic = []int{6, 6, 6, 4, 4, 5, 5, 5}

func banner(p palette, word string) []string {
	rows := make([]string, 5)
	for i, letter := range word {
		glyph := font[letter]
		for r := 0; r < 5; r++ {
			if i > 0 {
				rows[r] += "  "
			}
			rows[r] += p.paint(glyph[r], gradient256[i], gradientBasic[i], true)
		}
	}
	return rows
}

// ---------------------------------------------------------------
// Layout
// ---------------------------------------------------------------

// terminalWidth asks the terminal how wide it is. Falls back to 80.
func terminalWidth() int {
	if cols, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && cols > 0 {
		return cols
	}
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	if out, err := cmd.Output(); err == nil {
		fields := strings.Fields(string(out))
		if len(fields) == 2 {
			if cols, err := strconv.Atoi(fields[1]); err == nil && cols > 0 {
				return cols
			}
		}
	}
	return 80
}

// visibleWidth counts the columns a string takes on screen,
// ignoring color escape codes.
func visibleWidth(s string) int {
	width := 0
	inEscape := false
	for _, r := range s {
		switch {
		case inEscape:
			if r == 'm' {
				inEscape = false
			}
		case r == '\x1b':
			inEscape = true
		default:
			width++
		}
	}
	return width
}

func center(s string, width int) string {
	pad := (width - visibleWidth(s)) / 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat(" ", pad) + s
}

// ---------------------------------------------------------------
// Small helpers
// ---------------------------------------------------------------

func userName() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	if name := os.Getenv("USER"); name != "" {
		return name
	}
	return "operator"
}

func greetingFor(t time.Time) string {
	switch h := t.Hour(); {
	case h >= 5 && h < 12:
		return "Good morning"
	case h >= 12 && h < 17:
		return "Good afternoon"
	case h >= 17 && h < 21:
		return "Good evening"
	default:
		return "Good night"
	}
}

// ---------------------------------------------------------------
// Main
// ---------------------------------------------------------------

func main() {
	p := palette{level: detectColor()}
	width := terminalWidth()

	// The Linux text console cannot draw the arrow character.
	arrow := "→"
	if os.Getenv("TERM") == "linux" {
		arrow = "->"
	}

	var screen []string
	add := func(s string) { screen = append(screen, center(s, width)) }
	blank := func() { screen = append(screen, "") }

	blank()

	// Big lettering, or a compact title on narrow terminals.
	if width >= 60 {
		for _, row := range banner(p, "EKATHVAM") {
			add(row)
		}
	} else {
		add(p.paint("E K A T H V A M", 51, 6, true))
	}

	blank()
	add(p.paint("Unity through interoperability", 117, 6, false))
	blank()

	// The greeting.
	add(p.paint(greetingFor(time.Now())+", "+userName()+".", 229, 3, true))
	blank()

	// A few keywords.
	keywords := []string{"Open", "Interoperable", "Self-managing", "Secure"}
	if width < 54 {
		keywords = []string{"Open", "Interoperable", "Secure"} // fits small screens
	}
	add(p.dim(strings.Join(keywords, "  ·  ")))
	blank()

	// The development philosophy.
	sep := " " + p.dim(arrow) + " "
	add(p.paint("Build", 51, 6, true) + sep +
		p.paint("Validate", 45, 6, true) + sep +
		p.paint("Deliver", 99, 5, true) + sep +
		p.paint("Improve", 171, 5, true))
	blank()

	// Version.
	add(p.paint(version, 45, 6, true) + " " + p.paint(stage, 221, 3, true))
	blank()

	fmt.Println(strings.Join(screen, "\n"))
	
	outerbootstrap.Bootstrap()
}
