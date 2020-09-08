package main

import (
    "os"
    
    "tcloud-sdk/cli/cmd"
)

func main() {
	// home := homeDIR()
	cmd.Execute()
}

func homeDIR() string {
	// if runtime.GOOS == "windows" {
	// 	return os.Getenv("userprofile")
	// }
	return os.Getenv("HOME")
}
