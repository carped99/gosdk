package aclgate

import (
	"context"
	"fmt"
)

// Common relation names for standard permissions
const (
	RelationCanCreate = "can_create"
	RelationCanRead   = "can_read"
	RelationCanWrite  = "can_write"
	RelationCanDelete = "can_delete"
)

func checkTupleRelation(clientService ClientService, ctx context.Context, resource *Resource, subject *Subject, relationName string) (bool, error) {
	if resource == nil {
		return false, fmt.Errorf("resource cannot be nil")
	}
	if subject == nil {
		return false, fmt.Errorf("subject cannot be nil")
	}

	relation, err := NewRelation(relationName)
	if err != nil {
		return false, err
	}

	return checkTuple(clientService, ctx, &Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: relation,
	})
}

func checkTuple(clientService ClientService, ctx context.Context, tuple *Tuple) (bool, error) {
	if tuple == nil {
		return false, nil // No tuple to check, consider it false
	}

	return clientService.Check(ctx, &CheckRequest{Tuple: tuple})
}
