package internal

func init() {
	GetIrisDir()
	InitLogger()
	CleanupLogs()
	setupUUID()
}
