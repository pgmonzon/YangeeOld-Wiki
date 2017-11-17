package core

import (
  "log"
  "net/http"
  "time"
)

func RespuestaJSON(w http.ResponseWriter, req *http.Request, start time.Time, respuesta []byte, code int) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
  if string(respuesta) != "" {
		w.Write(respuesta)
	}

  log.Printf("%s\t%s\t%s\t%s\t%d\t%d\t%s",
		req.RemoteAddr,
		req.Method,
		req.RequestURI,
		req.Proto,
		code,
		len(respuesta),
		time.Since(start),
	)
}
