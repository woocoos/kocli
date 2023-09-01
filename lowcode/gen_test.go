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
}

func Test_GenSuite(t *testing.T) {
	suite.Run(t, new(genSuite))
}

func (s *genSuite) SetupSuite() {
	assert := gen.Assets{}
	assert.AddDir("internal/integration/lowcodetest")
	assert.Add("internal/integration/lowcodetest/.editorconfig", []byte(`# http://editorconfig.org
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
	kodir, err := filepath.Abs("internal/integration/lowcodetest")
	s.Require().NoError(err)
	cfg := &graph.Config{
		Schema: "./testdata/protable/schema.json",
		Target: kodir,
	}
	err = Generate(cfg)
	s.Require().NoError(err)
}
