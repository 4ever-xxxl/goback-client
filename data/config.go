package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

type ConfigStruct struct {
	BackupDir         string
	RestoreDir        string
	Key               string
	RestoreToOriginal bool
	TimedBackup       bool
	FsNotify          bool
	Cloud             string
	Port              string
	Username          string
	Password          string
	JWTTOKEN          string
	RootDir           string
}

var Config = ConfigStruct{}

const LocalConfigPath = "static/config.json"

func (c *ConfigStruct) Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	c.BackupDir = "D:/TESTDIR/backup/"
	c.RestoreDir = "D:/TESTDIR/restore/"
	c.Key = "Man Always Remember Love Because Of Romance Only"
	c.RestoreToOriginal = false
	c.TimedBackup = false
	c.FsNotify = false
	c.Cloud = getEnv("CLOUD", "")
	c.Port = getEnv("PORT", "")
	c.Username = "Red"
	c.Password = "123456"
}

func (c *ConfigStruct) Save() error {
	file, err := os.Create(LocalConfigPath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(Config); err != nil {
		return fmt.Errorf("无法编码文件列表: %v", err)
	}

	return nil
}

func (c *ConfigStruct) Load() error {
	file, err := os.Open(LocalConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		return fmt.Errorf("无法解码文件列表: %v", err)
	}

	return nil
}

func init() {
	if err := Config.Load(); err != nil {
		Config.Init()
	}
}
