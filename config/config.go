package config

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// envPrefix 는 설정을 덮어쓸 환경변수의 접두사다.
// 예) secret.jwtsecret -> PRODUCT_SECRET_JWTSECRET
const envPrefix = "PRODUCT"

// secretKeys 는 yaml 에 값을 두지 않고 환경변수로만 주입받는 키 목록이다.
var secretKeys = []string{
	"database.password",
	"redis.password",
}

func LoadConfig() Config {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("error load .env file: %v", err)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./files/config")

	// 환경변수가 yaml 값보다 우선한다.
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, key := range secretKeys {
		if err := viper.BindEnv(key); err != nil {
			log.Fatalf("error bind env %q: %v", key, err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error read config file: %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshal config: %v", err)
	}

	if err := validator.New().Struct(cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	return cfg
}

