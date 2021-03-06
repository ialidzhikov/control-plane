package runtime

import (
	"strings"

	pkg "github.com/kyma-project/control-plane/components/kyma-environment-broker/common/runtime"
	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/pkg/errors"
)

type converter struct {
	defaultSubaccountRegion string
}

func newConverter(platformRegion string) *converter {
	return &converter{
		defaultSubaccountRegion: platformRegion,
	}
}

func (c *converter) setRegionOrDefault(instance internal.Instance, runtime *pkg.RuntimeDTO) error {
	pp, err := instance.GetProvisioningParameters()
	if err != nil {
		return errors.Wrap(err, "while getting provisioning parameters")
	}

	if pp.PlatformRegion == "" {
		runtime.SubAccountRegion = c.defaultSubaccountRegion
	} else {
		runtime.SubAccountRegion = pp.PlatformRegion
	}
	return nil
}

func (c *converter) ApplyProvisioningOperation(dto *pkg.RuntimeDTO, pOpr *internal.ProvisioningOperation) {
	if pOpr != nil {
		c.applyOperation(&pOpr.Operation, dto.Status.Provisioning)
	}
}

func (c *converter) ApplyDeprovisioningOperation(dto *pkg.RuntimeDTO, dOpr *internal.DeprovisioningOperation) {
	if dOpr != nil {
		dto.Status.Deprovisioning = &pkg.Operation{}
		c.applyOperation(&dOpr.Operation, dto.Status.Deprovisioning)
	}
}

func (c *converter) applyOperation(source *internal.Operation, target *pkg.Operation) {
	if source != nil {
		target.OperationID = source.ID
		target.CreatedAt = source.CreatedAt
		target.State = string(source.State)
		target.Description = source.Description
		if source.OrchestrationID != "" {
			target.OrchestrationID = &source.OrchestrationID
		}
	}
}

func (c *converter) NewDTO(instance internal.Instance) (pkg.RuntimeDTO, error) {
	toReturn := pkg.RuntimeDTO{
		InstanceID:       instance.InstanceID,
		RuntimeID:        instance.RuntimeID,
		GlobalAccountID:  instance.GlobalAccountID,
		SubAccountID:     instance.SubAccountID,
		ServiceClassID:   instance.ServiceID,
		ServiceClassName: instance.ServiceName,
		ServicePlanID:    instance.ServicePlanID,
		ServicePlanName:  instance.ServicePlanName,
		ProviderRegion:   instance.ProviderRegion,
		Status: pkg.RuntimeStatus{
			CreatedAt:    instance.CreatedAt,
			ModifiedAt:   instance.UpdatedAt,
			Provisioning: &pkg.Operation{},
		},
	}

	err := c.setRegionOrDefault(instance, &toReturn)
	if err != nil {
		return pkg.RuntimeDTO{}, errors.Wrap(err, "while setting region")
	}

	urlSplitted := strings.Split(instance.DashboardURL, ".")
	if len(urlSplitted) > 1 {
		toReturn.ShootName = urlSplitted[1]
	}

	return toReturn, nil
}

func (c *converter) ApplyUpgradingKymaOperations(dto *pkg.RuntimeDTO, oprs []internal.UpgradeKymaOperation, totalCount int) {
	dto.Status.UpgradingKyma.TotalCount = totalCount
	dto.Status.UpgradingKyma.Count = len(oprs)
	dto.Status.UpgradingKyma.Data = make([]pkg.Operation, 0)
	for _, o := range oprs {
		op := pkg.Operation{}
		c.applyOperation(&o.Operation, &op)
		dto.Status.UpgradingKyma.Data = append(dto.Status.UpgradingKyma.Data, op)
	}
}
