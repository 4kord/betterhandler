package betterhandler

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type stringIntFloat64 struct {
	Key1 string  `json:"key1" xml:"key1" form:"key1"`
	Key2 int     `json:"key2" xml:"key2" form:"key2"`
	Key3 float64 `json:"key3" xml:"key3" form:"key3"`
}

func TestString(t *testing.T) {
	want := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit..."

	request := httptest.NewRequest("POST", "/string", nil)
	responseRecorder := httptest.NewRecorder()

	handler := BetterHandler(func(c *Ctx) {
		c.String(want)
	})

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Body.String() != want {
		t.Errorf("Want '%s', got '%s'", want, responseRecorder.Body.String())
	}
}

func TestJSON(t *testing.T) {
	give := stringIntFloat64{
		Key1: "Value1",
		Key2: 123,
		Key3: 123.123,
	}
	want := `{"key1":"Value1","key2":123,"key3":123.123}`

	request := httptest.NewRequest("POST", "/json", nil)
	responseRecorder := httptest.NewRecorder()

	handler := BetterHandler(func(c *Ctx) {
		c.JSON(give)
	})

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Body.String() != want {
		t.Error("error")
	}
}

func TestXML(t *testing.T) {
	give := stringIntFloat64{
		Key1: "Value1",
		Key2: 123,
		Key3: 123.123,
	}
	want := `<stringIntFloat64><key1>Value1</key1><key2>123</key2><key3>123.123</key3></stringIntFloat64>`

	request := httptest.NewRequest("POST", "/json", nil)
	responseRecorder := httptest.NewRecorder()

	handler := BetterHandler(func(c *Ctx) {
		c.XML(give)
	})

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Body.String() != want {
		t.Error("error")
	}
}

func TestBodyParser(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		contentType string
		body        string
		want        interface{}
	}{
		{
			name:        "application/json_string_int_float64",
			method:      "GET",
			contentType: "application/json",
			body:        `{"key1":"Value1","key2":10,"key3":12.12}`,
			want: stringIntFloat64{
				Key1: "Value1",
				Key2: 10,
				Key3: 12.12,
			},
		},
		{
			name:        "application/xml_string_int_float64",
			method:      "GET",
			contentType: "application/xml",
			body:        `<stringIntFloat64><key1>Value1</key1><key2>20</key2><key3>22.12</key3></stringIntFloat64>`,
			want: stringIntFloat64{
				Key1: "Value1",
				Key2: 20,
				Key3: 22.12,
			},
		},
		{
			name:        "multipart/form-data_string_int_float64",
			method:      "GET",
			contentType: `multipart/form-data; boundary="boundary"`,
			body: `--boundary
Content-Disposition: form-data; name="key1"

Value1
--boundary
Content-Disposition: form-data; name="key2"

30
--boundary
Content-Disposition: form-data; name="key3"

32.12
--boundary--`,
			want: stringIntFloat64{
				Key1: "Value1",
				Key2: 30,
				Key3: 32.12,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/bodyparser", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			handler := BetterHandler(func(c *Ctx) {
				var got stringIntFloat64

				err := c.BodyParser(&got)
				if err != nil {
					t.Error(err)
				}

				equal := reflect.DeepEqual(tt.want, got)
				if !equal {
					t.Errorf("Expected %v, got %v", tt.want, got)
				}
			})

			handler.ServeHTTP(responseRecorder, request)
		})
	}
}
