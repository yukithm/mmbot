package shell

import (
	"strings"

	"gopkg.in/readline.v1"
)

// Readline is extended version of chzyer/readline.
type Readline struct {
	*readline.Instance
}

// NewReadline returns new Readline with prompt.
func NewReadline(prompt string) (*Readline, error) {
	rl, err := readline.New(prompt)
	if err != nil {
		return nil, err
	}
	return &Readline{rl}, nil
}

// Readline reads STDIN and return it.
func (rl *Readline) Readline() (string, error) {
	line, err := rl.Instance.Readline()
	if err != nil {
		return "", err
	}

	origPrompt := rl.Config.Prompt
	defer rl.SetPrompt(origPrompt)

	for hasNextLine(line) {
		rl.SetPrompt("> ")
		next, err := rl.Instance.Readline()
		if err != nil {
			return "", err
		}
		if hasEscapedBreak(line) {
			line = strings.TrimSuffix(line, `\`)
		}
		line = line + "\n" + next
	}

	return line, nil
}

func hasEscapedBreak(line string) bool {
	r := strings.NewReplacer(`\\`, ``)
	tmp := r.Replace(line)

	return strings.HasSuffix(tmp, `\`)
}

func hasNextLine(line string) bool {
	if hasEscapedBreak(line) {
		return true
	}

	return inQuotes(line)
}

func inQuotes(line string) bool {
	var prev, current rune
	for _, c := range line {
		switch c {
		case '\'', '"', '`':
			if prev == '\\' {
				break
			}
			if current == 0 {
				current = c
			} else if current == c {
				current = 0
			}
		}
		prev = c
	}
	return current != 0
}
