package main

import (
	"fmt"

	"github.com/gonuts/commander"
)

const Version = "0.0.1"

func version(cmd *commander.Command, args []string) error {
	fmt.Println(Version)
	return nil
}

func init() {
	fs := newFlagSet("tcp")

	cmd.Subcommands = append(cmd.Subcommands, &commander.Command{
		UsageLine: "version",
		Short:     "display version information",
		Flag:      *fs,
		Run:       version,
	})
}
