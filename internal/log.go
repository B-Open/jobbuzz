package internal

import (
	"os"

	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Logger {
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)

	// still figuring out how to add file logger
	// filename := "/var/log/jobbuzz/jobbuzz.log"

	// os.MkdirAll(filename, 0755)
	// f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// logger.SetOutput(f)

	return logger
}
