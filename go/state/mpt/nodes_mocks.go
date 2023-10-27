// Code generated by MockGen. DO NOT EDIT.
// Source: nodes.go

// Package mpt is a generated GoMock package.
package mpt

import (
	reflect "reflect"

	common "github.com/Fantom-foundation/Carmen/go/common"
	shared "github.com/Fantom-foundation/Carmen/go/state/mpt/shared"
	gomock "go.uber.org/mock/gomock"
)

// MockNode is a mock of Node interface.
type MockNode struct {
	ctrl     *gomock.Controller
	recorder *MockNodeMockRecorder
}

// MockNodeMockRecorder is the mock recorder for MockNode.
type MockNodeMockRecorder struct {
	mock *MockNode
}

// NewMockNode creates a new mock instance.
func NewMockNode(ctrl *gomock.Controller) *MockNode {
	mock := &MockNode{ctrl: ctrl}
	mock.recorder = &MockNodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNode) EXPECT() *MockNodeMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockNode) Check(source NodeSource, thisId NodeId, path []Nibble) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", source, thisId, path)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockNodeMockRecorder) Check(source, thisId, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockNode)(nil).Check), source, thisId, path)
}

// ClearStorage mocks base method.
func (m *MockNode) ClearStorage(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], address common.Address, path []Nibble) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearStorage", manager, thisId, this, address, path)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ClearStorage indicates an expected call of ClearStorage.
func (mr *MockNodeMockRecorder) ClearStorage(manager, thisId, this, address, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearStorage", reflect.TypeOf((*MockNode)(nil).ClearStorage), manager, thisId, this, address, path)
}

// Dump mocks base method.
func (m *MockNode) Dump(source NodeSource, thisId NodeId, indent string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dump", source, thisId, indent)
}

// Dump indicates an expected call of Dump.
func (mr *MockNodeMockRecorder) Dump(source, thisId, indent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dump", reflect.TypeOf((*MockNode)(nil).Dump), source, thisId, indent)
}

// Freeze mocks base method.
func (m *MockNode) Freeze(manager NodeManager, this shared.WriteHandle[Node]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Freeze", manager, this)
	ret0, _ := ret[0].(error)
	return ret0
}

// Freeze indicates an expected call of Freeze.
func (mr *MockNodeMockRecorder) Freeze(manager, this interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Freeze", reflect.TypeOf((*MockNode)(nil).Freeze), manager, this)
}

// GetAccount mocks base method.
func (m *MockNode) GetAccount(source NodeSource, address common.Address, path []Nibble) (AccountInfo, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", source, address, path)
	ret0, _ := ret[0].(AccountInfo)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockNodeMockRecorder) GetAccount(source, address, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockNode)(nil).GetAccount), source, address, path)
}

// GetHash mocks base method.
func (m *MockNode) GetHash() (common.Hash, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHash")
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetHash indicates an expected call of GetHash.
func (mr *MockNodeMockRecorder) GetHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHash", reflect.TypeOf((*MockNode)(nil).GetHash))
}

// GetSlot mocks base method.
func (m *MockNode) GetSlot(source NodeSource, address common.Address, path []Nibble, key common.Key) (common.Value, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlot", source, address, path, key)
	ret0, _ := ret[0].(common.Value)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSlot indicates an expected call of GetSlot.
func (mr *MockNodeMockRecorder) GetSlot(source, address, path, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlot", reflect.TypeOf((*MockNode)(nil).GetSlot), source, address, path, key)
}

// GetValue mocks base method.
func (m *MockNode) GetValue(source NodeSource, key common.Key, path []Nibble) (common.Value, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", source, key, path)
	ret0, _ := ret[0].(common.Value)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetValue indicates an expected call of GetValue.
func (mr *MockNodeMockRecorder) GetValue(source, key, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockNode)(nil).GetValue), source, key, path)
}

// IsFrozen mocks base method.
func (m *MockNode) IsFrozen() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFrozen")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFrozen indicates an expected call of IsFrozen.
func (mr *MockNodeMockRecorder) IsFrozen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFrozen", reflect.TypeOf((*MockNode)(nil).IsFrozen))
}

// MarkFrozen mocks base method.
func (m *MockNode) MarkFrozen() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkFrozen")
}

// MarkFrozen indicates an expected call of MarkFrozen.
func (mr *MockNodeMockRecorder) MarkFrozen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkFrozen", reflect.TypeOf((*MockNode)(nil).MarkFrozen))
}

// Release mocks base method.
func (m *MockNode) Release(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release", manager, thisId, this)
	ret0, _ := ret[0].(error)
	return ret0
}

// Release indicates an expected call of Release.
func (mr *MockNodeMockRecorder) Release(manager, thisId, this interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockNode)(nil).Release), manager, thisId, this)
}

// SetAccount mocks base method.
func (m *MockNode) SetAccount(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], address common.Address, path []Nibble, info AccountInfo) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAccount", manager, thisId, this, address, path, info)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetAccount indicates an expected call of SetAccount.
func (mr *MockNodeMockRecorder) SetAccount(manager, thisId, this, address, path, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAccount", reflect.TypeOf((*MockNode)(nil).SetAccount), manager, thisId, this, address, path, info)
}

// SetHash mocks base method.
func (m *MockNode) SetHash(arg0 common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHash", arg0)
}

// SetHash indicates an expected call of SetHash.
func (mr *MockNodeMockRecorder) SetHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHash", reflect.TypeOf((*MockNode)(nil).SetHash), arg0)
}

// SetSlot mocks base method.
func (m *MockNode) SetSlot(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], address common.Address, path []Nibble, key common.Key, value common.Value) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSlot", manager, thisId, this, address, path, key, value)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetSlot indicates an expected call of SetSlot.
func (mr *MockNodeMockRecorder) SetSlot(manager, thisId, this, address, path, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSlot", reflect.TypeOf((*MockNode)(nil).SetSlot), manager, thisId, this, address, path, key, value)
}

// SetValue mocks base method.
func (m *MockNode) SetValue(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], key common.Key, path []Nibble, value common.Value) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetValue", manager, thisId, this, key, path, value)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetValue indicates an expected call of SetValue.
func (mr *MockNodeMockRecorder) SetValue(manager, thisId, this, key, path, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValue", reflect.TypeOf((*MockNode)(nil).SetValue), manager, thisId, this, key, path, value)
}

// Visit mocks base method.
func (m *MockNode) Visit(source NodeSource, thisId NodeId, depth int, visitor NodeVisitor) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Visit", source, thisId, depth, visitor)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Visit indicates an expected call of Visit.
func (mr *MockNodeMockRecorder) Visit(source, thisId, depth, visitor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Visit", reflect.TypeOf((*MockNode)(nil).Visit), source, thisId, depth, visitor)
}

// MockNodeSource is a mock of NodeSource interface.
type MockNodeSource struct {
	ctrl     *gomock.Controller
	recorder *MockNodeSourceMockRecorder
}

// MockNodeSourceMockRecorder is the mock recorder for MockNodeSource.
type MockNodeSourceMockRecorder struct {
	mock *MockNodeSource
}

// NewMockNodeSource creates a new mock instance.
func NewMockNodeSource(ctrl *gomock.Controller) *MockNodeSource {
	mock := &MockNodeSource{ctrl: ctrl}
	mock.recorder = &MockNodeSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeSource) EXPECT() *MockNodeSourceMockRecorder {
	return m.recorder
}

// getConfig mocks base method.
func (m *MockNodeSource) getConfig() MptConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getConfig")
	ret0, _ := ret[0].(MptConfig)
	return ret0
}

// getConfig indicates an expected call of getConfig.
func (mr *MockNodeSourceMockRecorder) getConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getConfig", reflect.TypeOf((*MockNodeSource)(nil).getConfig))
}

// getHashFor mocks base method.
func (m *MockNodeSource) getHashFor(arg0 NodeId) (common.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getHashFor", arg0)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getHashFor indicates an expected call of getHashFor.
func (mr *MockNodeSourceMockRecorder) getHashFor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getHashFor", reflect.TypeOf((*MockNodeSource)(nil).getHashFor), arg0)
}

// getNode mocks base method.
func (m *MockNodeSource) getNode(arg0 NodeId) (shared.ReadHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getNode", arg0)
	ret0, _ := ret[0].(shared.ReadHandle[Node])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getNode indicates an expected call of getNode.
func (mr *MockNodeSourceMockRecorder) getNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getNode", reflect.TypeOf((*MockNodeSource)(nil).getNode), arg0)
}

// hashAddress mocks base method.
func (m *MockNodeSource) hashAddress(address common.Address) common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "hashAddress", address)
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// hashAddress indicates an expected call of hashAddress.
func (mr *MockNodeSourceMockRecorder) hashAddress(address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "hashAddress", reflect.TypeOf((*MockNodeSource)(nil).hashAddress), address)
}

// hashKey mocks base method.
func (m *MockNodeSource) hashKey(arg0 common.Key) common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "hashKey", arg0)
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// hashKey indicates an expected call of hashKey.
func (mr *MockNodeSourceMockRecorder) hashKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "hashKey", reflect.TypeOf((*MockNodeSource)(nil).hashKey), arg0)
}

// MockNodeManager is a mock of NodeManager interface.
type MockNodeManager struct {
	ctrl     *gomock.Controller
	recorder *MockNodeManagerMockRecorder
}

// MockNodeManagerMockRecorder is the mock recorder for MockNodeManager.
type MockNodeManagerMockRecorder struct {
	mock *MockNodeManager
}

// NewMockNodeManager creates a new mock instance.
func NewMockNodeManager(ctrl *gomock.Controller) *MockNodeManager {
	mock := &MockNodeManager{ctrl: ctrl}
	mock.recorder = &MockNodeManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeManager) EXPECT() *MockNodeManagerMockRecorder {
	return m.recorder
}

// createAccount mocks base method.
func (m *MockNodeManager) createAccount() (NodeId, shared.WriteHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createAccount")
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(shared.WriteHandle[Node])
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// createAccount indicates an expected call of createAccount.
func (mr *MockNodeManagerMockRecorder) createAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createAccount", reflect.TypeOf((*MockNodeManager)(nil).createAccount))
}

// createBranch mocks base method.
func (m *MockNodeManager) createBranch() (NodeId, shared.WriteHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createBranch")
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(shared.WriteHandle[Node])
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// createBranch indicates an expected call of createBranch.
func (mr *MockNodeManagerMockRecorder) createBranch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createBranch", reflect.TypeOf((*MockNodeManager)(nil).createBranch))
}

// createExtension mocks base method.
func (m *MockNodeManager) createExtension() (NodeId, shared.WriteHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createExtension")
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(shared.WriteHandle[Node])
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// createExtension indicates an expected call of createExtension.
func (mr *MockNodeManagerMockRecorder) createExtension() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createExtension", reflect.TypeOf((*MockNodeManager)(nil).createExtension))
}

// createValue mocks base method.
func (m *MockNodeManager) createValue() (NodeId, shared.WriteHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createValue")
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(shared.WriteHandle[Node])
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// createValue indicates an expected call of createValue.
func (mr *MockNodeManagerMockRecorder) createValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createValue", reflect.TypeOf((*MockNodeManager)(nil).createValue))
}

// getConfig mocks base method.
func (m *MockNodeManager) getConfig() MptConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getConfig")
	ret0, _ := ret[0].(MptConfig)
	return ret0
}

// getConfig indicates an expected call of getConfig.
func (mr *MockNodeManagerMockRecorder) getConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getConfig", reflect.TypeOf((*MockNodeManager)(nil).getConfig))
}

// getHashFor mocks base method.
func (m *MockNodeManager) getHashFor(arg0 NodeId) (common.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getHashFor", arg0)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getHashFor indicates an expected call of getHashFor.
func (mr *MockNodeManagerMockRecorder) getHashFor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getHashFor", reflect.TypeOf((*MockNodeManager)(nil).getHashFor), arg0)
}

// getMutableNode mocks base method.
func (m *MockNodeManager) getMutableNode(arg0 NodeId) (shared.WriteHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getMutableNode", arg0)
	ret0, _ := ret[0].(shared.WriteHandle[Node])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getMutableNode indicates an expected call of getMutableNode.
func (mr *MockNodeManagerMockRecorder) getMutableNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getMutableNode", reflect.TypeOf((*MockNodeManager)(nil).getMutableNode), arg0)
}

// getNode mocks base method.
func (m *MockNodeManager) getNode(arg0 NodeId) (shared.ReadHandle[Node], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getNode", arg0)
	ret0, _ := ret[0].(shared.ReadHandle[Node])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getNode indicates an expected call of getNode.
func (mr *MockNodeManagerMockRecorder) getNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getNode", reflect.TypeOf((*MockNodeManager)(nil).getNode), arg0)
}

// hashAddress mocks base method.
func (m *MockNodeManager) hashAddress(address common.Address) common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "hashAddress", address)
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// hashAddress indicates an expected call of hashAddress.
func (mr *MockNodeManagerMockRecorder) hashAddress(address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "hashAddress", reflect.TypeOf((*MockNodeManager)(nil).hashAddress), address)
}

// hashKey mocks base method.
func (m *MockNodeManager) hashKey(arg0 common.Key) common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "hashKey", arg0)
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// hashKey indicates an expected call of hashKey.
func (mr *MockNodeManagerMockRecorder) hashKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "hashKey", reflect.TypeOf((*MockNodeManager)(nil).hashKey), arg0)
}

// release mocks base method.
func (m *MockNodeManager) release(arg0 NodeId) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "release", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// release indicates an expected call of release.
func (mr *MockNodeManagerMockRecorder) release(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "release", reflect.TypeOf((*MockNodeManager)(nil).release), arg0)
}

// update mocks base method.
func (m *MockNodeManager) update(arg0 NodeId, arg1 shared.WriteHandle[Node]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// update indicates an expected call of update.
func (mr *MockNodeManagerMockRecorder) update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "update", reflect.TypeOf((*MockNodeManager)(nil).update), arg0, arg1)
}

// MockleafNode is a mock of leafNode interface.
type MockleafNode struct {
	ctrl     *gomock.Controller
	recorder *MockleafNodeMockRecorder
}

// MockleafNodeMockRecorder is the mock recorder for MockleafNode.
type MockleafNodeMockRecorder struct {
	mock *MockleafNode
}

// NewMockleafNode creates a new mock instance.
func NewMockleafNode(ctrl *gomock.Controller) *MockleafNode {
	mock := &MockleafNode{ctrl: ctrl}
	mock.recorder = &MockleafNodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockleafNode) EXPECT() *MockleafNodeMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockleafNode) Check(source NodeSource, thisId NodeId, path []Nibble) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", source, thisId, path)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockleafNodeMockRecorder) Check(source, thisId, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockleafNode)(nil).Check), source, thisId, path)
}

// ClearStorage mocks base method.
func (m *MockleafNode) ClearStorage(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], address common.Address, path []Nibble) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearStorage", manager, thisId, this, address, path)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ClearStorage indicates an expected call of ClearStorage.
func (mr *MockleafNodeMockRecorder) ClearStorage(manager, thisId, this, address, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearStorage", reflect.TypeOf((*MockleafNode)(nil).ClearStorage), manager, thisId, this, address, path)
}

// Dump mocks base method.
func (m *MockleafNode) Dump(source NodeSource, thisId NodeId, indent string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dump", source, thisId, indent)
}

// Dump indicates an expected call of Dump.
func (mr *MockleafNodeMockRecorder) Dump(source, thisId, indent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dump", reflect.TypeOf((*MockleafNode)(nil).Dump), source, thisId, indent)
}

// Freeze mocks base method.
func (m *MockleafNode) Freeze(manager NodeManager, this shared.WriteHandle[Node]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Freeze", manager, this)
	ret0, _ := ret[0].(error)
	return ret0
}

// Freeze indicates an expected call of Freeze.
func (mr *MockleafNodeMockRecorder) Freeze(manager, this interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Freeze", reflect.TypeOf((*MockleafNode)(nil).Freeze), manager, this)
}

// GetAccount mocks base method.
func (m *MockleafNode) GetAccount(source NodeSource, address common.Address, path []Nibble) (AccountInfo, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", source, address, path)
	ret0, _ := ret[0].(AccountInfo)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockleafNodeMockRecorder) GetAccount(source, address, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockleafNode)(nil).GetAccount), source, address, path)
}

// GetHash mocks base method.
func (m *MockleafNode) GetHash() (common.Hash, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHash")
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetHash indicates an expected call of GetHash.
func (mr *MockleafNodeMockRecorder) GetHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHash", reflect.TypeOf((*MockleafNode)(nil).GetHash))
}

// GetSlot mocks base method.
func (m *MockleafNode) GetSlot(source NodeSource, address common.Address, path []Nibble, key common.Key) (common.Value, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlot", source, address, path, key)
	ret0, _ := ret[0].(common.Value)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSlot indicates an expected call of GetSlot.
func (mr *MockleafNodeMockRecorder) GetSlot(source, address, path, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlot", reflect.TypeOf((*MockleafNode)(nil).GetSlot), source, address, path, key)
}

// GetValue mocks base method.
func (m *MockleafNode) GetValue(source NodeSource, key common.Key, path []Nibble) (common.Value, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", source, key, path)
	ret0, _ := ret[0].(common.Value)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetValue indicates an expected call of GetValue.
func (mr *MockleafNodeMockRecorder) GetValue(source, key, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockleafNode)(nil).GetValue), source, key, path)
}

// IsFrozen mocks base method.
func (m *MockleafNode) IsFrozen() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFrozen")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFrozen indicates an expected call of IsFrozen.
func (mr *MockleafNodeMockRecorder) IsFrozen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFrozen", reflect.TypeOf((*MockleafNode)(nil).IsFrozen))
}

// MarkFrozen mocks base method.
func (m *MockleafNode) MarkFrozen() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkFrozen")
}

// MarkFrozen indicates an expected call of MarkFrozen.
func (mr *MockleafNodeMockRecorder) MarkFrozen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkFrozen", reflect.TypeOf((*MockleafNode)(nil).MarkFrozen))
}

// Release mocks base method.
func (m *MockleafNode) Release(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node]) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Release", manager, thisId, this)
	ret0, _ := ret[0].(error)
	return ret0
}

// Release indicates an expected call of Release.
func (mr *MockleafNodeMockRecorder) Release(manager, thisId, this interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockleafNode)(nil).Release), manager, thisId, this)
}

// SetAccount mocks base method.
func (m *MockleafNode) SetAccount(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], address common.Address, path []Nibble, info AccountInfo) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAccount", manager, thisId, this, address, path, info)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetAccount indicates an expected call of SetAccount.
func (mr *MockleafNodeMockRecorder) SetAccount(manager, thisId, this, address, path, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAccount", reflect.TypeOf((*MockleafNode)(nil).SetAccount), manager, thisId, this, address, path, info)
}

// SetHash mocks base method.
func (m *MockleafNode) SetHash(arg0 common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHash", arg0)
}

// SetHash indicates an expected call of SetHash.
func (mr *MockleafNodeMockRecorder) SetHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHash", reflect.TypeOf((*MockleafNode)(nil).SetHash), arg0)
}

// SetSlot mocks base method.
func (m *MockleafNode) SetSlot(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], address common.Address, path []Nibble, key common.Key, value common.Value) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSlot", manager, thisId, this, address, path, key, value)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetSlot indicates an expected call of SetSlot.
func (mr *MockleafNodeMockRecorder) SetSlot(manager, thisId, this, address, path, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSlot", reflect.TypeOf((*MockleafNode)(nil).SetSlot), manager, thisId, this, address, path, key, value)
}

// SetValue mocks base method.
func (m *MockleafNode) SetValue(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], key common.Key, path []Nibble, value common.Value) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetValue", manager, thisId, this, key, path, value)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SetValue indicates an expected call of SetValue.
func (mr *MockleafNodeMockRecorder) SetValue(manager, thisId, this, key, path, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValue", reflect.TypeOf((*MockleafNode)(nil).SetValue), manager, thisId, this, key, path, value)
}

// Visit mocks base method.
func (m *MockleafNode) Visit(source NodeSource, thisId NodeId, depth int, visitor NodeVisitor) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Visit", source, thisId, depth, visitor)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Visit indicates an expected call of Visit.
func (mr *MockleafNodeMockRecorder) Visit(source, thisId, depth, visitor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Visit", reflect.TypeOf((*MockleafNode)(nil).Visit), source, thisId, depth, visitor)
}

// setPathLength mocks base method.
func (m *MockleafNode) setPathLength(manager NodeManager, thisId NodeId, this shared.WriteHandle[Node], length byte) (NodeId, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "setPathLength", manager, thisId, this, length)
	ret0, _ := ret[0].(NodeId)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// setPathLength indicates an expected call of setPathLength.
func (mr *MockleafNodeMockRecorder) setPathLength(manager, thisId, this, length interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setPathLength", reflect.TypeOf((*MockleafNode)(nil).setPathLength), manager, thisId, this, length)
}
