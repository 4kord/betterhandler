package betterhandler

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type stringIntFloat64 struct {
	Key1 string  `json:"key1" xml:"key1" form:"key1"`
	Key2 int     `json:"key2" xml:"key2" form:"key2"`
	Key3 float64 `json:"key3" xml:"key3" form:"key3"`
}

type file struct {
	Key1 []*multipart.FileHeader `form:"key1"`
}

func TestString(t *testing.T) {
	cases := []struct {
		name string
		str  string
	}{
		{
			name: "test_1",
			str:  "String",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", nil)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				c.String(tt.str)
			})

			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Body.String() != tt.str {
				t.Errorf("Want '%s', got '%s'", tt.str, responseRecorder.Body.String())
			}
		})
	}
}

func TestJSON(t *testing.T) {
	cases := []struct {
		name string
		give interface{}
		want string
	}{
		{
			name: "test_1",
			give: stringIntFloat64{
				Key1: "Value1",
				Key2: 123,
				Key3: 123.123,
			},
			want: `{"key1":"Value1","key2":123,"key3":123.123}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", nil)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				c.JSON(tt.give)
			})

			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Body.String() != tt.want {
				t.Errorf("want %s, got %s", responseRecorder.Body.String(), tt.want)
			}
		})
	}
}

func TestXML(t *testing.T) {
	cases := []struct {
		name string
		give interface{}
		want string
	}{
		{
			name: "test_1",
			give: stringIntFloat64{
				Key1: "Value1",
				Key2: 123,
				Key3: 123.123,
			},
			want: `<stringIntFloat64><key1>Value1</key1><key2>123</key2><key3>123.123</key3></stringIntFloat64>`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", nil)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				c.XML(tt.give)
			})

			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Body.String() != tt.want {
				t.Errorf("want %s, got %s", responseRecorder.Body.String(), tt.want)
			}
		})
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
--boundary
Content-Disposition: form-data; name="key4"; filename="example.txt"

Value2
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
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				var got stringIntFloat64

				err := c.BodyParser(&got)
				if err != nil {
					t.Error(err)
				}

				equal := reflect.DeepEqual(tt.want, got)
				if !equal {
					t.Errorf("want %v, got %v", tt.want, got)
				}
			})

			handler.ServeHTTP(responseRecorder, request)
		})
	}
	t.Run("file_test", func(t *testing.T) {
		body := `--boundary
Content-Disposition: form-data; name="key1"; filename="example.txt"

example
--boundary--`
		request := httptest.NewRequest("GET", "/", strings.NewReader(body))
		request.Header.Set("Content-Type", "multipart/form-data; boundary=\"boundary\"")
		responseRecorder := httptest.NewRecorder()

		handler := BH(func(c *Ctx) {
			var got file

			err := c.BodyParser(&got)
			if err != nil {
				t.Error(err)
			}

			if got.Key1 == nil {
				t.Error("Expected []*multipart.FileHeader, got nil")
			}
		})

		handler.ServeHTTP(responseRecorder, request)
	})
}

func TestBaseURL(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "test_1",
			url:  "https://example.com/test",
			want: "https://example.com",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			request.URL, _ = url.Parse(tt.url)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				if c.BaseURL() != tt.want {
					t.Errorf("want %s, got %s", tt.want, c.BaseURL())
				}
			})

			handler.ServeHTTP(responseRecorder, request)
		})
	}
}

func TestSetCookie(t *testing.T) {
	cases := []struct {
		name   string
		cookie *http.Cookie
	}{
		{
			name: "test_1",
			cookie: &http.Cookie{
				Name:  "Cookie1",
				Value: "Value1",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				c.SetCookie(tt.cookie)
			})

			handler.ServeHTTP(responseRecorder, request)

			if tt.cookie.Name != responseRecorder.Result().Cookies()[0].Name || tt.cookie.Value != responseRecorder.Result().Cookies()[0].Value {
				t.Errorf("%v and %v are not equal", tt.cookie, responseRecorder.Result().Cookies()[0])
			}
		})
	}
}

func TestGetCookie(t *testing.T) {
	cases := []struct {
		name   string
		cookie *http.Cookie
	}{
		{
			name: "test_1",
			cookie: &http.Cookie{
				Name:  "Cookie1",
				Value: "Value1",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			request.AddCookie(tt.cookie)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				got, err := c.GetCookie(tt.cookie.Name)
				if err != nil {
					t.Error(err)
				}

				if tt.cookie.Name != got.Name || tt.cookie.Value != got.Value {
					t.Errorf("want %v, got %v", tt.cookie, got)
				}
			})

			handler.ServeHTTP(responseRecorder, request)
		})
	}
}

func TestGetCookieValue(t *testing.T) {
	cases := []struct {
		name   string
		cookie *http.Cookie
		want   string
	}{
		{
			name: "test_1",
			cookie: &http.Cookie{
				Name:  "Cookie1",
				Value: "Value1",
			},
			want: "Value1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			request.AddCookie(tt.cookie)
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				got, err := c.GetCookieValue(tt.cookie.Name)
				if err != nil {
					t.Error(err)
				}

				if tt.cookie.Value != got {
					t.Errorf("want %v, got %v", tt.cookie.Value, got)
				}
			})

			handler.ServeHTTP(responseRecorder, request)
		})
	}
}

func TestClearCookie(t *testing.T) {
	cases := []struct {
		name    string
		cookies []*http.Cookie
		clear   []string
	}{
		{
			name: "test_1",
			cookies: []*http.Cookie{
				{
					Name:  "Cookie1",
					Value: "Value1",
				},
				{
					Name:  "Cookie2",
					Value: "Value2",
				},
			},
			clear: []string{"Cookie1"},
		},
		{
			name:    "test_2",
			cookies: []*http.Cookie{},
			clear:   []string{},
		},
		{
			name: "test_3",
			cookies: []*http.Cookie{
				{
					Name:  "Cookie1",
					Value: "Value1",
				},
				{
					Name:  "Cookie2",
					Value: "Value2",
				},
			},
			clear: []string{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			for _, cookie := range tt.cookies {
				request.AddCookie(cookie)
			}
			responseRecorder := httptest.NewRecorder()

			handler := BH(func(c *Ctx) {
				c.ClearCookie(tt.clear...)
			})

			handler.ServeHTTP(responseRecorder, request)

			if len(tt.clear) != 0 && len(responseRecorder.Result().Cookies()) != len(tt.clear) {
				t.Errorf("Response cookies len != clear cookies len")
			}

			if len(tt.clear) == 0 && len(responseRecorder.Result().Cookies()) != len(tt.cookies) {
				t.Errorf("Response cookies len != cookies len")
			}

			for _, cookie := range responseRecorder.Result().Cookies() {
				if cookie.Value != "" {
					t.Errorf("| Value | Want empty, got %s", cookie.Value)
				}
				if !cookie.Expires.Before(time.Now().Add(-100*time.Hour)) || cookie.Expires.Equal(time.Now().Add(-100*time.Hour)) {
					t.Errorf("| Expires | Want < time.Now, got %s", cookie.Expires)
				}
				if cookie.MaxAge != -1 {
					t.Errorf("| MaxAge | Want < 0, got %d", cookie.MaxAge)
				}
			}
		})
	}
}
