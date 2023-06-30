package appres

import (
	"context"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/mitchellh/mapstructure"
	"github.com/tsingsun/woocoo"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/woocoos/entco/ecx"
	"github.com/woocoos/entco/pkg/snowflake"
	"github.com/woocoos/entco/schemax"
	"github.com/woocoos/knockout/ent"
	"github.com/woocoos/knockout/ent/app"
	"github.com/woocoos/knockout/ent/appaction"
	"github.com/woocoos/knockout/ent/appres"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	arnSplit = ":"
)

var (
	fieldFormat = func(field string) string {
		return fmt.Sprintf("%s/*%s", field, arnSplit)
	}
)

type Config struct {
	KnockoutConfig string
	EntConfig      string
	GQLConfig      string
	Dialect        string
	DSN            string
	AppCode        string
	PortalClient   *ent.Client
}

func initApp(cfg *Config) (*woocoo.App, error) {
	frame := woocoo.New(woocoo.WithAppConfiguration(
		conf.New(
			conf.WithLocalPath(cfg.KnockoutConfig),
		).Load()),
	)
	err := snowflake.SetDefaultNode(frame.AppConfiguration().Sub("snowflake"))
	if err != nil {
		return nil, err
	}
	if cfg.PortalClient == nil {
		cfg.PortalClient, err = ent.Open(frame.AppConfiguration().String("store.portal.driverName"),
			frame.AppConfiguration().String("store.portal.dsn"), ent.Debug())
		if err != nil {
			return nil, err
		}
	}

	return frame, nil
}

// GenEntSchemaRes generate resource from ent schema
func GenEntSchemaRes(cfg Config) error {
	_, err := initApp(&cfg)
	if err != nil {
		return err
	}
	entg, err := entc.LoadGraph(cfg.EntConfig, &gen.Config{})
	if err != nil {
		return err
	}

	appid := cfg.PortalClient.App.Query().Where(app.Code(cfg.AppCode)).OnlyX(context.Background()).ID

	builder := make([]*ent.AppResCreate, 0)
	for _, schema := range entg.Schemas {
		has, err, arnfn := checkResAnnotation(schema)
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		inst, err := cfg.PortalClient.AppRes.Query().Where(appres.AppID(appid), appres.TypeName(schema.Name)).First(context.Background())
		if err != nil && !ent.IsNotFound(err) {
			return err
		}

		arn := strings.Builder{}
		if checkTenant(schema) {
			arn.WriteString(arnSplit + schemax.FieldTenantID + arnSplit)
		}
		arn.WriteString(schema.Name + arnSplit)
		arn.WriteString(arnPicker(schema, arnfn))
		apc := cfg.PortalClient.AppRes.Create().SetName(schema.Name).SetTypeName(schema.Name).
			SetArnPattern(arn.String()).SetCreatedBy(0).SetAppID(appid)
		builder = append(builder, apc)
		if inst != nil {
			apc.SetID(inst.ID)
		}
	}
	err = ecx.WithTx(context.Background(), func(ctx context.Context) (ecx.Transactor, error) {
		return cfg.PortalClient.Tx(ctx)
	}, func(itx ecx.Transactor) error {
		tx := itx.(*ent.Tx)
		err := tx.AppRes.CreateBulk(builder...).OnConflict().UpdateNewValues().Exec(context.Background())
		return err
	})
	if err == nil {
		log.Print("done")
	}
	return nil
}

// 检查schema是否为根对象,如果只有tenant_id,则不需要生成,因为tenant_id在引入tenant组件时,强制使用的
func checkResAnnotation(sch *load.Schema) (has bool, err error, fn arnFieldFunc) {
	var ann schemax.Annotation
	for ak, vals := range sch.Annotations {
		if ak == schemax.EntCoAnnotationName {
			err = mapstructure.Decode(vals, &ann)
			if err != nil {
				return
			}
		}
	}
	if len(ann.Resources) > 0 {
		has = true
		fn = func(name string) bool {
			for _, res := range ann.Resources {
				// exclude tenant_id
				if res == schemax.FieldTenantID || res == ann.TenantField {
					continue
				}
				if res == name {
					return true
				}
			}
			return false
		}
	}
	return
}

// 检查是否具有租户字段
func checkTenant(sch *load.Schema) bool {
	for _, field := range sch.Fields {
		if field.Name == schemax.FieldTenantID {
			return true
		}
	}
	var ann schemax.Annotation
	for ak, vals := range sch.Annotations {
		if ak == schemax.EntCoAnnotationName {
			err := mapstructure.Decode(vals, &ann)
			if err != nil {
				panic(err)
			}
		}
	}
	if ann.TenantField != "" {
		return true
	}
	return false
}

type arnFieldFunc func(name string) bool

func arnPicker(sch *load.Schema, fns ...arnFieldFunc) string {
	arn := strings.Builder{}
	for _, field := range sch.Fields {
		for _, fn := range fns {
			if fn(field.Name) {
				arn.WriteString(fieldFormat(field.Name))
			}
		}
	}
	return arn.String()
}

func GenGqlActions(cfg Config) error {
	_, err := initApp(&cfg)
	if err != nil {
		return err
	}
	if err := os.Chdir(filepath.Dir(cfg.GQLConfig)); err != nil {
		return err
	}

	gcfg, err := config.LoadConfig(filepath.Base(cfg.GQLConfig))
	if err != nil {
		return err
	}
	err = gcfg.LoadSchema()
	if err != nil {
		return err
	}

	appid := cfg.PortalClient.App.Query().Where(app.Code(cfg.AppCode)).OnlyX(context.Background()).ID
	builder := make([]*ent.AppActionCreate, 0)
	for _, field := range gcfg.Schema.Query.Fields {
		if strings.HasPrefix(field.Name, "__") {
			continue
		}
		inst, err := cfg.PortalClient.AppAction.Query().Where(appaction.AppID(appid), appaction.Name(field.Name)).First(context.Background())
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		apc := cfg.PortalClient.AppAction.Create().SetName(field.Name).SetKind(appaction.KindGraphql).
			SetMethod(appaction.MethodRead).SetComments(field.Description).SetCreatedBy(0).SetAppID(appid)
		if inst != nil {
			apc.SetID(inst.ID)
		}
		builder = append(builder, apc)
	}
	if gcfg.Schema.Mutation != nil {
		for _, field := range gcfg.Schema.Mutation.Fields {
			inst, err := cfg.PortalClient.AppAction.Query().Where(appaction.AppID(appid), appaction.Name(field.Name)).First(context.Background())
			if err != nil && !ent.IsNotFound(err) {
				return err
			}
			apc := cfg.PortalClient.AppAction.Create().SetName(field.Name).SetKind(appaction.KindGraphql).
				SetMethod(appaction.MethodWrite).SetComments(field.Description).SetCreatedBy(0).SetAppID(appid)
			if inst != nil {
				apc.SetID(inst.ID)
			}
			builder = append(builder, apc)
		}
	}
	err = ecx.WithTx(context.Background(), func(ctx context.Context) (ecx.Transactor, error) {
		return cfg.PortalClient.Tx(ctx)
	}, func(itx ecx.Transactor) error {
		tx := itx.(*ent.Tx)
		err := tx.AppAction.CreateBulk(builder...).OnConflict().UpdateNewValues().Exec(context.Background())
		return err
	})
	if err == nil {
		log.Print("done")
	}
	return err
}
