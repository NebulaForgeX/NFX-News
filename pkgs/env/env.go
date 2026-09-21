package env

type Env string

const (
	Dev    Env = "dev"
	Secure Env = "secure"
)

func (e Env) IsSecure() bool {
	return e == Secure
}

func (e Env) IsDev() bool {
	return e == Dev
}

func (e Env) String() string {
	return string(e)
}
