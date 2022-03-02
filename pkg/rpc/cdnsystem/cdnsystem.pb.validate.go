// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/rpc/cdnsystem/cdnsystem.proto

package cdnsystem

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
)

// Validate checks the field values on SeedRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SeedRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SeedRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SeedRequestMultiError, or
// nil if none found.
func (m *SeedRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *SeedRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		err := SeedRequestValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if uri, err := url.Parse(m.GetUrl()); err != nil {
		err = SeedRequestValidationError{
			field:  "Url",
			reason: "value must be a valid URI",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	} else if !uri.IsAbs() {
		err := SeedRequestValidationError{
			field:  "Url",
			reason: "value must be absolute",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetUrlMeta()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SeedRequestValidationError{
					field:  "UrlMeta",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SeedRequestValidationError{
					field:  "UrlMeta",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUrlMeta()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SeedRequestValidationError{
				field:  "UrlMeta",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SeedRequestMultiError(errors)
	}
	return nil
}

// SeedRequestMultiError is an error wrapping multiple validation errors
// returned by SeedRequest.ValidateAll() if the designated constraints aren't met.
type SeedRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SeedRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SeedRequestMultiError) AllErrors() []error { return m }

// SeedRequestValidationError is the validation error returned by
// SeedRequest.Validate if the designated constraints aren't met.
type SeedRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SeedRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SeedRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SeedRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SeedRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SeedRequestValidationError) ErrorName() string { return "SeedRequestValidationError" }

// Error satisfies the builtin error interface
func (e SeedRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSeedRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SeedRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SeedRequestValidationError{}

// Validate checks the field values on PieceSeed with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PieceSeed) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PieceSeed with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PieceSeedMultiError, or nil
// if none found.
func (m *PieceSeed) ValidateAll() error {
	return m.validate(true)
}

func (m *PieceSeed) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetPeerId()) < 1 {
		err := PieceSeedValidationError{
			field:  "PeerId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetHostUuid()) < 1 {
		err := PieceSeedValidationError{
			field:  "HostUuid",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetPieceInfo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PieceSeedValidationError{
					field:  "PieceInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PieceSeedValidationError{
					field:  "PieceInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPieceInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PieceSeedValidationError{
				field:  "PieceInfo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Done

	// no validation rules for ContentLength

	// no validation rules for TotalPieceCount

	// no validation rules for BeginTime

	// no validation rules for EndTime

	if len(errors) > 0 {
		return PieceSeedMultiError(errors)
	}
	return nil
}

// PieceSeedMultiError is an error wrapping multiple validation errors returned
// by PieceSeed.ValidateAll() if the designated constraints aren't met.
type PieceSeedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PieceSeedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PieceSeedMultiError) AllErrors() []error { return m }

// PieceSeedValidationError is the validation error returned by
// PieceSeed.Validate if the designated constraints aren't met.
type PieceSeedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PieceSeedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PieceSeedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PieceSeedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PieceSeedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PieceSeedValidationError) ErrorName() string { return "PieceSeedValidationError" }

// Error satisfies the builtin error interface
func (e PieceSeedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPieceSeed.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PieceSeedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PieceSeedValidationError{}
