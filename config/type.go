package config

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	// Jaeger JaegerConfig
	// Kafka  KafkaConfig
}

type AppConfig struct {
	Name    string
	Version string
	Debug   bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// type JaegerConfig struct {
// 	Host string
// 	Port string
// }

// type KafkaConfig struct {
// 	Brokers []string
// 	GroupID string
// }
