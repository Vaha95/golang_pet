package config

import (
	"flag"
	"os"
	"strconv"
)

// ConfigENV holds configuration read from environment variables.
type ConfigENV struct {
	ListenHost    string
	URLHost       string
	FilePath      string
	DbDSN         string
	AuditFilePath string
	AuditURL      string
	EnableHttps   *bool // nil if not set or invalid
}

// ConfigFlag holds configuration read from command-line flags.
type ConfigFlag struct {
	ListenHost    string
	URLHost       string
	FilePath      string
	DbDSN         string
	AuditFilePath string
	AuditURL      string
	EnableHttps   bool
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
	flag.Parse()

	return ConfigFlag{
		ListenHost:    *listenHost,
		URLHost:       *urlHost,
		FilePath:      *filePath,
		DbDSN:         *dsn,
		AuditFilePath: *auditFile,
		AuditURL:      *auditURL,
		EnableHttps:   *enableHttps,
	}
}

func GetMainConfig() Config {
	env := GetConfigENV()
	flagCfg := GetConfigFlag()

	enableHttps := flagCfg.EnableHttps
	if env.EnableHttps != nil {
		enableHttps = *env.EnableHttps
	}

	return Config{
		ListenHost:    fallback(env.ListenHost, flagCfg.ListenHost),
		URLHost:       fallback(env.URLHost, flagCfg.URLHost),
		FilePath:      fallback(env.FilePath, flagCfg.FilePath),
		DbDSN:         fallback(env.DbDSN, flagCfg.DbDSN),
		AuditFilePath: fallback(env.AuditFilePath, flagCfg.AuditFilePath),
		AuditURL:      fallback(env.AuditURL, flagCfg.AuditURL),
		EnableHttps:   enableHttps,
	}
}

func fallback(envVal string, flagVal string) string {
	if envVal != "" {
		return envVal
	}
	return flagVal
}
