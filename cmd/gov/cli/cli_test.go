package cli_test

import (
	"testing"

	"github.com/qba73/gov/cmd/gov/cli"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"gov": cli.Run,
	})
}

func Test(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
	})
}
