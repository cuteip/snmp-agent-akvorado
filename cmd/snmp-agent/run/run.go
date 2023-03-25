package run

import (
	"github.com/cuteip/snmp-agent-akvorado/mibimpls/network_interface"
	"github.com/cuteip/snmp-agent-akvorado/mibimpls/system"

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

	master := GoSNMPServer.MasterAgent{
		// memo: 標準ではない別の logger を使ってしまっている。できれば統一したい
		Logger: GoSNMPServer.NewDefaultLogger(),
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
