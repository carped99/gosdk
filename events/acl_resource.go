package events

import (
	"fmt"
	"regexp"
)

var (
	resourceTypeRegex = regexp.MustCompile(`^[^:#\s]+$`)
	resourceIDRegex   = regexp.MustCompile(`^[^:#\s]+$`)
)

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
