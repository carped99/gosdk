package aclgate

import (
	"context"
	"fmt"

	v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"
	"google.golang.org/grpc"
)

// ClientService defines the interface for ACL operations
type ClientService interface {
	// Check verifying if the given user has the required permission
	Check(ctx context.Context, req *CheckRequest) (bool, error)

	// BatchCheck verifies multiple permissions at once
	BatchCheck(ctx context.Context, reqs []*CheckRequest) ([]*BatchCheckResult, error)

	// Mutate adds or removes permissions (advanced usage)
	Mutate(ctx context.Context, writes, deletes []*Tuple) (bool, error)

	// ListResources retrieves resources based on filters
	ListResources(ctx context.Context, req *ListResourcesRequest) (*ListResourcesResponse, error)

	// ListSubjects retrieves subjects based on filters
	ListSubjects(ctx context.Context, req *ListSubjectsRequest) (*ListSubjectsResponse, error)

	// Audit retrieves audit logs based on filters
	Audit(ctx context.Context, req *AuditRequest) (*AuditResponse, error)
}

// clientServiceImpl implements the ClientService interface
type clientServiceImpl struct {
	client v1.AclGateServiceClient
}

// Check verifying if the given user has the required permission
func (s *clientServiceImpl) Check(ctx context.Context, req *CheckRequest) (bool, error) {
	protoReq := &v1.CheckRequest{Tuple: toProtoTuple(req.Tuple)}
	resp, err := s.client.Check(ctx, protoReq)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return resp.GetAllowed(), nil
}

// BatchCheck verifies multiple permissions at once
func (s *clientServiceImpl) BatchCheck(ctx context.Context, reqs []*CheckRequest) ([]*BatchCheckResult, error) {
	if len(reqs) == 0 {
		return []*BatchCheckResult{}, nil
	}

	items := make([]*v1.CheckRequest, 0, len(reqs))
	for _, r := range reqs {
		items = append(items, &v1.CheckRequest{Tuple: toProtoTuple(r.Tuple)})
	}

	resp, err := s.client.BatchCheck(ctx, &v1.BatchCheckRequest{Items: items})
	if err != nil {
		return nil, err
	}

	results := make([]*BatchCheckResult, 0, len(resp.GetResults()))
	for _, r := range resp.GetResults() {
		tuple, err := toDomainTuple(r.GetRequest().GetTuple())
		if err != nil {
			return nil, err
		}

		results = append(results, &BatchCheckResult{
			Request: &CheckRequest{Tuple: tuple},
			Allowed: r.GetAllowed(),
		})
	}
	return results, nil
}

// Mutate adds or removes permissions
func (s *clientServiceImpl) Mutate(ctx context.Context, writes, deletes []*Tuple) (bool, error) {
	protoWrites := make([]*v1.Tuple, 0, len(writes))
	for _, t := range writes {
		protoWrites = append(protoWrites, toProtoTuple(t))
	}

	protoDeletes := make([]*v1.Tuple, 0, len(deletes))
	for _, t := range deletes {
		protoDeletes = append(protoDeletes, toProtoTuple(t))
	}

	resp, err := s.client.Mutate(ctx, &v1.MutateRequest{
		Writes:  protoWrites,
		Deletes: protoDeletes,
	})
	if err != nil {
		return false, err
	}
	return resp.GetSuccess(), nil
}

// ListResources retrieves resources based on filters
func (s *clientServiceImpl) ListResources(ctx context.Context, req *ListResourcesRequest) (*ListResourcesResponse, error) {
	protoReq := &v1.ListResourcesRequest{
		Type:     req.Type,
		Subject:  &v1.Subject{Type: req.Subject.Type, Id: req.Subject.ID},
		Relation: &v1.Relation{Name: req.Relation.Name},
	}

	resp, err := s.client.ListResources(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %v, %w", req, err)
	}

	resources, err := toDomainResources(resp.GetResources())
	if err != nil {
		return nil, fmt.Errorf("failed to convert resources: %w", err)
	}

	return &ListResourcesResponse{Resources: resources}, nil
}

// ListSubjects retrieves subjects based on filters
func (s *clientServiceImpl) ListSubjects(ctx context.Context, req *ListSubjectsRequest) (*ListSubjectsResponse, error) {
	protoReq := &v1.ListSubjectsRequest{
		Type:     req.Type,
		Resource: &v1.Resource{Type: req.Resource.Type, Id: req.Resource.ID},
		Relation: &v1.Relation{Name: req.Relation.Name},
	}

	resp, err := s.client.ListSubjects(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list subjects: %v, %w", req, err)
	}

	subjects, err := toDomainSubjects(resp.GetSubjects())
	if err != nil {
		return nil, fmt.Errorf("failed to convert subjects: %w", err)
	}

	return &ListSubjectsResponse{Subjects: subjects}, nil
}

// Audit retrieves audit logs based on filters
func (s *clientServiceImpl) Audit(ctx context.Context, req *AuditRequest) (*AuditResponse, error) {
	protoReq := &v1.AuditRequest{
		Resource: &v1.Resource{Type: req.Resource.Type, Id: req.Resource.ID},
		Subject:  &v1.Subject{Type: req.Subject.Type, Id: req.Subject.ID},
		Relation: &v1.Relation{Name: req.Relation.Name},
		PageSize: req.PageSize,
		Cursor:   req.Cursor,
	}

	resp, err := s.client.Audit(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	logs := make([]AuditLog, 0, len(resp.GetLogs()))
	for _, log := range resp.GetLogs() {
		tuple, err := toDomainTuple(log.GetTuple())
		if err != nil {
			return nil, err
		}
		logs = append(logs, AuditLog{
			ID:        log.GetId(),
			Action:    log.GetAction(),
			Tuple:     tuple,
			Actor:     log.GetActor(),
			Timestamp: log.GetTimestamp().String(),
			Reason:    log.GetReason(),
		})
	}

	return &AuditResponse{Logs: logs}, nil
}

func NewClientService(cc grpc.ClientConnInterface) (ClientService, error) {
	client := v1.NewAclGateServiceClient(cc)
	return &clientServiceImpl{client: client}, nil
}
