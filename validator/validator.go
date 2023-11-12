package validator

import (
	"reflect"
	"regexp"

	"strings"

	"github.com/cgalvisleon/elvis/console"
	. "github.com/cgalvisleon/elvis/utility"
	"gopkg.in/validator.v2"
)

func validateChannel(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return console.AlertF("validateChannel only validates strings")
	}
	channelLower := strings.ToLower(st.String())
	if channelLower != "whatsapp" && channelLower != "ivr" {
		return console.AlertF("Invalid (Only accept WhatsApp or IVR values)")
	}

	return nil
}

func validateOnlyDigit(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return console.AlertF("validateChannel only validates strings")
	}
	matched, _ := regexp.MatchString(`\D+`, st.String())
	if matched {
		return console.AlertF("Invalid (Only accept digits)")
	}

	return nil
}

func validatePhoneNumber(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return console.AlertF("validateChannel only validates strings")
	}
	matched, _ := regexp.MatchString(`^[0-9]{10,15}$`, st.String())
	if !matched {
		return console.AlertF("Invalid (Only accept between 10 and 15 digits)")
	}

	return nil
}

func ValidateOnlyDigit(value string, min, max int, errorLabel string) error {
	if reflect.TypeOf(value).Kind() != reflect.String {
		return console.AlertF("Invalid phone_number must be a string")
	}

	matched, _ := regexp.MatchString(`\D+`, value)
	if matched {
		return console.AlertF("Invalid (%s Only accept digits)", errorLabel)
	}

	if min != -1 && len(value) < min {
		return console.AlertF("Invalid (%s is less than %d digits)", errorLabel, min)
	}

	if max != -1 && len(value) > max {
		return console.AlertF("Invalid (%s is upper than %d digits)", errorLabel, max)
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	if reflect.TypeOf(phoneNumber).Kind() != reflect.String {
		return console.AlertF("Invalid phone_number must be a string")
	}
	matched, _ := regexp.MatchString(`^[0-9]{10,15}$`, phoneNumber)
	if !matched {
		return console.AlertF("Invalid phone_number (Only accept between 10 and 15 digits)")
	}

	return nil
}

func ValidateAccount(phoneNumber string) error {
	if reflect.TypeOf(phoneNumber).Kind() != reflect.String {
		return console.AlertF("Invalid phone_number must be a string")
	}
	matched, _ := regexp.MatchString(`^[0-9]{8,15}$`, phoneNumber)
	if !matched {
		return console.AlertF("Invalid account (Only accept between 8 and 15 digits)")
	}

	return nil
}

func validateUUID(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if param == "empty" && st.String() == "" {
		return nil
	}

	if st.Kind() != reflect.String {
		return console.AlertF("validateChannel only validates strings")
	}
	matched, _ := regexp.MatchString(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, st.String())
	if !matched {
		return console.AlertF("Invalid UUID")
	}

	return nil
}

func ValidateID(id string) error {
	if reflect.TypeOf(id).Kind() != reflect.String {
		return console.AlertF("Invalid ID")
	}
	matched, _ := regexp.MatchString(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, id)
	if !matched {
		return console.AlertF("Invalid ID")
	}

	return nil
}

func validRequired(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if !ValidStr(st.String(), 0, []string{"", "-1"}) {
		return console.AlertF("Required value")
	}
	return nil
}

func validName(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return console.AlertF("only accept string values")
	}
	matched, _ := regexp.MatchString(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]+$`, st.String())
	if !matched {
		return console.AlertF("Invalid (Not accept special characters)")
	}

	return nil
}

func validateIntType(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.Int {
		return console.AlertF("must be integer type")
	}

	return nil
}

func InitValidators() {
	validator.SetValidationFunc("onlyDigit", validateOnlyDigit)
	validator.SetValidationFunc("validatePhoneNumber", validatePhoneNumber)
	validator.SetValidationFunc("validateChannel", validateChannel)
	validator.SetValidationFunc("validateUUID", validateUUID)
	validator.SetValidationFunc("validateName", validName)
	validator.SetValidationFunc("required", validRequired)
	validator.SetValidationFunc("intType", validateIntType)
}
