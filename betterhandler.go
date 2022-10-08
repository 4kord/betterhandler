package betterhandler

import (
	"net/http"
)

// BetterHandler signature
type BH func(*Ctx)

// ServeHTTP method to implement Handler interface
func (h BH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(newCtx(w, r))
}

type Map map[string]interface{}
