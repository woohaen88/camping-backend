package enums

import "errors"

type Role string

const (
	Admin  Role = "admin"
	Client Role = "client"
)

func (r *Role) Check() error {
	value := *r
	switch value {
	case Admin, Client:
		return nil
	}
	return errors.New("invalid product type value")

}
