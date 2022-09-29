package betterhandler

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestWriteString(t *testing.T) {
    for i := 1; i <= 10; i++ {
        randomString := RandStringBytes(5, 50)
        testName := fmt.Sprintf("test_%v_with_test_data_'%v'", i, randomString)
        t.Run(testName, func(t *testing.T) {
            request := httptest.NewRequest("POST", "/writestring", nil)
            responseRecorder := httptest.NewRecorder()

            handler := BetterHandler(func(c Ctx) {
                c.Res.WriteString(randomString)
            })
            handler.ServeHTTP(responseRecorder, request)

            if responseRecorder.Body.String() != randomString {
                t.Errorf("Want '%s', got '%s'", randomString, responseRecorder.Body.String())
            }
        })
    }
}
