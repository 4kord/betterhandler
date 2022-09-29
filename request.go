package betterhandler

import (
	"encoding/json"
	"io"
)

// UnmarshalJson unmarhals reqeust body into v
func (r request) UnmarshalJson(v any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
