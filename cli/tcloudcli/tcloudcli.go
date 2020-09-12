package tcloudcli

// import (
// 	"fmt"
// 	"os"
// )

// TODO(Add more attributes)
type TcloudCli struct {
	userConfig *UserConfig
}

func NewTcloudCli(userConfig *UserConfig) *TcloudCli {
	tcloudcli := &TcloudCli{
		userConfig: userConfig,
	}
	return tcloudcli
}
