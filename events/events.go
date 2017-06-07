package events

import (
    "log"
    "os"
)

var (
    logger *log.Logger
)

func init() {
    logger = log.New(os.Stdout, "[events] ", 0)
}
