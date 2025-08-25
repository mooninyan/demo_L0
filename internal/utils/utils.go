package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

type MyLogger struct{}

func (l MyLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func GetDlqFlagOrDefault() bool {
	var isDlq, propertyError = strconv.ParseBool(os.Getenv("LISTEN_DLQ"))
	if propertyError != nil {
		isDlq = false
	}
	return isDlq
}

func AtoiMust(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func MapValues[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func WrapAndLog(err error) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "Unknown"
		line = 0
	}
	err = fmt.Errorf("%s:%d Error:%v", file, line, err.Error())
	log.Println(err)
	return err
}
