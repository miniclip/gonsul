package config
//
//import (
//	. "github.com/onsi/gomega"
//	"testing"
//
//	"fmt"
//	"github.com/miniclip/gonsul/tests/mocks"
//	"github.com/miniclip/gonsul/util"
//)
//
//func TestGetConfigSuccess(t *testing.T) {
//	RegisterTestingT(t)
//
//	// Instantiate our mocks
//	flagsMock := &mocks.IConfigFlags{}
//	// Get our mocked flags
//	configFlags := getConfigFlagsFor(
//		util.LogDebug,
//		StrategyOnce,
//		"",
//		"",
//		"",
//		"",
//		"",
//		"/",
//		"./..",
//		"http://consul.com",
//		"some-acl-1234567890-qwerty",
//		"",
//		false,
//		"tests/test-secrets-file-success.json",
//		false,
//		60,
//		"json,txt,ini",
//		10,
//	)
//
//	// Setup expectations
//	flagsMock.On("Parse").Return(configFlags)
//
//	configs, err := GetConfig(flagsMock)
//
//	Expect(flagsMock.AssertExpectations(t)).Should(BeTrue(), "Mocked method must be called")
//	Expect(err).To(BeNil())
//	Expect(configs).To(Not(BeNil()))
//}
//
//func TestGetConfigMultipleFail(t *testing.T) {
//	RegisterTestingT(t)
//
//	// Instantiate our mocks
//	flagsMock := &mocks.IConfigFlags{}
//	// Get our mocked flags
//	configFlags := getMultipleWrongConfigs()
//
//	for i, badConfigFlags := range configFlags {
//		// Reset our singleton config
//		DestroyConfig()
//
//		// Setup expectations
//		flagsMock.On("Parse").Return(badConfigFlags)
//
//		// Run our tested function injecting our mock
//		configs, err := GetConfig(flagsMock)
//
//		// Assert our expectations
//		Expect(flagsMock.AssertExpectations(t)).Should(BeTrue(), fmt.Sprintf("Mocked method must be called (%d)", i))
//		Expect(err).To(Not(BeNil()), fmt.Sprintf("Error must no be nil (%d)", i))
//		Expect(configs).To(BeNil(), fmt.Sprintf("Configs must be nil (%d)", i))
//	}
//}
//
//func getMultipleWrongConfigs() []ConfigFlags {
//	return []ConfigFlags{
//		getConfigFlagsFor("WRONG_LOG_LEVEL", StrategyOnce, "", "", "", "", "", "/", "./..", "http://consul.com", "some-acl-1234567890-qwerty", "", false, "tests/test-secrets-file-success.json", false, 60, "json,txt,ini", 10),
//		getConfigFlagsFor(util.LogDebug, "WRONG_STRATEGY", "", "", "", "", "", "/", "./..", "http://consul.com", "some-acl-1234567890-qwerty", "", false, "tests/test-secrets-file-success.json", false, 60, "json,txt,ini", 10),
//		getConfigFlagsFor(util.LogDebug, StrategyOnce, "", "", "", "", "", "/", "./..", "", "some-acl-1234567890-qwerty", "", false, "tests/test-secrets-file-success.json", false, 60, "json,txt,ini", 10),
//		getConfigFlagsFor(util.LogDebug, StrategyOnce, "", "", "", "", "", "/", "./..", "http://consul.com", "", "", false, "tests/test-secrets-file-success.json", false, 60, "json,txt,ini", 10),
//		getConfigFlagsFor(util.LogDebug, StrategyOnce, "", "", "", "", "", "/", "./..", "http://consul.com", "some-acl-1234567890-qwerty", "", false, "tests/test-secrets-file-success.json", false, 60, "", 10),
//		getConfigFlagsFor(util.LogDebug, StrategyOnce, "", "", "", "", "", "/", "./..", "http://consul.com", "some-acl-1234567890-qwerty", "", false, "tests/test-secrets-file-fail.json", false, 60, "json,txt,ini", 10),
//		getConfigFlagsFor(util.LogDebug, StrategyOnce, "", "", "", "", "", "/", "./..", "http://consul.com", "some-acl-1234567890-qwerty", "", false, "tests/test-secrets-file-non-existent.json", false, 60, "json,txt,ini", 10),
//	}
//}
//
//func getConfigFlagsFor(
//	ll, s, ru, rsk, rsu, rb, rrn, rbp, rr, cu, ca, cbp string,
//	ej bool,
//	sf string,
//	ad bool,
//	pi int,
//	ie string,
//	ti int,
//) ConfigFlags {
//	configFlags := ConfigFlags{
//		LogLevel:        &ll,
//		Strategy:        &s,
//		RepoURL:         &ru,
//		RepoSSHKey:      &rsk,
//		RepoSSHUser:     &rsu,
//		RepoBranch:      &rb,
//		RepoRemoteName:  &rrn,
//		RepoBasePath:    &rbp,
//		RepoRootDir:     &rr,
//		ConsulURL:       &cu,
//		ConsulACL:       &ca,
//		ConsulBasePath:  &cbp,
//		ExpandJSON:      &ej,
//		SecretsFile:     &sf,
//		AllowDeletes:    &ad,
//		PollInterval:    &pi,
//		ValidExtensions: &ie,
//		Timeout:         &ti,
//	}
//
//	return configFlags
//}
