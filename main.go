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

import (
	"context"
	"dextor/catcherr"
	"dextor/commands"
	"dextor/messages"
	"fmt"

	"github.com/spf13/pflag"
)

func main() {
	defer catcherr.Recover(`main()`)
	ctx := context.Background()

	var (
		helpFlag       bool
		connectFlag    bool
		reconnectFlag  bool
		disconnectFlag bool
		showIPFlag     bool
		fixDNSFlag     bool
		fixCFGFlag     bool
		versionFlag    bool
	)

	pflag.BoolVarP(&helpFlag, "help", "h", false, messages.HelpFlag)
	pflag.BoolVarP(&connectFlag, "connect", "c", false, messages.ConnectFlag)
	pflag.BoolVarP(&reconnectFlag, "reconnect", "r", false, messages.ReconnectFlag)
	pflag.BoolVarP(&disconnectFlag, "disconnect", "d", false, messages.DisconnectFlag)
	pflag.BoolVarP(&showIPFlag, "showip", "i", false, messages.ShowIPFlag)
	pflag.BoolVarP(&fixDNSFlag, "fixdns", "s", false, messages.FixDNSFlag)
	pflag.BoolVarP(&fixCFGFlag, "fixcfg", "g", false, messages.FixCFGFlag)
	pflag.BoolVarP(&versionFlag, "version", "v", false, messages.VersionFlag)
	pflag.Parse()

	if len(pflag.Args()) > 1 {
		catcherr.HandleErrorString(messages.ErrOnlyOneArgument)
	}

	switch {
	case connectFlag:
		commands.Connect(ctx)

	case disconnectFlag:
		commands.Disconnect(ctx)

	case reconnectFlag:
		commands.Reconnect(ctx)

	case showIPFlag:
		commands.ShowIP(ctx)

	case fixDNSFlag:
		commands.FixDNS()

	case fixCFGFlag:
		fmt.Print(messages.FixCFGWarn)

		var c string
		fmt.Scan(&c)

		switch c {
		case "Y", "y":
			catcherr.HandleError(commands.RestoreConfigs())
			fmt.Println(messages.ConfigsRestored)
		default:
			fmt.Println(messages.ConfigsRestoringCancelled)
		}

	case versionFlag:
		commands.ShowVersion()

	case helpFlag:
		pflag.Usage()

	default:
		pflag.Usage()
	}
}
