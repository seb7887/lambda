package lambda

import (
	"encoding/json"
)

func adaptJSON[I any](raw json.RawMessage) (*I, error) {
	var in I
	if err := json.Unmarshal(raw, &in); err != nil {
		return nil, err
	}
	return &in, nil
}
