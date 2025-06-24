package aclgate

import (
	"errors"
)

type Tuple struct {
	Resource Resource
	Subject  Subject
	Relation Relation
}

type CheckRequest struct {
	Tuple Tuple
}

type BatchCheckResult struct {
	Request CheckRequest
	Allowed bool
	Error   error
}

var (
	ErrInvalidTuple = errors.New("invalid tuple")
)

// NewTuple creates a new Tuple with the given parameters
func NewTuple(resourceType, resourceId, subjectType, subjectId, relation string) Tuple {
	return Tuple{
		Resource: NewResource(resourceType, resourceId),
		Subject:  NewSubject(subjectType, subjectId),
		Relation: NewRelation(relation),
	}
}

// NewCheckRequest creates a new CheckRequest with the given tuple
func NewCheckRequest(tuple Tuple) CheckRequest {
	return CheckRequest{Tuple: tuple}
}

// NewCheckRequestFromParams creates a new CheckRequest with individual parameters
func NewCheckRequestFromParams(resourceType, resourceId, subjectType, subjectId, relation string) CheckRequest {
	return CheckRequest{
		Tuple: NewTuple(resourceType, resourceId, subjectType, subjectId, relation),
	}
}

// ListResourcesRequest represents a request to list permissions
type ListResourcesRequest struct {
	Type     string
	Subject  Subject
	Relation Relation
}

type ListResourcesResponse struct {
	Resources []Resource
}

type ListSubjectsRequest struct {
	Type     string
	Resource Resource
	Relation Relation
	Limit    int32
	Offset   int32
}

// ListSubjectsResponse represents a response containing a list of permissions
type ListSubjectsResponse struct {
	Subjects []Subject
}

// AuditRequest represents a request to list audit logs
type AuditRequest struct {
	Resource Resource
	Subject  Subject
	Relation Relation
	PageSize int32
	Cursor   string
}

// AuditResponse represents a response containing a list of audit logs
type AuditResponse struct {
	Logs []AuditLog
}

// AuditLog represents a single audit log entry
type AuditLog struct {
	ID        string
	Action    string // e.g., "WRITE", "DELETE"
	Tuple     Tuple
	Actor     string
	Timestamp string
	Reason    string
}

// Resource represents a resource in the ACL system
type Resource struct {
	Type string
	ID   string
}

// Subject represents a subject in the ACL system
type Subject struct {
	Type string
	ID   string
}

// Relation represents a relation in the ACL system
type Relation struct {
	Name string
}

// NewResource creates a new Resource
func NewResource(resourceType, resourceId string) Resource {
	return Resource{
		Type: resourceType,
		ID:   resourceId,
	}
}

// NewSubject creates a new Subject
func NewSubject(subjectType, subjectId string) Subject {
	return Subject{
		Type: subjectType,
		ID:   subjectId,
	}
}

// NewRelation creates a new Relation
func NewRelation(name string) Relation {
	return Relation{
		Name: name,
	}
}
