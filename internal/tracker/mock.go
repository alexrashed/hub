package tracker

import (
	"github.com/artifacthub/hub/internal/hub"
	"github.com/stretchr/testify/mock"
)

// ErrorsCollectorMock is mock ErrorsCollector implementation.
type ErrorsCollectorMock struct {
	mock.Mock
}

// Append implements the ErrorsCollector interface.
func (m *ErrorsCollectorMock) Append(repositoryID string, err error) {
	m.Called(repositoryID, err)
}

// Flush implements the ErrorsCollector interface.
func (m *ErrorsCollectorMock) Flush() {
	m.Called()
}

// Init implements the ErrorsCollector interface.
func (m *ErrorsCollectorMock) Init(repositoryID string) {
	m.Called(repositoryID)
}

// SourceMock is mock TrackerSource implementation.
type SourceMock struct {
	mock.Mock
}

// GetPackagesAvailable implements the TrackerSource interface.
func (m *SourceMock) GetPackagesAvailable() (map[string]*hub.Package, error) {
	args := m.Called()
	data, _ := args.Get(0).(map[string]*hub.Package)
	return data, args.Error(1)
}
