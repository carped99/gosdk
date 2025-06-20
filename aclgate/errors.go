package aclgate

import (
	"errors"
	"fmt"
)

// Common ACL errors
var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrServiceNotFound  = errors.New("acl service not found in context")
	ErrPermissionDenied = errors.New("permission denied")
	ErrResourceNotFound = errors.New("resource not found")
	ErrSubjectNotFound  = errors.New("subject not found")
)

// AclError represents an ACL-specific error
type AclError struct {
	Code    string
	Message string
	Cause   error
}

func (e *AclError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("acl error [%s]: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("acl error [%s]: %s", e.Code, e.Message)
}

func (e *AclError) Unwrap() error {
	return e.Cause
}

// NewAclError creates a new ACL error
func NewAclError(code, message string, cause error) *AclError {
	return &AclError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsAclError checks if an error is an ACL error
func IsAclError(err error) bool {
	var aclErr *AclError
	return errors.As(err, &aclErr)
}

// GetAclErrorCode returns the error code if it's an ACL error
func GetAclErrorCode(err error) string {
	var aclErr *AclError
	if errors.As(err, &aclErr) {
		return aclErr.Code
	}
	return ""
}
