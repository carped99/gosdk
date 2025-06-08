package events

import (
	"fmt"
	"regexp"
)

var (
	subjectTypeRegex = regexp.MustCompile(`^[^:#\s]+$`)
	subjectIDRegex   = regexp.MustCompile(`^[^:#\s]+$`)
)

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
