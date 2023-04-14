package config

import (
	"github.com/miniclip/gonsul/internal/util"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/namsral/flag"
)

const StrategyDry = "DRYRUN"
const StrategyOnce = "ONCE"
const StrategyPoll = "POLL"
const StrategyHook = "HOOK"

type config struct {
	shouldClone     bool
	logLevel        int
	strategy        string
	repoUrl         string
	repoSSHKey      string
	repoSSHUser     string
	repoBranch      string
	repoRemoteName  string
	repoBasePath    string
	repoRootDir     string
	consulURL       string
	consulACL       string
	consulBasePath  string
	keyFile         string
	certFile        string
	caFile          string
	expandJSON      bool
	expandYAML      bool
	doSecrets       bool
	secretsMap      map[string]string
	allowDeletes    string
	pollInterval    int
	Working         chan bool
	validExtensions []string
	keepFileExt     bool
	timeout         int
	version         bool
}

// IConfig is our config interface, implemented by our config struct above. It allows
// us to pass along an interface so we can mock and test any function that receives it
type IConfig interface {
	IsCloning() bool
	GetLogLevel() int
	GetStrategy() string
	GetRepoURL() string
	GetRepoSSHKey() string
	GetRepoSSHUser() string
	GetRepoBranch() string
	GetRepoRemoteName() string
	GetRepoBasePath() string
	GetRepoRootDir() string
	GetConsulURL() string
	GetConsulACL() string
	GetConsulBasePath() string
	GetKeyFile() string
	GetCaFile() string
	GetCertFile() string
	ShouldExpandJSON() bool
	ShouldExpandYAML() bool
	DoSecrets() bool
	GetSecretsMap() map[string]string
	AllowDeletes() string
	GetPollInterval() int
	WorkingChan() chan bool
	GetValidExtensions() []string
	KeepFileExt() bool
	GetTimeout() int
	IsShowVersion() bool
}

// NewConfig is our config struct constructor.
func NewConfig() (IConfig, error) {
	// Parse our flags
	flags := parseFlags()
	// Build our configuration
	return buildConfig(flags)
}

func buildConfig(flags ConfigFlags) (*config, error) {
	// Set some local variable and some others defaulted
	var secrets map[string]string
	var err error
	var clone = true
	var doSecrets = false

	// If we were passed a -v (version) flag, nothing else matters
	if *flags.Version {
		return &config{
			version: true,
		}, nil
	}

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

	// Make sure delete method is properly given
	allowDeletes := strings.ToLower(*flags.AllowDeletes)
	if allowDeletes != "true" && allowDeletes != "false" && allowDeletes != "skip" {
		return nil, errors.New(fmt.Sprintf("AllowDelete method is invalid, please define one of the following valid options as argument: true, false, skip"))
	}

	// Shall we use a local copy of the repository instead of cloning ourselves
	// This should be useful if we use Gonsul on a CI stack (such as Bamboo)
	// And the repo is checked out already, alleviating Gonsul work
	if *flags.RepoURL == "" && *flags.RepoRootDir != "" {
		clone = false
	}

	// Make sure log level is properly set
	errorLevel := util.ErrorLevels[strings.ToUpper(*flags.LogLevel)]
	if errorLevel < util.LogLevelErr {
		return nil, errors.New(fmt.Sprintf("log level invalid, must be one of: %s, %s, %s", util.LogErr, util.LogInfo, util.LogDebug))
	}

	// Should we build a secrets map for on-the-fly mustache replacement
	if *flags.SecretsFile != "" {
		secrets, err = buildSecretsMap(*flags.SecretsFile, *flags.RepoRootDir)
		if err != nil {
			return nil, err
		}
		doSecrets = true
	}

	return &config{
		shouldClone:     clone,
		logLevel:        errorLevel,
		strategy:        strategy,
		repoUrl:         *flags.RepoURL,
		repoSSHKey:      *flags.RepoSSHKey,
		repoSSHUser:     *flags.RepoSSHUser,
		repoBranch:      *flags.RepoBranch,
		repoRemoteName:  *flags.RepoRemoteName,
		repoBasePath:    *flags.RepoBasePath,
		repoRootDir:     *flags.RepoRootDir,
		consulURL:       *flags.ConsulURL,
		consulACL:       *flags.ConsulACL,
		consulBasePath:  *flags.ConsulBasePath,
		keyFile:         *flags.KeyFile,
		caFile:          *flags.CaFile,
		certFile:        *flags.CertFile,
		expandJSON:      *flags.ExpandJSON,
		expandYAML:      *flags.ExpandYAML,
		doSecrets:       doSecrets,
		secretsMap:      secrets,
		allowDeletes:    *flags.AllowDeletes,
		pollInterval:    *flags.PollInterval,
		Working:         make(chan bool, 1),
		validExtensions: extensions,
		keepFileExt:     *flags.KeepFileExt,
		timeout:         *flags.Timeout,
		version:         *flags.Version,
	}, nil
}

func (config *config) IsCloning() bool {
	return config.shouldClone
}

func (config *config) GetLogLevel() int {
	return config.logLevel
}

func (config *config) GetStrategy() string {
	return config.strategy
}

func (config *config) GetRepoURL() string {
	return config.repoUrl
}

func (config *config) GetRepoSSHKey() string {
	return config.repoSSHKey
}

func (config *config) GetRepoSSHUser() string {
	return config.repoSSHUser
}

func (config *config) GetRepoBranch() string {
	return config.repoBranch
}

func (config *config) GetRepoRemoteName() string {
	return config.repoRemoteName
}

func (config *config) GetRepoBasePath() string {
	return config.repoBasePath
}

func (config *config) GetKeyFile() string {
	return config.keyFile
}

func (config *config) GetCaFile() string {
	return config.caFile
}

func (config *config) GetCertFile() string {
	return config.certFile
}

func (config *config) GetRepoRootDir() string {
	return config.repoRootDir
}

func (config *config) GetConsulURL() string {
	return config.consulURL
}

func (config *config) GetConsulACL() string {
	return config.consulACL
}

func (config *config) GetConsulBasePath() string {
	return config.consulBasePath
}

func (config *config) ShouldExpandJSON() bool {
	return config.expandJSON
}

func (config *config) ShouldExpandYAML() bool {
	return config.expandYAML
}

func (config *config) DoSecrets() bool {
	return config.doSecrets
}

func (config *config) GetSecretsMap() map[string]string {
	return config.secretsMap
}

func (config *config) AllowDeletes() string {
	return strings.ToLower(config.allowDeletes)
}

func (config *config) GetPollInterval() int {
	return config.pollInterval
}

func (config *config) WorkingChan() chan bool {
	return config.Working
}

func (config *config) GetValidExtensions() []string {
	return config.validExtensions
}

func (config *config) KeepFileExt() bool {
	return config.keepFileExt
}

func (config *config) GetTimeout() int {
	return config.timeout
}

func (config *config) IsShowVersion() bool {
	return config.version
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
