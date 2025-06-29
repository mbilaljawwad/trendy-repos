package utils

import (
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func ConvertToUTF8(str string) (string, error) {
	reader := transform.NewReader(strings.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	convertedStr, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(convertedStr), nil
}
