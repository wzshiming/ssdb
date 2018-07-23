package ssdb

import (
	"fmt"
)

var errIsEmpty = fmt.Errorf("error: respone is empty")

func ResultProcessing(resp Values, err error) (Values, error) {

	if err != nil {
		return nil, fmt.Errorf("error: %v", err.Error())
	}

	if len(resp) < 1 {
		return nil, errIsEmpty
	}

	if resp[0].Equal(ok) {
		return resp[1:], nil
	}

	if resp[0].Equal(notFound) {
		return nil, nil
	}

	if resp[0].Equal(clientError) {
		return nil, fmt.Errorf("error: %v", resp[1].String())
	}

	return nil, fmt.Errorf("error: parameter %v", resp[0])
}