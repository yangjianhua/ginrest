package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const (
	COOKIE_AUTH_KEY = "ginrest_auth_key"
	VERSION         = "1.0.0"
)

var CONFIG = &Config{
	ServerPort: 0,
	DbPort:     0,
	DbHost:     "",
	DbName:     "",
	DbUser:     "",
	DbPassword: "",
	DbUrl:      "",
	Debug:      false,
	DbPrefix:   "",
}

type Config struct {
	ServerPort int
	DbPort     int
	DbHost     string
	DbName     string
	DbUser     string
	DbPassword string
	DbUrl      string
	Debug      bool
	DbPrefix   string
}

func (this *Config) validate() string {
	strRet := ""

	if this.ServerPort == 0 {
		strRet = "ServerPort Must Config"
	}
	if this.DbHost == "" {
		strRet = "Database Host Must Config"
	}
	if this.DbPort == 0 {
		strRet = "Database Port Must Config"
	}
	if this.DbName == "" {
		strRet = "Database Name Must Config"
	}
	if this.DbUser == "" {
		strRet = "Database User Must Config"
	}
	if this.DbPassword == "" {
		strRet = "Database Password Must Config"
	}

	if this.DbPrefix != "" {
		b, _ := regexp.MatchString("[A-Za-z0-9_-]", this.DbPrefix)
		if !b {
			panic("Config Error: DbPrefix Not Valid String, only a-z,A-Z,_-")
		}
		if len(this.DbPrefix) > 10 {
			panic("Config Error: Dbprefix max 10 chars")
		}
	}

	return strRet
}

// Load Config From Local File which located in current path.
func LoadConfig() {
	ex, err := os.Executable()
	if err != nil {
		panic("Get Path Error.")
	}
	path := filepath.Dir(ex)
	file := path + "/conf.json"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Cannot load config file")
	} else {
		err := json.Unmarshal(content, CONFIG)
		if err != nil {
			fmt.Println(err.Error())
			panic("config file ERROR.")
		}
		strRet := CONFIG.validate()
		if strRet != "" {
			panic("CONFIG Error: " + strRet)
		}

		// Format MySQL URL, for connecting use in every controller
		CONFIG.DbUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", CONFIG.DbUser, CONFIG.DbPassword, CONFIG.DbHost, CONFIG.DbPort, CONFIG.DbName)
	}
}
