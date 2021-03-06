// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	authProto "stlab.itechart-group.com/go/food_delivery/authentication_service/GRPC"
	model "stlab.itechart-group.com/go/food_delivery/authentication_service/model"
)

// MockAppUser is a mock of AppUser interface.
type MockAppUser struct {
	ctrl     *gomock.Controller
	recorder *MockAppUserMockRecorder
}

// MockAppUserMockRecorder is the mock recorder for MockAppUser.
type MockAppUserMockRecorder struct {
	mock *MockAppUser
}

// NewMockAppUser creates a new mock instance.
func NewMockAppUser(ctrl *gomock.Controller) *MockAppUser {
	mock := &MockAppUser{ctrl: ctrl}
	mock.recorder = &MockAppUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppUser) EXPECT() *MockAppUserMockRecorder {
	return m.recorder
}

// AuthUser mocks base method.
func (m *MockAppUser) AuthUser(email, password string) (*authProto.GeneratedTokens, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUser", email, password)
	ret0, _ := ret[0].(*authProto.GeneratedTokens)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AuthUser indicates an expected call of AuthUser.
func (mr *MockAppUserMockRecorder) AuthUser(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUser", reflect.TypeOf((*MockAppUser)(nil).AuthUser), email, password)
}

// CheckInputRole mocks base method.
func (m *MockAppUser) CheckInputRole(role string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInputRole", role)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckInputRole indicates an expected call of CheckInputRole.
func (mr *MockAppUserMockRecorder) CheckInputRole(role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInputRole", reflect.TypeOf((*MockAppUser)(nil).CheckInputRole), role)
}

// CheckPasswordHash mocks base method.
func (m *MockAppUser) CheckPasswordHash(password, hash string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPasswordHash", password, hash)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckPasswordHash indicates an expected call of CheckPasswordHash.
func (mr *MockAppUserMockRecorder) CheckPasswordHash(password, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPasswordHash", reflect.TypeOf((*MockAppUser)(nil).CheckPasswordHash), password, hash)
}

// CheckRights mocks base method.
func (m *MockAppUser) CheckRights(neededPerms []string, givenPerms string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRights", neededPerms, givenPerms)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckRights indicates an expected call of CheckRights.
func (mr *MockAppUserMockRecorder) CheckRights(neededPerms, givenPerms interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRights", reflect.TypeOf((*MockAppUser)(nil).CheckRights), neededPerms, givenPerms)
}

// CheckRole mocks base method.
func (m *MockAppUser) CheckRole(neededRoles []string, givenRole string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRole", neededRoles, givenRole)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckRole indicates an expected call of CheckRole.
func (mr *MockAppUserMockRecorder) CheckRole(neededRoles, givenRole interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRole", reflect.TypeOf((*MockAppUser)(nil).CheckRole), neededRoles, givenRole)
}

// CreateCustomer mocks base method.
func (m *MockAppUser) CreateCustomer(user *model.CreateCustomer) (*authProto.GeneratedTokens, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", user)
	ret0, _ := ret[0].(*authProto.GeneratedTokens)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateCustomer indicates an expected call of CreateCustomer.
func (mr *MockAppUserMockRecorder) CreateCustomer(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockAppUser)(nil).CreateCustomer), user)
}

// CreateStaff mocks base method.
func (m *MockAppUser) CreateStaff(user *model.CreateStaff) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStaff", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStaff indicates an expected call of CreateStaff.
func (mr *MockAppUserMockRecorder) CreateStaff(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStaff", reflect.TypeOf((*MockAppUser)(nil).CreateStaff), user)
}

// DeleteUserByID mocks base method.
func (m *MockAppUser) DeleteUserByID(id int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserByID", id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUserByID indicates an expected call of DeleteUserByID.
func (mr *MockAppUserMockRecorder) DeleteUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserByID", reflect.TypeOf((*MockAppUser)(nil).DeleteUserByID), id)
}

// GetUser mocks base method.
func (m *MockAppUser) GetUser(id int) (*model.ResponseUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*model.ResponseUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAppUserMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAppUser)(nil).GetUser), id)
}

// GetUsers mocks base method.
func (m *MockAppUser) GetUsers(page, limit int, filters *model.RequestFilters) ([]model.ResponseUser, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", page, limit, filters)
	ret0, _ := ret[0].([]model.ResponseUser)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockAppUserMockRecorder) GetUsers(page, limit, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockAppUser)(nil).GetUsers), page, limit, filters)
}

// HashPassword mocks base method.
func (m *MockAppUser) HashPassword(password string, rounds int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPassword", password, rounds)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HashPassword indicates an expected call of HashPassword.
func (mr *MockAppUserMockRecorder) HashPassword(password, rounds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPassword", reflect.TypeOf((*MockAppUser)(nil).HashPassword), password, rounds)
}

// ParseToken mocks base method.
func (m *MockAppUser) ParseToken(token string) (*authProto.UserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(*authProto.UserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAppUserMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAppUser)(nil).ParseToken), token)
}

// RestorePassword mocks base method.
func (m *MockAppUser) RestorePassword(restore *model.RestorePassword) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestorePassword", restore)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestorePassword indicates an expected call of RestorePassword.
func (mr *MockAppUserMockRecorder) RestorePassword(restore interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestorePassword", reflect.TypeOf((*MockAppUser)(nil).RestorePassword), restore)
}

// UpdateUser mocks base method.
func (m *MockAppUser) UpdateUser(user *model.UpdateUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAppUserMockRecorder) UpdateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAppUser)(nil).UpdateUser), user)
}
