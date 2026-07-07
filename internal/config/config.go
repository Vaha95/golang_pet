package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
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
}

// ConfigFile holds configuration read from a JSON file (lowest priority).
type ConfigFile struct {
	ListenHost  string `json:"server_address"`
	URLHost     string `json:"base_url"`
	FilePath    string `json:"file_storage_path"`
	DbDSN       string `json:"database_dsn"`
	EnableHttps *bool  `json:"enable_https"`
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
}

func GetConfigENV() ConfigENV {
	enableHttpsRaw := os.Getenv("ENABLE_HTTPS")
	enableHttps := (*bool)(nil)
	if parsed, err := strconv.ParseBool(enableHttpsRaw); err == nil {
		enableHttps = &parsed
	}

	return ConfigENV{
		ListenHost:    os.Getenv("SERVER_ADDRESS"),
		URLHost:       os.Getenv("BASE_URL"),
		FilePath:      os.Getenv("FILE_STORAGE_PATH"),
		DbDSN:         os.Getenv("DATABASE_DSN"),
		AuditFilePath: os.Getenv("AUDIT_FILE"),
		AuditURL:      os.Getenv("AUDIT_URL"),
		EnableHttps:   enableHttps,
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

	configFilePath := fallback(env.ConfigFilePath, flagCfg.ConfigFilePath, nil)
	configFileData := getConfigFile(configFilePath)

	enableHttps := flagCfg.EnableHttps
	if env.EnableHttps != nil {
		enableHttps = *env.EnableHttps
	}

	return Config{
		ListenHost:    fallback(env.ListenHost, flagCfg.ListenHost, &(configFileData.ListenHost)),
		URLHost:       fallback(env.URLHost, flagCfg.URLHost, &(configFileData.URLHost)),
		FilePath:      fallback(env.FilePath, flagCfg.FilePath, &(configFileData.FilePath)),
		DbDSN:         fallback(env.DbDSN, flagCfg.DbDSN, &(configFileData.DbDSN)),
		AuditFilePath: fallback(env.AuditFilePath, flagCfg.AuditFilePath, nil),
		AuditURL:      fallback(env.AuditURL, flagCfg.AuditURL, nil),
		EnableHttps:   enableHttps,
	}
}

func fallback(envVal string, flagVal string, fileVal *string) string {
	if envVal != "" {
		return envVal
	}
	if flagVal != "" {
		return flagVal
	}

	return *fileVal
}
