package events

type AclTuples struct {
	Tuples []*AclTuple `json:"tuples,omitempty"`
}

func (p *AclTuples) Validate() error {
	for _, tuple := range p.Tuples {
		if err := tuple.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type AclTuplesBuilder struct {
	tuples *AclTuples
}

func (b *AclTuplesBuilder) AddTuples(tuples ...*AclTuple) *AclTuplesBuilder {
	b.tuples.Tuples = append(b.tuples.Tuples, tuples...)
	return b
}

// Build returns the constructed AclTuples
func (b *AclTuplesBuilder) Build() (*AclTuples, error) {
	return b.tuples, b.tuples.Validate()
}

func NewAclTuplesBuilder() *AclTuplesBuilder {
	return &AclTuplesBuilder{
		tuples: &AclTuples{},
	}
}
