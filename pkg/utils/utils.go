package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func PrettyJSONString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func ResolveArgs(args []string) map[string]string {
	m := make(map[string]string)
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) < 2 {
			panic(fmt.Sprintf("Invalid argument (format: <key>=<value>): %s", arg))
		}
		m[parts[0]] = parts[1]
	}
	return m
}

func IntFromStr(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func StrInArr(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func Exit(val int) {
	os.Exit(val)
}
