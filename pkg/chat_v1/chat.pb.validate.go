// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: chat.proto

package chat_v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
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
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Message with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Message) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Message with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in MessageMultiError, or nil if none found.
func (m *Message) ValidateAll() error {
	return m.validate(true)
}

func (m *Message) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for From

	// no validation rules for Text

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MessageValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MessageValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MessageValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return MessageMultiError(errors)
	}

	return nil
}

// MessageMultiError is an error wrapping multiple validation errors returned
// by Message.ValidateAll() if the designated constraints aren't met.
type MessageMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MessageMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MessageMultiError) AllErrors() []error { return m }

// MessageValidationError is the validation error returned by Message.Validate
// if the designated constraints aren't met.
type MessageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MessageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MessageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MessageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MessageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MessageValidationError) ErrorName() string { return "MessageValidationError" }

// Error satisfies the builtin error interface
func (e MessageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMessage.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MessageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MessageValidationError{}

// Validate checks the field values on Chat with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Chat) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Chat with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ChatMultiError, or nil if none found.
func (m *Chat) ValidateAll() error {
	return m.validate(true)
}

func (m *Chat) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetUsernames()) < 1 {
		err := ChatValidationError{
			field:  "Usernames",
			reason: "value must contain at least 1 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	_Chat_Usernames_Unique := make(map[string]struct{}, len(m.GetUsernames()))

	for idx, item := range m.GetUsernames() {
		_, _ = idx, item

		if _, exists := _Chat_Usernames_Unique[item]; exists {
			err := ChatValidationError{
				field:  fmt.Sprintf("Usernames[%v]", idx),
				reason: "repeated value must contain unique items",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {
			_Chat_Usernames_Unique[item] = struct{}{}
		}

		if l := utf8.RuneCountInString(item); l < 1 || l > 50 {
			err := ChatValidationError{
				field:  fmt.Sprintf("Usernames[%v]", idx),
				reason: "value length must be between 1 and 50 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return ChatMultiError(errors)
	}

	return nil
}

// ChatMultiError is an error wrapping multiple validation errors returned by
// Chat.ValidateAll() if the designated constraints aren't met.
type ChatMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChatMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChatMultiError) AllErrors() []error { return m }

// ChatValidationError is the validation error returned by Chat.Validate if the
// designated constraints aren't met.
type ChatValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChatValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChatValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChatValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChatValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChatValidationError) ErrorName() string { return "ChatValidationError" }

// Error satisfies the builtin error interface
func (e ChatValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChat.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChatValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChatValidationError{}

// Validate checks the field values on CreateResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateResponseMultiError,
// or nil if none found.
func (m *CreateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return CreateResponseMultiError(errors)
	}

	return nil
}

// CreateResponseMultiError is an error wrapping multiple validation errors
// returned by CreateResponse.ValidateAll() if the designated constraints
// aren't met.
type CreateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateResponseMultiError) AllErrors() []error { return m }

// CreateResponseValidationError is the validation error returned by
// CreateResponse.Validate if the designated constraints aren't met.
type CreateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateResponseValidationError) ErrorName() string { return "CreateResponseValidationError" }

// Error satisfies the builtin error interface
func (e CreateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateResponseValidationError{}

// Validate checks the field values on DeleteRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DeleteRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DeleteRequestMultiError, or
// nil if none found.
func (m *DeleteRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return DeleteRequestMultiError(errors)
	}

	return nil
}

// DeleteRequestMultiError is an error wrapping multiple validation errors
// returned by DeleteRequest.ValidateAll() if the designated constraints
// aren't met.
type DeleteRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteRequestMultiError) AllErrors() []error { return m }

// DeleteRequestValidationError is the validation error returned by
// DeleteRequest.Validate if the designated constraints aren't met.
type DeleteRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteRequestValidationError) ErrorName() string { return "DeleteRequestValidationError" }

// Error satisfies the builtin error interface
func (e DeleteRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteRequestValidationError{}

// Validate checks the field values on ConnectRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ConnectRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ConnectRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ConnectRequestMultiError,
// or nil if none found.
func (m *ConnectRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ConnectRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ChatId

	// no validation rules for Username

	if len(errors) > 0 {
		return ConnectRequestMultiError(errors)
	}

	return nil
}

// ConnectRequestMultiError is an error wrapping multiple validation errors
// returned by ConnectRequest.ValidateAll() if the designated constraints
// aren't met.
type ConnectRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ConnectRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ConnectRequestMultiError) AllErrors() []error { return m }

// ConnectRequestValidationError is the validation error returned by
// ConnectRequest.Validate if the designated constraints aren't met.
type ConnectRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConnectRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConnectRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConnectRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConnectRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConnectRequestValidationError) ErrorName() string { return "ConnectRequestValidationError" }

// Error satisfies the builtin error interface
func (e ConnectRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConnectRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConnectRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConnectRequestValidationError{}

// Validate checks the field values on SendMessageRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SendMessageRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SendMessageRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SendMessageRequestMultiError, or nil if none found.
func (m *SendMessageRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SendMessageRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ChatId

	if all {
		switch v := interface{}(m.GetMessage()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SendMessageRequestValidationError{
					field:  "Message",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SendMessageRequestValidationError{
					field:  "Message",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMessage()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SendMessageRequestValidationError{
				field:  "Message",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SendMessageRequestMultiError(errors)
	}

	return nil
}

// SendMessageRequestMultiError is an error wrapping multiple validation errors
// returned by SendMessageRequest.ValidateAll() if the designated constraints
// aren't met.
type SendMessageRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SendMessageRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SendMessageRequestMultiError) AllErrors() []error { return m }

// SendMessageRequestValidationError is the validation error returned by
// SendMessageRequest.Validate if the designated constraints aren't met.
type SendMessageRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SendMessageRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SendMessageRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SendMessageRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SendMessageRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SendMessageRequestValidationError) ErrorName() string {
	return "SendMessageRequestValidationError"
}

// Error satisfies the builtin error interface
func (e SendMessageRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSendMessageRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SendMessageRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SendMessageRequestValidationError{}
