package events

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
