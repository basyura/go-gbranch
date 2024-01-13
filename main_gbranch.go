package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

const commitMsgLength = 70 // Define the length of the commit message as a constant
const specialSymbol = "⏰"  // Define the special symbol to replace '+'

type Branch struct {
	Symbol  string
	Name    string
	Message string
}

func main() {
	// Execute the git command
	cmd := exec.Command("git", "branch", "-vv")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing git command:", err)
		return
	}

	// Find the longest branch name for alignment
	var maxLength int
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		branchLength := len(getBranchName(line))
		if branchLength > maxLength {
			maxLength = branchLength
		}
	}

	// Reset scanner to start
	scanner = bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parseAndPrintLine(line, maxLength)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading git output:", err)
	}
}

func getBranchName(line string) string {
	re := regexp.MustCompile(`^(\s*[\+\*]?\s*)([\w/+\-\.]+)`)
	matches := re.FindStringSubmatch(line)
	if matches != nil {
		return matches[2]
	}
	return ""
}

func parseAndPrintLine(line string, maxLength int) {
	// Improved regular expression to match the required parts of each line accurately
	re := regexp.MustCompile(`^(\s*[\+\*]?\s*)([\w/+\-\.]+)\s+([a-f0-9]+)\s+(\[.*?\])?\s*(.*)$`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return
	}

	isCurrent := false
	symbol := matches[1]
	if symbol == "* " {
		symbol = specialSymbol // Replace '+' with the special symbol and add a space
		isCurrent = true
	} else if symbol == "" {
		symbol = "  " // Two spaces for missing symbol
	}

	branch := matches[2]
	remote := ""
	if matches[4] != "" {
		// Extract remote branch name from the square brackets
		remoteBracketContents := matches[4]
		remoteParts := strings.SplitN(remoteBracketContents[1:len(remoteBracketContents)-1], ":", 2)
		remote = strings.TrimSpace(remoteParts[0])
	}
	msg := adjustSpace(matches[5]) // Adjust spaces in commit message

	padding := branch + strings.Repeat(" ", maxLength-len(branch)+1)

	out := fmt.Sprintf("%s %s: %s | %s\n", symbol, padding, msg, remote)
	if isCurrent {
		// red := color.New(color.BgRed)
		// red.Add(color.FgWhite)
		red := color.New(color.FgRed)
		out = red.Sprintf(out)
		// out = color.RedString(out)
	}

	fmt.Print(out)
}

// adjustSpace adjusts the string to have an effective length of 70 spaces, considering ASCII and non-ASCII characters
func adjustSpace(s string) string {
	effectiveLength := 0
	runes := []rune(s)
	for _, r := range runes {
		if unicode.IsPrint(r) && r < 128 {
			effectiveLength += 1 // ASCII characters count as 1
		} else {
			effectiveLength += 2 // Non-ASCII characters count as 2
		}
	}

	if effectiveLength > commitMsgLength {
		adjustedLength := 0
		for i, r := range runes {
			if adjustedLength >= commitMsgLength {
				return string(runes[:i])
			}
			if unicode.IsPrint(r) && r < 128 {
				adjustedLength += 1
			} else {
				adjustedLength += 2
			}
		}
	}
	return s + strings.Repeat(" ", commitMsgLength-effectiveLength)
}
