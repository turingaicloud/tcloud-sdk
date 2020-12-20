package tcloudcli

type TACCGlobalEnv struct {
	RepoName      string
	LocalWorkDir  string
	LocalConfDir  string
	RemoteWorkDir string
	RemoteUserDir string

	SlurmUserlog string
}

func NewGlobalEnv() *TACCGlobalEnv {
	var env TACCGlobalEnv
	env.SlurmUserlog = "slurm_log"
	return &env
}
