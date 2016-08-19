/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	// Mongo
	MongoHost     string
	MongoPort     int
	MongoDatabase string

	// Influx
	InfluxHost     string
	InfluxPort     int
	InfluxDatabase string

	// NATS
	NatsHost string
	NatsPort int
}

func (this *Config) Parse() {
	/**
	 * Config
	 */
	/** Viper setup */
	// We can use config.yml from different locations,
	// depending if we run from
	cfgDir := os.Getenv("MAINFLUX_CORE_SERVER_CONFIG_DIR")
	if cfgDir == "" {
		// default cfg path to source dir, as we keep cfg.yml there
		cfgDir = os.Getenv("GOPATH") + "/src/github.com/mainflux/mainflux-core-server/config"
	}
	viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(cfgDir)   // path to look for the config file in
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	this.MongoHost = viper.GetString("mongo.host")
	this.MongoPort = viper.GetInt("mongo.port")
	this.MongoDatabase = viper.GetString("mongo.db")

	this.InfluxHost = viper.GetString("influx.host")
	this.InfluxPort = viper.GetInt("influx.port")
	this.InfluxDatabase = viper.GetString("influx.db")

	this.NatsHost = viper.GetString("nats.host")
	this.NatsPort = viper.GetInt("nats.port")
}
