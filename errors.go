package cuttly

import (
	"errors"
)

// checkErrorCode associates the respons error code to a specific message
// Important: the messages are the original ones from the cuttly-api page
func checkErrorCode(errorCode int, isURL bool) error {
	switch errorCode {
	case 0:
		return errors.New("This shortened link does not exist")
	case 1:
		if isURL {
			return errors.New("The shortened link comes from the domain that shortens the link, i.e. the link has already been shortened")
		}
		break
	case 2:
		if isURL {
			return errors.New("The entered link is not a link")
		}
		return errors.New("Invalid API key")
	case 3:
		return errors.New("The preferred link name is already taken")
	case 4:
		return errors.New("Invalid API key")
	case 5:
		return errors.New("The link has not passed the validation. Includes invalid characters")
	case 6:
		return errors.New("The link provided is from a blocked domain")
	}

	return nil
}
