package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/uget/uget/core"
)

type cliPrompter struct {
	prefix    string
	overrides map[string]string
}

func newCliPrompter(prefix string, overrides map[string]string) *cliPrompter {
	return &cliPrompter{prefix, overrides}
}

func (c cliPrompter) Get(fields []core.Field) (map[string]string, error) {
	reader := bufio.NewReader(os.Stdin)
	values := map[string]string{}
	for _, field := range fields {
		if value, ok := c.overrides[field.Key]; ok {
			values[field.Key] = value
		} else {
			var deftext string
			if field.Value != "" {
				deftext = fmt.Sprintf(" (%s)", field.Value)
			}
			fmt.Printf("[%s] %s%s: ", c.prefix, field.Display, deftext)
			var entered string
			if field.Sensitive {
				t, err := gopass.GetPasswdMasked()
				if err != nil {
					return nil, err
				}
				entered = string(t)
			} else {
				line, err := reader.ReadString('\n')
				if err != nil {
					c.Error(err.Error())
					return nil, err
				}
				entered = strings.TrimSpace(string(line))
			}
			if entered == "" {
				entered = field.Value
			}
			values[field.Key] = entered
		}
	}
	return values, nil
}

func (c cliPrompter) Error(display string) {
	fmt.Fprintf(os.Stderr, "[%s] Error: %s\n", c.prefix, display)
}

func (c cliPrompter) Success() {
	fmt.Fprintf(os.Stderr, "[%s] Success.\n", c.prefix)
}
