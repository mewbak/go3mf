// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qmuntal/go3mf/internal/mesh (interfaces: mergeableNodes,mergeableFaces,mergeableBeams,MergeableMesh)

// Package mesh is a generated GoMock package.
package mesh

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meshinfo "github.com/qmuntal/go3mf/internal/meshinfo"
)

// MockmergeableNodes is a mock of mergeableNodes interface
type MockmergeableNodes struct {
	ctrl     *gomock.Controller
	recorder *MockmergeableNodesMockRecorder
}

// MockmergeableNodesMockRecorder is the mock recorder for MockmergeableNodes
type MockmergeableNodesMockRecorder struct {
	mock *MockmergeableNodes
}

// NewMockmergeableNodes creates a new mock instance
func NewMockmergeableNodes(ctrl *gomock.Controller) *MockmergeableNodes {
	mock := &MockmergeableNodes{ctrl: ctrl}
	mock.recorder = &MockmergeableNodesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockmergeableNodes) EXPECT() *MockmergeableNodesMockRecorder {
	return m.recorder
}

// Node mocks base method
func (m *MockmergeableNodes) Node(arg0 uint32) *Node {
	ret := m.ctrl.Call(m, "Node", arg0)
	ret0, _ := ret[0].(*Node)
	return ret0
}

// Node indicates an expected call of Node
func (mr *MockmergeableNodesMockRecorder) Node(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Node", reflect.TypeOf((*MockmergeableNodes)(nil).Node), arg0)
}

// NodeCount mocks base method
func (m *MockmergeableNodes) NodeCount() uint32 {
	ret := m.ctrl.Call(m, "NodeCount")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// NodeCount indicates an expected call of NodeCount
func (mr *MockmergeableNodesMockRecorder) NodeCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeCount", reflect.TypeOf((*MockmergeableNodes)(nil).NodeCount))
}

// MockmergeableFaces is a mock of mergeableFaces interface
type MockmergeableFaces struct {
	ctrl     *gomock.Controller
	recorder *MockmergeableFacesMockRecorder
}

// MockmergeableFacesMockRecorder is the mock recorder for MockmergeableFaces
type MockmergeableFacesMockRecorder struct {
	mock *MockmergeableFaces
}

// NewMockmergeableFaces creates a new mock instance
func NewMockmergeableFaces(ctrl *gomock.Controller) *MockmergeableFaces {
	mock := &MockmergeableFaces{ctrl: ctrl}
	mock.recorder = &MockmergeableFacesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockmergeableFaces) EXPECT() *MockmergeableFacesMockRecorder {
	return m.recorder
}

// Face mocks base method
func (m *MockmergeableFaces) Face(arg0 uint32) *Face {
	ret := m.ctrl.Call(m, "Face", arg0)
	ret0, _ := ret[0].(*Face)
	return ret0
}

// Face indicates an expected call of Face
func (mr *MockmergeableFacesMockRecorder) Face(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Face", reflect.TypeOf((*MockmergeableFaces)(nil).Face), arg0)
}

// FaceCount mocks base method
func (m *MockmergeableFaces) FaceCount() uint32 {
	ret := m.ctrl.Call(m, "FaceCount")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// FaceCount indicates an expected call of FaceCount
func (mr *MockmergeableFacesMockRecorder) FaceCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FaceCount", reflect.TypeOf((*MockmergeableFaces)(nil).FaceCount))
}

// InformationHandler mocks base method
func (m *MockmergeableFaces) InformationHandler() *meshinfo.Handler {
	ret := m.ctrl.Call(m, "InformationHandler")
	ret0, _ := ret[0].(*meshinfo.Handler)
	return ret0
}

// InformationHandler indicates an expected call of InformationHandler
func (mr *MockmergeableFacesMockRecorder) InformationHandler() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InformationHandler", reflect.TypeOf((*MockmergeableFaces)(nil).InformationHandler))
}

// MockmergeableBeams is a mock of mergeableBeams interface
type MockmergeableBeams struct {
	ctrl     *gomock.Controller
	recorder *MockmergeableBeamsMockRecorder
}

// MockmergeableBeamsMockRecorder is the mock recorder for MockmergeableBeams
type MockmergeableBeamsMockRecorder struct {
	mock *MockmergeableBeams
}

// NewMockmergeableBeams creates a new mock instance
func NewMockmergeableBeams(ctrl *gomock.Controller) *MockmergeableBeams {
	mock := &MockmergeableBeams{ctrl: ctrl}
	mock.recorder = &MockmergeableBeamsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockmergeableBeams) EXPECT() *MockmergeableBeamsMockRecorder {
	return m.recorder
}

// Beam mocks base method
func (m *MockmergeableBeams) Beam(arg0 uint32) *Beam {
	ret := m.ctrl.Call(m, "Beam", arg0)
	ret0, _ := ret[0].(*Beam)
	return ret0
}

// Beam indicates an expected call of Beam
func (mr *MockmergeableBeamsMockRecorder) Beam(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Beam", reflect.TypeOf((*MockmergeableBeams)(nil).Beam), arg0)
}

// BeamCount mocks base method
func (m *MockmergeableBeams) BeamCount() uint32 {
	ret := m.ctrl.Call(m, "BeamCount")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// BeamCount indicates an expected call of BeamCount
func (mr *MockmergeableBeamsMockRecorder) BeamCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeamCount", reflect.TypeOf((*MockmergeableBeams)(nil).BeamCount))
}

// MockMergeableMesh is a mock of MergeableMesh interface
type MockMergeableMesh struct {
	ctrl     *gomock.Controller
	recorder *MockMergeableMeshMockRecorder
}

// MockMergeableMeshMockRecorder is the mock recorder for MockMergeableMesh
type MockMergeableMeshMockRecorder struct {
	mock *MockMergeableMesh
}

// NewMockMergeableMesh creates a new mock instance
func NewMockMergeableMesh(ctrl *gomock.Controller) *MockMergeableMesh {
	mock := &MockMergeableMesh{ctrl: ctrl}
	mock.recorder = &MockMergeableMeshMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMergeableMesh) EXPECT() *MockMergeableMeshMockRecorder {
	return m.recorder
}

// Beam mocks base method
func (m *MockMergeableMesh) Beam(arg0 uint32) *Beam {
	ret := m.ctrl.Call(m, "Beam", arg0)
	ret0, _ := ret[0].(*Beam)
	return ret0
}

// Beam indicates an expected call of Beam
func (mr *MockMergeableMeshMockRecorder) Beam(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Beam", reflect.TypeOf((*MockMergeableMesh)(nil).Beam), arg0)
}

// BeamCount mocks base method
func (m *MockMergeableMesh) BeamCount() uint32 {
	ret := m.ctrl.Call(m, "BeamCount")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// BeamCount indicates an expected call of BeamCount
func (mr *MockMergeableMeshMockRecorder) BeamCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeamCount", reflect.TypeOf((*MockMergeableMesh)(nil).BeamCount))
}

// Face mocks base method
func (m *MockMergeableMesh) Face(arg0 uint32) *Face {
	ret := m.ctrl.Call(m, "Face", arg0)
	ret0, _ := ret[0].(*Face)
	return ret0
}

// Face indicates an expected call of Face
func (mr *MockMergeableMeshMockRecorder) Face(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Face", reflect.TypeOf((*MockMergeableMesh)(nil).Face), arg0)
}

// FaceCount mocks base method
func (m *MockMergeableMesh) FaceCount() uint32 {
	ret := m.ctrl.Call(m, "FaceCount")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// FaceCount indicates an expected call of FaceCount
func (mr *MockMergeableMeshMockRecorder) FaceCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FaceCount", reflect.TypeOf((*MockMergeableMesh)(nil).FaceCount))
}

// InformationHandler mocks base method
func (m *MockMergeableMesh) InformationHandler() *meshinfo.Handler {
	ret := m.ctrl.Call(m, "InformationHandler")
	ret0, _ := ret[0].(*meshinfo.Handler)
	return ret0
}

// InformationHandler indicates an expected call of InformationHandler
func (mr *MockMergeableMeshMockRecorder) InformationHandler() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InformationHandler", reflect.TypeOf((*MockMergeableMesh)(nil).InformationHandler))
}

// Node mocks base method
func (m *MockMergeableMesh) Node(arg0 uint32) *Node {
	ret := m.ctrl.Call(m, "Node", arg0)
	ret0, _ := ret[0].(*Node)
	return ret0
}

// Node indicates an expected call of Node
func (mr *MockMergeableMeshMockRecorder) Node(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Node", reflect.TypeOf((*MockMergeableMesh)(nil).Node), arg0)
}

// NodeCount mocks base method
func (m *MockMergeableMesh) NodeCount() uint32 {
	ret := m.ctrl.Call(m, "NodeCount")
	ret0, _ := ret[0].(uint32)
	return ret0
}

// NodeCount indicates an expected call of NodeCount
func (mr *MockMergeableMeshMockRecorder) NodeCount() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeCount", reflect.TypeOf((*MockMergeableMesh)(nil).NodeCount))
}
