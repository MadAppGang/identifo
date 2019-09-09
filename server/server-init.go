package server

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/madappgang/identifo/model"
	"gopkg.in/yaml.v2"
)

const initConfigFilename = "init-config.yaml"

func init() {
	warningMsg := `WARNING! Config file could not be read, so the default server-config.yaml will be used for the server configuration.
Note that when using Docker container, changes made to this file won't survive the container restart.`

	configFlag := flag.String("config", "", "Path to the file that describes the location of a server configuration file")
	flag.Parse()

	if configFlag == nil || len(*configFlag) == 0 {
		log.Println("Config file path not specified.\n", warningMsg)
		loadDefaultServerConfiguration(&ServerSettings)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get current working directory: %s\n", err)
	}

	initConfigBytes, err := ioutil.ReadFile(filepath.Join(wd, initConfigFilename))
	if err != nil {
		log.Println("Cannot read init configuration file: ", err, warningMsg)
		loadDefaultServerConfiguration(&ServerSettings)
		return
	}

	ic := new(initialConfig)
	if err = yaml.Unmarshal(initConfigBytes, ic); err != nil {
		log.Println("Cannot unmarshal init configuration file: ", err, warningMsg)
		loadDefaultServerConfiguration(&ServerSettings)
		return
	}

	if err = ic.Validate(); err != nil {
		log.Println("Cannot load initial config: ", err, warningMsg)
		loadDefaultServerConfiguration(&ServerSettings)
		return
	}

	switch ic.Location {
	case "local":
		err = loadConfigFromFile(ic)
	case "etcd":
		err = loadConfigFromEtcd(ic)
	case "s3":
		err = loadConfigFromS3(ic)
	}
	if err != nil {
		log.Fatalf("Cannot load config. %s", err)
	}
}

func loadConfigFromFile(ic *initialConfig) error {
	return nil // TODO: implement
}

func loadConfigFromEtcd(ic *initialConfig) error {
	return nil // TODO: implement
}

func loadConfigFromS3(ic *initialConfig) error {
	return nil // TODO: implement
}

// initialConfig is for settings required by the server on the start.
type initialConfig struct {
	Location string `yaml:"location"`
	Folder   string `yaml:"folder"`
	Filename string `yaml:"filename"`
	Bucket   string `yaml:"bucket"`
	Key      string `yaml:"key"`
}

func (ic *initialConfig) Validate() error {
	subject := "Initial config"

	if ic == nil {
		return fmt.Errorf("Nil initial server config")
	}

	switch ic.Location {
	case "local":
		if len(ic.Filename) == 0 {
			return fmt.Errorf("%s. Empty filename", subject)
		}
	case "s3":
		if len(ic.Filename) == 0 {
			return fmt.Errorf("%s. Empty filename", subject)
		}
		if len(ic.Bucket) == 0 {
			log.Panicf("%s. Empty bucket", subject)
		}
	case "etcd":
		if len(ic.Key) == 0 {
			log.Panicf("%s. Empty key", subject)
		}
	default:
		return fmt.Errorf("Unknown location '%s'", ic.Location)
	}
	return nil
}

// loadDefaultServerConfiguration loads configuration from the yaml file and writes it to out variable.
func loadDefaultServerConfiguration(out *model.ServerSettings) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Cannot get server configuration file:", err)
	}

	// Iterate through possible config paths until we find the valid one.
	configPaths := []string{
		os.Getenv(serverConfigPathEnvName),
		"./server-config.yaml",
		"../../server/server-config.yaml",
	}

	var configFile []byte

	for _, p := range configPaths {
		if p == "" {
			continue
		}
		configFile, err = ioutil.ReadFile(filepath.Join(dir, p))
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Fatalln("Cannot read server configuration file:", err)
	}

	if err = yaml.Unmarshal(configFile, out); err != nil {
		log.Fatalln("Cannot unmarshal server configuration file:", err)
	}

	if len(os.Getenv(out.AdminAccount.LoginEnvName)) == 0 {
		if err := os.Setenv(out.AdminAccount.LoginEnvName, defaultAdminLogin); err != nil {
			log.Fatalf("Could not set default %s: %s\n", out.AdminAccount.LoginEnvName, err)
		}
		log.Printf("WARNING! %s not set. Default value will be used.\n", out.AdminAccount.LoginEnvName)
	}
	if len(os.Getenv(out.AdminAccount.PasswordEnvName)) == 0 {
		if err := os.Setenv(out.AdminAccount.PasswordEnvName, defaultAdminPassword); err != nil {
			log.Fatalf("Could not set default %s: %s\n", out.AdminAccount.PasswordEnvName, err)
		}
		log.Printf("WARNING! %s not set. Default value will be used.\n", out.AdminAccount.PasswordEnvName)
	}

	if err := out.Validate(); err != nil {
		log.Fatalln(err)
	}
}
