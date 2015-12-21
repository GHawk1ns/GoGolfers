package test
import (
	"os"
	"testing"
	"github.com/ghawk1ns/golf/util"
	"github.com/ghawk1ns/golf/blah"
	"io/ioutil"
)


func TestLog(t *testing.T) {
	blah.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func TestConfig(t *testing.T) {
	config := util.GetConfig()

	if config.HBaseConfig.Host == "" {
		t.Error("No hostname")
	} else if config.HBaseConfig.Root == "" {
		t.Error("No zookeeper root")
	} else if config.HBaseConfig.Table == "" {
		t.Error("No table")
	}
}
