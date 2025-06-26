package aclgate

import (
	v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"
)

func toProtoTuple(t *Tuple) *v1.Tuple {
	if t == nil {
		return nil
	}

	var (
		resource *v1.Resource
		subject  *v1.Subject
		relation *v1.Relation
	)

	if t.Resource != nil {
		resource = &v1.Resource{
			Type: t.Resource.Type,
			Id:   t.Resource.ID,
		}
	}

	if t.Subject != nil {
		subject = &v1.Subject{
			Type: t.Subject.Type,
			Id:   t.Subject.ID,
		}
	}

	if t.Relation != nil {
		relation = &v1.Relation{
			Name: t.Relation.Name,
		}
	}

	return &v1.Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: relation,
	}
}

func toDomainTuple(t *v1.Tuple) (*Tuple, error) {
	if t == nil {
		return nil, nil
	}

	var (
		resource *Resource
		subject  *Subject
		relation *Relation
		err      error
	)

	if resource, err = toDomainResource(t.GetResource()); err != nil {
		return nil, err
	}

	if subject, err = toDomainSubject(t.GetSubject()); err != nil {
		return nil, err
	}

	if relation, err = toDomainRelation(t.GetRelation()); err != nil {
		return nil, err
	}

	return &Tuple{
		Resource: resource,
		Subject:  subject,
		Relation: relation,
	}, nil
}

func toDomainResource(r *v1.Resource) (*Resource, error) {
	if r == nil {
		return nil, nil
	}
	return NewResource(r.GetType(), r.GetId())
}

// toDomainResources optimizes memory allocation and provides early return on errors
func toDomainResources(values []*v1.Resource) ([]*Resource, error) {
	if len(values) == 0 {
		return nil, nil
	}

	result := make([]*Resource, 0, len(values))
	for _, it := range values {
		if res, err := toDomainResource(it); err != nil {
			return nil, err
		} else if res != nil {
			result = append(result, res)
		}
	}
	return result, nil
}

func toDomainSubject(s *v1.Subject) (*Subject, error) {
	if s == nil {
		return nil, nil
	}
	return NewSubject(s.GetType(), s.GetId())
}

// toDomainSubjects optimizes memory allocation and provides early return on errors
func toDomainSubjects(values []*v1.Subject) ([]*Subject, error) {
	if len(values) == 0 {
		return nil, nil
	}

	result := make([]*Subject, 0, len(values))
	for _, it := range values {
		if res, err := toDomainSubject(it); err != nil {
			return nil, err
		} else if res != nil {
			result = append(result, res)
		}
	}
	return result, nil
}

func toDomainRelation(r *v1.Relation) (*Relation, error) {
	if r == nil {
		return nil, nil
	}
	return NewRelation(r.GetName())
}

// toDomainRelations optimizes memory allocation and provides early return on errors
func toDomainRelations(values []*v1.Relation) ([]*Relation, error) {
	if len(values) == 0 {
		return nil, nil
	}

	result := make([]*Relation, 0, len(values))
	for _, it := range values {
		if res, err := toDomainRelation(it); err != nil {
			return nil, err
		} else if res != nil {
			result = append(result, res)
		}
	}
	return result, nil
}
