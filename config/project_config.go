package config
// 项目配置结构体
type ProjectConfig struct {
	ProjectName       string   `yaml:"project_name"`
	BindIP            string   `yaml:"bind_ip"`
	Port              int      `yaml:"port"`
	StaticDir         string   `yaml:"static_dir"`
	ErrorDir          string   `yaml:"error_dir"`
	DefaultDocument   string   `yaml:"default_document"`
	MySQLEnabled      bool     `yaml:"mysql_enabled"`
	TablePrefix       string   `yaml:"table_prefix"`
	MySQL             MySQLConfig `yaml:"mysql"`
	DebugMode         bool     `yaml:"debug_mode"`
	ForbiddenSuffixes []string `yaml:"forbidden_suffixes"`
	ForbiddenPaths    []string `yaml:"forbidden_paths"`
	ForbiddenDirs     []string `yaml:"forbidden_dirs"`
	AllowedHosts      []string `yaml:"allowed_hosts"` // 域名白名单
	MaxConcurrentRequests int `yaml:"max_concurrent_requests"`//最大总并发数
	MaxConcurrentPerIP int `yaml:"max_concurrent_per_ip"` // 每个 IP 最多并发请求数
	MaxBodySize int64 `yaml:"max_body_size"` // 请求体最大大小（单位：字节），如 2MB = 2*1024*1024
	MaxMemoryMB int `yaml:"max_memory_mb"` // 单位：MB，超出后拒绝新请求（如 256）
	IPBlacklistFile string `yaml:"ip_blacklist_file"` // 黑名单路径
	RequestTimeoutSeconds int `yaml:"request_timeout_seconds"` // 请求处理超时时间
	CORS CorsConfig `yaml:"cors"`
	MaintenanceMode bool `yaml:"maintenance_mode"` // 维护模式开关
	GzipEnabled bool `yaml:"gzip_enabled"` // 是否启用 Gzip 压缩
}
type CorsConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}
type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Dbname          string `yaml:"dbname"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // 秒
}


