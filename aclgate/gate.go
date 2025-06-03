package aclgate

import (
	"context"
	"fmt"
)

// AclGateService implements the AclService interface
type AclGateService struct {
	// Add any necessary fields here
}

// Check verifies if the given user has the required permission
func (s *AclGateService) Check(ctx context.Context, resourceType, resourceID, relation interface{}) error {
	// Convert any type to string using fmt.Sprint
	resourceTypeStr := fmt.Sprint(resourceType)
	resourceIDStr := fmt.Sprint(resourceID)
	relationStr := fmt.Sprint(relation)

	// TODO: Implement actual permission check logic using the string values
	_ = resourceTypeStr
	_ = resourceIDStr
	_ = relationStr

	return nil
}
