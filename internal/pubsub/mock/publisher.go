// Code generated by MockGen. DO NOT EDIT.
// Source: publisher.go
//
// Generated by this command:
//
//	mockgen -package=mock -source=publisher.go -destination=./mock/publisher.go
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	pubsub "github.com/ucpr/mongo-streamer/internal/pubsub"
	gomock "go.uber.org/mock/gomock"
)

// MockPublishResult is a mock of PublishResult interface.
type MockPublishResult struct {
	ctrl     *gomock.Controller
	recorder *MockPublishResultMockRecorder
}

// MockPublishResultMockRecorder is the mock recorder for MockPublishResult.
type MockPublishResultMockRecorder struct {
	mock *MockPublishResult
}

// NewMockPublishResult creates a new mock instance.
func NewMockPublishResult(ctrl *gomock.Controller) *MockPublishResult {
	mock := &MockPublishResult{ctrl: ctrl}
	mock.recorder = &MockPublishResultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublishResult) EXPECT() *MockPublishResultMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPublishResult) Get(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPublishResultMockRecorder) Get(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPublishResult)(nil).Get), ctx)
}

// Ready mocks base method.
func (m *MockPublishResult) Ready() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ready")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Ready indicates an expected call of Ready.
func (mr *MockPublishResultMockRecorder) Ready() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ready", reflect.TypeOf((*MockPublishResult)(nil).Ready))
}

// MockPublisher is a mock of Publisher interface.
type MockPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockPublisherMockRecorder
}

// MockPublisherMockRecorder is the mock recorder for MockPublisher.
type MockPublisherMockRecorder struct {
	mock *MockPublisher
}

// NewMockPublisher creates a new mock instance.
func NewMockPublisher(ctrl *gomock.Controller) *MockPublisher {
	mock := &MockPublisher{ctrl: ctrl}
	mock.recorder = &MockPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublisher) EXPECT() *MockPublisherMockRecorder {
	return m.recorder
}

// AsyncPublish mocks base method.
func (m *MockPublisher) AsyncPublish(ctx context.Context, msg pubsub.Message) pubsub.PublishResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AsyncPublish", ctx, msg)
	ret0, _ := ret[0].(pubsub.PublishResult)
	return ret0
}

// AsyncPublish indicates an expected call of AsyncPublish.
func (mr *MockPublisherMockRecorder) AsyncPublish(ctx, msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsyncPublish", reflect.TypeOf((*MockPublisher)(nil).AsyncPublish), ctx, msg)
}
