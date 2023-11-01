package common

import "log"

func CheckErr(e error) error {
	if e != nil {
		log.Println(e.Error())
	}
	return nil
}
