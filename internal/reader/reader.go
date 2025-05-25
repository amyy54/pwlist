package reader

import (
	"os"
	"strings"
)

func ReadFiles(paths []string) ([]string, error) {
	var res []string
	for _, path := range paths {
		file_out, err := os.ReadFile(path)
		if err != nil {
			return []string{}, err
		}
		res = append(res, strings.TrimSpace(string(file_out)))
	}
	return res, nil
}
