package main

import (
	"fmt"
	"os"

	"github.com/cuteip/snmp-agent-akvorado/cmd/snmp-agent-akvorado/run"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "snmp-agent",
		Short:   "SNMP Agent for Akvorado",
		Version: version,
	}

	rootCmd.SetVersionTemplate(versionTemplate())

	rootCmd.AddCommand(run.RootCmd())

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+V", err)
		os.Exit(1)
	}
}
