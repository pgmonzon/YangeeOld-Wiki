package core

import (
  "crypto/sha512"
  "encoding/binary"
	"bytes"
)

func HashSha512(clave string) (int64) {
	var claveInt64 int64

	h512 := sha512.New()
	h512.Write([]byte(clave))

	buf := bytes.NewBuffer(h512.Sum(nil))
	binary.Read(buf, binary.LittleEndian, &claveInt64)

	return claveInt64
}
