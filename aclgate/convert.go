package aclgate

import v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"

func toProtoTuple(t Tuple) *v1.Tuple {
	return &v1.Tuple{
		Resource: &v1.Resource{Type: t.ResourceType, Id: t.ResourceId},
		Subject:  &v1.Subject{Type: t.SubjectType, Id: t.SubjectId},
		Relation: &v1.Relation{Name: t.Relation},
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
		ResourceType: resource.GetType(),
		ResourceId:   resource.GetId(),
		SubjectType:  subject.GetType(),
		SubjectId:    subject.GetId(),
		Relation:     relation.GetName(),
	}
}
