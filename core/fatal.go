package core

import (
  "log"
)

func FatalErr(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
