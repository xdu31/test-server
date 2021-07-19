// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: github.com/xdu31/test-server/pkg/server/pb/service.proto

package pb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// Validate checks the field values on Ip with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Ip) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetId() < 0 {
		return IpValidationError{
			field:  "Id",
			reason: "value must be greater than or equal to 0",
		}
	}

	// no validation rules for IpAddress

	if v, ok := interface{}(m.GetCreatedTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IpValidationError{
				field:  "CreatedTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetUpdatedTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IpValidationError{
				field:  "UpdatedTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// IpValidationError is the validation error returned by Ip.Validate if the
// designated constraints aren't met.
type IpValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IpValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IpValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IpValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IpValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IpValidationError) ErrorName() string { return "IpValidationError" }

// Error satisfies the builtin error interface
func (e IpValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IpValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IpValidationError{}

// Validate checks the field values on CreateIpRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *CreateIpRequest) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetPayload()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateIpRequestValidationError{
				field:  "Payload",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// CreateIpRequestValidationError is the validation error returned by
// CreateIpRequest.Validate if the designated constraints aren't met.
type CreateIpRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateIpRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateIpRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateIpRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateIpRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateIpRequestValidationError) ErrorName() string { return "CreateIpRequestValidationError" }

// Error satisfies the builtin error interface
func (e CreateIpRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateIpRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateIpRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateIpRequestValidationError{}

// Validate checks the field values on CreateIpResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *CreateIpResponse) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetResult()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateIpResponseValidationError{
				field:  "Result",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// CreateIpResponseValidationError is the validation error returned by
// CreateIpResponse.Validate if the designated constraints aren't met.
type CreateIpResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateIpResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateIpResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateIpResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateIpResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateIpResponseValidationError) ErrorName() string { return "CreateIpResponseValidationError" }

// Error satisfies the builtin error interface
func (e CreateIpResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateIpResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateIpResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateIpResponseValidationError{}

// Validate checks the field values on ReadIpRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ReadIpRequest) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetId() < 0 {
		return ReadIpRequestValidationError{
			field:  "Id",
			reason: "value must be greater than or equal to 0",
		}
	}

	return nil
}

// ReadIpRequestValidationError is the validation error returned by
// ReadIpRequest.Validate if the designated constraints aren't met.
type ReadIpRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadIpRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadIpRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadIpRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadIpRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadIpRequestValidationError) ErrorName() string { return "ReadIpRequestValidationError" }

// Error satisfies the builtin error interface
func (e ReadIpRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadIpRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadIpRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadIpRequestValidationError{}

// Validate checks the field values on ReadIpResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ReadIpResponse) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetResult()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ReadIpResponseValidationError{
				field:  "Result",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// ReadIpResponseValidationError is the validation error returned by
// ReadIpResponse.Validate if the designated constraints aren't met.
type ReadIpResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadIpResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadIpResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadIpResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadIpResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadIpResponseValidationError) ErrorName() string { return "ReadIpResponseValidationError" }

// Error satisfies the builtin error interface
func (e ReadIpResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadIpResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadIpResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadIpResponseValidationError{}

// Validate checks the field values on UpdateIpRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *UpdateIpRequest) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetPayload()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateIpRequestValidationError{
				field:  "Payload",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetFields()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateIpRequestValidationError{
				field:  "Fields",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// UpdateIpRequestValidationError is the validation error returned by
// UpdateIpRequest.Validate if the designated constraints aren't met.
type UpdateIpRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateIpRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateIpRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateIpRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateIpRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateIpRequestValidationError) ErrorName() string { return "UpdateIpRequestValidationError" }

// Error satisfies the builtin error interface
func (e UpdateIpRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateIpRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateIpRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateIpRequestValidationError{}

// Validate checks the field values on UpdateIpResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *UpdateIpResponse) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetResult()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateIpResponseValidationError{
				field:  "Result",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// UpdateIpResponseValidationError is the validation error returned by
// UpdateIpResponse.Validate if the designated constraints aren't met.
type UpdateIpResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateIpResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateIpResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateIpResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateIpResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateIpResponseValidationError) ErrorName() string { return "UpdateIpResponseValidationError" }

// Error satisfies the builtin error interface
func (e UpdateIpResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateIpResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateIpResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateIpResponseValidationError{}

// Validate checks the field values on DeleteIpRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *DeleteIpRequest) Validate() error {
	if m == nil {
		return nil
	}

	_DeleteIpRequest_Ids_Unique := make(map[int32]struct{}, len(m.GetIds()))

	for idx, item := range m.GetIds() {
		_, _ = idx, item

		if _, exists := _DeleteIpRequest_Ids_Unique[item]; exists {
			return DeleteIpRequestValidationError{
				field:  fmt.Sprintf("Ids[%v]", idx),
				reason: "repeated value must contain unique items",
			}
		} else {
			_DeleteIpRequest_Ids_Unique[item] = struct{}{}
		}

		if item < 0 {
			return DeleteIpRequestValidationError{
				field:  fmt.Sprintf("Ids[%v]", idx),
				reason: "value must be greater than or equal to 0",
			}
		}

	}

	return nil
}

// DeleteIpRequestValidationError is the validation error returned by
// DeleteIpRequest.Validate if the designated constraints aren't met.
type DeleteIpRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteIpRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteIpRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteIpRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteIpRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteIpRequestValidationError) ErrorName() string { return "DeleteIpRequestValidationError" }

// Error satisfies the builtin error interface
func (e DeleteIpRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteIpRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteIpRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteIpRequestValidationError{}

// Validate checks the field values on DeleteIpResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *DeleteIpResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// DeleteIpResponseValidationError is the validation error returned by
// DeleteIpResponse.Validate if the designated constraints aren't met.
type DeleteIpResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteIpResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteIpResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteIpResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteIpResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteIpResponseValidationError) ErrorName() string { return "DeleteIpResponseValidationError" }

// Error satisfies the builtin error interface
func (e DeleteIpResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteIpResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteIpResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteIpResponseValidationError{}

// Validate checks the field values on ListIpsRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ListIpsRequest) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetFilter()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListIpsRequestValidationError{
				field:  "Filter",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetFields()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListIpsRequestValidationError{
				field:  "Fields",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// ListIpsRequestValidationError is the validation error returned by
// ListIpsRequest.Validate if the designated constraints aren't met.
type ListIpsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListIpsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListIpsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListIpsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListIpsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListIpsRequestValidationError) ErrorName() string { return "ListIpsRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListIpsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListIpsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListIpsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListIpsRequestValidationError{}

// Validate checks the field values on ListIpsResponse with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ListIpsResponse) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetResults() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListIpsResponseValidationError{
					field:  fmt.Sprintf("Results[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListIpsResponseValidationError is the validation error returned by
// ListIpsResponse.Validate if the designated constraints aren't met.
type ListIpsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListIpsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListIpsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListIpsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListIpsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListIpsResponseValidationError) ErrorName() string { return "ListIpsResponseValidationError" }

// Error satisfies the builtin error interface
func (e ListIpsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListIpsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListIpsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListIpsResponseValidationError{}
