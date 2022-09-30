package betterhandler

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalJson(t *testing.T) {
    give := `{"Test": "test"}`
    want := struct{Test string}{
        Test: "test",
    }
    
    request := httptest.NewRequest("GET", "/unmarshaljson", strings.NewReader(give))
    responseRecorder := httptest.NewRecorder()

    handler := BetterHandler(func(c Ctx) {
        var got struct{Test string}
        c.Req.UnmarshalJson(&got)
        fmt.Println()
        
        if cmp.Equal(want, got) == false {
            t.Errorf("Expected %v, got %v", want, got)
        }
    })
    handler.ServeHTTP(responseRecorder, request)
}
