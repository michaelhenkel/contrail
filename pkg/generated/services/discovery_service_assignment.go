package services

import (
	"context"
	"database/sql"
	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/db"
	"github.com/Juniper/contrail/pkg/generated/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//RESTCreateDiscoveryServiceAssignment handle a Create REST service.
func (service *ContrailService) RESTCreateDiscoveryServiceAssignment(c echo.Context) error {
	requestData := &models.CreateDiscoveryServiceAssignmentRequest{}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "discovery_service_assignment",
		}).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	ctx := c.Request().Context()
	response, err := service.CreateDiscoveryServiceAssignment(ctx, requestData)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusCreated, response)
}

//CreateDiscoveryServiceAssignment handle a Create API
func (service *ContrailService) CreateDiscoveryServiceAssignment(
	ctx context.Context,
	request *models.CreateDiscoveryServiceAssignmentRequest) (*models.CreateDiscoveryServiceAssignmentResponse, error) {
	model := request.DiscoveryServiceAssignment
	if model.UUID == "" {
		model.UUID = uuid.NewV4().String()
	}
	auth := common.GetAuthCTX(ctx)
	if auth == nil {
		return nil, common.ErrorUnauthenticated
	}

	if model.FQName == nil {
		if model.DisplayName != "" {
			model.FQName = []string{auth.DomainID(), auth.ProjectID(), model.DisplayName}
		} else {
			model.FQName = []string{auth.DomainID(), auth.ProjectID(), model.UUID}
		}
	}
	model.Perms2 = &models.PermType2{}
	model.Perms2.Owner = auth.ProjectID()
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.CreateDiscoveryServiceAssignment(ctx, tx, request)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "discovery_service_assignment",
		}).Debug("db create failed on create")
		return nil, common.ErrorInternal
	}
	return &models.CreateDiscoveryServiceAssignmentResponse{
		DiscoveryServiceAssignment: request.DiscoveryServiceAssignment,
	}, nil
}

//RESTUpdateDiscoveryServiceAssignment handles a REST Update request.
func (service *ContrailService) RESTUpdateDiscoveryServiceAssignment(c echo.Context) error {
	//id := c.Param("id")
	request := &models.UpdateDiscoveryServiceAssignmentRequest{}
	if err := c.Bind(request); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "discovery_service_assignment",
		}).Debug("bind failed on update")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	ctx := c.Request().Context()
	response, err := service.UpdateDiscoveryServiceAssignment(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//UpdateDiscoveryServiceAssignment handles a Update request.
func (service *ContrailService) UpdateDiscoveryServiceAssignment(
	ctx context.Context,
	request *models.UpdateDiscoveryServiceAssignmentRequest) (*models.UpdateDiscoveryServiceAssignmentResponse, error) {
	model := request.DiscoveryServiceAssignment
	if model == nil {
		return nil, common.ErrorBadRequest("Update body is empty")
	}
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.UpdateDiscoveryServiceAssignment(ctx, tx, request)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "discovery_service_assignment",
		}).Debug("db update failed")
		return nil, common.ErrorInternal
	}
	return &models.UpdateDiscoveryServiceAssignmentResponse{
		DiscoveryServiceAssignment: model,
	}, nil
}

//RESTDeleteDiscoveryServiceAssignment delete a resource using REST service.
func (service *ContrailService) RESTDeleteDiscoveryServiceAssignment(c echo.Context) error {
	id := c.Param("id")
	request := &models.DeleteDiscoveryServiceAssignmentRequest{
		ID: id,
	}
	ctx := c.Request().Context()
	_, err := service.DeleteDiscoveryServiceAssignment(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//DeleteDiscoveryServiceAssignment delete a resource.
func (service *ContrailService) DeleteDiscoveryServiceAssignment(ctx context.Context, request *models.DeleteDiscoveryServiceAssignmentRequest) (*models.DeleteDiscoveryServiceAssignmentResponse, error) {
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.DeleteDiscoveryServiceAssignment(ctx, tx, request)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return nil, common.ErrorInternal
	}
	return &models.DeleteDiscoveryServiceAssignmentResponse{
		ID: request.ID,
	}, nil
}

//RESTGetDiscoveryServiceAssignment a REST Get request.
func (service *ContrailService) RESTGetDiscoveryServiceAssignment(c echo.Context) error {
	id := c.Param("id")
	request := &models.GetDiscoveryServiceAssignmentRequest{
		ID: id,
	}
	ctx := c.Request().Context()
	response, err := service.GetDiscoveryServiceAssignment(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//GetDiscoveryServiceAssignment a Get request.
func (service *ContrailService) GetDiscoveryServiceAssignment(ctx context.Context, request *models.GetDiscoveryServiceAssignmentRequest) (response *models.GetDiscoveryServiceAssignmentResponse, err error) {
	spec := &models.ListSpec{
		Limit: 1,
		Filters: []*models.Filter{
			&models.Filter{
				Key:    "uuid",
				Values: []string{request.ID},
			},
		},
	}
	listRequest := &models.ListDiscoveryServiceAssignmentRequest{
		Spec: spec,
	}
	var result *models.ListDiscoveryServiceAssignmentResponse
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListDiscoveryServiceAssignment(ctx, tx, listRequest)
			return err
		}); err != nil {
		return nil, common.ErrorInternal
	}
	if len(result.DiscoveryServiceAssignments) == 0 {
		return nil, common.ErrorNotFound
	}
	response = &models.GetDiscoveryServiceAssignmentResponse{
		DiscoveryServiceAssignment: result.DiscoveryServiceAssignments[0],
	}
	return response, nil
}

//RESTListDiscoveryServiceAssignment handles a List REST service Request.
func (service *ContrailService) RESTListDiscoveryServiceAssignment(c echo.Context) error {
	var err error
	spec := common.GetListSpec(c)
	request := &models.ListDiscoveryServiceAssignmentRequest{
		Spec: spec,
	}
	ctx := c.Request().Context()
	response, err := service.ListDiscoveryServiceAssignment(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//ListDiscoveryServiceAssignment handles a List service Request.
func (service *ContrailService) ListDiscoveryServiceAssignment(
	ctx context.Context,
	request *models.ListDiscoveryServiceAssignmentRequest) (response *models.ListDiscoveryServiceAssignmentResponse, err error) {
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			response, err = db.ListDiscoveryServiceAssignment(ctx, tx, request)
			return err
		}); err != nil {
		return nil, common.ErrorInternal
	}
	return response, nil
}
