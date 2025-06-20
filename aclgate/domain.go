package aclgate

import (
	"errors"
)

type Tuple struct {
	ResourceType string
	ResourceId   string
	SubjectType  string
	SubjectId    string
	Relation     string
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
		ResourceType: resourceType,
		ResourceId:   resourceId,
		SubjectType:  subjectType,
		SubjectId:    subjectId,
		Relation:     relation,
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

// ListRequest represents a request to list permissions
type ListRequest struct {
	Resource Resource
	Subject  Subject
	Relation Relation
	Limit    int32
	Offset   int32
}

// ListResponse represents a response containing a list of permissions
type ListResponse struct {
	Tuples []Tuple
}

// AuditRequest represents a request to list audit logs
type AuditRequest struct {
	Resource Resource
	Subject  Subject
	Relation Relation
	Limit    int32
	Offset   int32
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

// NewListRequest creates a new ListRequest
func NewListRequest(resource Resource, subject Subject, relation Relation, limit, offset int32) ListRequest {
	return ListRequest{
		Resource: resource,
		Subject:  subject,
		Relation: relation,
		Limit:    limit,
		Offset:   offset,
	}
}

// NewAuditRequest creates a new AuditRequest
func NewAuditRequest(resource Resource, subject Subject, relation Relation, limit, offset int32) AuditRequest {
	return AuditRequest{
		Resource: resource,
		Subject:  subject,
		Relation: relation,
		Limit:    limit,
		Offset:   offset,
	}
}
