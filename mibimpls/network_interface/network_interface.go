// GoSNMPServer をベースに、Akvorado が取得する OID に関するデータのみを返すように変更
// https://github.com/slayercat/GoSNMPServer/blob/50d4684aabd065e8a94c17bd9e1d9a766f77bff7/mibImps/ifMib/networks.go

package network_interface

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/net"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

func OIDs() ([]*GoSNMPServer.PDUValueControlItem, error) {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	valInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "network ifs read failed")
	}

	for _, val := range valInterfaces {
		ifName := val.Name
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.1.%d", val.Index),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(val.Index), nil },
				Document: "ifIndex",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.2.%d", val.Index),
				Type:     gosnmp.OctetString,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(ifName), nil },
				Document: "ifDescr",
			},
			{
				// ifDescr が存在すれば十分んだと思われるが、Akvorado は ifAlias を取得しに来るので空として返す
				OID:      fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.18.%d", val.Index),
				Type:     gosnmp.OctetString,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
				Document: "ifAlias",
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.15.%d", val.Index),
				Type: gosnmp.Gauge32,
				// 単位は Mbps https://oidref.com/1.3.6.1.2.1.31.1.1.1.15
				// TODO: https://github.com/safchain/ethtool を利用して Speed を取得できるようにする
				// https://www.freedesktop.org/software/systemd/man/systemd.link.html#BitsPerSecond=
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(1000), nil },
				// https://github.com/akvorado/akvorado/blob/v1.8.1/inlet/snmp/poller.go#L174 には ifSpeed と書かれているが、おそらく ifHighSpeed
				Document: "ifHighSpeed",
			},
		}
		toRet = append(toRet, currentIf...)
	}
	return toRet, nil
}
