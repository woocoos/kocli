package appres

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/woocoos/entco/schemax/typex"
	"github.com/woocoos/knockout/ent"
	"github.com/woocoos/knockout/ent/app"
	"github.com/woocoos/knockout/ent/appaction"
	"github.com/woocoos/knockout/ent/migrate"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type testSuite struct {
	suite.Suite
	client *ent.Client
}

// set up
func (s *testSuite) SetupSuite() {
	client, err := ent.Open("sqlite3", "file:portal?mode=memory&cache=shared&_fk=1")
	s.Require().NoError(err)
	s.client = client
	err = client.Schema.Create(context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		migrate.WithForeignKeys(false))
	s.Require().NoError(err)
	client.App.Create().SetCreatedBy(1).SetID(1).SetName("test").
		SetStatus(typex.SimpleStatusInactive).SetCode("resource").SetKind(app.Kind("web")).
		SaveX(context.Background())
}

func TestAppres(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) TestGenEntSchemaRes() {
	cfg := Config{
		KnockoutConfig: "../integration/resource/knockout.yaml",
		EntConfig:      "../integration/resource/ent/schema",
		AppCode:        "resource",
		PortalClient:   s.client,
	}
	err := GenEntSchemaRes(cfg)
	s.Require().NoError(err)
	all, err := s.client.AppRes.Query().All(context.Background())
	s.Require().NoError(err)
	s.Len(all, 1)
	s.Equal("Resource", all[0].TypeName)
	s.Equal(":tenant_id:Resource:name/*:", all[0].ArnPattern)
}

func (s *testSuite) TestGenQqlAction() {
	cfg := Config{
		KnockoutConfig: "../integration/resource/knockout.yaml",
		GQLConfig:      "../integration/resource/gqlgen.yml",
		AppCode:        "resource",
		PortalClient:   s.client,
	}
	err := GenGqlActions(cfg)
	s.Require().NoError(err)
	all, err := s.client.AppAction.Query().All(context.Background())
	s.Require().NoError(err)
	s.Equal("node", all[0].Name)
	s.Equal(appaction.MethodRead, all[0].Method)
	s.Equal(appaction.KindGraphql, all[0].Kind)
	s.Equal("nodes", all[1].Name)
	s.Equal(appaction.MethodRead, all[1].Method)
	s.Equal(appaction.KindGraphql, all[1].Kind)

}
