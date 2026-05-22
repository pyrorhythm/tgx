package inline

import (
	"strings"

	"pyrorhythm.dev/fn/opt"
)

// ParsedQuery is a command name plus arguments (everything after the first token).
type ParsedQuery struct {
	Command string
	Args    string
}

// CommandRule describes a known inline command and minimum argument length.
type CommandRule struct {
	Name       string
	MinArgsLen int
}

// Parser splits "command args..." into command and args.
// If the first token is not a known command, the full query is routed to the default command.
type Parser struct {
	byName         map[string]CommandRule
	defaultCommand CommandRule
}

// NewParser registers known commands and a default for unrecognized first tokens.
func NewParser(defaultCommand CommandRule, commands ...CommandRule) *Parser {
	p := &Parser{
		byName:         make(map[string]CommandRule, len(commands)),
		defaultCommand: defaultCommand,
	}
	for _, cmd := range commands {
		p.byName[strings.ToLower(cmd.Name)] = cmd
	}
	return p
}

// Parse returns Nil for empty queries or when argument length rules are not met.
func (p *Parser) Parse(query string) opt.Of[ParsedQuery] {
	query = strings.TrimSpace(query)
	if query == "" {
		return opt.Nil[ParsedQuery]()
	}

	slots := strings.Fields(query)
	name := strings.ToLower(slots[0])
	args := ""
	if len(slots) > 1 {
		args = strings.Join(slots[1:], " ")
	}

	if rule, ok := p.byName[name]; ok {
		if len(args) < rule.MinArgsLen {
			return opt.Nil[ParsedQuery]()
		}
		return opt.SomeAny(ParsedQuery{Command: rule.Name, Args: args})
	}

	if len(query) < p.defaultCommand.MinArgsLen {
		return opt.Nil[ParsedQuery]()
	}
	return opt.SomeAny(ParsedQuery{Command: p.defaultCommand.Name, Args: query})
}
