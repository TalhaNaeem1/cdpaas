// Code generated by MockGen. DO NOT EDIT.
// Source: pipelineService/services/db (interfaces: Store)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	models "pipelineService/models/v1"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddPipeline mocks base method.
func (m *MockStore) AddPipeline(arg0 uuid.UUID, arg1 []models.ProductsPipelines) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPipeline", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPipeline indicates an expected call of AddPipeline.
func (mr *MockStoreMockRecorder) AddPipeline(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPipeline", reflect.TypeOf((*MockStore)(nil).AddPipeline), arg0, arg1)
}

// CreateConnectionAndSourceAgainstAPipeline mocks base method.
func (m *MockStore) CreateConnectionAndSourceAgainstAPipeline(arg0 models.Source, arg1 models.Connection) (models.Source, models.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateConnectionAndSourceAgainstAPipeline", arg0, arg1)
	ret0, _ := ret[0].(models.Source)
	ret1, _ := ret[1].(models.Connection)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateConnectionAndSourceAgainstAPipeline indicates an expected call of CreateConnectionAndSourceAgainstAPipeline.
func (mr *MockStoreMockRecorder) CreateConnectionAndSourceAgainstAPipeline(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateConnectionAndSourceAgainstAPipeline", reflect.TypeOf((*MockStore)(nil).CreateConnectionAndSourceAgainstAPipeline), arg0, arg1)
}

// CreateDataProduct mocks base method.
func (m *MockStore) CreateDataProduct(arg0 models.DataProduct) (models.DataProduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDataProduct", arg0)
	ret0, _ := ret[0].(models.DataProduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDataProduct indicates an expected call of CreateDataProduct.
func (mr *MockStoreMockRecorder) CreateDataProduct(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDataProduct", reflect.TypeOf((*MockStore)(nil).CreateDataProduct), arg0)
}

// CreateDestination mocks base method.
func (m *MockStore) CreateDestination(arg0 models.Destination) (models.Destination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDestination", arg0)
	ret0, _ := ret[0].(models.Destination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDestination indicates an expected call of CreateDestination.
func (mr *MockStoreMockRecorder) CreateDestination(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDestination", reflect.TypeOf((*MockStore)(nil).CreateDestination), arg0)
}

// CreatePipeline mocks base method.
func (m *MockStore) CreatePipeline(arg0 models.Pipeline) (models.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePipeline", arg0)
	ret0, _ := ret[0].(models.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePipeline indicates an expected call of CreatePipeline.
func (mr *MockStoreMockRecorder) CreatePipeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePipeline", reflect.TypeOf((*MockStore)(nil).CreatePipeline), arg0)
}

// CreatePipelineAssets mocks base method.
func (m *MockStore) CreatePipelineAssets(arg0 []models.PipelineAssets) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePipelineAssets", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePipelineAssets indicates an expected call of CreatePipelineAssets.
func (mr *MockStoreMockRecorder) CreatePipelineAssets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePipelineAssets", reflect.TypeOf((*MockStore)(nil).CreatePipelineAssets), arg0)
}

// CreatePipelineSchema mocks base method.
func (m *MockStore) CreatePipelineSchema(arg0 models.PipelineSchemas) (models.PipelineSchemas, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePipelineSchema", arg0)
	ret0, _ := ret[0].(models.PipelineSchemas)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePipelineSchema indicates an expected call of CreatePipelineSchema.
func (mr *MockStoreMockRecorder) CreatePipelineSchema(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePipelineSchema", reflect.TypeOf((*MockStore)(nil).CreatePipelineSchema), arg0)
}

// CreateTransformationPipeline mocks base method.
func (m *MockStore) CreateTransformationPipeline(arg0 models.TransformationPipelines) (models.TransformationPipelines, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransformationPipeline", arg0)
	ret0, _ := ret[0].(models.TransformationPipelines)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransformationPipeline indicates an expected call of CreateTransformationPipeline.
func (mr *MockStoreMockRecorder) CreateTransformationPipeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransformationPipeline", reflect.TypeOf((*MockStore)(nil).CreateTransformationPipeline), arg0)
}

// DeletePipeline mocks base method.
func (m *MockStore) DeletePipeline(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePipeline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePipeline indicates an expected call of DeletePipeline.
func (mr *MockStoreMockRecorder) DeletePipeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePipeline", reflect.TypeOf((*MockStore)(nil).DeletePipeline), arg0)
}

// DeletePipelineSchema mocks base method.
func (m *MockStore) DeletePipelineSchema(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePipelineSchema", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePipelineSchema indicates an expected call of DeletePipelineSchema.
func (mr *MockStoreMockRecorder) DeletePipelineSchema(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePipelineSchema", reflect.TypeOf((*MockStore)(nil).DeletePipelineSchema), arg0)
}

// EnablePipelineAssets mocks base method.
func (m *MockStore) EnablePipelineAssets(arg0 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnablePipelineAssets", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnablePipelineAssets indicates an expected call of EnablePipelineAssets.
func (mr *MockStoreMockRecorder) EnablePipelineAssets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnablePipelineAssets", reflect.TypeOf((*MockStore)(nil).EnablePipelineAssets), arg0)
}

// GetAllConnections mocks base method.
func (m *MockStore) GetAllConnections() ([]models.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllConnections")
	ret0, _ := ret[0].([]models.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllConnections indicates an expected call of GetAllConnections.
func (mr *MockStoreMockRecorder) GetAllConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllConnections", reflect.TypeOf((*MockStore)(nil).GetAllConnections))
}

// GetAllDataProducts mocks base method.
func (m *MockStore) GetAllDataProducts(arg0 int) ([]models.GetAllDataProductsView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDataProducts", arg0)
	ret0, _ := ret[0].([]models.GetAllDataProductsView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDataProducts indicates an expected call of GetAllDataProducts.
func (mr *MockStoreMockRecorder) GetAllDataProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDataProducts", reflect.TypeOf((*MockStore)(nil).GetAllDataProducts), arg0)
}

// GetAllPipelines mocks base method.
func (m *MockStore) GetAllPipelines(arg0 int) ([]models.PipelinesMetaData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPipelines", arg0)
	ret0, _ := ret[0].([]models.PipelinesMetaData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPipelines indicates an expected call of GetAllPipelines.
func (mr *MockStoreMockRecorder) GetAllPipelines(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPipelines", reflect.TypeOf((*MockStore)(nil).GetAllPipelines), arg0)
}

// GetAssetDetails mocks base method.
func (m *MockStore) GetAssetDetails(arg0 uuid.UUID) (models.AssetDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssetDetails", arg0)
	ret0, _ := ret[0].(models.AssetDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssetDetails indicates an expected call of GetAssetDetails.
func (mr *MockStoreMockRecorder) GetAssetDetails(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssetDetails", reflect.TypeOf((*MockStore)(nil).GetAssetDetails), arg0)
}

// GetConfiguredDestination mocks base method.
func (m *MockStore) GetConfiguredDestination(arg0 int) ([]models.ConfiguredDestination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfiguredDestination", arg0)
	ret0, _ := ret[0].([]models.ConfiguredDestination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfiguredDestination indicates an expected call of GetConfiguredDestination.
func (mr *MockStoreMockRecorder) GetConfiguredDestination(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfiguredDestination", reflect.TypeOf((*MockStore)(nil).GetConfiguredDestination), arg0)
}

// GetConnection mocks base method.
func (m *MockStore) GetConnection(arg0 string) (models.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnection", arg0)
	ret0, _ := ret[0].(models.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnection indicates an expected call of GetConnection.
func (mr *MockStoreMockRecorder) GetConnection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockStore)(nil).GetConnection), arg0)
}

// GetDataProduct mocks base method.
func (m *MockStore) GetDataProduct(arg0 uuid.UUID) (models.DataProductView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataProduct", arg0)
	ret0, _ := ret[0].(models.DataProductView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataProduct indicates an expected call of GetDataProduct.
func (mr *MockStoreMockRecorder) GetDataProduct(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataProduct", reflect.TypeOf((*MockStore)(nil).GetDataProduct), arg0)
}

// GetDataProductInfo mocks base method.
func (m *MockStore) GetDataProductInfo(arg0 uuid.UUID) (models.DataProduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataProductInfo", arg0)
	ret0, _ := ret[0].(models.DataProduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataProductInfo indicates an expected call of GetDataProductInfo.
func (mr *MockStoreMockRecorder) GetDataProductInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataProductInfo", reflect.TypeOf((*MockStore)(nil).GetDataProductInfo), arg0)
}

// GetDestination mocks base method.
func (m *MockStore) GetDestination(arg0 uuid.UUID) (models.Destination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDestination", arg0)
	ret0, _ := ret[0].(models.Destination)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDestination indicates an expected call of GetDestination.
func (mr *MockStoreMockRecorder) GetDestination(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDestination", reflect.TypeOf((*MockStore)(nil).GetDestination), arg0)
}

// GetDestinationSummary mocks base method.
func (m *MockStore) GetDestinationSummary(arg0 uuid.UUID) (models.DestinationSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDestinationSummary", arg0)
	ret0, _ := ret[0].(models.DestinationSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDestinationSummary indicates an expected call of GetDestinationSummary.
func (mr *MockStoreMockRecorder) GetDestinationSummary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDestinationSummary", reflect.TypeOf((*MockStore)(nil).GetDestinationSummary), arg0)
}

// GetPipeline mocks base method.
func (m *MockStore) GetPipeline(arg0 uuid.UUID) (models.PipelineView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipeline", arg0)
	ret0, _ := ret[0].(models.PipelineView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipeline indicates an expected call of GetPipeline.
func (mr *MockStoreMockRecorder) GetPipeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipeline", reflect.TypeOf((*MockStore)(nil).GetPipeline), arg0)
}

// GetPipelineAssets mocks base method.
func (m *MockStore) GetPipelineAssets(arg0 uuid.UUID) ([]models.PipelineAssets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineAssets", arg0)
	ret0, _ := ret[0].([]models.PipelineAssets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineAssets indicates an expected call of GetPipelineAssets.
func (mr *MockStoreMockRecorder) GetPipelineAssets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineAssets", reflect.TypeOf((*MockStore)(nil).GetPipelineAssets), arg0)
}

// GetPipelineConnection mocks base method.
func (m *MockStore) GetPipelineConnection(arg0 string) (models.PipelineConnection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineConnection", arg0)
	ret0, _ := ret[0].(models.PipelineConnection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineConnection indicates an expected call of GetPipelineConnection.
func (mr *MockStoreMockRecorder) GetPipelineConnection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineConnection", reflect.TypeOf((*MockStore)(nil).GetPipelineConnection), arg0)
}

// GetPipelineSchema mocks base method.
func (m *MockStore) GetPipelineSchema(arg0 uuid.UUID) (models.PipelineSchemas, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineSchema", arg0)
	ret0, _ := ret[0].(models.PipelineSchemas)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineSchema indicates an expected call of GetPipelineSchema.
func (mr *MockStoreMockRecorder) GetPipelineSchema(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineSchema", reflect.TypeOf((*MockStore)(nil).GetPipelineSchema), arg0)
}

// GetPipelineSourceAndConnectionID mocks base method.
func (m *MockStore) GetPipelineSourceAndConnectionID(arg0 uuid.UUID) (models.PipelineSourceAndConnectionID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineSourceAndConnectionID", arg0)
	ret0, _ := ret[0].(models.PipelineSourceAndConnectionID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineSourceAndConnectionID indicates an expected call of GetPipelineSourceAndConnectionID.
func (mr *MockStoreMockRecorder) GetPipelineSourceAndConnectionID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineSourceAndConnectionID", reflect.TypeOf((*MockStore)(nil).GetPipelineSourceAndConnectionID), arg0)
}

// GetProductConnection mocks base method.
func (m *MockStore) GetProductConnection(arg0 uuid.UUID) (models.TransformationPipelines, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductConnection", arg0)
	ret0, _ := ret[0].(models.TransformationPipelines)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductConnection indicates an expected call of GetProductConnection.
func (mr *MockStoreMockRecorder) GetProductConnection(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductConnection", reflect.TypeOf((*MockStore)(nil).GetProductConnection), arg0)
}

// GetProductDetails mocks base method.
func (m *MockStore) GetProductDetails() ([]models.ProductDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductDetails")
	ret0, _ := ret[0].([]models.ProductDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductDetails indicates an expected call of GetProductDetails.
func (mr *MockStoreMockRecorder) GetProductDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductDetails", reflect.TypeOf((*MockStore)(nil).GetProductDetails))
}

// GetSource mocks base method.
func (m *MockStore) GetSource(arg0 string) (models.Source, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSource", arg0)
	ret0, _ := ret[0].(models.Source)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSource indicates an expected call of GetSource.
func (mr *MockStoreMockRecorder) GetSource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSource", reflect.TypeOf((*MockStore)(nil).GetSource), arg0)
}

// GetSourceAndConnectionDetails mocks base method.
func (m *MockStore) GetSourceAndConnectionDetails(arg0 string) (models.ConnectionSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceAndConnectionDetails", arg0)
	ret0, _ := ret[0].(models.ConnectionSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSourceAndConnectionDetails indicates an expected call of GetSourceAndConnectionDetails.
func (mr *MockStoreMockRecorder) GetSourceAndConnectionDetails(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceAndConnectionDetails", reflect.TypeOf((*MockStore)(nil).GetSourceAndConnectionDetails), arg0)
}

// GetSourceAndDestinationAirbyteInfo mocks base method.
func (m *MockStore) GetSourceAndDestinationAirbyteInfo(arg0, arg1 string) (models.AirbyteSourceAndDestinations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceAndDestinationAirbyteInfo", arg0, arg1)
	ret0, _ := ret[0].(models.AirbyteSourceAndDestinations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSourceAndDestinationAirbyteInfo indicates an expected call of GetSourceAndDestinationAirbyteInfo.
func (mr *MockStoreMockRecorder) GetSourceAndDestinationAirbyteInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceAndDestinationAirbyteInfo", reflect.TypeOf((*MockStore)(nil).GetSourceAndDestinationAirbyteInfo), arg0, arg1)
}

// GetSupportedDestinations mocks base method.
func (m *MockStore) GetSupportedDestinations() ([]models.SupportedDestinations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupportedDestinations")
	ret0, _ := ret[0].([]models.SupportedDestinations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupportedDestinations indicates an expected call of GetSupportedDestinations.
func (mr *MockStoreMockRecorder) GetSupportedDestinations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupportedDestinations", reflect.TypeOf((*MockStore)(nil).GetSupportedDestinations))
}

// GetSupportedSources mocks base method.
func (m *MockStore) GetSupportedSources() ([]models.SupportedSources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupportedSources")
	ret0, _ := ret[0].([]models.SupportedSources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupportedSources indicates an expected call of GetSupportedSources.
func (mr *MockStoreMockRecorder) GetSupportedSources() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupportedSources", reflect.TypeOf((*MockStore)(nil).GetSupportedSources))
}

// GetTransformationPipeline mocks base method.
func (m *MockStore) GetTransformationPipeline(arg0 uuid.UUID) (models.TransformationPipelines, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransformationPipeline", arg0)
	ret0, _ := ret[0].(models.TransformationPipelines)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransformationPipeline indicates an expected call of GetTransformationPipeline.
func (mr *MockStoreMockRecorder) GetTransformationPipeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransformationPipeline", reflect.TypeOf((*MockStore)(nil).GetTransformationPipeline), arg0)
}

// GetTransformedAssetDetails mocks base method.
func (m *MockStore) GetTransformedAssetDetails(arg0 uuid.UUID) (models.TransformedAssetDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransformedAssetDetails", arg0)
	ret0, _ := ret[0].(models.TransformedAssetDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransformedAssetDetails indicates an expected call of GetTransformedAssetDetails.
func (mr *MockStoreMockRecorder) GetTransformedAssetDetails(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransformedAssetDetails", reflect.TypeOf((*MockStore)(nil).GetTransformedAssetDetails), arg0)
}

// GetTransformedAssets mocks base method.
func (m *MockStore) GetTransformedAssets(arg0 uuid.UUID) ([]models.ProductAssets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransformedAssets", arg0)
	ret0, _ := ret[0].([]models.ProductAssets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransformedAssets indicates an expected call of GetTransformedAssets.
func (mr *MockStoreMockRecorder) GetTransformedAssets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransformedAssets", reflect.TypeOf((*MockStore)(nil).GetTransformedAssets), arg0)
}

// PreviewData mocks base method.
func (m *MockStore) PreviewData(arg0 *gorm.DB, arg1, arg2 string) ([]map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreviewData", arg0, arg1, arg2)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreviewData indicates an expected call of PreviewData.
func (mr *MockStoreMockRecorder) PreviewData(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreviewData", reflect.TypeOf((*MockStore)(nil).PreviewData), arg0, arg1, arg2)
}

// SyncTransformedAssets mocks base method.
func (m *MockStore) SyncTransformedAssets(arg0 []models.ProductAssetDetails) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncTransformedAssets", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncTransformedAssets indicates an expected call of SyncTransformedAssets.
func (mr *MockStoreMockRecorder) SyncTransformedAssets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncTransformedAssets", reflect.TypeOf((*MockStore)(nil).SyncTransformedAssets), arg0)
}

// UpdateConnectionInfo mocks base method.
func (m *MockStore) UpdateConnectionInfo(arg0 models.Connection, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConnectionInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConnectionInfo indicates an expected call of UpdateConnectionInfo.
func (mr *MockStoreMockRecorder) UpdateConnectionInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConnectionInfo", reflect.TypeOf((*MockStore)(nil).UpdateConnectionInfo), arg0, arg1)
}

// UpdateConnectionSchedule mocks base method.
func (m *MockStore) UpdateConnectionSchedule(arg0 models.Connection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConnectionSchedule", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConnectionSchedule indicates an expected call of UpdateConnectionSchedule.
func (mr *MockStoreMockRecorder) UpdateConnectionSchedule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConnectionSchedule", reflect.TypeOf((*MockStore)(nil).UpdateConnectionSchedule), arg0)
}

// UpdateConnections mocks base method.
func (m *MockStore) UpdateConnections(arg0 []models.Connection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConnections", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConnections indicates an expected call of UpdateConnections.
func (mr *MockStoreMockRecorder) UpdateConnections(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConnections", reflect.TypeOf((*MockStore)(nil).UpdateConnections), arg0)
}

// UpdateDataProduct mocks base method.
func (m *MockStore) UpdateDataProduct(arg0 models.DataProduct) (models.DataProduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDataProduct", arg0)
	ret0, _ := ret[0].(models.DataProduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateDataProduct indicates an expected call of UpdateDataProduct.
func (mr *MockStoreMockRecorder) UpdateDataProduct(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDataProduct", reflect.TypeOf((*MockStore)(nil).UpdateDataProduct), arg0)
}

// UpdatePipeline mocks base method.
func (m *MockStore) UpdatePipeline(arg0 models.UpdatePipeline) (models.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePipeline", arg0)
	ret0, _ := ret[0].(models.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePipeline indicates an expected call of UpdatePipeline.
func (mr *MockStoreMockRecorder) UpdatePipeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePipeline", reflect.TypeOf((*MockStore)(nil).UpdatePipeline), arg0)
}

// UpdatePipelineStatus mocks base method.
func (m *MockStore) UpdatePipelineStatus(arg0 uuid.UUID, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePipelineStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePipelineStatus indicates an expected call of UpdatePipelineStatus.
func (mr *MockStoreMockRecorder) UpdatePipelineStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePipelineStatus", reflect.TypeOf((*MockStore)(nil).UpdatePipelineStatus), arg0, arg1)
}
