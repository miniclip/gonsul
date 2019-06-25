package configuration

import (
	"github.com/miniclip/gonsul/errorutil"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/miniclip/gonsul/interfaces"
)

const StrategyDry = "DRYRUN"
const StrategyOnce = "ONCE"
const StrategyPoll = "POLL"
const StrategyHook = "HOOK"

var config *Config

type Config struct {
	shouldClone    	bool
	logLevel       	int
	strategy       	string
	repoUrl        	string
	repoSSHKey     	string
	repoSSHUser    	string
	repoBranch     	string
	repoRemoteName 	string
	repoBasePath   	string
	repoRootDir    	string
	consulURL      	string
	consulACL      	string
	consulBasePath 	string
	expandJSON     	bool
	doSecrets      	bool
	secretsMap     	map[string]string
	allowDeletes   	bool
	pollInterval   	int
	Working			chan bool
	validExtensions	[]string
	timeout   			int
}

func GetConfig(flagParser interfaces.IConfigFlags) (*Config, error) {
	// Set our local error var
	var err error
	// Singleton check
	if config == nil {
		// Parse our flags
		flags := flagParser.Parse()
		// Build our configuration
		config, err = buildConfig(flags)
	}

	// Return the config and error (whichever state they are)
	return config, err
}

func DestroyConfig() {
	config = nil
}

func buildConfig(flags interfaces.ConfigFlags) (*Config, error) {

	// Set some local variable and some others defaulted
	var secrets map[string]string
	var err error
	clone := true
	doSecrets := false

	// Make sure we have the mandatory flags set
	if *flags.ConsulURL == "" || *flags.ValidExtensions == "" {
		flag.PrintDefaults()
		return nil, errors.New("required flags not set")
	}

	// Set our valid extensions
	extensions, err := setValidExtensions(*flags.ValidExtensions)
	if err != nil {
		return nil, err
	}

	// Make sure strategy is properly given
	strategy := strings.ToUpper(*flags.Strategy)
	if strategy != StrategyDry && strategy != StrategyOnce && strategy != StrategyPoll && strategy != StrategyHook {
		return nil, errors.New(fmt.Sprintf("strategy invalid, must be one of: %s, %s, %s, %s", StrategyDry, StrategyOnce, StrategyPoll, StrategyHook))
	}

	// Shall we use a local copy of the repository instead of cloning ourselves
	// This should be useful if we use Gonsul on a CI stack (such as Bamboo)
	// And the repo is checked out already, alleviating Gonsul work
	if *flags.RepoURL == "" && *flags.RepoRootDir != "" {
		clone = false
	}

	// Make sure log level is properly set
	errorLevel := errorutil.ErrorLevels[strings.ToUpper(*flags.LogLevel)]
	if errorLevel < errorutil.LogLevelErr {
		return nil, errors.New(fmt.Sprintf("log level invalid, must be one of: %s, %s, %s", errorutil.LogErr, errorutil.LogInfo, errorutil.LogDebug))
	}

	// Should we build a secrets map for on-the-fly mustache replacement
	if *flags.SecretsFile != "" {
		secrets, err = buildSecretsMap(*flags.SecretsFile, *flags.RepoRootDir)
		if err != nil {
			return nil, err
		}
		doSecrets = true
	}

	return &Config{
		shouldClone:    	clone,
		logLevel:       	errorLevel,
		strategy:       	strategy,
		repoUrl:        	*flags.RepoURL,
		repoSSHKey:     	*flags.RepoSSHKey,
		repoSSHUser:    	*flags.RepoSSHUser,
		repoBranch:     	*flags.RepoBranch,
		repoRemoteName: 	*flags.RepoRemoteName,
		repoBasePath:   	*flags.RepoBasePath,
		repoRootDir:    	*flags.RepoRootDir,
		consulURL:      	*flags.ConsulURL,
		consulACL:      	*flags.ConsulACL,
		consulBasePath: 	*flags.ConsulBasePath,
		expandJSON:     	*flags.ExpandJSON,
		doSecrets:      	doSecrets,
		secretsMap:     	secrets,
		allowDeletes:   	*flags.AllowDeletes,
		pollInterval:   	*flags.PollInterval,
		Working: 			make(chan bool, 1),
		validExtensions: 	extensions,
		timeout:   				*flags.Timeout,
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

func (config *Config) GetPollInterval() int {
	return config.pollInterval
}

func (config *Config) GetValidExtensions() []string {
	return config.validExtensions
}

func (config *Config) GetTimeout() int {
	return config.timeout
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
	content, err := ioutil.ReadFile(file) // just pass the file name
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

func setValidExtensions(validExtensions string) ([]string, error) {
	var extensionsArr []string

	// Try to explode the string
	extensions := strings.Split(validExtensions, ",")
	if len(extensions) < 1 {
		return nil, errors.New(fmt.Sprintf("could not open get valid extensions from flag (%s). Value given: %s", "--input-ext", validExtensions))
	}

	for _, extension := range extensions {
		extensionsArr = append(extensionsArr, extension)
	}

	return extensionsArr, nil
}
