package betterhandler

import (
	"net/http/httptest"
	"testing"
)

func TestString(t *testing.T) {
	want := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."

	request := httptest.NewRequest("POST", "/writestring", nil)
	responseRecorder := httptest.NewRecorder()

	handler := BetterHandler(func(c *Ctx) {
		c.String(want)
	})

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Body.String() != want {
		t.Errorf("Want '%s', got '%s'", want, responseRecorder.Body.String())
	}
}

// func TestBodyParser(t *testing.T) {
// 	handler := BetterHandler(func(c *Ctx) {
// 		var values struct {
// 			Key1 string  `form:"key1"`
// 			Key2 int32   `form:"key2"`
// 			Key3 float32 `form:"key3"`
// 		}

// 		err := c.BodyParser(&values)
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		fmt.Println(values)
// 	})
// 	http.ListenAndServe(":8080", handler)
// }
