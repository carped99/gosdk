package aclgate

import (
	"context"
	"google.golang.org/grpc"
)

type SimpleClientService interface {
	ClientService

	// CanCreate verifies if the given user can create the resource
	CanCreate(ctx context.Context, resource *Resource, subject *Subject) (bool, error)

	// CanRead verifying if the given user has the required permission
	CanRead(ctx context.Context, resource *Resource, subject *Subject) (bool, error)

	// CanWrite verifies if the given user can write to the resource
	CanWrite(ctx context.Context, resource *Resource, subject *Subject) (bool, error)

	// CanDelete verifies if the given user can delete the resource
	CanDelete(ctx context.Context, resource *Resource, subject *Subject) (bool, error)

	// Write adds permissions
	Write(ctx context.Context, tuples []*Tuple) (bool, error)

	// Delete removes permissions
	Delete(ctx context.Context, tuples []*Tuple) (bool, error)

	// DeleteResource removes all permissions for a specific resource
	DeleteResource(ctx context.Context, resource *Resource) (bool, error)

	// DeleteSubject removes all permissions for a specific subject
	DeleteSubject(ctx context.Context, subject *Subject) (bool, error)
}

// simpleClientServiceImpl implements the SimpleClientService interface
type simpleClientServiceImpl struct {
	ClientService
}

func (s *simpleClientServiceImpl) CanCreate(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return canCreate(s, ctx, resource, subject)
}

func (s *simpleClientServiceImpl) CanRead(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return canRead(s, ctx, resource, subject)
}

func (s *simpleClientServiceImpl) CanWrite(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return canWrite(s, ctx, resource, subject)
}

func (s *simpleClientServiceImpl) CanDelete(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return canDelete(s, ctx, resource, subject)
}

// Write adds permissions
func (s *simpleClientServiceImpl) Write(ctx context.Context, tuples []*Tuple) (bool, error) {
	return writeTuples(s, ctx, tuples)
}

// Delete removes permissions
func (s *simpleClientServiceImpl) Delete(ctx context.Context, tuples []*Tuple) (bool, error) {
	return deleteTuples(s, ctx, tuples)
}

// DeleteResource removes all permissions for a specific resource
func (s *simpleClientServiceImpl) DeleteResource(ctx context.Context, resource *Resource) (bool, error) {
	return deleteResource(s, ctx, resource)
}

// DeleteSubject removes all permissions for a specific subject
func (s *simpleClientServiceImpl) DeleteSubject(ctx context.Context, subject *Subject) (bool, error) {
	return deleteSubject(s, ctx, subject)
}

func NewSimpleClientService(cc grpc.ClientConnInterface) (SimpleClientService, error) {
	base, err := NewClientService(cc)
	if err != nil {
		return nil, err
	}
	return &simpleClientServiceImpl{ClientService: base}, nil
}
