package events

import (
	"fmt"
	"regexp"
)

var (
	relationRegex = regexp.MustCompile(`^[^:#@\s]+$`)
)

type AclRelation struct {
	Name string `json:"name"`
}

func (p *AclRelation) Validate() error {
	if !relationRegex.MatchString(p.Name) {
		return fmt.Errorf("invalid relation name: '%s'", p.Name)
	}
	return nil
}
