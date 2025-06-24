package aclgate

import "context"

func canCreate(clientService ClientService, ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return checkTuple(clientService, ctx, Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: NewRelation("can_create"),
	})
}

func canRead(clientService ClientService, ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return checkTuple(clientService, ctx, Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: NewRelation("can_read"),
	})
}

func canWrite(clientService ClientService, ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return checkTuple(clientService, ctx, Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: NewRelation("can_write"),
	})
}

func canDelete(clientService ClientService, ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return checkTuple(clientService, ctx, Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: NewRelation("can_delete"),
	})
}

func checkTuple(clientService ClientService, ctx context.Context, tuple Tuple) (bool, error) {
	return clientService.Check(ctx, CheckRequest{Tuple: tuple})
}

func writeTuples(clientService ClientService, ctx context.Context, tuples []Tuple) (bool, error) {
	return clientService.Mutate(ctx, tuples, nil)
}

func deleteTuples(clientService ClientService, ctx context.Context, tuples []Tuple) (bool, error) {
	return clientService.Mutate(ctx, nil, tuples)
}

func deleteResource(clientService ClientService, ctx context.Context, resource Resource) (bool, error) {
	return clientService.Mutate(ctx, nil, []Tuple{
		{
			Resource: resource,
		},
	})
}

func deleteSubject(clientService ClientService, ctx context.Context, subject Subject) (bool, error) {
	return clientService.Mutate(ctx, nil, []Tuple{
		{
			Subject: subject,
		},
	})
}
