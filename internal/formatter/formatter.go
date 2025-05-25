package formatter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/QMUL/ntlmgen"
)

type HashType int

const (
	NTLM HashType = iota
	MD5
	SHA1
	SHA256
	SHA512
)

const FormatMatcher = `\|(\w+)(?:-([a-zA-Z]+))?\|`

func GenMatchList(format_str string, file_contents []string) ([]RegReplace, [][]string) {
	re := regexp.MustCompile(FormatMatcher)
	matches := re.FindAllStringSubmatch(format_str, -1)

	var match_list []RegReplace
	var iter_args [][]string

	for _, match := range matches {
		full := match[0]
		index := match[1]
		mod := match[2]
		slog.Info("Checking regex:", "full", full, "index", index, "mod", mod)

		index_int, index_err := strconv.Atoi(index)
		if index_err == nil { // It's a number
			if len(file_contents) > index_int {
				file_content := file_contents[index_int]
				if mod != "" {
					switch mod {
					case "lower":
						file_content = strings.ToLower(file_content)
					case "upper":
						file_content = strings.ToUpper(file_content)
					case "title":
						file_content = strings.ToTitle(file_content)
					}
				}
				match_list = append(match_list, RegReplace{FormatBlock: full, Items: strings.Split(file_content, "\n")})
				iter_args = append(iter_args, strings.Split(file_content, "\n"))
			} else {
				slog.Info("File index was referenced, but length of file contents didn't match", "index_int", index_int, "len(file_contents)", len(file_contents))
				log.Fatal("A file index was referenced that wasn't passed!")
				os.Exit(1)
			}
		} else { // It's not a number
			var chosen_charset string
			switch index {
			case "l":
				chosen_charset = "a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z"
			case "d":
				chosen_charset = "0,1,2,3,4,5,6,7,8,9"
			case "h":
				chosen_charset = "0,1,2,3,4,5,6,7,8,9,a,b,c,d,e,f"
			default:
				chosen_charset = ""
			}
			if mod == "upper" {
				chosen_charset = strings.ToUpper(chosen_charset)
			}

			match_list = append(match_list, RegReplace{FormatBlock: full, Items: strings.Split(chosen_charset, ",")})
			iter_args = append(iter_args, strings.Split(chosen_charset, ","))
		}
	}
	return match_list, iter_args
}

func HashGen(hashType HashType, text string) string {
	switch hashType {
	case NTLM:
		return ntlmgen.Ntlmgen(text)
	case MD5:
		md5 := md5.Sum([]byte(text))
		return fmt.Sprintf("%x", md5)
	case SHA1:
		sha1 := sha1.Sum([]byte(text))
		return fmt.Sprintf("%x", sha1)
	case SHA256:
		sha256 := sha256.New().Sum([]byte(text))
		return fmt.Sprintf("%x", sha256)
	case SHA512:
		sha512 := sha512.New().Sum([]byte(text))
		return fmt.Sprintf("%x", sha512)
	}
	return ""
}
