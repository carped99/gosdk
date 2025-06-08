package events

type AclTupleBuilder struct {
	payload *AclTuple
}

func NewAclTupleBuilder() *AclTupleBuilder {
	return &AclTupleBuilder{
		payload: &AclTuple{},
	}
}

func (b *AclTupleBuilder) SetResourceType(resourceType string) *AclTupleBuilder {
	if b.payload.Resource == nil {
		b.payload.Resource = &AclResource{}
	}
	b.payload.Resource.Type = resourceType
	return b
}

func (b *AclTupleBuilder) SetResourceID(resourceId string) *AclTupleBuilder {
	if b.payload.Resource == nil {
		b.payload.Resource = &AclResource{}
	}
	b.payload.Resource.ID = resourceId
	return b
}

func (b *AclTupleBuilder) SetSubjectType(subjectType string) *AclTupleBuilder {
	if b.payload.Subject == nil {
		b.payload.Subject = &AclSubject{}
	}
	b.payload.Subject.Type = subjectType
	return b
}

func (b *AclTupleBuilder) SetSubjectID(subjectId string) *AclTupleBuilder {
	if b.payload.Subject == nil {
		b.payload.Subject = &AclSubject{}
	}
	b.payload.Subject.ID = subjectId
	return b
}

func (b *AclTupleBuilder) SetRelationName(relationName string) *AclTupleBuilder {
	if b.payload.Relation == nil {
		b.payload.Relation = &AclRelation{}
	}
	b.payload.Relation.Name = relationName
	return b
}

func (b *AclTupleBuilder) Build() (*AclTuple, error) {
	if b.payload.Resource != nil {
		if err := b.payload.Resource.Validate(); err != nil {
			return nil, err
		}
	}

	if b.payload.Subject != nil {
		if err := b.payload.Subject.Validate(); err != nil {
			return nil, err
		}
	}

	if b.payload.Relation != nil {
		if err := b.payload.Relation.Validate(); err != nil {
			return nil, err
		}
	}
	return b.payload, nil
}
