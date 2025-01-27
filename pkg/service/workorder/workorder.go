package workorder

import (
	"github.com/pkg/errors"
	apiResource "github.com/yametech/devops/pkg/api/resource/workorder"
	"github.com/yametech/devops/pkg/common"
	"github.com/yametech/devops/pkg/core"
	"github.com/yametech/devops/pkg/resource/workorder"
	"github.com/yametech/devops/pkg/service"
	"github.com/yametech/devops/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	service.IService
}

func NewWorkOrderService(i service.IService) *Service {
	return &Service{i}
}

func (s *Service) List(orderType int, orderStatus int, search string, page, pageSize int64) ([]interface{}, error) {
	offset := (page - 1) * pageSize

	filter := make(map[string]interface{})

	if orderStatus == -1 {
		filter["$or"] = []map[string]interface{}{
			{
				"spec.order_type": orderType,
				"spec.number":     bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
			{
				"spec.order_type": orderType,
				"spec.title":      bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
			//{
			//	"spec.order_type": orderType,
			//	"spec.creator":      bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			//},
		}
	} else {
		filter["$or"] = []map[string]interface{}{
			{
				"spec.order_type":   orderType,
				"spec.order_status": orderStatus,
				"spec.number":       bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
			{
				"spec.order_type":   orderType,
				"spec.order_status": orderStatus,
				"spec.title":        bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			},
			//{
			//	"spec.order_type": orderType,
			//	"spec.order_status": orderStatus,
			//	"spec.creator":      bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			//},
		}
	}

	sort := map[string]interface{}{
		"metadata.version": -1,
	}

	return s.IService.ListByFilter(common.DefaultNamespace, common.WorkOrder, filter, sort, offset, pageSize)
}

func (s *Service) Create(request *apiResource.Request) (core.IObject, error) {
	req := &workorder.WorkOrder{
		Spec: workorder.Spec{
			OrderType: request.OrderType,
			Title:     request.Title,
			Relation:  request.Relation,
			Attribute: request.Attribute,
			Apply:     request.Apply,
			Check:     request.Check,
			Result:    request.Result,
			OrderStatus: 1,
		},
	}

	req.GenerateNumber()
	req.GenerateVersion()
	return s.IService.Create(common.DefaultNamespace, common.WorkOrder, req)
}

func (s *Service) Get(uuid string) (core.IObject, error) {
	order := &workorder.WorkOrder{}
	if err := s.IService.GetByUUID(common.DefaultNamespace, common.WorkOrder, uuid, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) Update(uuid string, request *apiResource.Request) (core.IObject, bool, error) {
	dbObj := &workorder.WorkOrder{}
	if err := s.GetByUUID(common.DefaultNamespace, common.WorkOrder, uuid, dbObj); err != nil {
		return nil, false, errors.New("The workorder is not exist")
	}

	dbObj.Spec.OrderStatus = request.OrderStatus
	dbObj.Spec.Title = request.Title
	dbObj.Spec.Attribute = request.Attribute
	dbObj.Spec.Apply = request.Apply
	dbObj.Spec.Check = request.Check
	dbObj.Spec.Result = request.Result

	dbObj.GenerateVersion()
	return s.IService.Apply(common.DefaultNamespace, common.WorkOrder, dbObj.UUID, dbObj, false)
}

func (s *Service) Delete(uuid string) (bool, error) {
	if err := s.IService.Delete(common.DefaultNamespace, common.WorkOrder, uuid); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) GetWorkOrderStatus(relation string, orderType int) (workorder.OrderStatus, error) {

	filter := map[string]interface{}{
		"spec.relation": relation,
		"spec.order_type": orderType,
	}
	sort := map[string]interface{}{
		"metadata.created_time": -1,
	}

	data, _ := s.IService.ListByFilter(common.DefaultNamespace, common.WorkOrder, filter, sort, 0 , 1)
	if len(data) == 0{
		return workorder.None, errors.New("The workorder is not exist")
	}

	order := make([]*workorder.WorkOrder, 0)
	if err := utils.UnstructuredObjectToInstanceObj(data, &order); err != nil {
		return workorder.None, err
	}

	return order[0].Spec.OrderStatus, nil
}
