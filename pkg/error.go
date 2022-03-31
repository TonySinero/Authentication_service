package pkg

import "errors"

const (
	EmailDoesNotExist = "user with this email does not exist"
)

var ErrorEmailDoesNotExist = errors.New(EmailDoesNotExist)
