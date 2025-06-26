package aclgate

import (
	"regexp"
)

var (
	resourceTypeRegex = regexp.MustCompile(`^[^:#@\s]{1,254}$`)
	resourceIDRegex   = regexp.MustCompile(`^[^#:\s]+$`)
	subjectTypeRegex  = regexp.MustCompile(`^[^:#@\s]{1,254}$`)
	subjectIDRegex    = regexp.MustCompile(`^[^:#\s]+$`)
	relationNameRegex = regexp.MustCompile(`^[^:#@\s]{1,50}$`)
)

type Tuple struct {
	Resource *Resource
	Subject  *Subject
	Relation *Relation
}

type CheckRequest struct {
	Tuple *Tuple
}

type BatchCheckResult struct {
	Request *CheckRequest
	Allowed bool
	Error   error
}

// NewTuple creates a new Tuple with the given parameters
func NewTuple(resourceType, resourceId, subjectType, subjectId, relationName string) (*Tuple, error) {
	resource, err := NewResource(resourceType, resourceId)
	if err != nil {
		return nil, err
	}

	subject, err := NewSubject(subjectType, subjectId)
	if err != nil {
		return nil, err
	}

	relation, err := NewRelation(relationName)
	if err != nil {
		return nil, err
	}

	return &Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: relation,
	}, nil
}

// ListResourcesRequest represents a request to list permissions
type ListResourcesRequest struct {
	Type     string
	Subject  *Subject
	Relation *Relation
}

type ListResourcesResponse struct {
	Resources []*Resource
}

type ListSubjectsRequest struct {
	Type     string
	Resource *Resource
	Relation *Relation
	Limit    int32
	Offset   int32
}

// ListSubjectsResponse represents a response containing a list of permissions
type ListSubjectsResponse struct {
	Subjects []*Subject
}

// AuditRequest represents a request to list audit logs
type AuditRequest struct {
	Resource *Resource
	Subject  *Subject
	Relation *Relation
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
	Tuple     *Tuple
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
func NewResource(resourceType, resourceId string) (*Resource, error) {
	if !resourceTypeRegex.MatchString(resourceType) {
		return nil, ErrInvalidResourceType
	}

	if !resourceIDRegex.MatchString(resourceId) {
		return nil, ErrInvalidResourceId
	}

	return &Resource{
		Type: resourceType,
		ID:   resourceId,
	}, nil
}

// NewSubject creates a new Subject
func NewSubject(subjectType, subjectId string) (*Subject, error) {
	if !subjectTypeRegex.MatchString(subjectType) {
		return nil, ErrInvalidSubjectType
	}

	if !subjectIDRegex.MatchString(subjectId) {
		return nil, ErrInvalidSubjectId
	}
	return &Subject{
		Type: subjectType,
		ID:   subjectId,
	}, nil
}

// NewRelation creates a new Relation
func NewRelation(name string) (*Relation, error) {
	if !relationNameRegex.MatchString(name) {
		return nil, ErrInvalidRelationName
	}

	return &Relation{
		Name: name,
	}, nil
}
