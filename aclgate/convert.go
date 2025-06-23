package aclgate

import v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"

func toProtoTuple(t Tuple) *v1.Tuple {
	return &v1.Tuple{
		Resource: &v1.Resource{Type: t.Resource.Type, Id: t.Resource.ID},
		Subject:  &v1.Subject{Type: t.Subject.Type, Id: t.Subject.ID},
		Relation: &v1.Relation{Name: t.Relation.Name},
	}
}

func toDomainTuple(t *v1.Tuple) Tuple {
	if t == nil {
		return Tuple{}
	}

	resource := t.GetResource()
	subject := t.GetSubject()
	relation := t.GetRelation()

	return Tuple{
		Resource: Resource{Type: resource.GetType(), ID: resource.GetId()},
		Subject:  Subject{Type: subject.GetType(), ID: subject.GetId()},
		Relation: Relation{Name: relation.GetName()},
	}
}
