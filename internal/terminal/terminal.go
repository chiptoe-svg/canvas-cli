package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// IsStdoutTerminal returns true if stdout is a terminal.
func IsStdoutTerminal() bool {
	return IsTerminal(int(os.Stdout.Fd()))
}

// IsStderrTerminal returns true if stderr is a terminal.
func IsStderrTerminal() bool {
	return IsTerminal(int(os.Stderr.Fd()))
}

// IsStdinTerminal returns true if stdin is a terminal.
func IsStdinTerminal() bool {
	return IsTerminal(int(os.Stdin.Fd()))
}

// SupportsColor returns true if the terminal likely supports color output.
// Checks stdout TTY and the NO_COLOR convention (https://no-color.org/).
func SupportsColor() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	return IsStdoutTerminal()
}

// ReadSecret prints prompt to stderr and reads one line WITHOUT echoing when stdin is a
// terminal, so a pasted token/secret never lands in scrollback. On a pipe it falls back to a
// normal line read so scripts still work. It reads a full line (long API tokens / OAuth codes
// are fine) and trims surrounding whitespace — unlike fmt.Scanln, which echoes the input and
// stalls on long terminal pastes.
func ReadSecret(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	if IsStdinTerminal() {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// Width returns the width of the terminal connected to stdout.
// Returns 80 as a fallback if the width cannot be determined.
func Width() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
