package print

import (
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"strings"
)

const (
	ColorReset      = "\x1b[0m\x1b[K"
	ColorLineNumber = "\x1b[1;33m"  /* yellow with black background */
	ColorPath       = "\x1b[1;32m"  /* bold green */
	ColorMatch      = "\x1b[30;43m" /* black with yellow background */
)

type Modifier interface {
	Path(path string) string
	LineNumber(lineNum int) string
	Match(pattern, match string) string
}

type Plain struct {}

func (m *Plain) Path(path string) string {
	return path
}

func (m *Plain) LineNumber(lineNum int) string {
	return fmt.Sprintf("%d:", lineNum)
}

func (m *Plain) Match(pattern, match string) string {
	return match
}

type Color struct {}

func (m *Color) Path(path string) string {
	return ColorPath + path + ColorReset
}

func (m *Color) LineNumber(lineNum int) string {
	return ColorLineNumber + fmt.Sprintf("%d", lineNum) + ColorReset + ":"
}

func (m *Color) Match(pattern, match string) string {
	return strings.Replace(match, pattern, ColorMatch + pattern + ColorReset, -1)
}

type Match struct {
	LineNum int
	Match   string
}

type Params struct {
	Pattern string
	Path    string
	Matches []*Match
}

type Printer struct {
	In     chan *Params
	Done   chan bool
	Option *option.Option
}

func (self *Printer) Print() {
	var modifier Modifier
	if self.Option.NoColor {
		modifier = &Plain{}
	} else {
		modifier = &Color{}
	}

	for arg := range self.In {

		if len(arg.Matches) == 0 {
			continue
		}

		path := modifier.Path(arg.Path)
		if self.Option.FilesWithMatches {
			fmt.Println(path)
			continue
		}
		if self.Option.NoGroup {
			path += ":"
		} else {
			fmt.Println(path)
			path = ""
		}

		for _, v := range arg.Matches {
			if v == nil {
				continue
			}
			fmt.Print(path)
			fmt.Print(modifier.LineNumber(v.LineNum))
			fmt.Println(modifier.Match(arg.Pattern, v.Match))
		}
		if !self.Option.NoGroup {
			fmt.Println()
		}
	}
	self.Done <- true
}
