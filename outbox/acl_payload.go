package outbox

type AclMessageBuilder struct {
	*MessageBuilder
}

func NewAclMessageBuilder() *AclMessageBuilder {
	builder := &AclMessageBuilder{
		MessageBuilder: NewMessageBuilder(),
	}

	builder.SetEventTopic("acls.events")

	return builder
}

func (b *AclMessageBuilder) Build() (*Message, error) {
	return b.Build()
}

type AclPayload struct {
	Tuples []*AclTuple `json:"tuples,omitempty"`
}

func (p *AclPayload) Validate() error {
	for _, tuple := range p.Tuples {
		if err := tuple.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type AclPayloadBuilder struct {
	payload *AclPayload
}

func (b *AclPayloadBuilder) AddTuples(tuples ...*AclTuple) *AclPayloadBuilder {
	b.payload.Tuples = append(b.payload.Tuples, tuples...)
	return b
}

// Build returns the constructed AclPayload
func (b *AclPayloadBuilder) Build() (*AclPayload, error) {
	return b.payload, b.payload.Validate()
}

func NewAclPayloadBuilder() *AclPayloadBuilder {
	return &AclPayloadBuilder{
		payload: &AclPayload{},
	}
}
