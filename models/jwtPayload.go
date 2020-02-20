package models

type JwtPayload struct {
	ProfileKey string
}

type JwtValidate struct {
	JwtCiphered string
}
