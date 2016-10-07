package tracker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/powerman/narada-go/narada/staging"
)

func TestMain(m *testing.M) { os.Exit(staging.TearDown(m.Run())) }

func readTestFile(t *testing.T, path string) []byte {
	b, err := ioutil.ReadFile(filepath.Join("var", "testdata", path))
	if err != nil {
		pwd, err2 := os.Getwd()
		t.Fatal(err, pwd, err2)
	}
	return b
}
