package integration

import (
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/namespaces"
	"github.com/rancher/rancher/tests/framework/extensions/users"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	namespaceName = "testnamespace"
)

type ProjectUserTestSuite struct {
	suite.Suite
	testUser *management.User
	client   *rancher.Client
	project  *management.Project
	session  *session.Session
}

func (p *ProjectUserTestSuite) TearDownSuite() {
	p.session.Cleanup()
}

func (p *ProjectUserTestSuite) SetupSuite() {
	testSession := session.NewSession(p.T())
	p.session = testSession

	client, err := rancher.NewClient("", testSession)
	require.NoError(p.T(), err)

	p.client = client

	projectConfig := &management.Project{
		ClusterID: "local",
		Name:      "TestProject",
	}

	testProject, err := client.Management.Project.Create(projectConfig)
	require.NoError(p.T(), err)

	p.project = testProject

	enabled := true
	user := &management.User{
		Username: "testusername",
		Password: "passwordpasswordd",
		Name:     "displayname",
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	require.NoError(p.T(), err)

	p.testUser = newUser
}

func (p *ProjectUserTestSuite) TestCreateNamespaceProjectMember() {
	subSession := p.session.NewSession()
	defer subSession.Cleanup()

	client, err := p.client.WithSession(subSession)
	require.NoError(p.T(), err)

	err = users.AddProjectMember(client, p.project, p.testUser, "project-member")
	require.NoError(p.T(), err)

	testUser, err := client.AsUser(p.testUser)
	require.NoError(p.T(), err)

	createdNamespace, err := namespaces.CreateNamespace(testUser, namespaceName, "{}", map[string]string{}, map[string]string{}, p.project)
	assert.NoError(p.T(), err)
	assert.Equal(p.T(), namespaceName, createdNamespace.Name)
}

func (p *ProjectUserTestSuite) TestCreateNamespaceProjectOwner() {
	subSession := p.session.NewSession()
	defer subSession.Cleanup()

	client, err := p.client.WithSession(subSession)
	require.NoError(p.T(), err)

	err = users.AddProjectMember(client, p.project, p.testUser, "project-owner")
	require.NoError(p.T(), err)

	testUser, err := client.AsUser(p.testUser)
	require.NoError(p.T(), err)

	createdNamespace, err := namespaces.CreateNamespace(testUser, namespaceName, "{}", map[string]string{}, map[string]string{}, p.project)
	assert.NoError(p.T(), err)
	assert.Equal(p.T(), namespaceName, createdNamespace.Name)
}

func TestProjectUserTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectUserTestSuite))
}
