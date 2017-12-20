package core

import (
  "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
  "strings"
  "net/http"
)

func Encriptar(llave string, texto string) (string, error, int) {
  key := []byte(llave)
  plaintext := []byte(texto)

  block, err := aes.NewCipher(key)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
    return "", fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError
  }

  nonce := make([]byte, 12)
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
    return "", fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError
  }

  aesgcm, err := cipher.NewGCM(block)
  if err != nil {
    s := []string{"INTERNAL_SERVER_ERROR:", err.Error()}
    return "", fmt.Errorf(strings.Join(s, " ")), http.StatusInternalServerError
  }

  ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
  return string(ciphertext), nil, http.StatusOK
}
