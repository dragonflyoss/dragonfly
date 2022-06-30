// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/rpc/scheduler/scheduler.proto

package scheduler

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

	base "d7y.io/dragonfly/v2/pkg/rpc/base"
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

	_ = base.Pattern(0)

	_ = base.TaskType(0)

	_ = base.SizeScope(0)

	_ = base.Code(0)

	_ = base.Code(0)

	_ = base.Code(0)

	_ = base.TaskType(0)

	_ = base.TaskType(0)
)

// Validate checks the field values on PeerTaskRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *PeerTaskRequest) Validate() error {
	if m == nil {
		return nil
	}

	if uri, err := url.Parse(m.GetUrl()); err != nil {
		return PeerTaskRequestValidationError{
			field:  "Url",
			reason: "value must be a valid URI",
			cause:  err,
		}
	} else if !uri.IsAbs() {
		return PeerTaskRequestValidationError{
			field:  "Url",
			reason: "value must be absolute",
		}
	}

	if m.GetUrlMeta() == nil {
		return PeerTaskRequestValidationError{
			field:  "UrlMeta",
			reason: "value is required",
		}
	}

	if v, ok := any(m.GetUrlMeta()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PeerTaskRequestValidationError{
				field:  "UrlMeta",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if utf8.RuneCountInString(m.GetPeerId()) < 1 {
		return PeerTaskRequestValidationError{
			field:  "PeerId",
			reason: "value length must be at least 1 runes",
		}
	}

	if v, ok := any(m.GetPeerHost()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PeerTaskRequestValidationError{
				field:  "PeerHost",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := any(m.GetHostLoad()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PeerTaskRequestValidationError{
				field:  "HostLoad",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for IsMigrating

	// no validation rules for Pattern

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return PeerTaskRequestValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// PeerTaskRequestValidationError is the validation error returned by
// PeerTaskRequest.Validate if the designated constraints aren't met.
type PeerTaskRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PeerTaskRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PeerTaskRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PeerTaskRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PeerTaskRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PeerTaskRequestValidationError) ErrorName() string { return "PeerTaskRequestValidationError" }

// Error satisfies the builtin error interface
func (e PeerTaskRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPeerTaskRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PeerTaskRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PeerTaskRequestValidationError{}

// Validate checks the field values on RegisterResult with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *RegisterResult) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for TaskType

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return RegisterResultValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	if _, ok := base.SizeScope_name[int32(m.GetSizeScope())]; !ok {
		return RegisterResultValidationError{
			field:  "SizeScope",
			reason: "value must be one of the defined enum values",
		}
	}

	if v, ok := any(m.GetExtendAttribute()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RegisterResultValidationError{
				field:  "ExtendAttribute",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.DirectPiece.(type) {

	case *RegisterResult_SinglePiece:

		if v, ok := any(m.GetSinglePiece()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RegisterResultValidationError{
					field:  "SinglePiece",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *RegisterResult_PieceContent:
		// no validation rules for PieceContent

	}

	return nil
}

// RegisterResultValidationError is the validation error returned by
// RegisterResult.Validate if the designated constraints aren't met.
type RegisterResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterResultValidationError) ErrorName() string { return "RegisterResultValidationError" }

// Error satisfies the builtin error interface
func (e RegisterResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterResultValidationError{}

// Validate checks the field values on SinglePiece with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *SinglePiece) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetDstPid()) < 1 {
		return SinglePieceValidationError{
			field:  "DstPid",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetDstAddr()) < 1 {
		return SinglePieceValidationError{
			field:  "DstAddr",
			reason: "value length must be at least 1 runes",
		}
	}

	if v, ok := any(m.GetPieceInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SinglePieceValidationError{
				field:  "PieceInfo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// SinglePieceValidationError is the validation error returned by
// SinglePiece.Validate if the designated constraints aren't met.
type SinglePieceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SinglePieceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SinglePieceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SinglePieceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SinglePieceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SinglePieceValidationError) ErrorName() string { return "SinglePieceValidationError" }

// Error satisfies the builtin error interface
func (e SinglePieceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSinglePiece.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SinglePieceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SinglePieceValidationError{}

// Validate checks the field values on PeerHost with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *PeerHost) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetId()) < 1 {
		return PeerHostValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
	}

	if ip := net.ParseIP(m.GetIp()); ip == nil {
		return PeerHostValidationError{
			field:  "Ip",
			reason: "value must be a valid IP address",
		}
	}

	if val := m.GetRpcPort(); val < 1024 || val >= 65535 {
		return PeerHostValidationError{
			field:  "RpcPort",
			reason: "value must be inside range [1024, 65535)",
		}
	}

	if val := m.GetDownPort(); val < 1024 || val >= 65535 {
		return PeerHostValidationError{
			field:  "DownPort",
			reason: "value must be inside range [1024, 65535)",
		}
	}

	if err := m._validateHostname(m.GetHostName()); err != nil {
		return PeerHostValidationError{
			field:  "HostName",
			reason: "value must be a valid hostname",
			cause:  err,
		}
	}

	// no validation rules for SecurityDomain

	// no validation rules for Location

	// no validation rules for Idc

	// no validation rules for NetTopology

	return nil
}

func (m *PeerHost) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

// PeerHostValidationError is the validation error returned by
// PeerHost.Validate if the designated constraints aren't met.
type PeerHostValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PeerHostValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PeerHostValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PeerHostValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PeerHostValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PeerHostValidationError) ErrorName() string { return "PeerHostValidationError" }

// Error satisfies the builtin error interface
func (e PeerHostValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPeerHost.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PeerHostValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PeerHostValidationError{}

// Validate checks the field values on PieceResult with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *PieceResult) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return PieceResultValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetSrcPid()) < 1 {
		return PieceResultValidationError{
			field:  "SrcPid",
			reason: "value length must be at least 1 runes",
		}
	}

	// no validation rules for DstPid

	if v, ok := any(m.GetPieceInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PieceResultValidationError{
				field:  "PieceInfo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for BeginTime

	// no validation rules for EndTime

	// no validation rules for Success

	// no validation rules for Code

	if v, ok := any(m.GetHostLoad()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PieceResultValidationError{
				field:  "HostLoad",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for FinishedCount

	if v, ok := any(m.GetExtendAttribute()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PieceResultValidationError{
				field:  "ExtendAttribute",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// PieceResultValidationError is the validation error returned by
// PieceResult.Validate if the designated constraints aren't met.
type PieceResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PieceResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PieceResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PieceResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PieceResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PieceResultValidationError) ErrorName() string { return "PieceResultValidationError" }

// Error satisfies the builtin error interface
func (e PieceResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPieceResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PieceResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PieceResultValidationError{}

// Validate checks the field values on PeerPacket with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *PeerPacket) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return PeerPacketValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetSrcPid()) < 1 {
		return PeerPacketValidationError{
			field:  "SrcPid",
			reason: "value length must be at least 1 runes",
		}
	}

	if m.GetParallelCount() < 1 {
		return PeerPacketValidationError{
			field:  "ParallelCount",
			reason: "value must be greater than or equal to 1",
		}
	}

	if v, ok := any(m.GetMainPeer()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PeerPacketValidationError{
				field:  "MainPeer",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetStealPeers() {
		_, _ = idx, item

		if v, ok := any(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PeerPacketValidationError{
					field:  fmt.Sprintf("StealPeers[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Code

	return nil
}

// PeerPacketValidationError is the validation error returned by
// PeerPacket.Validate if the designated constraints aren't met.
type PeerPacketValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PeerPacketValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PeerPacketValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PeerPacketValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PeerPacketValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PeerPacketValidationError) ErrorName() string { return "PeerPacketValidationError" }

// Error satisfies the builtin error interface
func (e PeerPacketValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPeerPacket.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PeerPacketValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PeerPacketValidationError{}

// Validate checks the field values on PeerResult with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *PeerResult) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return PeerResultValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetPeerId()) < 1 {
		return PeerResultValidationError{
			field:  "PeerId",
			reason: "value length must be at least 1 runes",
		}
	}

	if ip := net.ParseIP(m.GetSrcIp()); ip == nil {
		return PeerResultValidationError{
			field:  "SrcIp",
			reason: "value must be a valid IP address",
		}
	}

	// no validation rules for SecurityDomain

	// no validation rules for Idc

	if uri, err := url.Parse(m.GetUrl()); err != nil {
		return PeerResultValidationError{
			field:  "Url",
			reason: "value must be a valid URI",
			cause:  err,
		}
	} else if !uri.IsAbs() {
		return PeerResultValidationError{
			field:  "Url",
			reason: "value must be absolute",
		}
	}

	if m.GetContentLength() < -1 {
		return PeerResultValidationError{
			field:  "ContentLength",
			reason: "value must be greater than or equal to -1",
		}
	}

	// no validation rules for Traffic

	// no validation rules for Cost

	// no validation rules for Success

	// no validation rules for Code

	if m.GetTotalPieceCount() < -1 {
		return PeerResultValidationError{
			field:  "TotalPieceCount",
			reason: "value must be greater than or equal to -1",
		}
	}

	return nil
}

// PeerResultValidationError is the validation error returned by
// PeerResult.Validate if the designated constraints aren't met.
type PeerResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PeerResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PeerResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PeerResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PeerResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PeerResultValidationError) ErrorName() string { return "PeerResultValidationError" }

// Error satisfies the builtin error interface
func (e PeerResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPeerResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PeerResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PeerResultValidationError{}

// Validate checks the field values on PeerTarget with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *PeerTarget) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return PeerTargetValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetPeerId()) < 1 {
		return PeerTargetValidationError{
			field:  "PeerId",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// PeerTargetValidationError is the validation error returned by
// PeerTarget.Validate if the designated constraints aren't met.
type PeerTargetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PeerTargetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PeerTargetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PeerTargetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PeerTargetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PeerTargetValidationError) ErrorName() string { return "PeerTargetValidationError" }

// Error satisfies the builtin error interface
func (e PeerTargetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPeerTarget.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PeerTargetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PeerTargetValidationError{}

// Validate checks the field values on StatTaskRequest with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *StatTaskRequest) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return StatTaskRequestValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// StatTaskRequestValidationError is the validation error returned by
// StatTaskRequest.Validate if the designated constraints aren't met.
type StatTaskRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StatTaskRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StatTaskRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StatTaskRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StatTaskRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StatTaskRequestValidationError) ErrorName() string { return "StatTaskRequestValidationError" }

// Error satisfies the builtin error interface
func (e StatTaskRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStatTaskRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StatTaskRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StatTaskRequestValidationError{}

// Validate checks the field values on Task with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Task) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetId()) < 1 {
		return TaskValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
	}

	// no validation rules for Type

	if m.GetContentLength() < 1 {
		return TaskValidationError{
			field:  "ContentLength",
			reason: "value must be greater than or equal to 1",
		}
	}

	if m.GetTotalPieceCount() < 1 {
		return TaskValidationError{
			field:  "TotalPieceCount",
			reason: "value must be greater than or equal to 1",
		}
	}

	if utf8.RuneCountInString(m.GetState()) < 1 {
		return TaskValidationError{
			field:  "State",
			reason: "value length must be at least 1 runes",
		}
	}

	if m.GetPeerCount() < 0 {
		return TaskValidationError{
			field:  "PeerCount",
			reason: "value must be greater than or equal to 0",
		}
	}

	// no validation rules for HasAvailablePeer

	return nil
}

// TaskValidationError is the validation error returned by Task.Validate if the
// designated constraints aren't met.
type TaskValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TaskValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TaskValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TaskValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TaskValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TaskValidationError) ErrorName() string { return "TaskValidationError" }

// Error satisfies the builtin error interface
func (e TaskValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTask.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TaskValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TaskValidationError{}

// Validate checks the field values on AnnounceTaskRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *AnnounceTaskRequest) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetTaskId()) < 1 {
		return AnnounceTaskRequestValidationError{
			field:  "TaskId",
			reason: "value length must be at least 1 runes",
		}
	}

	if m.GetUrl() != "" {

		if uri, err := url.Parse(m.GetUrl()); err != nil {
			return AnnounceTaskRequestValidationError{
				field:  "Url",
				reason: "value must be a valid URI",
				cause:  err,
			}
		} else if !uri.IsAbs() {
			return AnnounceTaskRequestValidationError{
				field:  "Url",
				reason: "value must be absolute",
			}
		}

	}

	if m.GetUrlMeta() == nil {
		return AnnounceTaskRequestValidationError{
			field:  "UrlMeta",
			reason: "value is required",
		}
	}

	if v, ok := any(m.GetUrlMeta()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AnnounceTaskRequestValidationError{
				field:  "UrlMeta",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := any(m.GetPeerHost()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AnnounceTaskRequestValidationError{
				field:  "PeerHost",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.GetPiecePacket() == nil {
		return AnnounceTaskRequestValidationError{
			field:  "PiecePacket",
			reason: "value is required",
		}
	}

	if v, ok := any(m.GetPiecePacket()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AnnounceTaskRequestValidationError{
				field:  "PiecePacket",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for TaskType

	return nil
}

// AnnounceTaskRequestValidationError is the validation error returned by
// AnnounceTaskRequest.Validate if the designated constraints aren't met.
type AnnounceTaskRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AnnounceTaskRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AnnounceTaskRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AnnounceTaskRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AnnounceTaskRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AnnounceTaskRequestValidationError) ErrorName() string {
	return "AnnounceTaskRequestValidationError"
}

// Error satisfies the builtin error interface
func (e AnnounceTaskRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAnnounceTaskRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AnnounceTaskRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AnnounceTaskRequestValidationError{}

// Validate checks the field values on PeerPacket_DestPeer with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *PeerPacket_DestPeer) Validate() error {
	if m == nil {
		return nil
	}

	if ip := net.ParseIP(m.GetIp()); ip == nil {
		return PeerPacket_DestPeerValidationError{
			field:  "Ip",
			reason: "value must be a valid IP address",
		}
	}

	if val := m.GetRpcPort(); val < 1024 || val >= 65535 {
		return PeerPacket_DestPeerValidationError{
			field:  "RpcPort",
			reason: "value must be inside range [1024, 65535)",
		}
	}

	if utf8.RuneCountInString(m.GetPeerId()) < 1 {
		return PeerPacket_DestPeerValidationError{
			field:  "PeerId",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// PeerPacket_DestPeerValidationError is the validation error returned by
// PeerPacket_DestPeer.Validate if the designated constraints aren't met.
type PeerPacket_DestPeerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PeerPacket_DestPeerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PeerPacket_DestPeerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PeerPacket_DestPeerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PeerPacket_DestPeerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PeerPacket_DestPeerValidationError) ErrorName() string {
	return "PeerPacket_DestPeerValidationError"
}

// Error satisfies the builtin error interface
func (e PeerPacket_DestPeerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPeerPacket_DestPeer.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PeerPacket_DestPeerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PeerPacket_DestPeerValidationError{}
