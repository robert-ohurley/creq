package args

import (
	"log"
	"net/url"
	"strings"
	"unicode"
)

type UserFlag struct {
	F         string
	Parameter string
}

type flag struct {
	Flag          string
	Short         string
	ParamRequired bool
	DefaultValue  any
	Description   string
	ReplOnly      bool
}

var ProgramFlags []flag

func init() {
	setProgramFlags()
}

func ParseArgs(cmdLineArgs []string) []UserFlag {
	var Flags []UserFlag

	for i := 1; i < len(cmdLineArgs); i++ {
		var f UserFlag
		flagExist, parsedFlag := isFlag(cmdLineArgs[i])
		shortExist, parsedShort := isShort(cmdLineArgs[i])

		//check for optional first URL
		if !flagExist && !shortExist && i == 1 {
			_, err := url.ParseRequestURI(cmdLineArgs[i])

			if err != nil {
				log.Fatalln("Given URL is invalid. Example: http://github.com")
			}
			Flags = append(Flags, UserFlag{F: "url", Parameter: cmdLineArgs[i]})
		}

		if shortExist {
			f.F = parsedShort
			Flags = append(Flags, f)
			continue
		} else if flagExist {
			f.F = parsedFlag
			Flags = append(Flags, f)
			continue
		} else if len(Flags) >= 1 {
			//todo: check that last added flag actually requires a parameter otherwise return an error
			Flags[len(Flags)-1].Parameter = strings.Trim(cmdLineArgs[i], " ")
		}
	}

	return Flags
}

func isFlag(s string) (bool, string) {
	if len(s) > 3 && strings.HasPrefix(s, "--") {
		for i := 2; i < len(s); i++ {
			if !unicode.IsLetter(rune(s[i])) && s[i] != byte('-') {
				return false, ""
			}
		}
		return true, strings.Trim(strings.ToLower(s[2:]), " ")
	} else {
		return false, ""
	}
}

func isShort(s string) (bool, string) {
	if len(s) == 2 && strings.HasPrefix(s, "-") && unicode.IsLetter(rune(s[1])) {
		return true, string(s[1])
	} else {
		return false, ""
	}
}
