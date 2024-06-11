package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var Log = zerolog.New(zerolog.ConsoleWriter{})

func SetupLogger(logger zerolog.Logger) error {
	var level string
	var filePath string

	level = "debug"
	fmt.Println(level)
	filePath = "F:\\Dean-ai\\users\\logger\\logger.log"
	if filePath == "" {
		return errors.New("logger file path not found")
	}
	basePath := filepath.Dir(filePath)
	created, err := CheckpathExists(basePath, 1)
	if !created {
		return errors.New(err)
	}

	file, err1 := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err1 != nil {
		return err1
	}

	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Log = zerolog.New(zerolog.ConsoleWriter{Out: file, NoColor: false, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	return nil

}

func CheckpathExists(path string, value int) (bool, string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if value == 1 {
				if err := os.Mkdir(path, 0700); err != nil {
					errmsg := fmt.Sprintf("Err: Path to backup \"%s\". %s", path, err.Error())
					return false, errmsg
				} else {
					return true, ""
				}
			} else {
				errmsg := fmt.Sprintf("Err: Path to backup not \"%s\" not exists", path)
				return false, errmsg
			}
		} else {
			errmsg := fmt.Sprintf("Err: Path to backup not \"%s\" not exists", path)
			return false, errmsg
		}
	} else {
		return true, ""
	}
}
