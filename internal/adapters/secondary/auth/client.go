package auth

type Adapter struct {
	JWTSalt string
}

func NewAdapter(JWTSalt string) *Adapter {
	return &Adapter{
		JWTSalt,
	}
}
