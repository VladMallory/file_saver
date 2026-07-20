package flags

import (
	"flag"
	"log"
	"os"
)

// RunFlags содержит значения флагов команды "run".
type RunFlags struct{}

func ParseRun() RunFlags {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	if err := runCmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Ошибка парсинга флагов: %v", err)
	}

	return RunFlags{}
}
