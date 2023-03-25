package system

import (
	"os"

	"github.com/pkg/errors"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

func OIDs() ([]*GoSNMPServer.PDUValueControlItem, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}

	items := []*GoSNMPServer.PDUValueControlItem{
		{
			// 1.3.6.1.2.1.1.5 ではなく 1.3.6.1.2.1.1.5.0 にする
			// Advorado が .0 にアクセスするため
			// https://github.com/akvorado/akvorado/blob/v1.8.1/inlet/snmp/poller.go#L169
			OID:      "1.3.6.1.2.1.1.5.0",
			Type:     gosnmp.OctetString,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(hostname), nil },
			Document: "sysName",
		},
	}
	return items, nil
}
