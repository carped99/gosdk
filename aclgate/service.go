package aclgate

import (
	"context"
	"fmt"

	v1 "github.com/carped99/gosdk/aclgate/api/gen/aclgate/v1"
	"google.golang.org/grpc"
)

// AclGateService defines the interface for ACL operations
type AclGateService interface {
	// Check verifying if the given user has the required permission
	Check(ctx context.Context, req CheckRequest) (bool, error)

	// BatchCheck verifies multiple permissions at once
	BatchCheck(ctx context.Context, reqs []CheckRequest) ([]BatchCheckResult, error)

	// Write adds permissions
	Write(ctx context.Context, tuples []Tuple) (bool, error)

	// Delete removes permissions
	Delete(ctx context.Context, tuples []Tuple) (bool, error)

	// DeleteResource removes all permissions for a specific resource
	DeleteResource(ctx context.Context, resourceType, resourceId string) (bool, error)

	// DeleteSubject removes all permissions for a specific subject
	DeleteSubject(ctx context.Context, subjectType, subjectId string) (bool, error)

	// Mutate adds or removes permissions (advanced usage)
	Mutate(ctx context.Context, writes, deletes []Tuple) (bool, error)

	// List retrieves permissions based on filters
	List(ctx context.Context, req ListRequest) (*ListResponse, error)

	// Audit retrieves audit logs based on filters
	Audit(ctx context.Context, req AuditRequest) (*AuditResponse, error)

	// StreamCheck streams permission checks in real-time
	StreamCheck(ctx context.Context) (StreamCheckClient, error)
}

// StreamCheckClient represents a client for streaming permission checks
type StreamCheckClient interface {
	// Send sends a permission check request
	Send(request CheckRequest) error
	// Recv receives a permission check response
	Recv() (*StreamCheckResponse, error)
	// Close closes the stream
	Close() error
}

// StreamCheckResponse represents a streaming permission check response
type StreamCheckResponse struct {
	Allowed bool
	Reason  string
	Error   string
}

// aclGateService implements the AclGateService interface
type aclGateService struct {
	client v1.AclGateServiceClient
}

// Check verifying if the given user has the required permission
func (s *aclGateService) Check(ctx context.Context, req CheckRequest) (bool, error) {
	protoReq := &v1.CheckRequest{Tuple: toProtoTuple(req.Tuple)}
	resp, err := s.client.Check(ctx, protoReq)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return resp.GetAllowed(), nil
}

// BatchCheck verifies multiple permissions at once
func (s *aclGateService) BatchCheck(ctx context.Context, reqs []CheckRequest) ([]BatchCheckResult, error) {
	if len(reqs) == 0 {
		return []BatchCheckResult{}, nil
	}

	items := make([]*v1.CheckRequest, 0, len(reqs))
	for _, r := range reqs {
		items = append(items, &v1.CheckRequest{Tuple: toProtoTuple(r.Tuple)})
	}

	resp, err := s.client.BatchCheck(ctx, &v1.BatchCheckRequest{Items: items})
	if err != nil {
		return nil, err
	}

	results := make([]BatchCheckResult, 0, len(resp.GetResults()))
	for _, r := range resp.GetResults() {
		results = append(results, BatchCheckResult{
			Request: CheckRequest{Tuple: toDomainTuple(r.GetRequest().GetTuple())},
			Allowed: r.GetAllowed(),
		})
	}
	return results, nil
}

// Write adds permissions
func (s *aclGateService) Write(ctx context.Context, tuples []Tuple) (bool, error) {
	return s.Mutate(ctx, tuples, nil)
}

// Delete removes permissions
func (s *aclGateService) Delete(ctx context.Context, tuples []Tuple) (bool, error) {
	return s.Mutate(ctx, nil, tuples)
}

// DeleteResource removes all permissions for a specific resource
func (s *aclGateService) DeleteResource(ctx context.Context, resourceType, resourceId string) (bool, error) {
	// Create a wildcard tuple for the resource (empty subject and relation)
	tuple := Tuple{
		ResourceType: resourceType,
		ResourceId:   resourceId,
		SubjectType:  "", // wildcard - all subjects
		SubjectId:    "", // wildcard - all subjects
		Relation:     "", // wildcard - all relations
	}

	return s.Mutate(ctx, nil, []Tuple{tuple})
}

// DeleteSubject removes all permissions for a specific subject
func (s *aclGateService) DeleteSubject(ctx context.Context, subjectType, subjectId string) (bool, error) {
	// Create a wildcard tuple for the subject (empty resource and relation)
	tuple := Tuple{
		ResourceType: "", // wildcard - all resources
		ResourceId:   "", // wildcard - all resources
		SubjectType:  subjectType,
		SubjectId:    subjectId,
		Relation:     "", // wildcard - all relations
	}

	return s.Mutate(ctx, nil, []Tuple{tuple})
}

// Mutate adds or removes permissions
func (s *aclGateService) Mutate(ctx context.Context, writes, deletes []Tuple) (bool, error) {
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

// StreamCheck streams permission checks in real-time
func (s *aclGateService) StreamCheck(ctx context.Context) (StreamCheckClient, error) {
	stream, err := s.client.StreamCheck(ctx)
	if err != nil {
		return nil, err
	}

	return &streamCheckClient{stream: stream}, nil
}

// streamCheckClient implements StreamCheckClient
type streamCheckClient struct {
	stream grpc.BidiStreamingClient[v1.StreamCheckRequest, v1.StreamCheckResponse]
}

// Send sends a permission check request
func (c *streamCheckClient) Send(request CheckRequest) error {
	return c.stream.Send(&v1.StreamCheckRequest{
		Tuple: toProtoTuple(request.Tuple),
	})
}

// Recv receives a permission check response
func (c *streamCheckClient) Recv() (*StreamCheckResponse, error) {
	response, err := c.stream.Recv()
	if err != nil {
		return nil, err
	}

	return &StreamCheckResponse{
		Allowed: response.GetAllowed(),
		Reason:  response.GetReason(),
		Error:   response.GetError(),
	}, nil
}

// Close closes the stream
func (c *streamCheckClient) Close() error {
	return c.stream.CloseSend()
}

// List retrieves permissions based on filters
func (s *aclGateService) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	protoReq := &v1.ListRequest{
		Resource: &v1.Resource{Type: req.Resource.Type, Id: req.Resource.ID},
		Subject:  &v1.Subject{Type: req.Subject.Type, Id: req.Subject.ID},
		Relation: &v1.Relation{Name: req.Relation.Name},
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	resp, err := s.client.List(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	tuples := make([]Tuple, 0, len(resp.GetTuples()))
	for _, t := range resp.GetTuples() {
		tuples = append(tuples, toDomainTuple(t))
	}

	return &ListResponse{Tuples: tuples}, nil
}

// Audit retrieves audit logs based on filters
func (s *aclGateService) Audit(ctx context.Context, req AuditRequest) (*AuditResponse, error) {
	protoReq := &v1.AuditRequest{
		Resource: &v1.Resource{Type: req.Resource.Type, Id: req.Resource.ID},
		Subject:  &v1.Subject{Type: req.Subject.Type, Id: req.Subject.ID},
		Relation: &v1.Relation{Name: req.Relation.Name},
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	resp, err := s.client.Audit(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	logs := make([]AuditLog, 0, len(resp.GetLogs()))
	for _, log := range resp.GetLogs() {
		logs = append(logs, AuditLog{
			ID:        log.GetId(),
			Action:    log.GetAction(),
			Tuple:     toDomainTuple(log.GetTuple()),
			Actor:     log.GetActor(),
			Timestamp: log.GetTimestamp().String(),
			Reason:    log.GetReason(),
		})
	}

	return &AuditResponse{Logs: logs}, nil
}

func NewAclGateService(cc grpc.ClientConnInterface) (AclGateService, error) {
	client := v1.NewAclGateServiceClient(cc)
	return &aclGateService{
		client: client,
	}, nil
}
