package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const ConfigFile = "config.json"

// Struct to hold all necessary configs to run the bot
type Config struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	OwnerID   string `json:"OwnerID"`
}

func ReadConfig() *Config {
	var configData Config

	if _, err := os.Stat(ConfigFile); err == nil {
		// Reading file and extracting values
		var byteFile []byte
		if byteFile, err = os.ReadFile(ConfigFile); err != nil {
			log.Fatal("Error reading config file: ", err)
		}
		json.Unmarshal(byteFile, &configData)

		// If the user did not add their token return error
		if configData.Token == "" {
			fmt.Println("Please provide your Bot Information in the config.json file!")
			os.Exit(1)
		}

	} else if os.IsNotExist(err) {
		// Creating config file
		var jsonConfig []byte
		if jsonConfig, err = json.MarshalIndent(configData, "", " "); err != nil {
			fmt.Println("Error creating config json: ", err)
		}

		if err = os.WriteFile(ConfigFile, jsonConfig, 0640); err != nil {
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
