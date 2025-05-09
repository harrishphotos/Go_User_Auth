package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateVerificationToken() (string,error){

	b:= make([]byte, 32)
	if _,err:= rand.Read(b); err!=nil{
		return "", err
	}

	return hex.EncodeToString(b), nil
}