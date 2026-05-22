package inline

import (
	"github.com/mymmrac/telego"

	"pyrorhythm.dev/tgx"
)

// QueryHandler handles a parsed inline query for one command.
type QueryHandler func(*tgx.Ctx, telego.InlineQuery, string) error

// Router dispatches inline queries to registered handlers.
type Router struct {
	parser   *Parser
	handlers map[string]QueryHandler
}

// NewRouter creates a router that uses parser for query parsing.
func NewRouter(parser *Parser) *Router {
	return &Router{
		parser:   parser,
		handlers: make(map[string]QueryHandler),
	}
}

// On registers a handler for a command name.
func (r *Router) On(command string, h QueryHandler) {
	r.handlers[command] = h
}

// Handle parses the query and runs the matching handler.
func (r *Router) Handle(c *tgx.Ctx, q telego.InlineQuery) error {
	parsed := r.parser.Parse(q.Query)
	if !parsed.Valid() {
		return nil
	}

	h, ok := r.handlers[parsed.Val().Command]
	if !ok {
		return nil
	}
	return h(c, q, parsed.Val().Args)
}
