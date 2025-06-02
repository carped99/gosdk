package outbox

import (
	"fmt"
	"regexp"
)

var (
	resourceTypeRegex = regexp.MustCompile(`^[^:#\s]+$`)
	resourceIDRegex   = regexp.MustCompile(`^[^:#\s]+$`)
	subjectTypeRegex  = regexp.MustCompile(`^[^:#\s]+$`)
	subjectIDRegex    = regexp.MustCompile(`^[^:#\s]+$`)
	relationRegex     = regexp.MustCompile(`^[^:#@\s]+$`)
)

type AclTuple struct {
	Resource *AclResource `json:"resource,omitempty"`
	Subject  *AclSubject  `json:"subject,omitempty"`
	Relation *AclRelation `json:"relation,omitempty"`
}

func (p *AclTuple) Validate() error {
	if p.Resource != nil {
		if err := p.Resource.Validate(); err != nil {
			return err
		}
	}

	if p.Subject != nil {
		if err := p.Subject.Validate(); err != nil {
			return err
		}
	}

	if p.Relation != nil {
		if err := p.Relation.Validate(); err != nil {
			return err
		}
	}
	return nil
}

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

type AclResource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (p *AclResource) Validate() error {
	if !resourceTypeRegex.MatchString(p.Type) {
		return fmt.Errorf("invalid resource type: '%s'", p.Type)
	}

	if !resourceIDRegex.MatchString(p.ID) {
		return fmt.Errorf("invalid resource id: '%s'", p.ID)
	}

	return nil
}

type AclSubject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (p *AclSubject) Validate() error {
	if !subjectTypeRegex.MatchString(p.Type) {
		return fmt.Errorf("invalid subject type: '%s'", p.Type)
	}

	if !subjectIDRegex.MatchString(p.ID) {
		return fmt.Errorf("invalid subject id: '%s'", p.ID)
	}

	return nil
}

type AclRelation struct {
	Name string `json:"name"`
}

func (p *AclRelation) Validate() error {
	if !relationRegex.MatchString(p.Name) {
		return fmt.Errorf("invalid relation name: '%s'", p.Name)
	}
	return nil
}
