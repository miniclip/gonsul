package configuration

import (
	"flag"
	"fmt"
	"github.com/miniclip/gonsul/errorutil"
)

type IConfigFlags interface {
	Parse() configFlags
}

type configFlags struct {
	logLevel        *string
	strategy        *string
	repoURL         *string
	repoSSHKey      *string
	repoSSHUser     *string
	repoBranch      *string
	repoRemoteName  *string
	repoBasePath    *string
	repoRootDir     *string
	consulURL       *string
	consulACL       *string
	consulBasePath  *string
	expandJSON      *bool
	secretsFile     *string
	allowDeletes    *bool
	pollInterval    *int
	validExtensions *string
}

type configFlagsParser struct {
	Flags configFlags
}

func (flags *configFlagsParser) Parse() configFlags {
	flags.Flags.logLevel 		= flag.String("log-level", errorutil.LogErr, fmt.Sprintf("The desired log level (%s, %s, %s)", errorutil.LogErr, errorutil.LogInfo, errorutil.LogDebug))
	flags.Flags.strategy 		= flag.String("strategy", StrategyOnce, fmt.Sprintf("The Gonsul operation mode (%s, %s, %s, %s)", StrategyDry, StrategyOnce, StrategyPoll, StrategyHook))
	flags.Flags.repoURL = flag.String("repo-url", "", "The repository URL (Full URL with scheme)")
	flags.Flags.repoSSHKey 		= flag.String("repo-ssh-key", "", "The SSH private key location (Full path)")
	flags.Flags.repoSSHUser 	= flag.String("repo-ssh-user", "git", "The SSH user name")
	flags.Flags.repoBranch 		= flag.String("repo-branch", "master", "Which branch should we look at")
	flags.Flags.repoRemoteName 	= flag.String("repo-remote-name", "origin", "The repository remote name")
	flags.Flags.repoBasePath 	= flag.String("repo-base-path", "/", "The base directory to look from inside the repo")
	flags.Flags.repoRootDir 	= flag.String("repo-root", "/tmp/gonsul/repo", "The path where the repo will be downloaded to")
	flags.Flags.consulURL 		= flag.String("consul-url", "", "(REQUIRED) The Consul URL REST API endpoint (Full URL with scheme)")
	flags.Flags.consulACL 		= flag.String("consul-acl", "", "(REQUIRED) The Consul ACL to use (Must have write on the KV following --consul-base path)")
	flags.Flags.consulBasePath 	= flag.String("consul-base-path", "", "The base KV path will be prefixed to dir path - DO NOT START NOR END WITH SLASH")
	flags.Flags.expandJSON 		= flag.Bool("expand-json", false, "Expand and parse JSON files as full paths?")
	flags.Flags.secretsFile 	= flag.String("secrets-file", "", "A key value json file with placeholders->secrets mapping, in order to do on the fly replace")
	flags.Flags.allowDeletes 	= flag.Bool("allow-deletes", false, "Show Gonsul issue deletes? (If not, nothing will be done and a report on conflicting deletes will be shown)")
	flags.Flags.pollInterval 	= flag.Int("poll-interval", 60, "The number of seconds for the repository polling interval")
	flags.Flags.validExtensions	= flag.String("input-ext", "json,txt,ini", "A comma separated list of file extensions valid as input")

	// Parse our command line flags
	flag.Parse()

	return flags.Flags
}
