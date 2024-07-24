package cmd

import (
	"github.com/itchyny/gojq"
)

func jq(input interface{}, filter string, output func(interface{}) error) error {
	query, err := gojq.Parse(filter)
	if err != nil {
		return err
	}

	iter := query.Run(input)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil { //nolint:errorlint
				break
			}

			return err
		}

		if err := output(v); err != nil {
			return err
		}
	}

	return nil
}
