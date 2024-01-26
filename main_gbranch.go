package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// Define the length of the commit message as a constant
const commitMsgLength = 70

// Define the special symbol to replace '+'
const specialSymbol = "⚡"

type Branch struct {
	IsCurrent bool
	Symbol    string
	Name      string
	Remote    string
	Message   string
}

func main() {
	fmt.Println("")
	// 引数が有る場合は git branch を引数付きで実行したいと判断
	if len(os.Args) != 1 {
		execNormal(os.Args)
		return
	}
	// Execute the git command
	cmd := exec.Command("git", "branch", "-vv")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing git command:", err)
		return
	}

	// Reset scanner to start
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	branches := []*Branch{}
	for scanner.Scan() {
		line := scanner.Text()
		branch := parseAndPrintLine(line)
		branches = append(branches, branch)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading git output:", err)
		return
	}

	// print to console
	print(branches)
}

func execNormal(osArgs []string) {
	args := osArgs
	args[0] = "branch"
	// コマンド生成
	cmd := exec.Command("git", os.Args...)
	// 出力先設定
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// 実行
	err := cmd.Run()
	// エラー時
	if err != nil {
		fmt.Println(stderr.String())
		return
	}
	// 通常時
	fmt.Println(stdout.String())
}

func getBranchName(line string) string {
	re := regexp.MustCompile(`^(\s*[\+\*]?\s*)([\w/+\-\.]+)`)
	matches := re.FindStringSubmatch(line)
	if matches != nil {
		return matches[2]
	}

	return ""
}

func parseAndPrintLine(line string) *Branch {
	// Improved regular expression to match the required parts of each line accurately
	re := regexp.MustCompile(`^(\s*[\+\*]?\s*)([\w/+\-\.]+)\s+([a-f0-9]+)\s+(\[.*?\])?\s*(.*)$`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	// replace symbol
	symbol := matches[1]
	isCurrent := false
	if symbol == "* " {
		symbol = specialSymbol
		isCurrent = true
	} else if symbol == "" {
		symbol = "  "
	}

	branch := matches[2]
	remote := ""
	// Extract remote branch name from the square brackets
	if matches[4] != "" {
		remoteBracketContents := matches[4]
		remoteParts := strings.SplitN(remoteBracketContents[1:len(remoteBracketContents)-1], ":", 2)
		remote = strings.TrimSpace(remoteParts[0])
	}

	// Adjust spaces in commit message
	//msg := adjustSpace(matches[5])
	msg := matches[5]

	return &Branch{
		IsCurrent: isCurrent,
		Symbol:    symbol,
		Name:      branch,
		Remote:    remote,
		Message:   msg,
	}
}

// adjustSpace adjusts the string to have an effective length of 70 spaces, considering ASCII and non-ASCII characters
func adjustSpace(s string, maxLen int) string {
	length := strLen(s)
	if maxLen > commitMsgLength {
		maxLen = commitMsgLength
	}

	if length > maxLen {
		runes := []rune(s)
		adjustedLength := 0
		for i, r := range []rune(s) {
			if adjustedLength >= maxLen {
				return string(runes[:i])
			}
			if unicode.IsPrint(r) && r < 128 {
				adjustedLength += 1
			} else {
				adjustedLength += 2
			}
		}
	}
	return s + strings.Repeat(" ", maxLen-length)
}

func print(branches []*Branch) {
	max_branch_len := 0
	max_commit_len := 0
	for _, b := range branches {
		if len(b.Name) > max_branch_len {
			max_branch_len = len(b.Name)
		}
		if len(b.Message) > max_commit_len {
			max_commit_len = strLen(b.Message)
		}
	}

	// IsCurrent が true の Branch を見つけて先頭に移動
	/*
		for i, branch := range branches {
			// 要素を先頭に移動
			if branch.IsCurrent {
				branches = append(branches[:i], branches[i+1:]...)
				branches = append([]*Branch{branch}, branches...)
				break
			}
		}
	*/

	// 順番に出力
	for _, b := range branches {
		name := b.Name
		name += strings.Repeat(" ", max_branch_len-len(name)+1)
		msg := adjustSpace(b.Message, max_commit_len)
		out := fmt.Sprintf("%s %s %s - %s\n", b.Symbol, name, msg, b.Remote)

		if b.IsCurrent {
			// red := color.New(color.BgRed)
			// red.Add(color.FgWhite)
			red := color.New(color.FgYellow)
			out = red.Sprintf(out)
			// out = color.RedString(out)
		}

		fmt.Print(out)
	}
}

func strLen(s string) int {
	length := 0
	runes := []rune(s)
	for _, r := range runes {
		if unicode.IsPrint(r) && r < 128 {
			length += 1 // ASCII characters count as 1
		} else {
			length += 2 // Non-ASCII characters count as 2
		}
	}

	return length
}
