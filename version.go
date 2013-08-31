package main

import (
    "fmt"
    "github.com/gonuts/commander"
)

const Version = "0.0.1"

func version(cmd *commander.Command, args []string) {
    fmt.Println(Version)
}

func init() {
    fs := newFlagSet("tcp")

    cmd.Commands = append(cmd.Commands, &commander.Command{
        UsageLine: "version",
        Short:     "display version information",
        Flag:      *fs,
        Run:       version,
    })
}
