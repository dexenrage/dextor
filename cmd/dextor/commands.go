/*
Copyright 2022 dexenrage

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import "fmt"

const (
	cmdStopTor   = `/bin/bash -c "systemctl stop tor"`
	cmdFuser     = `/bin/bash -c "fuser -k 5353/tcp 9040/tcp 9050/tcp 9051/tcp > /dev/null 2>&1"`
	cmdRestartNM = `/bin/bash -c "systemctl restart --now NetworkManager"`

	cmdSysctlReadValues = `/bin/bash -c "sysctl -p"`

	cmdResetIPTables = `/bin/bash -c "iptables -P INPUT ACCEPT
	iptables -P FORWARD ACCEPT
	iptables -P OUTPUT ACCEPT
	iptables -t nat -F
	iptables -t mangle -F
	iptables -F
	iptables -X"`
)

func cmdIPTablesConfig() string {
	_, userID := getTorUser()
	cfg := `/bin/bash -c "NON_TOR="192.168.0.0/24 192.168.1.0/24 192.168.31.0/24"
	TOR_UID=%v
	TRANS_PORT=9040
	iptables -F
	iptables -t nat -F
	iptables -t nat -A OUTPUT -m owner --uid-owner $TOR_UID -j RETURN
	iptables -t nat -A OUTPUT -p udp --dport 53 -j REDIRECT --to-ports 5353
	for NET in $NON_TOR 127.0.0.0/9 127.128.0.0/10; do
	 iptables -t nat -A OUTPUT -d $NET -j RETURN
	done
	iptables -t nat -A OUTPUT -p tcp --syn -j REDIRECT --to-ports $TRANS_PORT
	iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
	for NET in $NON_TOR 127.0.0.0/8; do
	 iptables -A OUTPUT -d $NET -j ACCEPT
	done
	iptables -A OUTPUT -m owner --uid-owner $TOR_UID -j ACCEPT
	iptables -A OUTPUT -j REJECT"`
	return fmt.Sprintf(cfg, userID)
}

func cmdConnectTor() string {
	user, _ := getTorUser()
	return fmt.Sprintf(`sudo -u %v tor -f %v`, user, torrcDIR)
}
