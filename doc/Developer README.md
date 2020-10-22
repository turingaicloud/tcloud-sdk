## Developer README.md

#### Directory description

+ Directory tree

  ~~~
  .
  ├── Makefile
  ├── README.md
  ├── cli
  │   ├── Makefile
  │   ├── cmd
  │   │   ├── add.go
  │   │   ├── config.go
  │   │   ├── download.go
  │   │   ├── init.go
  │   │   ├── install.go
  │   │   ├── log.go
  │   │   ├── ps.go
  │   │   └── submit.go
  │   ├── main.go
  │   └── tcloudcli
  │       ├── tcloudcli.go
  │       ├── tuxivconfig.go
  │       └── userconfig.go
  ├── doc
  └── example
  ~~~

+ `cli` directory

  + `cmd`
    + In this folder, each command will corresponds to a go file.
    + Each go file define the usage, args and other informations for the command, which are packaged into a function for main function to call.
  + `main.go`
    + In this file, we define the main function. 
    + New command will be added in the main function as the subcommand of `tcloud`
  + `tcloudcli`
    + In this folder, we will define the actual function of each command.
    + `tclcoudcli.go` is designed for the tcloud command line related functions, every tcloud cli functions are in this file.
    + `tuxivconfig.go` is designed for the tuxiv.conf related functions, which will be called in the tcloud cli functions.
    + `userconfig.go` is designed for the user configuration ralated functions, which will be called in tcloud cli functions.

#### How to add a new command(top down view)

Here we use an example of `tcloud submit` to illustrate the development process

+ Add command to `main.go`

  ~~~go
  unc newTcloudCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
  	var tcloudCmd = &cobra.Command{
  		Use:     "tcloud",
  		Short:   "TACC Command-line Interface v" + VERSION,
  		Version: VERSION,
  	}
  	tcloudCmd.AddCommand(cmd.NewSubmitCommand(cli)) // add new command
  	...
  	return tcloudCmd
  }
  ~~~

+ Add command to `cmd` directory

  ~~~go
  package cmd
  
  import (
  	"github.com/spf13/cobra"
  	"tcloud-sdk/cli/tcloudcli"
  )
  
  func NewSubmitCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
  	return &cobra.Command{
  		Use:   "submit",
  		Short: "Submit a job to TACC",
  		Args:  cobra.MaximumNArgs(1),
  		Run: func(cmd *cobra.Command, args []string) {
  			cli.XSubmit(args...)
  		},
  	}
  }
  ~~~

+ Add command to `tcloudcli/tcloudcli.go`

  ~~~go
  func (tcloudcli *TcloudCli) XSubmit(args ...string) bool {
  	...
  	return false
  }
  ~~~

  
