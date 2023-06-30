package project

import (
	"github.com/stretchr/testify/require"
	"github.com/tsingsun/woocoo/cmd/woco/project"
	"path/filepath"
	"testing"
)

func Test_Generate(t *testing.T) {
	kodir, err := filepath.Abs("internal/integration/knockouttest")
	require.NoError(t, err)
	type args struct {
		cfg  *project.Config
		opts []project.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "knockout",
			args: args{
				cfg: &project.Config{
					Package: "github.com/tsingsun/knockouttest",
					Header:  "//go:build ignore",
					Target:  kodir,
					Modules: []string{"knockout"},
				},
				opts: []project.Option{
					project.Extensions(New(WithTargetDir(kodir))),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := project.Generate(tt.args.cfg, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
