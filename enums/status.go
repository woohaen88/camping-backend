package enums

import "errors"

type Status string

const (
	StatusOk    Status = "OK"
	StatusCant  Status = "CANT"
	StatusOther Status = "OTHER"
)

func (status *Status) String() error {
	statusValue := *status
	switch *status {
	case StatusOk, StatusCant, StatusOther:
		*status = statusValue
		return nil
	}

	return errors.New("invalid product type value")
}
