package mocks

import (
	"context"
	"reflect"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"go.uber.org/mock/gomock"
)

// MockAppRepository is a mock of AppRepository interface
type MockAppRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAppRepositoryMockRecorder
}

// MockAppRepositoryMockRecorder is the mock recorder for MockAppRepository
type MockAppRepositoryMockRecorder struct {
	mock *MockAppRepository
}

// NewMockAppRepository creates a new mock instance
func NewMockAppRepository(ctrl *gomock.Controller) *MockAppRepository {
	mock := &MockAppRepository{ctrl: ctrl}
	mock.recorder = &MockAppRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppRepository) EXPECT() *MockAppRepositoryMockRecorder {
	return m.recorder
}

// CreateApp mocks base method
func (m *MockAppRepository) CreateApp(arg0 context.Context, arg1 *models.App) error {
	ret := m.ctrl.Call(m, "CreateApp", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateApp indicates an expected call of CreateApp.
func (mr *MockAppRepositoryMockRecorder) CreateApp(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApp", reflect.TypeOf((*MockAppRepository)(nil).CreateApp), arg0, arg1)
}

// FindAppByID mocks base method
func (m *MockAppRepository) FindAppByID(arg0 context.Context, arg1 string) (*models.App, error) {
	ret := m.ctrl.Call(m, "FindAppByID", arg0, arg1)
	ret0, _ := ret[0].(*models.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppByID indicates an expected call of FindAppByID.
func (mr *MockAppRepositoryMockRecorder) FindAppByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppByID", reflect.TypeOf((*MockAppRepository)(nil).FindAppByID), arg0, arg1)
}

// DeleteApp mocks base method
func (m *MockAppRepository) DeleteApp(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteApp", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApp indicates an expected call of DeleteApp.
func (mr *MockAppRepositoryMockRecorder) DeleteApp(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApp", reflect.TypeOf((*MockAppRepository)(nil).DeleteApp), arg0, arg1)
}

// FindAppsByUserID mocks base method
func (m *MockAppRepository) FindAppsByUserID(arg0 context.Context, arg1 string) ([]*models.App, error) {
	ret := m.ctrl.Call(m, "FindAppsByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*models.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppsByUserID indicates an expected call of FindAppsByUserID.
func (mr *MockAppRepositoryMockRecorder) FindAppsByUserID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppsByUserID", reflect.TypeOf((*MockAppRepository)(nil).FindAppsByUserID), arg0, arg1)
}

// UpdateApp mocks base method
func (m *MockAppRepository) UpdateApp(arg0 context.Context, arg1 *models.App) error {
	ret := m.ctrl.Call(m, "UpdateApp", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateApp indicates an expected call of UpdateApp.
func (mr *MockAppRepositoryMockRecorder) UpdateApp(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateApp", reflect.TypeOf((*MockAppRepository)(nil).UpdateApp), arg0, arg1)
}

// MockAppMemberRepository is a mock of AppMemberRepository interface
type MockAppMemberRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAppMemberRepositoryMockRecorder
}

// MockAppMemberRepositoryMockRecorder is the mock recorder for MockAppMemberRepository
type MockAppMemberRepositoryMockRecorder struct {
	mock *MockAppMemberRepository
}

// NewMockAppMemberRepository creates a new mock instance
func NewMockAppMemberRepository(ctrl *gomock.Controller) *MockAppMemberRepository {
	mock := &MockAppMemberRepository{ctrl: ctrl}
	mock.recorder = &MockAppMemberRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppMemberRepository) EXPECT() *MockAppMemberRepositoryMockRecorder {
	return m.recorder
}

// CreateAppMember mocks base method
func (m *MockAppMemberRepository) CreateAppMember(arg0 context.Context, arg1 *models.AppMember) error {
	ret := m.ctrl.Call(m, "CreateAppMember", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAppMember indicates an expected call of CreateAppMember
func (mr *MockAppMemberRepositoryMockRecorder) CreateAppMember(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAppMember", reflect.TypeOf((*MockAppMemberRepository)(nil).CreateAppMember), arg0, arg1)
}

// FindAppMemberByID mocks base method
func (m *MockAppMemberRepository) FindAppMemberByID(arg0 context.Context, arg1 string) (*models.AppMember, error) {
	ret := m.ctrl.Call(m, "FindAppMemberByID", arg0, arg1)
	ret0, _ := ret[0].(*models.AppMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppMemberByID indicates an expected call of FindAppMemberByID
func (mr *MockAppMemberRepositoryMockRecorder) FindAppMemberByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppMemberByID", reflect.TypeOf((*MockAppMemberRepository)(nil).FindAppMemberByID), arg0, arg1)
}

// DeleteAppMember mocks base method
func (m *MockAppMemberRepository) DeleteAppMember(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteAppMember", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAppMember indicates an expected call of DeleteAppMember
func (mr *MockAppMemberRepositoryMockRecorder) DeleteAppMember(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAppMember", reflect.TypeOf((*MockAppMemberRepository)(nil).DeleteAppMember), arg0, arg1)
}

// FindAppMembersByAppID mocks base method
func (m *MockAppMemberRepository) FindAppMembersByAppID(arg0 context.Context, arg1 string) ([]*models.AppMember, error) {
	ret := m.ctrl.Call(m, "FindAppMembersByAppID", arg0, arg1)
	ret0, _ := ret[0].([]*models.AppMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppMembersByAppID indicates an expected call of FindAppMembersByAppID
func (mr *MockAppMemberRepositoryMockRecorder) FindAppMembersByAppID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppMembersByAppID", reflect.TypeOf((*MockAppMemberRepository)(nil).FindAppMembersByAppID), arg0, arg1)
}

// UpdateAppMember mocks base method
func (m *MockAppMemberRepository) UpdateAppMember(arg0 context.Context, arg1 *models.AppMember) error {
	ret := m.ctrl.Call(m, "UpdateAppMember", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAppMember indicates an expected call of UpdateAppMember
func (mr *MockAppMemberRepositoryMockRecorder) UpdateAppMember(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAppMember", reflect.TypeOf((*MockAppMemberRepository)(nil).UpdateAppMember), arg0, arg1)
}

// FindAppMemberByAppIDAndUserId mocks base method
func (m *MockAppMemberRepository) FindAppMemberByAppIDAndUserId(arg0 context.Context, arg1, arg2 string) (*models.AppMember, error) {
	ret := m.ctrl.Call(m, "FindAppMemberByAppIDAndUserId", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.AppMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppMemberByAppIDAndUserId indicates an expected call of FindAppMemberByAppIDAndUserId
func (mr *MockAppMemberRepositoryMockRecorder) FindAppMemberByAppIDAndUserId(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppMemberByAppIDAndUserId", reflect.TypeOf((*MockAppMemberRepository)(nil).FindAppMemberByAppIDAndUserId), arg0, arg1, arg2)
}

// FindAppMembersByUserID mocks base method
func (m *MockAppMemberRepository) FindAppMembersByUserID(arg0 context.Context, arg1 string) ([]*models.AppMember, error) {
	ret := m.ctrl.Call(m, "FindAppMembersByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*models.AppMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppMembersByUserID indicates an expected call of FindAppMembersByUserID
func (mr *MockAppMemberRepositoryMockRecorder) FindAppMembersByUserID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppMembersByUserID", reflect.TypeOf((*MockAppMemberRepository)(nil).FindAppMembersByUserID), arg0, arg1)
}

// MockApiKeyRepository is a mock of ApiKeyRepository interface
type MockApiKeyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockApiKeyRepositoryMockRecorder
}

// MockApiKeyRepositoryMockRecorder is the mock recorder for MockApiKeyRepository
type MockApiKeyRepositoryMockRecorder struct {
	mock *MockApiKeyRepository
}

// NewMockApiKeyRepository creates a new mock instance
func NewMockApiKeyRepository(ctrl *gomock.Controller) *MockApiKeyRepository {
	mock := &MockApiKeyRepository{ctrl: ctrl}
	mock.recorder = &MockApiKeyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApiKeyRepository) EXPECT() *MockApiKeyRepositoryMockRecorder {
	return m.recorder
}

// CreateApiKey mocks base method
func (m *MockApiKeyRepository) CreateApiKey(arg0 context.Context, arg1 *models.ApiKey) error {
	ret := m.ctrl.Call(m, "CreateApiKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateApiKey indicates an expected call of CreateApiKey
func (mr *MockApiKeyRepositoryMockRecorder) CreateApiKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApiKey", reflect.TypeOf((*MockApiKeyRepository)(nil).CreateApiKey), arg0, arg1)
}

// FindApiKeyByID mocks base method
func (m *MockApiKeyRepository) FindApiKeyByID(arg0 context.Context, arg1 string) (*models.ApiKey, error) {
	ret := m.ctrl.Call(m, "FindApiKeyByID", arg0, arg1)
	ret0, _ := ret[0].(*models.ApiKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindApiKeyByID indicates an expected call of FindApiKeyByID
func (mr *MockApiKeyRepositoryMockRecorder) FindApiKeyByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindApiKeyByID", reflect.TypeOf((*MockApiKeyRepository)(nil).FindApiKeyByID), arg0, arg1)
}

// DeleteApiKey mocks base method
func (m *MockApiKeyRepository) DeleteApiKey(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteApiKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteApiKey indicates an expected call of DeleteApiKey
func (mr *MockApiKeyRepositoryMockRecorder) DeleteApiKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteApiKey", reflect.TypeOf((*MockApiKeyRepository)(nil).DeleteApiKey), arg0, arg1)
}

// FindApiKeysByAppID mocks base method
func (m *MockApiKeyRepository) FindApiKeysByAppID(arg0 context.Context, arg1 string) ([]*models.ApiKey, error) {
	ret := m.ctrl.Call(m, "FindApiKeysByAppID", arg0, arg1)
	ret0, _ := ret[0].([]*models.ApiKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindApiKeysByAppID indicates an expected call of FindApiKeysByAppID
func (mr *MockApiKeyRepositoryMockRecorder) FindApiKeysByAppID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindApiKeysByAppID", reflect.TypeOf((*MockApiKeyRepository)(nil).FindApiKeysByAppID), arg0, arg1)
}

// UpdateApiKey mocks base method
func (m *MockApiKeyRepository) UpdateApiKey(arg0 context.Context, arg1 *models.ApiKey) error {
	ret := m.ctrl.Call(m, "UpdateApiKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateApiKey indicates an expected call of UpdateApiKey
func (mr *MockApiKeyRepositoryMockRecorder) UpdateApiKey(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateApiKey", reflect.TypeOf((*MockApiKeyRepository)(nil).UpdateApiKey), arg0, arg1)
}

// MockUserRepository is a mock of UserRepository interface
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method
func (m *MockUserRepository) CreateUser(arg0 context.Context, arg1 *models.User) error {
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserRepositoryMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), arg0, arg1)
}

// FindUserByID mocks base method
func (m *MockUserRepository) FindUserByID(arg0 context.Context, arg1 string) (*models.User, error) {
	ret := m.ctrl.Call(m, "FindUserByID", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByID indicates an expected call of FindUserByID
func (mr *MockUserRepositoryMockRecorder) FindUserByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByID", reflect.TypeOf((*MockUserRepository)(nil).FindUserByID), arg0, arg1)
}

// FindUserByEmail mocks base method
func (m *MockUserRepository) FindUserByEmail(arg0 context.Context, arg1 string) (*models.User, error) {
	ret := m.ctrl.Call(m, "FindUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail
func (mr *MockUserRepositoryMockRecorder) FindUserByEmail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindUserByEmail), arg0, arg1)
}

// UpdateUser mocks base method
func (m *MockUserRepository) UpdateUser(arg0 context.Context, arg1 *models.User) error {
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser
func (mr *MockUserRepositoryMockRecorder) UpdateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), arg0, arg1)
}

// MockWebhookRepository is a mock of WebhookRepository interface
type MockWebhookRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWebhookRepositoryMockRecorder
}

// MockWebhookRepositoryMockRecorder is the mock recorder for MockWebhookRepository
type MockWebhookRepositoryMockRecorder struct {
	mock *MockWebhookRepository
}

// NewMockWebhookRepository creates a new mock instance
func NewMockWebhookRepository(ctrl *gomock.Controller) *MockWebhookRepository {
	mock := &MockWebhookRepository{ctrl: ctrl}
	mock.recorder = &MockWebhookRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWebhookRepository) EXPECT() *MockWebhookRepositoryMockRecorder {
	return m.recorder
}

// CreateWebhook mocks base method
func (m *MockWebhookRepository) CreateWebhook(arg0 context.Context, arg1 *models.Webhook) error {
	ret := m.ctrl.Call(m, "CreateWebhook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWebhook indicates an expected call of CreateWebhook
func (mr *MockWebhookRepositoryMockRecorder) CreateWebhook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebhook", reflect.TypeOf((*MockWebhookRepository)(nil).CreateWebhook), arg0, arg1)
}

// FindWebhookByID mocks base method
func (m *MockWebhookRepository) FindWebhookByID(arg0 context.Context, arg1 string) (*models.Webhook, error) {
	ret := m.ctrl.Call(m, "FindWebhookByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Webhook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWebhookByID indicates an expected call of FindWebhookByID
func (mr *MockWebhookRepositoryMockRecorder) FindWebhookByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWebhookByID", reflect.TypeOf((*MockWebhookRepository)(nil).FindWebhookByID), arg0, arg1)
}

// DeleteWebhook mocks base method
func (m *MockWebhookRepository) DeleteWebhook(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteWebhook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWebhook indicates an expected call of DeleteWebhook
func (mr *MockWebhookRepositoryMockRecorder) DeleteWebhook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWebhook", reflect.TypeOf((*MockWebhookRepository)(nil).DeleteWebhook), arg0, arg1)
}

// FindWebhooksByAppID mocks base method
func (m *MockWebhookRepository) FindWebhooksByAppID(arg0 context.Context, arg1 string) ([]*models.Webhook, error) {
	ret := m.ctrl.Call(m, "FindWebhooksByAppID", arg0, arg1)
	ret0, _ := ret[0].([]*models.Webhook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWebhooksByAppID indicates an expected call of FindWebhooksByAppID
func (mr *MockWebhookRepositoryMockRecorder) FindWebhooksByAppID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWebhooksByAppID", reflect.TypeOf((*MockWebhookRepository)(nil).FindWebhooksByAppID), arg0, arg1)
}

// UpdateWebhook mocks base method
func (m *MockWebhookRepository) UpdateWebhook(arg0 context.Context, arg1 *models.Webhook) error {
	ret := m.ctrl.Call(m, "UpdateWebhook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWebhook indicates an expected call of UpdateWebhook
func (mr *MockWebhookRepositoryMockRecorder) UpdateWebhook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWebhook", reflect.TypeOf((*MockWebhookRepository)(nil).UpdateWebhook), arg0, arg1)
}

// MockDomainRepository is a mock of DomainRepository interface
type MockDomainRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDomainRepositoryMockRecorder
}

// MockDomainRepositoryMockRecorder is the mock recorder for MockDomainRepository
type MockDomainRepositoryMockRecorder struct {
	mock *MockDomainRepository
}

// NewMockDomainRepository creates a new mock instance
func NewMockDomainRepository(ctrl *gomock.Controller) *MockDomainRepository {
	mock := &MockDomainRepository{ctrl: ctrl}
	mock.recorder = &MockDomainRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDomainRepository) EXPECT() *MockDomainRepositoryMockRecorder {
	return m.recorder
}

// CreateDomain mocks base method
func (m *MockDomainRepository) CreateDomain(arg0 context.Context, arg1 *models.Domain) error {
	ret := m.ctrl.Call(m, "CreateDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateDomain indicates an expected call of CreateDomain
func (mr *MockDomainRepositoryMockRecorder) CreateDomain(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDomain", reflect.TypeOf((*MockDomainRepository)(nil).CreateDomain), arg0, arg1)
}

// FindDomainByID mocks base method
func (m *MockDomainRepository) FindDomainByID(arg0 context.Context, arg1 string) (*models.Domain, error) {
	ret := m.ctrl.Call(m, "FindDomainByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDomainByID indicates an expected call of FindDomainByID
func (mr *MockDomainRepositoryMockRecorder) FindDomainByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDomainByID", reflect.TypeOf((*MockDomainRepository)(nil).FindDomainByID), arg0, arg1)
}

// DeleteDomain mocks base method
func (m *MockDomainRepository) DeleteDomain(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDomain indicates an expected call of DeleteDomain
func (mr *MockDomainRepositoryMockRecorder) DeleteDomain(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDomain", reflect.TypeOf((*MockDomainRepository)(nil).DeleteDomain), arg0, arg1)
}

// FindDomainByAppID mocks base method
func (m *MockDomainRepository) FindDomainByAppID(arg0 context.Context, arg1 string) (*models.Domain, error) {
	ret := m.ctrl.Call(m, "FindDomainByAppID", arg0, arg1)
	ret0, _ := ret[0].(*models.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDomainByAppID indicates an expected call of FindDomainsByAppID
func (mr *MockDomainRepositoryMockRecorder) FindDomainByAppID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDomainByAppID", reflect.TypeOf((*MockDomainRepository)(nil).FindDomainByAppID), arg0, arg1)
}

// UpdateDomain mocks base method
func (m *MockDomainRepository) UpdateDomain(arg0 context.Context, arg1 *models.Domain) error {
	ret := m.ctrl.Call(m, "UpdateDomain", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDomain indicates an expected call of UpdateDomain
func (mr *MockDomainRepositoryMockRecorder) UpdateDomain(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDomain", reflect.TypeOf((*MockDomainRepository)(nil).UpdateDomain), arg0, arg1)
}

// FindDomainByDomain mocks base method
func (m *MockDomainRepository) FindDomainByDomain(arg0 context.Context, arg1 string) (*models.Domain, error) {
	ret := m.ctrl.Call(m, "FindDomainByDomain", arg0, arg1)
	ret0, _ := ret[0].(*models.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDomainByDomain indicates an expected call of FindDomainByDomain
func (mr *MockDomainRepositoryMockRecorder) FindDomainByDomain(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDomainByDomain", reflect.TypeOf((*MockDomainRepository)(nil).FindDomainByDomain), arg0, arg1)
}

// MockFolderRepository is a mock of FolderRepository interface
type MockFolderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFolderRepositoryMockRecorder
}

// MockFolderRepositoryMockRecorder is the mock recorder for MockFolderRepository
type MockFolderRepositoryMockRecorder struct {
	mock *MockFolderRepository
}

// NewMockFolderRepository creates a new mock instance
func NewMockFolderRepository(ctrl *gomock.Controller) *MockFolderRepository {
	mock := &MockFolderRepository{ctrl: ctrl}
	mock.recorder = &MockFolderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFolderRepository) EXPECT() *MockFolderRepositoryMockRecorder {
	return m.recorder
}

// CreateFolder mocks base method
func (m *MockFolderRepository) CreateFolder(arg0 context.Context, arg1 *models.Folder) error {
	ret := m.ctrl.Call(m, "CreateFolder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFolder indicates an expected call of CreateFolder
func (mr *MockFolderRepositoryMockRecorder) CreateFolder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFolder", reflect.TypeOf((*MockFolderRepository)(nil).CreateFolder), arg0, arg1)
}

// FindFolderByID mocks base method
func (m *MockFolderRepository) FindFolderByID(arg0 context.Context, arg1 string) (*models.Folder, error) {
	ret := m.ctrl.Call(m, "FindFolderByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFolderByID indicates an expected call of FindFolderByID
func (mr *MockFolderRepositoryMockRecorder) FindFolderByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFolderByID", reflect.TypeOf((*MockFolderRepository)(nil).FindFolderByID), arg0, arg1)
}

// UpdateFolder mocks base method
func (m *MockFolderRepository) UpdateFolder(arg0 context.Context, arg1 *models.Folder) error {
	ret := m.ctrl.Call(m, "UpdateFolder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFolder indicates an expected call of UpdateFolder
func (mr *MockFolderRepositoryMockRecorder) UpdateFolder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFolder", reflect.TypeOf((*MockFolderRepository)(nil).UpdateFolder), arg0, arg1)
}

// DeleteFolder mocks base method
func (m *MockFolderRepository) DeleteFolder(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteFolder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFolder indicates an expected call of DeleteFolder
func (mr *MockFolderRepositoryMockRecorder) DeleteFolder(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFolder", reflect.TypeOf((*MockFolderRepository)(nil).DeleteFolder), arg0, arg1)
}

// FindFoldersByAppID mocks base method
func (m *MockFolderRepository) FindFoldersByAppID(arg0 context.Context, arg1 string) ([]*models.Folder, error) {
	ret := m.ctrl.Call(m, "FindFoldersByAppID", arg0, arg1)
	ret0, _ := ret[0].([]*models.Folder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFoldersByAppID indicates an expected call of FindFoldersByAppID
func (mr *MockFolderRepositoryMockRecorder) FindFoldersByAppID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFoldersByAppID", reflect.TypeOf((*MockFolderRepository)(nil).FindFoldersByAppID), arg0, arg1)
}

// MockFileRepository is a mock of FileRepository interface
type MockFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFileRepositoryMockRecorder
}

// MockFileRepositoryMockRecorder is the mock recorder for MockFileRepository
type MockFileRepositoryMockRecorder struct {
	mock *MockFileRepository
}

// NewMockFileRepository creates a new mock instance
func NewMockFileRepository(ctrl *gomock.Controller) *MockFileRepository {
	mock := &MockFileRepository{ctrl: ctrl}
	mock.recorder = &MockFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileRepository) EXPECT() *MockFileRepositoryMockRecorder {
	return m.recorder
}

// CreateFile mocks base method
func (m *MockFileRepository) CreateFile(arg0 context.Context, arg1 *models.File) error {
	ret := m.ctrl.Call(m, "CreateFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFile indicates an expected call of CreateFile
func (mr *MockFileRepositoryMockRecorder) CreateFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFile", reflect.TypeOf((*MockFileRepository)(nil).CreateFile), arg0, arg1)
}

// FindFileByID mocks base method
func (m *MockFileRepository) FindFileByID(arg0 context.Context, arg1 string) (*models.File, error) {
	ret := m.ctrl.Call(m, "FindFileByID", arg0, arg1)
	ret0, _ := ret[0].(*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFileByID indicates an expected call of FindFileByID
func (mr *MockFileRepositoryMockRecorder) FindFileByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFileByID", reflect.TypeOf((*MockFileRepository)(nil).FindFileByID), arg0, arg1)
}

// UpdateFile mocks base method
func (m *MockFileRepository) UpdateFile(arg0 context.Context, arg1 *models.File) error {
	ret := m.ctrl.Call(m, "UpdateFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFile indicates an expected call of UpdateFile
func (mr *MockFileRepositoryMockRecorder) UpdateFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFile", reflect.TypeOf((*MockFileRepository)(nil).UpdateFile), arg0, arg1)
}

// DeleteFile mocks base method
func (m *MockFileRepository) DeleteFile(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "DeleteFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile
func (mr *MockFileRepositoryMockRecorder) DeleteFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockFileRepository)(nil).DeleteFile), arg0, arg1)
}

// FindFilesByFolderID mocks base method
func (m *MockFileRepository) FindFilesByFolderID(arg0 context.Context, arg1 string) ([]*models.File, error) {
	ret := m.ctrl.Call(m, "FindFilesByFolderID", arg0, arg1)
	ret0, _ := ret[0].([]*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFilesByFolderID indicates an expected call of FindFilesByFolderID
func (mr *MockFileRepositoryMockRecorder) FindFilesByFolderID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFilesByFolderID", reflect.TypeOf((*MockFileRepository)(nil).FindFilesByFolderID), arg0, arg1)
}

// FindFilesByUserID  mocks base method
func (m *MockFileRepository) FindFilesByUserID(arg0 context.Context, arg1 string) ([]*models.File, error) {
	ret := m.ctrl.Call(m, "FindFilesByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFilesByUserID indicates an expected call of FindFilesByUserID
func (mr *MockFileRepositoryMockRecorder) FindFilesByUserID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFilesByUserID", reflect.TypeOf((*MockFileRepository)(nil).FindFilesByUserID), arg0, arg1)
}

// FindFilesByAppID mocks base method
func (m *MockFileRepository) FindFilesByAppID(arg0 context.Context, arg1 string) ([]*models.File, error) {
	ret := m.ctrl.Call(m, "FindFilesByAppID", arg0, arg1)
	ret0, _ := ret[0].([]*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFilesByAppID indicates an expected call of FindFilesByAppID
func (mr *MockFileRepositoryMockRecorder) FindFilesByAppID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFilesByAppID", reflect.TypeOf((*MockFileRepository)(nil).FindFilesByAppID), arg0, arg1)
}

// FindFilesByAppIDWithQuery mocks base method
func (m *MockFileRepository) FindFilesByAppIDWithQuery(arg0 context.Context, arg1 string, arg2 *dto.FileQueryDto) ([]*models.File, error) {
	ret := m.ctrl.Call(m, "FindFilesByAppIDWithQuery", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*models.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFilesByAppIDWithQuery indicates an expected call of FindFilesByAppIDWithQuery
func (mr *MockFileRepositoryMockRecorder) FindFilesByAppIDWithQuery(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFilesByAppIDWithQuery", reflect.TypeOf((*MockFileRepository)(nil).FindFilesByAppIDWithQuery), arg0, arg1, arg2)
}
