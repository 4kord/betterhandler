package betterhandler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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
func (ctx *Ctx) String(v string) error {
	ctx.rw.Header().Set("Content-Type", "text/plain")

	_, err := ctx.rw.Write([]byte(v))

	return err
}

// JSON writes JSON into responseWriter
func (ctx *Ctx) JSON(v any) error {
	ctx.rw.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = ctx.rw.Write(b)

	return err
}

// XML writes XML into responseWriter
func (ctx *Ctx) XML(v any) error {
	ctx.rw.Header().Set("Content-Type", "application/xml")

	b, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	_, err = ctx.rw.Write(b)

	return err
}

// BodyParser unmarhals request body into v
func (ctx *Ctx) BodyParser(v any) error {
	ctype := ctx.r.Header.Get("Content-Type")

	if strings.HasPrefix(ctype, "application/json") {
		b, err := io.ReadAll(ctx.r.Body)
		if err != nil {
			return err
		}

		return json.Unmarshal(b, v)
	} else if strings.HasPrefix(ctype, "application/xml") || strings.HasPrefix(ctype, "text/xml") {
		b, err := io.ReadAll(ctx.r.Body)
		if err != nil {
			return err
		}

		return xml.Unmarshal(b, v)
	} else if strings.HasPrefix(ctype, "multipart/form-data") {
		err := ctx.r.ParseMultipartForm(10 << 32)
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
			if len(ctx.r.MultipartForm.Value[field.Tag.Get("form")]) != 0 {
				formValue = ctx.r.MultipartForm.Value[field.Tag.Get("form")][0]
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
	}

	return nil
}
