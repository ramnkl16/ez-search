package services

import (
	"fmt"

	"github.com/ramnkl16/ez-search/abstractimpl"

	"github.com/ramnkl16/ez-search/models"
	"github.com/ramnkl16/ez-search/rest_errors"
	"github.com/ramnkl16/ez-search/utils/date_utils"
	"github.com/ramnkl16/ez-search/utils/uid_utils"
)

//Auto code generated with help of mysql table schema
// table : WidgetMeta

//WidgetMeta service as variable
var (
	WidgetMetaService widgetMetaServiceInterface = &widgetMetaService{}
)

type widgetMetaService struct{}

type widgetMetaServiceInterface interface {
	Create(models.WidgetMeta) rest_errors.RestErr
	Update(models.WidgetMeta) rest_errors.RestErr
	Get(string) (*models.WidgetMeta, rest_errors.RestErr)
	Delete(string) rest_errors.RestErr
	Search(string, string, string) (models.WidgetMetas, rest_errors.RestErr)
	//WebuiSearch(string, string) (models.WidgetMetas, rest_errors.RestErr)
}

func (srv *widgetMetaService) Create(wm models.WidgetMeta) rest_errors.RestErr {

	wm.IsActive = "t"
	wm.CreatedAt = date_utils.GetNowSearchFormat()
	wm.UpdatedAt = date_utils.GetNowSearchFormat()
	if len(wm.ID) == 0 {
		wm.ID = uid_utils.GetUid("rq", true)
	}
	if err := wm.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}

func (srv *widgetMetaService) Update(wm models.WidgetMeta) rest_errors.RestErr {
	wm.UpdatedAt = date_utils.GetNowSearchFormat()
	if err := wm.CreateOrUpdate(); err != nil {
		return err
	}
	return nil
}

func (srv *widgetMetaService) Get(id string) (*models.WidgetMeta, rest_errors.RestErr) {
	// dao := &models.WidgetMeta{ID: id}
	// if err := dao.Get(); err != nil {
	// 	return nil, err
	// }
	// return dao, nil
	return nil, nil

}

func (srv *widgetMetaService) Delete(id string) rest_errors.RestErr {
	dao := &models.WidgetMeta{ID: id}
	dao.UpdatedAt = date_utils.GetNowSearchFormat()

	if err := dao.Delete(id); err != nil {
		return err
	}
	return nil
}
func (srv *widgetMetaService) Search(start string, limit string, namespaceId string) (models.WidgetMetas, rest_errors.RestErr) {
	dao := &models.WidgetMeta{}
	if start == "" {
		start = "0"
	}
	if limit == "" {
		limit = "50"
	}

	q := fmt.Sprintf("select * from %s where +isActive:t,+division:%s limit %s,%s", abstractimpl.QueryMetaTable, namespaceId, start, limit)
	if namespaceId == "platform" || len(namespaceId) == 0 {
		q = fmt.Sprintf("select * from %s where +isActive:t limit %s,%s", abstractimpl.QueryMetaTable, start, limit)
	}
	list, err := dao.GetAll(q)
	return list, err
}
