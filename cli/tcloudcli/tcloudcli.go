package tcloudcli

// import (
// 	"fmt"
// 	"os"
// )

type TcloudCli struct {
	userConfig *UserConfig
}

func NewTcloudCli(userConfig *UserConfig) *TcloudCli {
	tcloudcli := &TcloudCli{
		userConfig: userConfig,
	}
	return tcloudcli
}
