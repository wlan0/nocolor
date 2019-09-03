// Copyright 2019 Sidhartha Mani <sidharthamn@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"os"
	"io"

	"github.com/spf13/cobra"
)

type colorstate int

const (
	colored colorstate = iota
	uncolored
)

var cmd = &cobra.Command{
	Use: "nocolor",
	Short: "strips terminal color from output",
	Long: `
commands | nocolor
`,
	Run: func(*cobra.Command, []string) {
		if err := noColor(os.Stdin, os.Stdout); err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
	},
}


func main() {
	cmd.Execute()
}

type colorParser func(r byte, out io.Writer) (interface{}, colorstate, error)

func noColor(input, output *os.File) error {
	buf := bufio.NewReader(input)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("nocolor: %v", err)
		}
		var colorFn colorParser
		colorFn = stripColor
		for i := range line {
			//			fmt.Fprintf(output, "%q", line[i])
			parseFnInterface, _, err := colorFn(line[i], output)
			if err != nil {
				return fmt.Errorf("nocolor: %v", err)
			}
			var ok bool
			if colorFn, ok = parseFnInterface.(func(r byte, out io.Writer) (interface{}, colorstate, error)); !ok {
				return fmt.Errorf("nocolor: internal error: colorParser expected")
			}
		}
	}
	return nil
}

func stripColor(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == '\\' {
		return stripColorBackslash, uncolored, nil 
	}
	if r == '\x1b' {
		return stripColorInColoredZoneX1b, colored, nil
	}
	fmt.Fprintf(out, "%c", r)
	return stripColor, uncolored, nil
}

func stripColorBackslash(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == 'x' {
		return stripColorInColoredZoneX, colored, nil
	}
	if r == 'e' {
		return stripColorInColoredZoneX1b, colored, nil
	}
	fmt.Fprintf(out, "\\%c", r)
	return stripColor, uncolored, nil
}

func stripColorInColoredZoneX(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == '1' {
		return stripColorInColoredZoneX1, colored, nil
	}
	return stripColor, colored, fmt.Errorf("nocolor: unexpected escape sequence")
}

func stripColorInColoredZoneX1(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == 'b' {
		return stripColorInColoredZoneX1b, colored, nil
	}
	return stripColor, colored, fmt.Errorf("nocolor: unexpected escape sequence")
}

func stripColorInColoredZoneX1b(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == '[' {
		return stripColorInColoredZoneX1bSQB, colored, nil
	}
	fmt.Fprintf(out, "\x1b%c", r)
	return stripColor, colored, nil
}

func stripColorInColoredZoneX1bSQB(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r >= '0' && r <= '9' {
		return stripColorInColoredZoneX1bSQB, colored, nil
	}
	if r == ';' {
		return stripColorInColoredZoneX1bSQB, colored, nil
	}
	if r == 'm' {
		return stripColorInColoredZoneX1bSQBIn, colored, nil
	}
	return stripColor, colored, fmt.Errorf("nocolor: unexpected color code")
}

func stripColorInColoredZoneX1bSQBIn(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == '\\' {
		return stripColorOutColoredZone, colored, nil
	}
	if r == '\x1b' {
		return stripColorOutColoredZoneX1b, colored, nil
	}
	fmt.Fprintf(out, "%c", r)
	return stripColorInColoredZoneX1bSQBIn, colored, nil
}

func stripColorOutColoredZone(r byte, out io.Writer) (interface{}, colorstate, error) {
	if  r == 'x' {
		return stripColorOutColoredZoneX, colored, nil
	}
	if r == 'e' {
		return stripColorOutColoredZoneX1b, colored, nil
	}
	fmt.Fprintf(out, "\\%c", r)
	return stripColorInColoredZoneX1bSQBIn, colored, nil
}

func stripColorOutColoredZoneX(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == '1' {
		return stripColorOutColoredZoneX1, colored, nil
	}
	return stripColorInColoredZoneX1bSQBIn, colored, fmt.Errorf("nocolor: unexpected escape sequence")
}

func stripColorOutColoredZoneX1(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == 'b' {
		return stripColorOutColoredZoneX1b, colored, nil
	}
	return stripColorInColoredZoneX1bSQBIn, colored, fmt.Errorf("nocolor: unexpected escape sequence")
}

func stripColorOutColoredZoneX1b(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r == '[' {
		return stripColorOutColoredZoneX1bSQB, colored, nil
	}
	fmt.Fprintf(out, "\x1b%c", r)
	return stripColor, colored, nil
}

func stripColorOutColoredZoneX1bSQB(r byte, out io.Writer) (interface{}, colorstate, error) {
	if r >= '0' && r <= '9' {
		return stripColorOutColoredZoneX1bSQB, colored, nil
	}
	if r == ';' {
		return stripColorOutColoredZoneX1bSQB, colored, nil
	}
	if r == 'm' {
		return stripColor, colored, nil
	}
	return stripColor, colored, fmt.Errorf("nocolor: unexpected color code")
}
