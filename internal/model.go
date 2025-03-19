package internal

import (
	"encoding/json"
	"errors"
)

type CheckingPolicy struct {
	ID         string          `json:"id"`
	Combination json.RawMessage `json:"combination"`
	Policy     PolicyData `json:"policy_data"`
}

type PolicyData struct {
	Type   string          `json:"type"`
	Params []string `json:"params"`
}

func (d *PolicyData) Scan(src interface{}) (err error) {
	var data PolicyData
	
	switch src.(type) {
	case string:
		err = json.Unmarshal([]byte(src.(string)), &data)
	case []byte:
		err = json.Unmarshal(src.([]byte), &data)
	default:
		return errors.New("Incompatible type for Skills")
	}

	if err != nil {
		return
	}
	*d = data
	return nil
}
