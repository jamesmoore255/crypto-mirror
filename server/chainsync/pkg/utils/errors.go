package utils

import "errors"

func GenericError() error {
	return errors.New("Internal server error")
}
