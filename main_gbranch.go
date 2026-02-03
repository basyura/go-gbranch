package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// Define the length of the commit message as a constant
const commitMsgLength = 70

// Define the special symbol to replace '+'
const specialSymbol = "⭕️"

type Branch struct {
	IsCurrent bool
	Symbol    string
	Name      string
	Remote    string
	Message   string
	Worktree  string
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
	re := regexp.MustCompile(`^(\s*[\+\*]?\s*)([\w/+\-\.]+)\s+([a-f0-9]+)(?:\s+\(([^)]*)\))?\s+(\[.*?\])?\s*(.*)$`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	// replace symbol
	symbolToken := matches[1]
	isCurrent := strings.Contains(symbolToken, "*")
	symbol := ""
	if isCurrent {
		symbol = os.Getenv("GBRANCH_SYMBOL")
		if symbol == "" {
			symbol = specialSymbol
		}
	} else if strings.Contains(symbolToken, "+") {
		symbol = "+"
	}

	branch := matches[2]
	remote := ""
	// Extract remote branch name from the square brackets
	if matches[5] != "" {
		remoteBracketContents := matches[5]
		remoteParts := strings.SplitN(remoteBracketContents[1:len(remoteBracketContents)-1], ":", 2)
		remote = strings.TrimSpace(remoteParts[0])
	}

	// Adjust spaces in commit message
	//msg := adjustSpace(matches[5])
	msg := matches[6]
	worktree := lastPathElement(matches[4])

	return &Branch{
		IsCurrent: isCurrent,
		Symbol:    symbol,
		Name:      branch,
		Remote:    remote,
		Message:   msg,
		Worktree:  worktree,
	}
}

// adjustSpace adjusts the string to have an effective length of 70 spaces, considering ASCII and non-ASCII characters
func adjustSpace(s string, maxLen int) string {
	length := strLen(s)
	if maxLen > commitMsgLength {
		maxLen = commitMsgLength
	}

	if maxLen <= 0 {
		return ""
	}
	if length >= maxLen {
		return chopRightByWidth(s, maxLen)
	}
	return padRightByWidth(s, maxLen)
}

func print(branches []*Branch) {
	max_symbol_width := 0
	max_branch_width := 0
	max_commit_len := 0
	for _, b := range branches {
		symbolWidth := symbolDisplayWidth(b)
		if symbolWidth > max_symbol_width {
			max_symbol_width = symbolWidth
		}
		branchName := buildBranchName(b)
		if strLen(branchName) > max_branch_width {
			max_branch_width = strLen(branchName)
		}
		msgWidth := strLen(b.Message)
		if msgWidth > max_commit_len {
			max_commit_len = msgWidth
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

	rawTermWidth := getTerminalWidth()
	termWidth := rawTermWidth - 1
	if termWidth < 0 {
		termWidth = 0
	}
	// 順番に出力
	for _, b := range branches {
		symbol := b.Symbol
		symbolWidth := symbolDisplayWidth(b)
		if symbolWidth < max_symbol_width {
			symbol += strings.Repeat(" ", max_symbol_width-symbolWidth)
		}
		name := padRightByWidth(buildBranchName(b), max_branch_width) + " "
		msg := adjustSpace(b.Message, max_commit_len)
		remote := b.Remote

		if termWidth > 0 {
			outLine := fmt.Sprintf("%s %s %s - %s", symbol, name, msg, remote)
			outWidth := strLen(outLine)
			if outWidth > termWidth {
				excess := outWidth - termWidth
				remoteWidth := strLen(remote)
				if remoteWidth > 0 {
					newRemoteWidth := remoteWidth - excess
					if newRemoteWidth < 0 {
						newRemoteWidth = 0
					}
					remote = chopRightByWidth(remote, newRemoteWidth)
				}
				outLine = fmt.Sprintf("%s %s %s - %s", symbol, name, msg, remote)
				outWidth = strLen(outLine)
				if outWidth > termWidth {
					excess = outWidth - termWidth
					msgWidth := strLen(msg)
					newMsgWidth := msgWidth - excess
					if newMsgWidth < 0 {
						newMsgWidth = 0
					}
					msg = chopRightByWidth(msg, newMsgWidth)
				}
			}
		}

		out := fmt.Sprintf("%s %s %s - %s\n", symbol, name, msg, remote)

		if b.IsCurrent {
			fg := getFg()
			c := color.New(fg)
			// c.Add(color.BgWhite)
			out = c.Sprint(out)
		}

		fmt.Print(out)
	}
}

func buildBranchName(b *Branch) string {
	if b.Worktree == "" {
		return b.Name
	}
	return fmt.Sprintf("%s (%s)", b.Name, b.Worktree)
}

func lastPathElement(path string) string {
	if path == "" {
		return ""
	}
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func getFg() color.Attribute {

	s := os.Getenv("GBRANCH_FG")
	switch s {
	case "red":
		return color.FgRed
	case "hired":
		return color.FgHiRed
	case "blue":
		return color.FgBlue
	case "hiblue":
		return color.FgHiBlue
	case "yellow":
		return color.FgYellow
	case "hiyellow":
		return color.FgHiYellow
	case "black":
		return color.FgBlack
	case "hiblack":
		return color.FgHiBlack
	}

	return color.FgRed
}

func getTerminalWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	s := os.Getenv("COLUMNS")
	if s == "" {
		return 0
	}
	width, err := strconv.Atoi(s)
	if err != nil || width <= 0 {
		return 0
	}

	return width
}

func strLen(s string) int {
	length := 0
	runes := []rune(s)
	for _, r := range runes {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Cf, r) {
			// Skip zero-width combining/format characters (e.g. variation selectors).
			continue
		}
		if unicode.IsPrint(r) && r < 128 {
			length += 1 // ASCII characters count as 1
		} else {
			length += 2 // Non-ASCII characters count as 2
		}
	}

	return length
}

func padRightByWidth(s string, width int) string {
	if width <= 0 {
		return s
	}
	currentWidth := strLen(s)
	if currentWidth >= width {
		return s
	}

	return s + strings.Repeat(" ", width-currentWidth)
}

func symbolDisplayWidth(b *Branch) int {
	if b.IsCurrent {
		// specialSymbol / GBRANCH_SYMBOL are treated as width 2
		if b.Symbol == "" {
			return 0
		}
		return 2
	}
	return strLen(b.Symbol)
}

func chopRightByWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	width := 0
	runes := []rune(s)
	for i, r := range runes {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Cf, r) {
			continue
		}
		if unicode.IsPrint(r) && r < 128 {
			width += 1
		} else {
			width += 2
		}
		if width > maxWidth {
			return string(runes[:i])
		}
	}

	return s
}
