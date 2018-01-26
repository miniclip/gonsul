package configuration

import (
	"flag"
	"errors"
	"github.com/miniclip/gonsul/errorutil"
	"strings"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

const StrategyDry 		= "DRYRUN"
const StrategyOnce 		= "ONCE"
const StrategyPoll 		= "POLL"
const StrategyHook 		= "HOOK"

var logLevel 			= flag.String("log-level", errorutil.LogErr, fmt.Sprintf("The desired log level (%s, %s, %s)", errorutil.LogErr, errorutil.LogInfo, errorutil.LogDebug))
var strategyFlag 		= flag.String("strategy", StrategyOnce, fmt.Sprintf("The Gonsul operation mode (%s, %s, %s, %s)", StrategyDry, StrategyOnce, StrategyPoll, StrategyHook))
var repoURLFlag 		= flag.String("repo-url", "", "The repository URL (Full URL with scheme)")
var repoSSHKeyFlag 		= flag.String("repo-ssh-key", "", "The SSH private key location (Full path)")
var repoSSHUserFlag 	= flag.String("repo-ssh-user", "git", "The SSH user name")
var repoBranchFlag 		= flag.String("repo-branch", "master", "Which branch should we look at")
var repoRemoteNameFlag 	= flag.String("repo-remote-name", "origin", "The repository remote name")
var repoBasePathFlag 	= flag.String("repo-base-path", "/", "The base directory to look from inside the repo")
var repoRootDirFlag 	= flag.String("repo-root", "/tmp/gonsul/repo", "The path where the repo will be downloaded to")
var consulURLFlag 		= flag.String("consul-url", "", "(REQUIRED) The Consul URL REST API endpoint (Full URL with scheme)")
var consulACLFlag 		= flag.String("consul-acl", "", "(REQUIRED) The Consul ACL to use (Must have write on the KV following --consul-base path)")
var consulBasePathFlag 	= flag.String("consul-base-path", "", "The base KV path will be prefixed to dir path - DO NOT START NOR END WITH SLASH")
var expandJSONFlag 		= flag.Bool("expand-json", false, "Expand and parse JSON files as full paths?")
var secretsFile 		= flag.String("secrets-file", "", "A key value json file with placeholders->secrets mapping, in order to do on the fly replace")
var allowDeletesFlag 	= flag.Bool("allow-deletes", false, "Show Gonsul issue deletes? (If not, nothing will be done and a report on conflicting deletes will be shown)")

var config				*Config

type Config struct {
	shouldClone     	bool
	logLevel        	int
	strategy        	string
	repoUrl         	string
	repoSSHKey      	string
	repoSSHUser     	string
	repoBranch      	string
	repoRemoteName  	string
	repoBasePath    	string
	repoRootDir     	string
	consulURL       	string
	consulACL       	string
	consulBasePath    	string
	expandJSON    		bool
	doSecrets    		bool
	secretsMap			map[string]string
	allowDeletes		bool
}

func GetConfig() (*Config, error) {

	var err error

	if config == nil {
		config, err = buildConfig()
		return config, err
	}

	return config, nil
}

func buildConfig() (*Config, error) {

	var secrets 	map[string]string
	var doSecrets 	bool
	var err 		error

	// Parse our command line flags
	flag.Parse()

	// Set some local variable defaults
	clone	:= true

	// Make sure we have the mandatory flags set
	if *consulURLFlag == "" || *consulACLFlag == "" {
		flag.PrintDefaults()
		return nil, errors.New("required flags not set")
	}

	// Make sure strategy is properly given
	strategy := strings.ToUpper(*strategyFlag)
	if strategy != StrategyDry && strategy != StrategyOnce && strategy != StrategyPoll && strategy != StrategyHook {
		return nil, errors.New(fmt.Sprintf("strategy invalid, must be one of: %s, %s, %s, %s", StrategyDry, StrategyOnce, StrategyPoll, StrategyHook))
	}

	// Shall we use a local copy of the repository instead of cloning ourselves
	// This should be useful if we use Gonsul on a CI stack (such as Bamboo)
	// And the repo is checked out already, alleviating Gonsul work
	if *repoURLFlag == "" && *repoRootDirFlag != "" {
		clone = false
	}

	// Make sure log level is properly set
	errorLevel 	:= errorutil.ErrorLevels[strings.ToUpper(*logLevel)]
	if errorLevel < errorutil.LogLevelErr {
		return nil, errors.New(fmt.Sprintf("log level invalid, must be one of: %s, %s, %s", errorutil.LogErr, errorutil.LogInfo, errorutil.LogDebug))
	}

	doSecrets = false
	if *secretsFile != "" {
		secrets, err = buildSecretsMap(*secretsFile, *repoRootDirFlag)
		if err != nil {
			return nil, err
		}
		doSecrets = true
	}

	return &Config{
		shouldClone:   	clone,
		logLevel:      	errorLevel,
		strategy:      	strategy,
		repoUrl:       	*repoURLFlag,
		repoSSHKey:    	*repoSSHKeyFlag,
		repoSSHUser:   	*repoSSHUserFlag,
		repoBranch:    	*repoBranchFlag,
		repoRemoteName:	*repoRemoteNameFlag,
		repoBasePath:  	*repoBasePathFlag,
		repoRootDir:   	*repoRootDirFlag,
		consulURL:     	*consulURLFlag,
		consulACL:     	*consulACLFlag,
		consulBasePath:	*consulBasePathFlag,
		expandJSON:		*expandJSONFlag,
		doSecrets:		doSecrets,
		secretsMap:		secrets,
		allowDeletes:	*allowDeletesFlag,
	}, nil
}

func (config *Config) IsCloning() bool {
	return config.shouldClone
}

func (config *Config) GetLogLevel() int {
	return config.logLevel
}

func (config *Config) GetStrategy() string {
	return config.strategy
}

func (config *Config) GetRepoURL() string {
	return config.repoUrl
}

func (config *Config) GetRepoSSHKey() string {
	return config.repoSSHKey
}

func (config *Config) GetRepoSSHUser() string {
	return config.repoSSHUser
}

func (config *Config) GetRepoBranch() string {
	return config.repoBranch
}

func (config *Config) GetRepoRemoteName() string {
	return config.repoRemoteName
}

func (config *Config) GetRepoBasePath() string {
	return config.repoBasePath
}

func (config *Config) GetRepoRootDir() string {
	return config.repoRootDir
}

func (config *Config) GetConsulURL() string {
	return config.consulURL
}

func (config *Config) GetConsulACL() string {
	return config.consulACL
}

func (config *Config) GetConsulbasePath() string {
	return config.consulBasePath
}

func (config *Config) ShouldExpandJSON() bool {
	return config.expandJSON
}

func (config *Config) DoSecrets() bool {
	return config.doSecrets
}

func (config *Config) GetSecretsMap() map[string]string {
	return config.secretsMap
}

func (config *Config) AllowDeletes() bool {
	return config.allowDeletes
}

func buildSecretsMap(secretsFile string, repoRootPath string) (map[string]string, error) {
	var file = secretsFile
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// The file path as is is not a valid file, let's try concatenate it with base path
		file = repoRootPath + "/" + secretsFile
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// Provided file nowhere to be seen
			return nil, errors.New(fmt.Sprintf("the provided secrets file (%s) cannot be found", secretsFile))
		}
	}

	// we're still here, we got a file, open it, try to parse JSON and return our map
	content, err 	:= ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not open file (%s). Error message: %s", secretsFile, err.Error()))
	}

	var secretsMap map[string]string

	// Decode data into "generic"
	err = json.Unmarshal([]byte(content), &secretsMap)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not parse keys JSON file (%s). Error message: %s", secretsFile, err.Error()))
	}

	return secretsMap, nil
}