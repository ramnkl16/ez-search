package ezsmtp

type Config struct {
	From     string `json:"from"`
	Password string `json:"password"`
	SmtpHost string `json:"smtpHost"`
	SmtpPort string `json:"smtpPort"`
}

var (
	Conf       Config
	subjectKey = "default.resetpassword.subject.key"
	bodyKey    = "default.resetpassword.body.key"
)

func SetConfig(c Config) {
	Conf = c
}
