// Code generated by MockGen. DO NOT EDIT.
// Source: ./question.go
//
// Generated by this command:
//
//	mockgen -source=./question.go -destination=../../mocks/question.mock.go -package=quemocks -typed=true Service
//

// Package quemocks is a generated GoMock package.
package quemocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/ecodeclub/webook/internal/question/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockService) Delete(ctx context.Context, qid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, qid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(ctx, qid any) *MockServiceDeleteCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), ctx, qid)
	return &MockServiceDeleteCall{Call: call}
}

// MockServiceDeleteCall wrap *gomock.Call
type MockServiceDeleteCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceDeleteCall) Return(arg0 error) *MockServiceDeleteCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceDeleteCall) Do(f func(context.Context, int64) error) *MockServiceDeleteCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceDeleteCall) DoAndReturn(f func(context.Context, int64) error) *MockServiceDeleteCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Detail mocks base method.
func (m *MockService) Detail(ctx context.Context, qid int64) (domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detail", ctx, qid)
	ret0, _ := ret[0].(domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Detail indicates an expected call of Detail.
func (mr *MockServiceMockRecorder) Detail(ctx, qid any) *MockServiceDetailCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detail", reflect.TypeOf((*MockService)(nil).Detail), ctx, qid)
	return &MockServiceDetailCall{Call: call}
}

// MockServiceDetailCall wrap *gomock.Call
type MockServiceDetailCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceDetailCall) Return(arg0 domain.Question, arg1 error) *MockServiceDetailCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceDetailCall) Do(f func(context.Context, int64) (domain.Question, error)) *MockServiceDetailCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceDetailCall) DoAndReturn(f func(context.Context, int64) (domain.Question, error)) *MockServiceDetailCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetPubByIDs mocks base method.
func (m *MockService) GetPubByIDs(ctx context.Context, ids []int64) ([]domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPubByIDs", ctx, ids)
	ret0, _ := ret[0].([]domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPubByIDs indicates an expected call of GetPubByIDs.
func (mr *MockServiceMockRecorder) GetPubByIDs(ctx, ids any) *MockServiceGetPubByIDsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPubByIDs", reflect.TypeOf((*MockService)(nil).GetPubByIDs), ctx, ids)
	return &MockServiceGetPubByIDsCall{Call: call}
}

// MockServiceGetPubByIDsCall wrap *gomock.Call
type MockServiceGetPubByIDsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceGetPubByIDsCall) Return(arg0 []domain.Question, arg1 error) *MockServiceGetPubByIDsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceGetPubByIDsCall) Do(f func(context.Context, []int64) ([]domain.Question, error)) *MockServiceGetPubByIDsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceGetPubByIDsCall) DoAndReturn(f func(context.Context, []int64) ([]domain.Question, error)) *MockServiceGetPubByIDsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockService) List(ctx context.Context, offset, limit int) ([]domain.Question, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, offset, limit)
	ret0, _ := ret[0].([]domain.Question)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockServiceMockRecorder) List(ctx, offset, limit any) *MockServiceListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), ctx, offset, limit)
	return &MockServiceListCall{Call: call}
}

// MockServiceListCall wrap *gomock.Call
type MockServiceListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceListCall) Return(arg0 []domain.Question, arg1 int64, arg2 error) *MockServiceListCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceListCall) Do(f func(context.Context, int, int) ([]domain.Question, int64, error)) *MockServiceListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceListCall) DoAndReturn(f func(context.Context, int, int) ([]domain.Question, int64, error)) *MockServiceListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PubDetail mocks base method.
func (m *MockService) PubDetail(ctx context.Context, qid int64) (domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PubDetail", ctx, qid)
	ret0, _ := ret[0].(domain.Question)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PubDetail indicates an expected call of PubDetail.
func (mr *MockServiceMockRecorder) PubDetail(ctx, qid any) *MockServicePubDetailCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PubDetail", reflect.TypeOf((*MockService)(nil).PubDetail), ctx, qid)
	return &MockServicePubDetailCall{Call: call}
}

// MockServicePubDetailCall wrap *gomock.Call
type MockServicePubDetailCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServicePubDetailCall) Return(arg0 domain.Question, arg1 error) *MockServicePubDetailCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServicePubDetailCall) Do(f func(context.Context, int64) (domain.Question, error)) *MockServicePubDetailCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServicePubDetailCall) DoAndReturn(f func(context.Context, int64) (domain.Question, error)) *MockServicePubDetailCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PubList mocks base method.
func (m *MockService) PubList(ctx context.Context, offset, limit int) (int64, []domain.Question, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PubList", ctx, offset, limit)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].([]domain.Question)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PubList indicates an expected call of PubList.
func (mr *MockServiceMockRecorder) PubList(ctx, offset, limit any) *MockServicePubListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PubList", reflect.TypeOf((*MockService)(nil).PubList), ctx, offset, limit)
	return &MockServicePubListCall{Call: call}
}

// MockServicePubListCall wrap *gomock.Call
type MockServicePubListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServicePubListCall) Return(arg0 int64, arg1 []domain.Question, arg2 error) *MockServicePubListCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServicePubListCall) Do(f func(context.Context, int, int) (int64, []domain.Question, error)) *MockServicePubListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServicePubListCall) DoAndReturn(f func(context.Context, int, int) (int64, []domain.Question, error)) *MockServicePubListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Publish mocks base method.
func (m *MockService) Publish(ctx context.Context, que *domain.Question) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, que)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockServiceMockRecorder) Publish(ctx, que any) *MockServicePublishCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockService)(nil).Publish), ctx, que)
	return &MockServicePublishCall{Call: call}
}

// MockServicePublishCall wrap *gomock.Call
type MockServicePublishCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServicePublishCall) Return(arg0 int64, arg1 error) *MockServicePublishCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServicePublishCall) Do(f func(context.Context, *domain.Question) (int64, error)) *MockServicePublishCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServicePublishCall) DoAndReturn(f func(context.Context, *domain.Question) (int64, error)) *MockServicePublishCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m *MockService) Save(ctx context.Context, question *domain.Question) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, question)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockServiceMockRecorder) Save(ctx, question any) *MockServiceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockService)(nil).Save), ctx, question)
	return &MockServiceSaveCall{Call: call}
}

// MockServiceSaveCall wrap *gomock.Call
type MockServiceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockServiceSaveCall) Return(arg0 int64, arg1 error) *MockServiceSaveCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockServiceSaveCall) Do(f func(context.Context, *domain.Question) (int64, error)) *MockServiceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockServiceSaveCall) DoAndReturn(f func(context.Context, *domain.Question) (int64, error)) *MockServiceSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
