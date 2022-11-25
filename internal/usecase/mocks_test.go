// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package usecase_test is a generated GoMock package.
package usecase_test

import (
	context "context"
	reflect "reflect"

	entity "github.com/PaulYakow/metrics-track/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockIClient is a mock of IClient interface.
type MockIClient struct {
	ctrl     *gomock.Controller
	recorder *MockIClientMockRecorder
}

// MockIClientMockRecorder is the mock recorder for MockIClient.
type MockIClientMockRecorder struct {
	mock *MockIClient
}

// NewMockIClient creates a new mock instance.
func NewMockIClient(ctrl *gomock.Controller) *MockIClient {
	mock := &MockIClient{ctrl: ctrl}
	mock.recorder = &MockIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClient) EXPECT() *MockIClientMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockIClient) GetAll() []entity.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]entity.Metric)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIClientMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIClient)(nil).GetAll))
}

// Poll mocks base method.
func (m *MockIClient) Poll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Poll")
}

// Poll indicates an expected call of Poll.
func (mr *MockIClientMockRecorder) Poll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poll", reflect.TypeOf((*MockIClient)(nil).Poll))
}

// MockIClientMemory is a mock of IClientMemory interface.
type MockIClientMemory struct {
	ctrl     *gomock.Controller
	recorder *MockIClientMemoryMockRecorder
}

// MockIClientMemoryMockRecorder is the mock recorder for MockIClientMemory.
type MockIClientMemoryMockRecorder struct {
	mock *MockIClientMemory
}

// NewMockIClientMemory creates a new mock instance.
func NewMockIClientMemory(ctrl *gomock.Controller) *MockIClientMemory {
	mock := &MockIClientMemory{ctrl: ctrl}
	mock.recorder = &MockIClientMemoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClientMemory) EXPECT() *MockIClientMemoryMockRecorder {
	return m.recorder
}

// ReadAll mocks base method.
func (m *MockIClientMemory) ReadAll() []entity.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAll")
	ret0, _ := ret[0].([]entity.Metric)
	return ret0
}

// ReadAll indicates an expected call of ReadAll.
func (mr *MockIClientMemoryMockRecorder) ReadAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAll", reflect.TypeOf((*MockIClientMemory)(nil).ReadAll))
}

// Store mocks base method.
func (m *MockIClientMemory) Store(arg0 map[string]*entity.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Store", arg0)
}

// Store indicates an expected call of Store.
func (mr *MockIClientMemoryMockRecorder) Store(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockIClientMemory)(nil).Store), arg0)
}

// MockIClientGather is a mock of IClientGather interface.
type MockIClientGather struct {
	ctrl     *gomock.Controller
	recorder *MockIClientGatherMockRecorder
}

// MockIClientGatherMockRecorder is the mock recorder for MockIClientGather.
type MockIClientGatherMockRecorder struct {
	mock *MockIClientGather
}

// NewMockIClientGather creates a new mock instance.
func NewMockIClientGather(ctrl *gomock.Controller) *MockIClientGather {
	mock := &MockIClientGather{ctrl: ctrl}
	mock.recorder = &MockIClientGatherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIClientGather) EXPECT() *MockIClientGatherMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockIClientGather) Update() map[string]*entity.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update")
	ret0, _ := ret[0].(map[string]*entity.Metric)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIClientGatherMockRecorder) Update() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIClientGather)(nil).Update))
}

// MockIServer is a mock of IServer interface.
type MockIServer struct {
	ctrl     *gomock.Controller
	recorder *MockIServerMockRecorder
}

// MockIServerMockRecorder is the mock recorder for MockIServer.
type MockIServerMockRecorder struct {
	mock *MockIServer
}

// NewMockIServer creates a new mock instance.
func NewMockIServer(ctrl *gomock.Controller) *MockIServer {
	mock := &MockIServer{ctrl: ctrl}
	mock.recorder = &MockIServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIServer) EXPECT() *MockIServerMockRecorder {
	return m.recorder
}

// CheckRepo mocks base method.
func (m *MockIServer) CheckRepo() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRepo")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckRepo indicates an expected call of CheckRepo.
func (mr *MockIServerMockRecorder) CheckRepo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRepo", reflect.TypeOf((*MockIServer)(nil).CheckRepo))
}

// Get mocks base method.
func (m *MockIServer) Get(ctx context.Context, metric entity.Metric) (*entity.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, metric)
	ret0, _ := ret[0].(*entity.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIServerMockRecorder) Get(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIServer)(nil).Get), ctx, metric)
}

// GetAll mocks base method.
func (m *MockIServer) GetAll(ctx context.Context) ([]entity.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]entity.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockIServerMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockIServer)(nil).GetAll), ctx)
}

// Save mocks base method.
func (m *MockIServer) Save(metric *entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIServerMockRecorder) Save(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIServer)(nil).Save), metric)
}

// SaveBatch mocks base method.
func (m *MockIServer) SaveBatch(metrics []entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatch", metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBatch indicates an expected call of SaveBatch.
func (mr *MockIServerMockRecorder) SaveBatch(metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatch", reflect.TypeOf((*MockIServer)(nil).SaveBatch), metrics)
}

// MockIServerRepo is a mock of IServerRepo interface.
type MockIServerRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIServerRepoMockRecorder
}

// MockIServerRepoMockRecorder is the mock recorder for MockIServerRepo.
type MockIServerRepoMockRecorder struct {
	mock *MockIServerRepo
}

// NewMockIServerRepo creates a new mock instance.
func NewMockIServerRepo(ctrl *gomock.Controller) *MockIServerRepo {
	mock := &MockIServerRepo{ctrl: ctrl}
	mock.recorder = &MockIServerRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIServerRepo) EXPECT() *MockIServerRepoMockRecorder {
	return m.recorder
}

// CheckConnection mocks base method.
func (m *MockIServerRepo) CheckConnection() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckConnection")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckConnection indicates an expected call of CheckConnection.
func (mr *MockIServerRepoMockRecorder) CheckConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckConnection", reflect.TypeOf((*MockIServerRepo)(nil).CheckConnection))
}

// Read mocks base method.
func (m *MockIServerRepo) Read(ctx context.Context, metric entity.Metric) (*entity.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, metric)
	ret0, _ := ret[0].(*entity.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockIServerRepoMockRecorder) Read(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockIServerRepo)(nil).Read), ctx, metric)
}

// ReadAll mocks base method.
func (m *MockIServerRepo) ReadAll(ctx context.Context) ([]entity.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAll", ctx)
	ret0, _ := ret[0].([]entity.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAll indicates an expected call of ReadAll.
func (mr *MockIServerRepoMockRecorder) ReadAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAll", reflect.TypeOf((*MockIServerRepo)(nil).ReadAll), ctx)
}

// Store mocks base method.
func (m *MockIServerRepo) Store(metric *entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockIServerRepoMockRecorder) Store(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockIServerRepo)(nil).Store), metric)
}

// StoreBatch mocks base method.
func (m *MockIServerRepo) StoreBatch(metrics []entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreBatch", metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreBatch indicates an expected call of StoreBatch.
func (mr *MockIServerRepoMockRecorder) StoreBatch(metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreBatch", reflect.TypeOf((*MockIServerRepo)(nil).StoreBatch), metrics)
}

// MockIHasher is a mock of IHasher interface.
type MockIHasher struct {
	ctrl     *gomock.Controller
	recorder *MockIHasherMockRecorder
}

// MockIHasherMockRecorder is the mock recorder for MockIHasher.
type MockIHasherMockRecorder struct {
	mock *MockIHasher
}

// NewMockIHasher creates a new mock instance.
func NewMockIHasher(ctrl *gomock.Controller) *MockIHasher {
	mock := &MockIHasher{ctrl: ctrl}
	mock.recorder = &MockIHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHasher) EXPECT() *MockIHasherMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockIHasher) Check(arg0 *entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockIHasherMockRecorder) Check(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockIHasher)(nil).Check), arg0)
}

// ProcessBatch mocks base method.
func (m *MockIHasher) ProcessBatch(arg0 []entity.Metric) []entity.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessBatch", arg0)
	ret0, _ := ret[0].([]entity.Metric)
	return ret0
}

// ProcessBatch indicates an expected call of ProcessBatch.
func (mr *MockIHasherMockRecorder) ProcessBatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessBatch", reflect.TypeOf((*MockIHasher)(nil).ProcessBatch), arg0)
}

// ProcessPointer mocks base method.
func (m *MockIHasher) ProcessPointer(arg0 *entity.Metric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProcessPointer", arg0)
}

// ProcessPointer indicates an expected call of ProcessPointer.
func (mr *MockIHasherMockRecorder) ProcessPointer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessPointer", reflect.TypeOf((*MockIHasher)(nil).ProcessPointer), arg0)
}

// ProcessSingle mocks base method.
func (m *MockIHasher) ProcessSingle(arg0 entity.Metric) entity.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessSingle", arg0)
	ret0, _ := ret[0].(entity.Metric)
	return ret0
}

// ProcessSingle indicates an expected call of ProcessSingle.
func (mr *MockIHasherMockRecorder) ProcessSingle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessSingle", reflect.TypeOf((*MockIHasher)(nil).ProcessSingle), arg0)
}
