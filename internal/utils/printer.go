package utils

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func PrintInfo(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	fmt.Printf("\033[0;36m[*]\033[0m %s\n", msg)
}

func PrintSuccess(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	fmt.Printf("\033[0;32m[+]\033[0m %s\n", msg)
}

func PrintError(msg string, err error) {
	if GlobalDebugFlag {
		log.Error().Err(err).Msg(msg)
		return
	}
	if err != nil {
		fmt.Printf("\033[0;31m[-]\033[0m %s: %v\n", msg, err)
	} else {
		fmt.Printf("\033[0;31m[-]\033[0m %s\n", msg)
	}
}

func PrintFatal(msg string, err error) {
	PrintError(msg, err)
	os.Exit(1)
}

func PrintWarn(msg string, err error) {
	if GlobalDebugFlag {
		log.Warn().Err(err).Msg(msg)
		return
	}
	if err != nil {
		fmt.Printf("\033[0;33m[!]\033[0m %s: %v\n", msg, err)
	} else {
		fmt.Printf("\033[0;33m[!]\033[0m %s\n", msg)
	}
}

func PrintGeneric(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	fmt.Println(msg)
}
