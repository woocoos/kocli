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
	client, err := ent.Open("sqlite3", "file:portal?mode=memory&cache=shared&_fk=1", ent.Debug())
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

func (s *testSuite) TestGenAppMenu() {
	s.client.AppAction.Delete().ExecX(context.Background())
	cfg := Config{
		KnockoutConfig: "../integration/resource/knockout.yaml",
		MenuConfig:     "../integration/resource/web/src/components/layout/menu.json",
		AppCode:        "resource",
		PortalClient:   s.client,
	}
	menuUpd := 0
	s.Run("init", func() {
		err := GenAppMenu(cfg)
		s.Require().NoError(err)
		menus, err := s.client.AppMenu.Query().All(context.Background())
		s.Require().NoError(err)
		s.Len(menus, 6)
		menuUpd = menus[0].ID
		actions, err := s.client.AppAction.Query().Where().All(context.Background())
		s.Require().NoError(err)
		s.Len(actions, 4)
	})
	upd, err := s.client.AppMenu.Get(context.Background(), menuUpd)
	s.Require().NoError(err)
	upd.Update().SetIcon("changed").SetUpdatedBy(0).SaveX(context.Background())
	s.Run("update", func() {
		err := GenAppMenu(cfg)
		s.Require().NoError(err)
		menus, err := s.client.AppMenu.Query().All(context.Background())
		s.Require().NoError(err)
		s.Len(menus, 6)
		actions, err := s.client.AppAction.Query().Where().All(context.Background())
		s.Require().NoError(err)
		s.Len(actions, 4)

		s.NotEqual(menus[0].Icon, "changed")
		s.Equal(menus[0].Icon, "fa")
	})
}
