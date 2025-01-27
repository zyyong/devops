package appservice

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/appproject"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/appproject"
	"github.com/yametech/devops/pkg/service"
)

type NamespaceConfigService struct {
	service.IService
}

func NewNamespaceConfigService(i service.IService) *NamespaceConfigService {
	return &NamespaceConfigService{IService: i}
}

func (n *NamespaceConfigService) GetByFilter(appid string) (core.IObject, error) {
	req := &appproject.Resource{
		Spec: appproject.ResourceSpec{
			App: appid,
		},
	}

	if err := n.IService.GetByFilter(common.DefaultNamespace, common.Resource, req, map[string]interface{}{
		"spec.app": req.Spec.App,
	}); err != nil {
		return nil, err
	}

	return req, nil
}

func (n *NamespaceConfigService) Update(data *apiResource.NameSpaceRequest) (core.IObject, bool, error) {

	app := &appproject.AppProject{}
	if err := n.GetByUUID(common.DefaultNamespace, common.Namespace, data.App, app); err != nil {
		return nil, false, errors.New("The namespace is not exist")
	}

	if app.Spec.AppType != appproject.Namespace {
		return nil, false, errors.New("This is not an Namespace type")
	}

	dbObj := &appproject.Resource{}
	n.IService.GetByFilter(common.DefaultNamespace, common.Resource, dbObj, map[string]interface{}{
		"spec.app": app.Metadata.UUID,
	})

	// create history
	// Get creator
	history := &appproject.ConfigHistory{
		Spec: appproject.HistorySpec{
			App: dbObj.Spec.App,
			History: map[string]interface{}{
				"creator": "",
				"cpu_before": dbObj.Spec.Cpu,
				"memory_before": dbObj.Spec.Memory,
			},
		},
	}

	dbObj.Spec.App = app.Metadata.UUID
	dbObj.Spec.Threshold = data.Threshold
	dbObj.Spec.Approval = data.Approval
	dbObj.Spec.Cpu = data.Cpu
	dbObj.Spec.Memory = data.Memory
	dbObj.Spec.Pod = data.Pod

	dbObj.GenerateVersion()

	result, update, err := n.IService.Apply(common.DefaultNamespace, common.Resource, dbObj.UUID, dbObj, false)
	if err != nil {
		return nil, false, err
	}

	history.Spec.History["cpu_now"] = dbObj.Spec.Cpu
	history.Spec.History["memory_now"] = dbObj.Spec.Memory


	if _, err = n.IService.Create(common.DefaultNamespace, common.History, history); err != nil {
		return nil, false, errors.New("the history create failed")
	}

	return result, update, nil
}