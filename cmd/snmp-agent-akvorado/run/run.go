package run

import (
	"github.com/cuteip/snmp-agent-akvorado/mibimpls/network_interface"
	"github.com/cuteip/snmp-agent-akvorado/mibimpls/system"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/slayercat/GoSNMPServer"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "run",
		Short: "Run SNMP Agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	rootCmd.Flags().String("listen-addr", "127.0.0.1:161", "Listen Address")
	rootCmd.Flags().String("community", "public", "Community")
	rootCmd.Flags().String("log-level", "warn", "Log Level (panic, fatal, error, warn, info, debug, trace)")

	return rootCmd
}

func run(cmd *cobra.Command, args []string) error {
	community, err := cmd.Flags().GetString("community")
	if err != nil {
		return errors.Wrap(err, "failed to get flag 'community'")
	}

	oidNetworkInterface, err := network_interface.OIDs()
	if err != nil {
		return errors.Wrap(err, "failed to get network interface oids")
	}

	oidSystem, err := system.OIDs()
	if err != nil {
		return errors.Wrap(err, "failed to get system oids")
	}

	oids := oidNetworkInterface
	oids = append(oids, oidSystem...)

	logLevelFlag, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return errors.Wrap(err, "failed to get flag 'log-level'")
	}

	logLevel, err := logrus.ParseLevel(logLevelFlag)
	if err != nil {
		return errors.Wrap(err, "failed to parse log level")
	}

	logger := logrus.New()
	logger.SetLevel(logLevel)

	master := GoSNMPServer.MasterAgent{
		// memo: 標準ではない別の logger を使ってしまっている。できれば統一したい
		Logger: logger,
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				CommunityIDs: []string{community},
				OIDs:         oids,
			},
		},
	}

	server := GoSNMPServer.NewSNMPServer(master)

	listenAddr, err := cmd.Flags().GetString("listen-addr")
	if err != nil {
		return errors.Wrap(err, "failed to get flag 'listen-addr")
	}

	err = server.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}

	return server.ServeForever()
}
