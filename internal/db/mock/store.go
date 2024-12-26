// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/danilshap/domains-generator/internal/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	gomock "github.com/golang/mock/gomock"
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

// CreateDomain mocks base method.
func (m *MockStore) CreateDomain(arg0 context.Context, arg1 db.CreateDomainParams) (db.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDomain", arg0, arg1)
	ret0, _ := ret[0].(db.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDomain indicates an expected call of CreateDomain.
func (mr *MockStoreMockRecorder) CreateDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDomain", reflect.TypeOf((*MockStore)(nil).CreateDomain), arg0, arg1)
}

// CreateMailbox mocks base method.
func (m *MockStore) CreateMailbox(arg0 context.Context, arg1 db.CreateMailboxParams) (db.Mailbox, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMailbox", arg0, arg1)
	ret0, _ := ret[0].(db.Mailbox)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMailbox indicates an expected call of CreateMailbox.
func (mr *MockStoreMockRecorder) CreateMailbox(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMailbox", reflect.TypeOf((*MockStore)(nil).CreateMailbox), arg0, arg1)
}

// DeleteDomain mocks base method.
func (m *MockStore) DeleteDomain(arg0 context.Context, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDomain indicates an expected call of DeleteDomain.
func (mr *MockStoreMockRecorder) DeleteDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDomain", reflect.TypeOf((*MockStore)(nil).DeleteDomain), arg0, arg1)
}

// DeleteMailbox mocks base method.
func (m *MockStore) DeleteMailbox(arg0 context.Context, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMailbox", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMailbox indicates an expected call of DeleteMailbox.
func (mr *MockStoreMockRecorder) DeleteMailbox(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMailbox", reflect.TypeOf((*MockStore)(nil).DeleteMailbox), arg0, arg1)
}

// GetAllDomains mocks base method.
func (m *MockStore) GetAllDomains(arg0 context.Context, arg1 db.GetAllDomainsParams) ([]db.GetAllDomainsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDomains", arg0, arg1)
	ret0, _ := ret[0].([]db.GetAllDomainsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDomains indicates an expected call of GetAllDomains.
func (mr *MockStoreMockRecorder) GetAllDomains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDomains", reflect.TypeOf((*MockStore)(nil).GetAllDomains), arg0, arg1)
}

// GetAllMailboxes mocks base method.
func (m *MockStore) GetAllMailboxes(arg0 context.Context, arg1 db.GetAllMailboxesParams) ([]db.GetAllMailboxesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMailboxes", arg0, arg1)
	ret0, _ := ret[0].([]db.GetAllMailboxesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMailboxes indicates an expected call of GetAllMailboxes.
func (mr *MockStoreMockRecorder) GetAllMailboxes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMailboxes", reflect.TypeOf((*MockStore)(nil).GetAllMailboxes), arg0, arg1)
}

// GetDomainByID mocks base method.
func (m *MockStore) GetDomainByID(arg0 context.Context, arg1 int32) (db.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainByID", arg0, arg1)
	ret0, _ := ret[0].(db.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainByID indicates an expected call of GetDomainByID.
func (mr *MockStoreMockRecorder) GetDomainByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainByID", reflect.TypeOf((*MockStore)(nil).GetDomainByID), arg0, arg1)
}

// GetDomainByName mocks base method.
func (m *MockStore) GetDomainByName(arg0 context.Context, arg1 string) (db.GetDomainByNameRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainByName", arg0, arg1)
	ret0, _ := ret[0].(db.GetDomainByNameRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainByName indicates an expected call of GetDomainByName.
func (mr *MockStoreMockRecorder) GetDomainByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainByName", reflect.TypeOf((*MockStore)(nil).GetDomainByName), arg0, arg1)
}

// GetDomainsCount mocks base method.
func (m *MockStore) GetDomainsCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomainsCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomainsCount indicates an expected call of GetDomainsCount.
func (mr *MockStoreMockRecorder) GetDomainsCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomainsCount", reflect.TypeOf((*MockStore)(nil).GetDomainsCount), arg0)
}

// GetMailboxByID mocks base method.
func (m *MockStore) GetMailboxByID(arg0 context.Context, arg1 int32) (db.Mailbox, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxByID", arg0, arg1)
	ret0, _ := ret[0].(db.Mailbox)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxByID indicates an expected call of GetMailboxByID.
func (mr *MockStoreMockRecorder) GetMailboxByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxByID", reflect.TypeOf((*MockStore)(nil).GetMailboxByID), arg0, arg1)
}

// GetMailboxesByDomain mocks base method.
func (m *MockStore) GetMailboxesByDomain(arg0 context.Context, arg1 string) ([]db.GetMailboxesByDomainRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxesByDomain", arg0, arg1)
	ret0, _ := ret[0].([]db.GetMailboxesByDomainRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxesByDomain indicates an expected call of GetMailboxesByDomain.
func (mr *MockStoreMockRecorder) GetMailboxesByDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxesByDomain", reflect.TypeOf((*MockStore)(nil).GetMailboxesByDomain), arg0, arg1)
}

// GetMailboxesByDomainID mocks base method.
func (m *MockStore) GetMailboxesByDomainID(arg0 context.Context, arg1 db.GetMailboxesByDomainIDParams) ([]db.Mailbox, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxesByDomainID", arg0, arg1)
	ret0, _ := ret[0].([]db.Mailbox)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxesByDomainID indicates an expected call of GetMailboxesByDomainID.
func (mr *MockStoreMockRecorder) GetMailboxesByDomainID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxesByDomainID", reflect.TypeOf((*MockStore)(nil).GetMailboxesByDomainID), arg0, arg1)
}

// GetMailboxesCount mocks base method.
func (m *MockStore) GetMailboxesCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxesCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxesCount indicates an expected call of GetMailboxesCount.
func (mr *MockStoreMockRecorder) GetMailboxesCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxesCount", reflect.TypeOf((*MockStore)(nil).GetMailboxesCount), arg0)
}

// GetMailboxesCountByDomainID mocks base method.
func (m *MockStore) GetMailboxesCountByDomainID(arg0 context.Context, arg1 int32) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxesCountByDomainID", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxesCountByDomainID indicates an expected call of GetMailboxesCountByDomainID.
func (mr *MockStoreMockRecorder) GetMailboxesCountByDomainID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxesCountByDomainID", reflect.TypeOf((*MockStore)(nil).GetMailboxesCountByDomainID), arg0, arg1)
}

// GetMailboxesStats mocks base method.
func (m *MockStore) GetMailboxesStats(arg0 context.Context) (db.GetMailboxesStatsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxesStats", arg0)
	ret0, _ := ret[0].(db.GetMailboxesStatsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxesStats indicates an expected call of GetMailboxesStats.
func (mr *MockStoreMockRecorder) GetMailboxesStats(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxesStats", reflect.TypeOf((*MockStore)(nil).GetMailboxesStats), arg0)
}

// GetMailboxesWithFilters mocks base method.
func (m *MockStore) GetMailboxesWithFilters(arg0 context.Context, arg1 db.GetMailboxesWithFiltersParams) ([]db.GetMailboxesWithFiltersRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMailboxesWithFilters", arg0, arg1)
	ret0, _ := ret[0].([]db.GetMailboxesWithFiltersRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMailboxesWithFilters indicates an expected call of GetMailboxesWithFilters.
func (mr *MockStoreMockRecorder) GetMailboxesWithFilters(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMailboxesWithFilters", reflect.TypeOf((*MockStore)(nil).GetMailboxesWithFilters), arg0, arg1)
}

// SetDomainStatus mocks base method.
func (m *MockStore) SetDomainStatus(arg0 context.Context, arg1 db.SetDomainStatusParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDomainStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDomainStatus indicates an expected call of SetDomainStatus.
func (mr *MockStoreMockRecorder) SetDomainStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDomainStatus", reflect.TypeOf((*MockStore)(nil).SetDomainStatus), arg0, arg1)
}

// SetMailboxStatus mocks base method.
func (m *MockStore) SetMailboxStatus(arg0 context.Context, arg1 db.SetMailboxStatusParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMailboxStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMailboxStatus indicates an expected call of SetMailboxStatus.
func (mr *MockStoreMockRecorder) SetMailboxStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMailboxStatus", reflect.TypeOf((*MockStore)(nil).SetMailboxStatus), arg0, arg1)
}

// UpdateDomain mocks base method.
func (m *MockStore) UpdateDomain(arg0 context.Context, arg1 db.UpdateDomainParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDomain indicates an expected call of UpdateDomain.
func (mr *MockStoreMockRecorder) UpdateDomain(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDomain", reflect.TypeOf((*MockStore)(nil).UpdateDomain), arg0, arg1)
}

// UpdateDomainAndMailboxesStatus mocks base method.
func (m *MockStore) UpdateDomainAndMailboxesStatus(arg0 context.Context, arg1, arg2 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDomainAndMailboxesStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDomainAndMailboxesStatus indicates an expected call of UpdateDomainAndMailboxesStatus.
func (mr *MockStoreMockRecorder) UpdateDomainAndMailboxesStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDomainAndMailboxesStatus", reflect.TypeOf((*MockStore)(nil).UpdateDomainAndMailboxesStatus), arg0, arg1, arg2)
}

// UpdateMailbox mocks base method.
func (m *MockStore) UpdateMailbox(arg0 context.Context, arg1 db.UpdateMailboxParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMailbox", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMailbox indicates an expected call of UpdateMailbox.
func (mr *MockStoreMockRecorder) UpdateMailbox(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMailbox", reflect.TypeOf((*MockStore)(nil).UpdateMailbox), arg0, arg1)
}
