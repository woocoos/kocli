package lowcode

import (
	"github.com/stretchr/testify/suite"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"github.com/woocoos/kocli/lowcode/graph"
	"path/filepath"
	"testing"
)

type genSuite struct {
	suite.Suite
	tmpDir string
}

func Test_GenSuite(t *testing.T) {
	suite.Run(t, &genSuite{
		tmpDir: "testdata/tmp/lowcodetest",
	})
}

func (s *genSuite) SetupSuite() {
	assert := gen.Assets{}
	assert.AddDir(s.tmpDir)
	assert.Add(filepath.Join(s.tmpDir, ".editorconfig"), []byte(`# http://editorconfig.org
root = true

[*]
indent_style = space
indent_size = 2
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true
`))
	s.Require().NoError(assert.Write())
}

func (s *genSuite) Test_Generate_ProTable() {
	kodir, err := filepath.Abs(s.tmpDir)
	s.Require().NoError(err)
	cfg := &graph.Config{
		Schema: "./testdata/protable/schema.json",
		Target: kodir,
	}
	err = Generate(cfg)
	s.Require().NoError(err)
}
