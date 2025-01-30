package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigFile = "config.yaml"

// Struct to hold all necessary configs to run the bot
type Config struct {
	Token     string `yaml:"Token"`
	BotPrefix string `yaml:"BotPrefix"`
	OwnerID   string `yaml:"OwnerID"`
}

func ReadConfig() *Config {
	var configData Config

	if _, err := os.Stat(ConfigFile); err == nil {
		// Reading file and extracting values
		var byteFile []byte
		if byteFile, err = os.ReadFile(ConfigFile); err != nil {
			log.Fatal("Error reading config file: ", err)
		}
		yaml.Unmarshal(byteFile, &configData)

		// If the user did not add their token return error
		if configData.Token == "" {
			fmt.Println("Please provide your Bot Information in the config.yaml file!")
			os.Exit(1)
		}

	} else if os.IsNotExist(err) {
		// Creating config file
		var yamlConfig []byte
		if yamlConfig, err = yaml.Marshal(configData); err != nil {
			fmt.Println("Error creating config json: ", err)
		}

		if err = os.WriteFile(ConfigFile, yamlConfig, 0640); err != nil {
			fmt.Println("Error writing file on disk, check if you have the right permissions!", err)
		}

		fmt.Println("Please make your configuration in " + ConfigFile + ".\nThen restart the app.")
		return nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		fmt.Println("Schrodinger: ", err)
		return nil
	}

	return &configData
}
