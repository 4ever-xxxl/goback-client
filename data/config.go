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
}

var Config = ConfigStruct{}

func (c *ConfigStruct) Init() {
	c.BackupDir = "D:/TESTDIR/backup/"
	c.RestoreDir = "D:/TESTDIR/restore/"
	c.Key = "Man Always Remember Love Because Of Romance Only"
	c.RestoreToOriginal = false
	c.TimedBackup = false
	c.FsNotify = false
	c.Cloud = "127.0.0.1"
	c.Port = "8080"
}

func init() {
	Config.Init()
}
