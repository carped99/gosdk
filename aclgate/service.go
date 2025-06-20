package aclgate

import (
	"context"
	v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"
	"google.golang.org/grpc"
)

// AclGateService defines the interface for ACL operations
type AclGateService interface {
	// Check verifies if the given user has the required permission
	Check(ctx context.Context, resourceType, resourceId, subjectType, subjectId, relation string) (bool, error)
}

// aclGateService implements the AclGateService interface
type aclGateService struct {
	client v1.AclGateServiceClient
}

// Check verifies if the given user has the required permission
func (s *aclGateService) Check(ctx context.Context, resourceType, resourceId, subjectType, subjectId, relation string) (bool, error) {
	response, err := s.client.Check(ctx, &v1.CheckRequest{
		Tuple: &v1.Tuple{
			Resource: &v1.Resource{
				Type: resourceType,
				Id:   resourceId,
			},
			Subject: &v1.Subject{
				Type: subjectType,
				Id:   subjectId,
			},
			Relation: &v1.Relation{
				Name: relation,
			},
		},
	})

	if err != nil {
		return false, err
	}

	return response.GetAllowed(), nil
}

func (s *aclGateService) MustCheck(ctx context.Context, resourceType, resourceId, subjectType, subjectId, relation string) bool {
	result, err := s.Check(ctx, resourceType, resourceId, subjectType, subjectId, relation)
	if err != nil {
		return false
	}
	return result
}

func NewAclGateService(cc grpc.ClientConnInterface) (AclGateService, error) {
	client := v1.NewAclGateServiceClient(cc)
	return &aclGateService{
		client: client,
	}, nil
}
