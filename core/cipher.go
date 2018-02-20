package core

import (
  "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
  "strings"
  "net/http"
  "encoding/base64"
)

// encrypt string to base64 crypto using AES
func Encriptar(llave string, texto string) (string, error, int) {
  key := []byte(llave)
  plaintext := []byte(texto)

  block, err := aes.NewCipher(key)
	if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil, http.StatusOK
}

func Desencriptar(llave string, cryptoText string) (string, error, int) {
  key := []byte(llave)

  ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR: ", err.Error()}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
	}

	if len(ciphertext) < aes.BlockSize {
    s := []string{"INTERNAL_SERVER_ERROR: ", "ciphertext es corto"}
    return "", fmt.Errorf(strings.Join(s, "")), http.StatusInternalServerError
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

  return string(ciphertext), nil, http.StatusOK
}
