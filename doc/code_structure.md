## Code Structure of TCLOUD-SDK

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
    + This folder stores the description of all tcloud commands.
    + Each file corresponds to a unique tcloud command. E.g., `submit.go` is related to command `tcloud submit`.
    + The function in each file defines the information of the corresponding command,
    and is called when user runs that command. The information of the command may include:  command description, requirement of args, and callable subcommand function.
  + `main.go`
    + In this file, we define the main function of CLI. 
    + New commands can be added in the main function as the subcommand of `tcloud`.
  + `tcloudcli`
    + In this folder, we define all operations of `tcloud command-line`.
    + `tcloudcli.go` defines the concrete operations of CLI subcommands. It packages each operation into `X<func>` functions for calling.
    + `tuxivconfig.go` defines the parsing and env operations upon `tuxiv.conf`, The operations vary in: tuxiv.conf parsing, `<conf_file>` generating and `TACC_env` configuring.
    + `userconfig.go` defines the user configuration functions for CLI, which are called when initializing CLI. `userconfig` includes: __UserName__, __SSHpath__ to TACC cluster, __Authfile__ for user to authenticate TACC cluster, __Dir__ defines user's directory in TACC. (By default, Dir[0]=`RepoDir`, Dir[1]=`UserDir`), __path__ is where the file locally stored.

#### How to add a new command(top down view)

Here we set `tcloud submit` as an example to illustrate the process for adding a new command.

+ Register new command in `main.go`:

  ~~~go
  func newTcloudCommand(cli *tcloudcli.TcloudCli) *cobra.Command {
  	var tcloudCmd = &cobra.Command{
  		Use:     "tcloud",
  		Short:   "TACC Command-line Interface v" + VERSION,
  		Version: VERSION,
  	}
  	tcloudCmd.AddCommand(cmd.NewSubmitCommand(cli)) // register submit command
  	...
  	return tcloudCmd
  }
  ~~~

+ Add callable command in `cmd` directory

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

+ Add operations of `tcloud submit` in `tcloudcli/tcloudcli.go`

  ~~~go
  func (tcloudcli *TcloudCli) XSubmit(args ...string) bool {
  	...
  	return false
  }
  ~~~

  
