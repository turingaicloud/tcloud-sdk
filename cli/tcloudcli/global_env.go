package tcloudcli

type TACCGlobalEnv struct {
	RepoName      string
	LocalWorkDir  string
	LocalConfDir  string
	RemoteWorkDir string
	RemoteUserDir string

	SlurmUserlog string
}

var DEFAULT_SLURMDIR = "slurm_log"

func NewGlobalEnv() *TACCGlobalEnv {
	var env TACCGlobalEnv
	env.SlurmUserlog = DEFAULT_SLURMDIR
	return &env
}
