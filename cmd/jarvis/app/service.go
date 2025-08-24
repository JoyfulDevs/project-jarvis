package app

type JarvisService struct {
	AppToken string
	BotToken string
}

func NewJarvisService(appToken string, botToken string) *JarvisService {
	return &JarvisService{
		AppToken: appToken,
		BotToken: botToken,
	}
}
