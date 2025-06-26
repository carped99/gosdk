package aclgate

import "context"

func canCreate(clientService ClientService, ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return checkTupleRelation(clientService, ctx, resource, subject, "can_create")
}

func canRead(clientService ClientService, ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return checkTupleRelation(clientService, ctx, resource, subject, "can_read")
}

func canWrite(clientService ClientService, ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return checkTupleRelation(clientService, ctx, resource, subject, "can_write")
}

func canDelete(clientService ClientService, ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	return checkTupleRelation(clientService, ctx, resource, subject, "can_delete")
}

func checkTupleRelation(clientService ClientService, ctx context.Context, resource *Resource, subject *Subject, relationName string) (bool, error) {
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
	return clientService.Check(ctx, &CheckRequest{Tuple: tuple})
}

func writeTuples(clientService ClientService, ctx context.Context, tuples []*Tuple) (bool, error) {
	return clientService.Mutate(ctx, tuples, nil)
}

func deleteTuples(clientService ClientService, ctx context.Context, tuples []*Tuple) (bool, error) {
	return clientService.Mutate(ctx, nil, tuples)
}

func deleteResource(clientService ClientService, ctx context.Context, resource *Resource) (bool, error) {
	return clientService.Mutate(ctx, nil, []*Tuple{
		{
			Resource: resource,
		},
	})
}

func deleteSubject(clientService ClientService, ctx context.Context, subject *Subject) (bool, error) {
	return clientService.Mutate(ctx, nil, []*Tuple{
		{
			Subject: subject,
		},
	})
}
