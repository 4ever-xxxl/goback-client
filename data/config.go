package data

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

func (c *ConfigStruct) Init() {
	c.BackupDir = "D:/TESTDIR/backup/"
	c.RestoreDir = "D:/TESTDIR/restore/"
	c.Key = "Man Always Remember Love Because Of Romance Only"
	c.RestoreToOriginal = false
	c.TimedBackup = false
	c.FsNotify = false
	c.Cloud = "124.222.42.111"
	c.Port = "8000"
	c.Username = "Red"
	c.Password = "123456"
}

func init() {
	Config.Init()
}
