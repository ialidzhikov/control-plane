package memory

import (
	"sort"
	"sync"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/common/pagination"

	"github.com/pivotal-cf/brokerapi/v7/domain"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage/dberr"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal/storage/dbsession/dbmodel"
)

type operations struct {
	mu sync.Mutex

	provisioningOperations   map[string]internal.ProvisioningOperation
	deprovisioningOperations map[string]internal.DeprovisioningOperation
	upgradeKymaOperations    map[string]internal.UpgradeKymaOperation
}

// NewOperation creates in-memory storage for OSB operations.
func NewOperation() *operations {
	return &operations{
		provisioningOperations:   make(map[string]internal.ProvisioningOperation, 0),
		deprovisioningOperations: make(map[string]internal.DeprovisioningOperation, 0),
		upgradeKymaOperations:    make(map[string]internal.UpgradeKymaOperation, 0),
	}
}

func (s *operations) InsertProvisioningOperation(operation internal.ProvisioningOperation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := operation.ID
	if _, exists := s.provisioningOperations[id]; exists {
		return dberr.AlreadyExists("instance operation with id %s already exist", id)
	}

	s.provisioningOperations[id] = operation
	return nil
}

func (s *operations) GetProvisioningOperationByID(operationID string) (*internal.ProvisioningOperation, error) {
	op, exists := s.provisioningOperations[operationID]
	if !exists {
		return nil, dberr.NotFound("instance provisioning operation with id %s not found", operationID)
	}
	return &op, nil
}

func (s *operations) GetProvisioningOperationByInstanceID(instanceID string) (*internal.ProvisioningOperation, error) {
	for _, op := range s.provisioningOperations {
		if op.InstanceID == instanceID {
			return &op, nil
		}
	}
	return nil, dberr.NotFound("instance provisioning operation with instanceID %s not found", instanceID)
}

func (s *operations) UpdateProvisioningOperation(op internal.ProvisioningOperation) (*internal.ProvisioningOperation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldOp, exists := s.provisioningOperations[op.ID]
	if !exists {
		return nil, dberr.NotFound("instance operation with id %s not found", op.ID)
	}
	if oldOp.Version != op.Version {
		return nil, dberr.Conflict("unable to update provisioning operation with id %s (for instance id %s) - conflict", op.ID, op.InstanceID)
	}
	op.Version = op.Version + 1
	s.provisioningOperations[op.ID] = op

	return &op, nil
}

func (s *operations) InsertDeprovisioningOperation(operation internal.DeprovisioningOperation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := operation.ID
	if _, exists := s.deprovisioningOperations[id]; exists {
		return dberr.AlreadyExists("instance operation with id %s already exist", id)
	}

	s.deprovisioningOperations[id] = operation
	return nil
}

func (s *operations) GetDeprovisioningOperationByID(operationID string) (*internal.DeprovisioningOperation, error) {
	op, exists := s.deprovisioningOperations[operationID]
	if !exists {
		return nil, dberr.NotFound("instance deprovisioning operation with id %s not found", operationID)
	}
	return &op, nil
}

func (s *operations) GetDeprovisioningOperationByInstanceID(instanceID string) (*internal.DeprovisioningOperation, error) {
	for _, op := range s.deprovisioningOperations {
		if op.InstanceID == instanceID {
			return &op, nil
		}
	}

	return nil, dberr.NotFound("instance deprovisioning operation with instanceID %s not found", instanceID)
}

func (s *operations) UpdateDeprovisioningOperation(op internal.DeprovisioningOperation) (*internal.DeprovisioningOperation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldOp, exists := s.deprovisioningOperations[op.ID]
	if !exists {
		return nil, dberr.NotFound("instance operation with id %s not found", op.ID)
	}
	if oldOp.Version != op.Version {
		return nil, dberr.Conflict("unable to update deprovisioning operation with id %s (for instance id %s) - conflict", op.ID, op.InstanceID)
	}
	op.Version = op.Version + 1
	s.deprovisioningOperations[op.ID] = op

	return &op, nil
}

func (s *operations) InsertUpgradeKymaOperation(operation internal.UpgradeKymaOperation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := operation.ID
	if _, exists := s.upgradeKymaOperations[id]; exists {
		return dberr.AlreadyExists("instance operation with id %s already exist", id)
	}

	s.upgradeKymaOperations[id] = operation
	return nil
}

func (s *operations) GetUpgradeKymaOperationByID(operationID string) (*internal.UpgradeKymaOperation, error) {
	op, exists := s.upgradeKymaOperations[operationID]
	if !exists {
		return nil, dberr.NotFound("instance upgradeKyma operation with id %s not found", operationID)
	}
	return &op, nil
}

func (s *operations) GetUpgradeKymaOperationByInstanceID(instanceID string) (*internal.UpgradeKymaOperation, error) {
	for _, op := range s.upgradeKymaOperations {
		if op.InstanceID == instanceID {
			return &op, nil
		}
	}

	return nil, dberr.NotFound("instance upgradeKyma operation with instanceID %s not found", instanceID)
}

func (s *operations) UpdateUpgradeKymaOperation(op internal.UpgradeKymaOperation) (*internal.UpgradeKymaOperation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldOp, exists := s.upgradeKymaOperations[op.ID]
	if !exists {
		return nil, dberr.NotFound("instance operation with id %s not found", op.ID)
	}
	if oldOp.Version != op.Version {
		return nil, dberr.Conflict("unable to update upgradeKyma operation with id %s (for instance id %s) - conflict", op.ID, op.InstanceID)
	}
	op.Version = op.Version + 1
	s.upgradeKymaOperations[op.ID] = op

	return &op, nil
}

func (s *operations) GetOperationByID(operationID string) (*internal.Operation, error) {
	var res *internal.Operation

	provisionOp, exists := s.provisioningOperations[operationID]
	if exists {
		res = &provisionOp.Operation
	}
	deprovisionOp, exists := s.deprovisioningOperations[operationID]
	if exists {
		res = &deprovisionOp.Operation
	}
	upgradeKymaOp, exists := s.upgradeKymaOperations[operationID]
	if exists {
		res = &upgradeKymaOp.Operation
	}
	if res == nil {
		return nil, dberr.NotFound("instance operation with id %s not found", operationID)
	}

	return res, nil
}

func (s *operations) GetOperationsInProgressByType(opType dbmodel.OperationType) ([]internal.Operation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ops := make([]internal.Operation, 0)
	switch opType {
	case dbmodel.OperationTypeProvision:
		for _, op := range s.provisioningOperations {
			if op.State == domain.InProgress {
				ops = append(ops, op.Operation)
			}
		}
	case dbmodel.OperationTypeDeprovision:
		for _, op := range s.deprovisioningOperations {
			if op.State == domain.InProgress {
				ops = append(ops, op.Operation)
			}
		}
	}

	return ops, nil
}

func (s *operations) GetOperationsForIDs(opIdList []string) ([]internal.Operation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ops := make([]internal.Operation, 0)
	for _, opID := range opIdList {
		for _, op := range s.upgradeKymaOperations {
			if op.Operation.ID == opID {
				ops = append(ops, op.Operation)
			}
		}
	}

	for _, opID := range opIdList {
		for _, op := range s.provisioningOperations {
			if op.Operation.ID == opID {
				ops = append(ops, op.Operation)
			}
		}
	}

	for _, opID := range opIdList {
		for _, op := range s.deprovisioningOperations {
			if op.Operation.ID == opID {
				ops = append(ops, op.Operation)
			}
		}
	}
	if len(ops) == 0 {
		return nil, dberr.NotFound("operations with ids from list %+q not exist", opIdList)
	}

	return ops, nil
}

func (s *operations) GetOperationStats() (internal.OperationStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := internal.OperationStats{
		Provisioning:   map[domain.LastOperationState]int{domain.InProgress: 0, domain.Succeeded: 0, domain.Failed: 0},
		Deprovisioning: map[domain.LastOperationState]int{domain.InProgress: 0, domain.Succeeded: 0, domain.Failed: 0},
	}

	for _, op := range s.provisioningOperations {
		result.Provisioning[op.State] = result.Provisioning[op.State] + 1
	}
	for _, op := range s.deprovisioningOperations {
		result.Deprovisioning[op.State] = result.Deprovisioning[op.State] + 1
	}
	return result, nil
}

func (s *operations) GetOperationStatsForOrchestration(orchestrationID string) (map[domain.LastOperationState]int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := map[domain.LastOperationState]int{
		domain.InProgress: 0,
		domain.Succeeded:  0,
		domain.Failed:     0,
	}
	for _, op := range s.upgradeKymaOperations {
		result[op.State] = result[op.State] + 1
	}
	return result, nil
}

func (s *operations) ListUpgradeKymaOperationsByOrchestrationID(orchestrationID string, pageSize, page int) ([]internal.UpgradeKymaOperation, int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]internal.UpgradeKymaOperation, 0)

	for _, op := range s.upgradeKymaOperations {
		if op.OrchestrationID == orchestrationID {
			result = append(result, op)
		}
	}
	offset := pagination.ConvertPageAndPageSizeToOffset(pageSize, page)

	sortedOperations := s.getUpgradeSortedByCreatedAt(s.upgradeKymaOperations)
	result = make([]internal.UpgradeKymaOperation, 0)

	for i := offset; i < offset+pageSize && i < len(sortedOperations)+offset; i++ {
		result = append(result, s.upgradeKymaOperations[sortedOperations[i].OrchestrationID])
	}

	return result,
		len(result),
		len(s.upgradeKymaOperations),
		nil
}

func (s *operations) ListUpgradeKymaOperationsByInstanceID(instanceID string) ([]internal.UpgradeKymaOperation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]internal.UpgradeKymaOperation, 0)

	for _, op := range s.upgradeKymaOperations {
		if op.InstanceID == instanceID {
			result = append(result, op)
		}
	}

	sortedOperations := s.getUpgradeSortedByCreatedAt(s.upgradeKymaOperations)
	result = make([]internal.UpgradeKymaOperation, 0)

	return sortedOperations, nil
}

func (s *operations) getUpgradeSortedByCreatedAt(operations map[string]internal.UpgradeKymaOperation) []internal.UpgradeKymaOperation {
	operationsList := make([]internal.UpgradeKymaOperation, 0, len(operations))
	for _, v := range operations {
		operationsList = append(operationsList, v)
	}
	sort.Slice(operationsList, func(i, j int) bool {
		return operationsList[i].CreatedAt.Before(operationsList[j].CreatedAt)
	})
	return operationsList
}
