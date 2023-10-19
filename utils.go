// version : V1
package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type DefaultValidator struct{}

const tagName = "validate"

func panicReceiver() []error {
	errs := []error{}
	if r := recover(); r != nil {
		fmt.Println("Panic occurred:", r)
		errs = append(errs, fmt.Errorf("Panic occurred: %v", r))
		_, file, line, ok := runtime.Caller(0)
		if ok {
			fmt.Println("File:", file, "Line:", line)
		}
	}
	return errs
}

func ValidateISO8583Message(s interface{}) []error {
	errs := []error{}
	v := reflect.ValueOf(s)

	defer panicReceiver()

	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		validator := GetValidatorFromTag(tag)
		fieldName := v.Type().Field(i).Name
		valid, err := validator.Validate(v.Field(i).Interface(), fieldName)
		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}
	return errs
}

func GetValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "number":
		validator := NumberValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "length=%d", &validator.Length)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "length=%d", &validator.Length)
		return validator
	}
	return DefaultValidator{}
}

type Validator interface {
	Validate(interface{}, string) (bool, error)
}

func (v DefaultValidator) Validate(val interface{}, fieldName string) (bool, error) {
	return true, nil
}

type NumberValidator struct {
	Length int
}

func RemovePoint(str string) string {
	return strings.Replace(str, ".", "", -1)
}

func (v NumberValidator) Validate(val interface{}, fieldName string) (bool, error) {
	num := val.(int)
	str := fieldName
	valueNumber := intToString(v.Length)
	valueRequest := intToString(num)

	if len(valueRequest) != len(valueNumber) {
		return false, fmt.Errorf("El campo %v debe tener %v", str, v.Length)
	}

	return true, nil
}

type StringValidator struct {
	Length int
}

func (v StringValidator) Validate(val interface{}, fieldName string) (bool, error) {

	valueRequest := val.(string)

	if fieldName == "AmountTransaction_004" {
		valueRequest = RemovePoint(valueRequest)
	}

	if len(valueRequest) == 0 {
		return false, fmt.Errorf("no puede estar vacio")
	}
	if len(valueRequest) != v.Length {
		return false, fmt.Errorf("debe tener %v", v.Length)
	}

	return true, nil
}

// floatToString converts a float64 number to string.
//
// It takes a float64 number as a parameter and returns the string representation of that number.
// The return type is string.
func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// hexToBinary converts a hexadecimal string to a binary string.
//
// hexStr: the hexadecimal string to convert.
// Returns the binary representation of the input hexadecimal string.
func HxToBinary(hexStr string) string {
	hexBytes := []byte(hexStr)
	binaryStr := ""
	for _, hexChar := range hexBytes {
		value, err := strconv.ParseInt(string(hexChar), 16, 0)
		if err != nil {
			panic(err)
		}
		binaryStr += fmt.Sprintf("%04b", value)
	}
	return binaryStr
}

// timeToString converts a time.Time value to a string representation.
//
// It takes a time.Time parameter and returns a string.
func TimeToString(time time.Time) string {
	return time.Format("HHmmss")
}

// validarTime validates a time string.
//
// timeStr: the string to be validated.
// bool: true if the time string is valid, false otherwise.
func ValidarTime(timeStr string) bool {
	regex := regexp.MustCompile("^([0-2]?[0-9]):([0-5]?[0-9]):([0-5]?[0-9])$")
	return regex.MatchString(timeStr)
}

// intToString converts an integer value to a string.
//
// It takes an integer value as a parameter and returns the corresponding string representation.
// The returned string will not contain any leading or trailing spaces.
func intToString(value int) string {
	return strconv.Itoa(value)
}

// convertir string a int
// convertirStringToInt converts a string to an integer.
//
// It takes a string as a parameter and returns the corresponding integer value.
func StringToInt(str string) (int, error) {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
