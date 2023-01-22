// Code generated by MockGen. DO NOT EDIT.
// Source: deployer/manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// MockDeploymentManager is a mock of DeploymentManager interface.
type MockDeploymentManager struct {
	ctrl     *gomock.Controller
	recorder *MockDeploymentManagerMockRecorder
}

// MockDeploymentManagerMockRecorder is the mock recorder for MockDeploymentManager.
type MockDeploymentManagerMockRecorder struct {
	mock *MockDeploymentManager
}

// NewMockDeploymentManager creates a new mock instance.
func NewMockDeploymentManager(ctrl *gomock.Controller) *MockDeploymentManager {
	mock := &MockDeploymentManager{ctrl: ctrl}
	mock.recorder = &MockDeploymentManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeploymentManager) EXPECT() *MockDeploymentManagerMockRecorder {
	return m.recorder
}

// CancelAll mocks base method.
func (m *MockDeploymentManager) CancelAll() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelAll")
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelAll indicates an expected call of CancelAll.
func (mr *MockDeploymentManagerMockRecorder) CancelAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelAll", reflect.TypeOf((*MockDeploymentManager)(nil).CancelAll))
}

// CancelWorkloads mocks base method.
func (m *MockDeploymentManager) CancelWorkloads(workloads map[uint32]map[string]bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelWorkloads", workloads)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelWorkloads indicates an expected call of CancelWorkloads.
func (mr *MockDeploymentManagerMockRecorder) CancelWorkloads(workloads interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelWorkloads", reflect.TypeOf((*MockDeploymentManager)(nil).CancelWorkloads), workloads)
}

// Commit mocks base method.
func (m *MockDeploymentManager) Commit(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockDeploymentManagerMockRecorder) Commit(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockDeploymentManager)(nil).Commit), ctx)
}

// GetContractIDs mocks base method.
func (m *MockDeploymentManager) GetContractIDs() map[uint32]uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractIDs")
	ret0, _ := ret[0].(map[uint32]uint64)
	return ret0
}

// GetContractIDs indicates an expected call of GetContractIDs.
func (mr *MockDeploymentManagerMockRecorder) GetContractIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractIDs", reflect.TypeOf((*MockDeploymentManager)(nil).GetContractIDs))
}

// GetDeployment mocks base method.
func (m *MockDeploymentManager) GetDeployment(nodeID uint32) (gridtypes.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeployment", nodeID)
	ret0, _ := ret[0].(gridtypes.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeployment indicates an expected call of GetDeployment.
func (mr *MockDeploymentManagerMockRecorder) GetDeployment(nodeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeployment", reflect.TypeOf((*MockDeploymentManager)(nil).GetDeployment), nodeID)
}

// GetWorkload mocks base method.
func (m *MockDeploymentManager) GetWorkload(nodeID uint32, name string) (gridtypes.Workload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkload", nodeID, name)
	ret0, _ := ret[0].(gridtypes.Workload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkload indicates an expected call of GetWorkload.
func (mr *MockDeploymentManagerMockRecorder) GetWorkload(nodeID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkload", reflect.TypeOf((*MockDeploymentManager)(nil).GetWorkload), nodeID, name)
}

// SetWorkloads mocks base method.
func (m *MockDeploymentManager) SetWorkloads(workloads map[uint32][]gridtypes.Workload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWorkloads", workloads)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWorkloads indicates an expected call of SetWorkloads.
func (mr *MockDeploymentManagerMockRecorder) SetWorkloads(workloads interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWorkloads", reflect.TypeOf((*MockDeploymentManager)(nil).SetWorkloads), workloads)
}

// Stage mocks base method.
func (m *MockDeploymentManager) Stage(workloadGenerator workloads.WorkloadGenerator, nodeID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stage", workloadGenerator, nodeID)
	ret0, _ := ret[1].(error)
	return ret0
}

// Stage indicates an expected call of Stage.
func (mr *MockDeploymentManagerMockRecorder) Stage(workloadGenerator, nodeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stage", reflect.TypeOf((*MockDeploymentManager)(nil).Stage), workloadGenerator, nodeID)
}
