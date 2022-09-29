package betterhandler

import (
	"net/http"
)

// Custom response struct
type response struct {
    http.ResponseWriter
}

// Custom request struct
type request struct {
    *http.Request
}

// Context struct
type Ctx struct {
    Res response
    Req request
}

// Creates context from ResponseWriter and Request
func newCtx(w http.ResponseWriter, r *http.Request) Ctx {
    return Ctx{
        Res: response{
            ResponseWriter: w,
        },
        Req: request{
            Request: r,
        },
    }
}

// BetterHandler signature
type BetterHandler func(Ctx)

// ServeHTTP method to implement Handler interface
func (h BetterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h(newCtx(w, r))
}
