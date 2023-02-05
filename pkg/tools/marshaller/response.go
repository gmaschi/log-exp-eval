package marshaller

import "encoding/json"

// Response marshals the obtained data to a response object.
func Response(data, res interface{}) error {
	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bData, &res)
	if err != nil {
		return err
	}

	return nil
}
