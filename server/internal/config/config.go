package config

type Config struct {
	Server   ServerConfig
	Media    MediaConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type MediaConfig struct {
	WorkerPort int
	RtcPort    int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host string
	Port int
}
