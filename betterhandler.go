package betterhandler

import (
	"net/http"
)

// BetterHandler signature
type BetterHandler func(*Ctx)

// ServeHTTP method to implement Handler interface
func (h BetterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(newCtx(w, r))
}

type Map map[string]interface{}
