package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

// TLS configuration constants for server.
const (
	TLS_ADDRESS = ":443"
	CERT_FILE   = "cert.pem"
	KEY_FILE    = "key.pem"
	GRPC_PORT   = ":50051"
)

// ConfigENV holds configuration read from environment variables.
type ConfigENV struct {
	ListenHost     string
	URLHost        string
	FilePath       string
	DbDSN          string
	AuditFilePath  string
	AuditURL       string
	EnableHttps    *bool // nil if not set or invalid
	ConfigFilePath string
	TLSAddress     string
	CertFile       string
	KeyFile        string
	TrustedSubnet  string
	GRPCPort       string
}

// ConfigFlag holds configuration read from command-line flags.
type ConfigFlag struct {
	ListenHost     string
	URLHost        string
	FilePath       string
	DbDSN          string
	AuditFilePath  string
	AuditURL       string
	EnableHttps    bool
	ConfigFilePath string
	TrustedSubnet  string
	GRPCPort       string
}

// ConfigFile holds configuration read from a JSON file (lowest priority).
type ConfigFile struct {
	ListenHost    string `json:"server_address"`
	URLHost       string `json:"base_url"`
	FilePath      string `json:"file_storage_path"`
	DbDSN         string `json:"database_dsn"`
	EnableHttps   *bool  `json:"enable_https"`
	TrustedSubnet string `json:"trusted_subnet"`
}

// Config holds merged application configuration (ENV takes precedence over flags).
type Config struct {
	ListenHost    string
	URLHost       string
	FilePath      string
	DbDSN         string
	AuditFilePath string
	AuditURL      string
	EnableHttps   bool
	TLSAddress    string
	CertFile      string
	KeyFile       string
	TrustedSubnet string
	GRPCPort      string
}

func GetConfigENV() ConfigENV {
	enableHttpsRaw, _ := os.LookupEnv("ENABLE_HTTPS")
	enableHttps := (*bool)(nil)
	if parsed, err := strconv.ParseBool(enableHttpsRaw); err == nil {
		enableHttps = &parsed
	}

	listenHost, _ := os.LookupEnv("SERVER_ADDRESS")
	urlHost, _ := os.LookupEnv("BASE_URL")
	filePath, _ := os.LookupEnv("FILE_STORAGE_PATH")
	dbDSN, _ := os.LookupEnv("DATABASE_DSN")
	auditFilePath, _ := os.LookupEnv("AUDIT_FILE")
	auditURL, _ := os.LookupEnv("AUDIT_URL")
	configFilePath, _ := os.LookupEnv("CONFIG")
	tlsAddress, _ := os.LookupEnv("TLS_ADDRESS")
	certFile, _ := os.LookupEnv("CERT_FILE")
	keyFile, _ := os.LookupEnv("KEY_FILE")
	trustedSubnet, _ := os.LookupEnv("TRUSTED_SUBNET")
	grpcPort, _ := os.LookupEnv("GRPC_PORT")

	return ConfigENV{
		ListenHost:     listenHost,
		URLHost:        urlHost,
		FilePath:       filePath,
		DbDSN:          dbDSN,
		AuditFilePath:  auditFilePath,
		AuditURL:       auditURL,
		ConfigFilePath: configFilePath,
		EnableHttps:    enableHttps,
		TLSAddress:     tlsAddress,
		CertFile:       certFile,
		KeyFile:        keyFile,
		TrustedSubnet:  trustedSubnet,
		GRPCPort:       grpcPort,
	}
}

func GetConfigFlag() ConfigFlag {
	listenHost := flag.String("a", `localhost:8080`, "Host for app")
	urlHost := flag.String("b", `http://localhost:8080`, "Host for url")
	filePath := flag.String("f", ``, "Storage file path")
	dsn := flag.String("d", ``, "Database dsn")
	auditFile := flag.String("audit-file", ``, "Audit file path")
	auditURL := flag.String("audit-url", ``, "Audit URL path")
	enableHttps := flag.Bool("s", false, "Enable Https")
	configPath := flag.String("c", ``, "Config file path (json)")
	flag.StringVar(configPath, "config", ``, "Config file path (json)")
	trustedSubnet := flag.String("t", "", "Trusted subnet")
	grpcPort := flag.String("g", GRPC_PORT, "gRPC server address")
	flag.Parse()

	return ConfigFlag{
		ListenHost:     *listenHost,
		URLHost:        *urlHost,
		FilePath:       *filePath,
		DbDSN:          *dsn,
		AuditFilePath:  *auditFile,
		AuditURL:       *auditURL,
		EnableHttps:    *enableHttps,
		ConfigFilePath: *configPath,
		TrustedSubnet:  *trustedSubnet,
		GRPCPort:       *grpcPort,
	}
}

func getConfigFile(configPath string) ConfigFile {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ConfigFile{}
	}

	var cfg ConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ConfigFile{}
	}

	return cfg
}

func GetMainConfig() Config {
	env := GetConfigENV()
	flagCfg := GetConfigFlag()

	configFilePath := coalesce(env.ConfigFilePath, flagCfg.ConfigFilePath, nil)
	configFileData := getConfigFile(configFilePath)

	enableHttps := flagCfg.EnableHttps
	if env.EnableHttps != nil {
		enableHttps = *env.EnableHttps
	}

	defaultPort := GRPC_PORT

	return Config{
		ListenHost:    coalesce(flagCfg.ListenHost, env.ListenHost, &(configFileData.ListenHost)),
		URLHost:       coalesce(flagCfg.URLHost, env.URLHost, &(configFileData.URLHost)),
		FilePath:      coalesce(flagCfg.FilePath, env.FilePath, &(configFileData.FilePath)),
		DbDSN:         coalesce(flagCfg.DbDSN, env.DbDSN, &(configFileData.DbDSN)),
		AuditFilePath: coalesce(flagCfg.AuditFilePath, env.AuditFilePath, nil),
		AuditURL:      coalesce(flagCfg.AuditURL, env.AuditURL, nil),
		EnableHttps:   enableHttps,
		TLSAddress:    coalesce(env.TLSAddress, TLS_ADDRESS, nil),
		CertFile:      coalesce(env.CertFile, CERT_FILE, nil),
		KeyFile:       coalesce(env.KeyFile, KEY_FILE, nil),
		TrustedSubnet: coalesce(flagCfg.TrustedSubnet, env.TrustedSubnet, &(configFileData.TrustedSubnet)),
		GRPCPort:      coalesce(flagCfg.GRPCPort, env.GRPCPort, &defaultPort),
	}
}

func coalesce(firstVal string, secondVal string, thirdVal *string) string {
	if firstVal != "" {
		return firstVal
	}
	if secondVal != "" {
		return secondVal
	}
	if thirdVal != nil {
		return *thirdVal
	}

	return ""
}
