package configuration

import (
	"flag"
	"fmt"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/interfaces"
)



type ConfigFlagsParser struct {
	Flags interfaces.ConfigFlags
}

func (flags *ConfigFlagsParser) Parse() interfaces.ConfigFlags {
	flags.Flags.LogLevel = flag.String("log-level", errorutil.LogErr, fmt.Sprintf("The desired log level (%s, %s, %s)", errorutil.LogErr, errorutil.LogInfo, errorutil.LogDebug))
	flags.Flags.Strategy 		= flag.String("strategy", StrategyOnce, fmt.Sprintf("The Gonsul operation mode (%s, %s, %s, %s)", StrategyDry, StrategyOnce, StrategyPoll, StrategyHook))
	flags.Flags.RepoURL = flag.String("repo-url", "", "The repository URL (Full URL with scheme)")
	flags.Flags.RepoSSHKey = flag.String("repo-ssh-key", "", "The SSH private key location (Full path)")
	flags.Flags.RepoSSHUser = flag.String("repo-ssh-user", "git", "The SSH user name")
	flags.Flags.RepoBranch = flag.String("repo-branch", "master", "Which branch should we look at")
	flags.Flags.RepoRemoteName = flag.String("repo-remote-name", "origin", "The repository remote name")
	flags.Flags.RepoBasePath = flag.String("repo-base-path", "/", "The base directory to look from inside the repo")
	flags.Flags.RepoRootDir = flag.String("repo-root", "/tmp/gonsul/repo", "The path where the repo will be downloaded to")
	flags.Flags.ConsulURL = flag.String("consul-url", "", "(REQUIRED) The Consul URL REST API endpoint (Full URL with scheme)")
	flags.Flags.ConsulACL = flag.String("consul-acl", "", "The Consul ACL to use (Must have write on the KV following --consul-base path)")
	flags.Flags.ConsulBasePath = flag.String("consul-base-path", "", "The base KV path will be prefixed to dir path - DO NOT START NOR END WITH SLASH")
	flags.Flags.ExpandJSON = flag.Bool("expand-json", false, "Expand and parse JSON files as full paths? (Default false)")
	flags.Flags.SecretsFile = flag.String("secrets-file", "", "A key value json file with placeholders->secrets mapping, in order to do on the fly replace")
	flags.Flags.AllowDeletes = flag.Bool("allow-deletes", false, "Show Gonsul issue deletes? (If not, nothing will be done and a report on conflicting deletes will be shown) (Default false)")
	flags.Flags.PollInterval = flag.Int("poll-interval", 60, "The number of seconds for the repository polling interval")
	flags.Flags.ValidExtensions = flag.String("input-ext", "json,txt,ini", "A comma separated list of file extensions valid as input")
	flags.Flags.Timeout = flag.Int("timeout", 5, "The number of seconds for the client to wait for a response from Consul")

	// Parse our command line flags
	flag.Parse()

	return flags.Flags
}

func NewFlagsParser() *ConfigFlagsParser {
	return &ConfigFlagsParser{
		Flags: interfaces.ConfigFlags{},
	}
}
