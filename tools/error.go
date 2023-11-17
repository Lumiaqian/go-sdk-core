package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type CommonError struct {
	apiName string
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (c *CommonError) Error() string {
	return fmt.Sprintf("%s Error, errcode=%d, errmsg=%s", c.apiName, c.ErrCode, c.ErrMsg)
}

// NewCommonError creates a new CommonError with the given API name, error code, and error message.
// This function can be used for returning a generic error when no specific error code and message are available.
func NewCommonError(apiName string, code int64, msg string) *CommonError {
	return &CommonError{
		apiName: apiName,
		ErrCode: code,
		ErrMsg:  msg,
	}
}

// DecodeWithCommonError decodes the response using the CommonError struct
func DecodeWithCommonError(response []byte, apiName string) error {
	var commError CommonError
	err := json.Unmarshal(response, &commError)
	if err != nil {
		return err
	}
	commError.apiName = apiName
	if commError.ErrCode != 0 {
		return &commError
	}
	return nil
}

// DecodeWithError decodes the response JSON into an object and checks for common error fields.
func DecodeWithError(response []byte, obj any, apiName string) error {
	err := json.Unmarshal(response, obj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	reflectObj := reflect.ValueOf(obj)
	if !reflectObj.IsValid() {
		return fmt.Errorf("invalid object")
	}
	commonError := reflectObj.Elem().FieldByName("CommonError")
	if !commonError.IsValid() || commonError.Kind() != reflect.Struct {
		return fmt.Errorf("invalid or non-struct commonError")
	}
	errCode := commonError.FieldByName("ErrCode")
	errMsg := commonError.FieldByName("ErrMsg")
	if !errCode.IsValid() || !errMsg.IsValid() {
		return fmt.Errorf("invalid errCode or errMsg")
	}
	if errCode.Int() != 0 {
		return &CommonError{
			apiName: apiName,
			ErrCode: errCode.Int(),
			ErrMsg:  errMsg.String(),
		}
	}
	return nil
}
