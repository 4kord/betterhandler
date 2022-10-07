package betterhandler

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Context struct
type Ctx struct {
	rw http.ResponseWriter
	r  *http.Request
}

// Creates context from ResponseWriter and Request
func newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		rw: w,
		r:  r,
	}
}

// String writes String into responseWriter
func (c *Ctx) String(v string) error {
	c.rw.Header().Set("Content-Type", "text/plain")

	_, err := c.rw.Write([]byte(v))

	return err
}

// JSON writes JSON into responseWriter
func (c *Ctx) JSON(v any) error {
	c.rw.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = c.rw.Write(b)

	return err
}

// XML writes XML into responseWriter
func (c *Ctx) XML(v any) error {
	c.rw.Header().Set("Content-Type", "application/xml")

	b, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	_, err = c.rw.Write(b)

	return err
}

// BodyParser unmarhals request body into v
func (c *Ctx) BodyParser(v any) error {
	ctype := c.r.Header.Get("Content-Type")

	if strings.HasPrefix(ctype, "application/json") {
		b, err := io.ReadAll(c.r.Body)
		if err != nil {
			return err
		}

		return json.Unmarshal(b, v)
	} else if strings.HasPrefix(ctype, "application/xml") || strings.HasPrefix(ctype, "text/xml") {
		b, err := io.ReadAll(c.r.Body)
		if err != nil {
			return err
		}

		return xml.Unmarshal(b, v)
	} else if strings.HasPrefix(ctype, "multipart/form-data") {
		err := c.r.ParseMultipartForm(10 << 32)
		if err != nil {
			return err
		}

		reflectionTypePtr := reflect.TypeOf(v)
		reflectionValuePtr := reflect.ValueOf(v)

		if reflectionTypePtr.Kind() != reflect.Pointer {
			return fmt.Errorf("Expected kind 'pointer', got kind %s", reflectionTypePtr.Kind().String())
		}

		reflectionType := reflectionTypePtr.Elem()
		reflectionValue := reflectionValuePtr.Elem()

		if reflectionType.Kind() != reflect.Struct {
			return fmt.Errorf("Expected kind 'struct', got kind %s", reflectionType.Kind().String())
		}

		for i := 0; i < reflectionType.NumField(); i++ {
			field := reflectionType.Field(i)
			fieldValue := reflectionValue.Field(i)
			var formValue string
			if len(c.r.MultipartForm.Value[field.Tag.Get("form")]) != 0 {
				formValue = c.r.MultipartForm.Value[field.Tag.Get("form")][0]
			}

			switch field.Type.Kind() {
			case reflect.String:
				fieldValue.SetString(formValue)
			case reflect.Int, reflect.Int32, reflect.Int64:
				i, _ := strconv.Atoi(formValue)
				fieldValue.SetInt(int64(i))
			case reflect.Float32, reflect.Float64:
				f, _ := strconv.ParseFloat(formValue, 64)
				fieldValue.SetFloat(f)
			}
		}
	} else {
		return errors.New("Unsupported content type")
	}

	return nil
}

// Returns the base URL (protocol + host) as a string
func (c *Ctx) BaseURL() string {
	return c.r.URL.Scheme + "://" + c.r.URL.Host
}

// Context returns the request's context. To change the context, use WithContext.
func (c *Ctx) Context() context.Context {
	return c.r.Context()
}

// Sets cookie
func (c *Ctx) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.rw, cookie)
}

// Get cookie by key
func (c *Ctx) GetCookie(key string) (*http.Cookie, error) {
	cookie, err := c.r.Cookie(key)
	return cookie, err
}

// Get cookie value by key
func (c *Ctx) GetCookieValue(key string) (string, error) {
	cookie, err := c.r.Cookie(key)
	return cookie.Value, err
}

// Expire a client cookie (or all cookies if left empty)
func (c *Ctx) ClearCookie(key ...string) {
	if len(key) == 0 {
		for _, cookie := range c.r.Cookies() {
			newCookie := http.Cookie{
				Name:    cookie.Name,
				Value:   "",
				MaxAge:  -1,
				Expires: time.Now().Add(-100 * time.Hour),
			}

			http.SetCookie(c.rw, &newCookie)
		}
		return
	}

	for _, k := range key {
		cookie, err := c.r.Cookie(k)
		if err != nil {
			continue
		}
		newCookie := http.Cookie{
			Name:    cookie.Name,
			Value:   "",
			MaxAge:  -1,
			Expires: time.Now().Add(-100 * time.Hour),
		}

		http.SetCookie(c.rw, &newCookie)
	}
}
