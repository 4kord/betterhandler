package betterhandler

import (
	"net/http/httptest"
	"testing"
)

func TestWriteString(t *testing.T) {
    cases := []struct {
        name string
        want string
    }{
        {
            name: "test_1",
            want: "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit...",
        },
        {
            name: "test_2",
            want: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. In semper hendrerit ligula eget posuere. Nullam non aliquam ante, eu interdum ipsum. Ut fringilla tempor semper. Donec ullamcorper pulvinar ante, in tristique sem porttitor id. Maecenas consectetur ipsum eget fringilla congue. Suspendisse et arcu tortor. Integer viverra quam metus, vitae viverra velit dictum eu. Sed id dapibus magna, ac porttitor lacus. In sed ligula vestibulum, tincidunt orci sed, pharetra nisi. Integer ultrices diam sapien, nec viverra magna consectetur ac. Morbi ac massa massa. Morbi eget purus quis enim pulvinar tincidunt luctus non augue. Suspendisse finibus lobortis ultrices. Phasellus ut imperdiet nulla, eget ornare quam.",
        },
    }

    for _, tt := range cases {
        t.Run(tt.name, func(t *testing.T) {
            request := httptest.NewRequest("POST", "/writestring", nil)
            responseRecorder := httptest.NewRecorder()

            handler := BetterHandler(func(c Ctx) {
                c.Res.WriteString(tt.want)
            })
            handler.ServeHTTP(responseRecorder, request)

            if responseRecorder.Body.String() != tt.want {
                t.Errorf("Want '%s', got '%s'", tt.want, responseRecorder.Body.String())
            }
        })
    }
}
