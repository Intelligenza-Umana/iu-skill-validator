package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type checkOutput struct {
	Path     string   `json:"path"`
	Valid    bool     `json:"valid"`
	Findings []string `json:"findings"`
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "check" {
		fmt.Fprintln(os.Stderr, "Usage: iu-skill-validator check [--json] <skillDir>")
		os.Exit(1)
	}

	flags := flag.NewFlagSet("check", flag.ExitOnError)
	jsonOutput := flags.Bool("json", false, "output JSON")
	_ = flags.Parse(os.Args[2:])

	args := flags.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: iu-skill-validator check [--json] <skillDir>")
		os.Exit(1)
	}

	skillDir := args[0]
	valid, findings := ValidateSkillDir(skillDir)

	if *jsonOutput {
		output := checkOutput{
			Path:     skillDir,
			Valid:    valid,
			Findings: findings,
		}
		if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode JSON output: %v\n", err)
			os.Exit(1)
		}
	} else {
		for _, finding := range findings {
			fmt.Println(finding)
		}
		if valid {
			fmt.Println("PASS (static checks only)")
		} else {
			fmt.Println("FAIL (static checks only)")
		}
	}

	if !valid {
		os.Exit(1)
	}
}
