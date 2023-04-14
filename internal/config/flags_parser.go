package config

import (
	"fmt"
	"os"

	"github.com/miniclip/gonsul/internal/util"
	"github.com/namsral/flag"
)

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
	KeyFile         *string
	CaFile          *string
	CertFile        *string
	ExpandJSON      *bool
	ExpandYAML      *bool
	SecretsFile     *string
	AllowDeletes    *string
	PollInterval    *int
	ValidExtensions *string
	KeepFileExt     *bool
	Timeout         *int
	Version         *bool
}

func parseFlags() ConfigFlags {
	flags := ConfigFlags{}

	flag.CommandLine = flag.NewFlagSetWithEnvPrefix(os.Args[0], "GONSUL", flag.ExitOnError)

	flag.String(flag.DefaultConfigFlagname, "", "The path to a configuration file")

	flags.LogLevel = flag.String("log-level", util.LogErr, fmt.Sprintf("The desired log level (%s, %s, %s)", util.LogErr, util.LogInfo, util.LogDebug))
	flags.Strategy = flag.String("strategy", StrategyOnce, fmt.Sprintf("The Gonsul operation mode (%s, %s, %s, %s)", StrategyDry, StrategyOnce, StrategyPoll, StrategyHook))
	flags.RepoURL = flag.String("repo-url", "", "The repository URL (Full URL with scheme)")
	flags.RepoSSHKey = flag.String("repo-ssh-key", "", "The SSH private key location (Full path)")
	flags.RepoSSHUser = flag.String("repo-ssh-user", "git", "The SSH user name")
	flags.RepoBranch = flag.String("repo-branch", "master", "Which branch should we look at")
	flags.RepoRemoteName = flag.String("repo-remote-name", "origin", "The repository remote name")
	flags.RepoBasePath = flag.String("repo-base-path", "/", "The base directory to look from inside the repo")
	flags.RepoRootDir = flag.String("repo-root", "/tmp/gonsul/repo", "The path where the repo will be downloaded to")
	flags.ConsulURL = flag.String("consul-url", "", "(REQUIRED) The Consul URL REST API endpoint (Full URL with scheme)")
	flags.ConsulACL = flag.String("consul-acl", "", "The Consul ACL to use (Must have write on the KV following --consul-base path)")
	flags.ConsulBasePath = flag.String("consul-base-path", "", "The base KV path will be prefixed to dir path")
	flags.KeyFile = flag.String("key-file", "", "The key path for mTls")
	flags.CaFile = flag.String("ca-file", "", "The ca certificat path for mTls")
	flags.CertFile = flag.String("cert-file", "", "The certificat path for mTls")
	flags.ExpandJSON = flag.Bool("expand-json", false, "Expand and parse JSON files as full paths? (Default false)")
	flags.ExpandYAML = flag.Bool("expand-yaml", false, "Expand and parse YAML files as full paths? (Default false)")
	flags.SecretsFile = flag.String("secrets-file", "", "A key value json file with placeholders->secrets mapping, in order to do on the fly replace")
	flags.AllowDeletes = flag.String("allow-deletes", "false", "false, nothing will be done and a report on conflicting deletes will be shown; true: deletes reported conflitcs and proceeds; skip: reportes conflitcs, does not performe any deletes and proceeds syncing remaining files.) (Default false)")
	flags.PollInterval = flag.Int("poll-interval", 60, "The number of seconds for the repository polling interval")
	flags.ValidExtensions = flag.String("input-ext", "json,txt,ini", "A comma separated list of file extensions valid as input")
	flags.KeepFileExt = flag.Bool("keep-ext", false, "Do we want to keep file name extensions ? (If not set to true defaults by ommiting the file name extension.) (Default false)")
	flags.Timeout = flag.Int("timeout", 5, "The number of seconds for the client to wait for a response from Consul")
	flags.Version = flag.Bool("v", false, "Will show the Gonsul version")

	// Parse our command line flags
	flag.Parse()

	return flags
}
