package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tsingsun/woocoo/cmd/woco/project"
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
			name: "knockout-skip-run-gen",
			args: args{
				cfg: &project.Config{
					Package: "github.com/tsingsun/knockouttest",
					Header:  "//go:build ignore",
					Target:  kodir,
					Modules: []string{"knockout"},
				},
				opts: []project.Option{
					project.Extensions(New(
						WithSkipRunGen(true),
						WithTargetDir(kodir),
						WithFrontend(),
					)),
				},
			},
			wantErr: true,
		},
		{
			name: "knockout-api-app",
			args: args{
				cfg: &project.Config{
					Package: "github.com/tsingsun/knockouttest",
					Target:  kodir,
					Modules: []string{"knockout"},
				},
				opts: []project.Option{
					project.Extensions(New(
						WithTargetDir(kodir),
					)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(tt.args.cfg.Target, os.ModePerm))
			require.NoError(t, os.Chdir(tt.args.cfg.Target))
			if err := project.Generate(tt.args.cfg, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
