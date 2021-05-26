package server

import (
	"log"
	"os"

	"github.com/madappgang/identifo/model"
)

const (
	defaultAdminLogin    = "admin@admin.com"
	defaultAdminPassword = "password"
)

const warningMsg = "WARNING! Config file could not be read, so the default server-config.yaml will be used for the server configuration. Note that when using Docker container, changes made to this file won't survive the container restart."

// func init() {
// 	configFlag := flag.String("config", "", "The location of a server configuration file (local file, s3 or etcd)")
// 	etcdKeyName := flag.String("etcd_key", "identifo", "Key for config settings in etcd folder")
// 	flag.Parse()

// 	if configFlag == nil || len(*configFlag) == 0 {
// 		log.Println("Config file path not specified.")
// 		loadDefaultServerConfiguration(&ServerSettings)
// 		return
// 	}

// 	// Parse the URL and ensure there are no errors.
// 	u, err := url.Parse(*configFlag)
// 	if err != nil {
// 		log.Println("unable to parse config flag:", err, warningMsg)
// 		loadDefaultServerConfiguration(&ServerSettings)
// 		return
// 	}

// 	ic := new(initialConfig)
// 	ic.Key = *etcdKeyName
// 	ic.raw = *configFlag
// 	switch u.Scheme {
// 	case "etcd":
// 		loadConfigFromEtcd(ic, &ServerSettings)
// 	case "s3":
// 		loadConfigFromS3(ic, &ServerSettings)
// 	case "file", "":
// 		loadConfigFromFile(ic, &ServerSettings)
// 	default:
// 		loadDefaultServerConfiguration(&ServerSettings)
// 	}
// }

// func loadConfigFromFile(ic *initialConfig, out *model.ServerSettings) {
// 	log.Println("Loading server configuration from specified file...")
// 	filename := ic.raw
// 	if strings.HasPrefix(strings.ToUpper(filename), "FILE://") {
// 		filename = filename[7:]
// 	}
// 	configFile, err := os.ReadFile(filename)
// 	if err != nil {
// 		log.Fatalln("Cannot read server configuration file:", err)
// 	}

// 	if err = yaml.Unmarshal(configFile, out); err != nil {
// 		log.Fatalln("Cannot unmarshal server configuration file:", err)
// 	}

// 	if err := out.Validate(); err != nil {
// 		log.Fatalln("Invalid settings.", err)
// 	}
// 	loadAdminEnvVars(out.AdminAccount)
// 	log.Println("Server configuration loaded from the file.")
// }

// func loadConfigFromEtcd(ic *initialConfig, out *model.ServerSettings) {
// 	log.Println("Loading server configuration from the etcd...")
// 	cfg := clientv3.Config{
// 		DialTimeout: 5 * time.Second,
// 	}

// 	components := strings.Split(ic.raw[7:], "@")
// 	if len(components) > 1 {
// 		cfg.Endpoints = strings.Split(components[1], ",")
// 		creds := strings.Split(components[0], ":")
// 		if len(creds) == 2 {
// 			cfg.Username = creds[0]
// 			cfg.Password = creds[1]
// 		}
// 	} else if len(components) == 1 {
// 		cfg.Endpoints = strings.Split(components[0], ",")
// 	} else {
// 		log.Fatalf("could not get etcd endpoints from config: %s", ic.raw)
// 	}

// 	etcdClient, err := clientv3.New(cfg)
// 	if err != nil {
// 		log.Fatalf("Cannot not connect to etcd config storage: %s", err)
// 	}

// 	res, err := etcdClient.Get(context.Background(), ic.Key)
// 	if err != nil {
// 		log.Fatalf("Cannot get value by key %s: %s", ic.Key, err)
// 	}
// 	if len(res.Kvs) == 0 {
// 		log.Fatalf("Etcd: No value for key %s", ic.Key)
// 	}

// 	if err = json.Unmarshal(res.Kvs[0].Value, out); err != nil {
// 		log.Fatalf("Cannot unmarshal value of key '%s'. %s", ic.Key, err)
// 	}
// }

// func loadConfigFromS3(ic *initialConfig, out *model.ServerSettings) {
// 	log.Println("Loading server configuration from the S3 bucket...")

// 	components := strings.Split(ic.raw[5:], "@")
// 	var pathComponents []string
// 	region := ""
// 	if len(components) == 2 {
// 		region = components[0]
// 		pathComponents = strings.Split(components[1], "/")
// 	} else if len(components) == 1 {
// 		pathComponents = strings.Split(components[0], "/")
// 	} else {
// 		log.Fatalf("could not get s3 file path from config: %s", ic.raw)
// 	}
// 	if len(pathComponents) < 2 {
// 		log.Fatalf("could not get s3 file path from config: %s", ic.raw)
// 	}
// 	bucket := pathComponents[0]
// 	path := strings.Join(pathComponents[1:], "/")

// 	s3client, err := s3Storage.NewS3Client(region)
// 	if err != nil {
// 		log.Fatalf("Cannot initialize S3 client: %s.", err)
// 	}
// 	getObjInput := &s3.GetObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(path),
// 	}

// 	resp, err := s3client.GetObject(getObjInput)
// 	if err != nil {
// 		log.Fatalf("Cannot get object from S3: %s", err)
// 	}
// 	defer resp.Body.Close()

// 	if err = yaml.NewDecoder(resp.Body).Decode(out); err != nil {
// 		log.Fatalf("Cannot decode S3 response: %s", err)
// 	}
// 	log.Println("Server configuration loaded from the S3 bucket.")
// }

// // initialConfig is for settings required by the server on the start.
// type initialConfig struct {
// 	raw string
// 	Key string `yaml:"key"`
// }

// const serverConfigPathEnvName = "SERVER_CONFIG_PATH"

// // loadDefaultServerConfiguration loads configuration from the yaml file and writes it to out variable.
// func loadDefaultServerConfiguration(out *model.ServerSettings) {
// 	log.Println(warningMsg, "\n", "Loading default server configuration...")
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		log.Fatalln("Cannot get server configuration file:", err)
// 	}

// 	// Iterate through possible config paths until we find the valid one.
// 	configPaths := []string{
// 		os.Getenv(serverConfigPathEnvName),
// 		"./server-config.yaml",
// 		"../../server/server-config.yaml",
// 	}

// 	var configFile []byte

// 	for _, p := range configPaths {
// 		if p == "" {
// 			continue
// 		}
// 		configFile, err = ioutil.ReadFile(filepath.Join(dir, p))
// 		if err == nil {
// 			break
// 		}
// 	}

// 	if err != nil {
// 		log.Fatalln("Cannot read server configuration file:", err)
// 	}

// 	if err = yaml.Unmarshal(configFile, out); err != nil {
// 		log.Fatalln("Cannot unmarshal server configuration file:", err)
// 	}

// 	if err := out.Validate(); err != nil {
// 		log.Fatalln(err)
// 	}
// 	loadAdminEnvVars(out.AdminAccount)
// 	log.Println("Default server configuration loaded.")
// }

func loadAdminEnvVars(vars model.AdminAccountSettings) {
	if len(os.Getenv(vars.LoginEnvName)) == 0 {
		if err := os.Setenv(vars.LoginEnvName, defaultAdminLogin); err != nil {
			log.Fatalf("Could not set default %s: %s\n", vars.LoginEnvName, err)
		}
		log.Printf("WARNING! %s not set. Default value will be used.\n", vars.LoginEnvName)
	}
	if len(os.Getenv(vars.PasswordEnvName)) == 0 {
		if err := os.Setenv(vars.PasswordEnvName, defaultAdminPassword); err != nil {
			log.Fatalf("Could not set default %s: %s\n", vars.PasswordEnvName, err)
		}
		log.Printf("WARNING! %s not set. Default value will be used.\n", vars.PasswordEnvName)
	}
}
