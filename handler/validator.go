package handler

import (
	"fmt"
	"reflect"
	"regexp"
	"stlab.itechart-group.com/go/food_delivery/authentication_service/model"
	"strings"
)

const tagName = "validate"

var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

type Validator interface {
	Validate(interface{}) error
}

type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) error {
	return nil
}

type EmailValidator struct {
}

func (v EmailValidator) Validate(val interface{}) error {
	if !mailRe.MatchString(val.(string)) {
		return fmt.Errorf("emailValidator: it is not a valid email address")
	}
	return nil
}

type RoleIdValidator struct{}

func (v RoleIdValidator) Validate(val interface{}) error {
	if val.(int) < 0 || val.(int) == 0 {
		return fmt.Errorf("roleIdValidator: roleId must be positive integer")
	}
	//check for superadmin creation
	if val.(int) == 6 {
		return fmt.Errorf("you don't have enaugh rights to create user with such a role")
	}
	return nil
}

type PasswordValidator struct {
}

func (v PasswordValidator) Validate(val interface{}) error {
	if (len(val.(string)) < 8 || len(val.(string)) > 15) && len(val.(string)) != 0 {
		return fmt.Errorf("passwordValidator: the length of the password should be between 8 to 15 characters")
	}
	if len(val.(string)) == 0 {
		return nil
	}
	for _, i := range val.(string) {
		if string(i) == " " {
			return fmt.Errorf("passwordValidator: password should not contain any space")
		}
	}
	var num int
	var upper int
	var lower int
	var special int
	for _, i := range val.(string) {
		var overlap = false
		for _, j := range model.PasswordNumber {
			if i == j {
				overlap = true
				num = num + 1
				break
			}
		}
		if overlap == false {
			for _, j := range model.PasswordLower {
				if i == j {
					overlap = true
					lower = lower + 1
					break
				}
			}
		}
		if overlap == false {
			for _, j := range model.PasswordUpper {
				if i == j {
					overlap = true
					upper = upper + 1
					break
				}
			}
		}
		if overlap == false {
			for _, j := range model.PasswordSpecial {
				if i == j {
					overlap = true
					special = special + 1
					break
				}
			}
		}
		if overlap == false {
			return fmt.Errorf("passwordValidator: the password must contain at least one digit(0-9), " +
				"one lowercase letter(a-z), one uppercase letter(A-Z), one special character (@,#,%%,&,!,$)")
		} else {
			overlap = false
		}
	}
	if num == 0 || lower == 0 || upper == 0 || special == 0 {
		num = 0
		lower = 0
		upper = 0
		special = 0
		return fmt.Errorf("passwordValidator: the password must contain at least one digit(0-9), " +
			"one lowercase letter(a-z), one uppercase letter(A-Z), one special character (@,#,%%,&,!,$)")
	} else {
		num = 0
		lower = 0
		upper = 0
		special = 0
	}
	return nil
}

func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "email":
		return EmailValidator{}
	case "password":
		return PasswordValidator{}
	case "roleId":
		return RoleIdValidator{}
	}

	return DefaultValidator{}
}

func ValidateStruct(s interface{}) map[string]string {
	var errs = make(map[string]string)
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		validator := getValidatorFromTag(tag)
		err := validator.Validate(v.Field(i).Interface())
		if err != nil {
			errs[v.Type().Field(i).Name] = err.Error()
		}
	}
	return errs
}
