package interfaces

type ConfigFlags struct {
	LogLevel        *string
	Strategy        *string
	RepoURL         *string
	RepoSSHKey      *string
	RepoSSHUser     *string
	RepoBranch      *string
	RepoRemoteName  *string
	RepoBasePath    *string
	RepoRootDir     *string
	ConsulURL       *string
	ConsulACL       *string
	ConsulBasePath  *string
	ExpandJSON      *bool
	SecretsFile     *string
	AllowDeletes    *bool
	PollInterval    *int
	ValidExtensions *string
}

type IConfigFlags interface {
	Parse() ConfigFlags
}