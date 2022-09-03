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

package messages

const (
	RunViaSudo         = "Please run program via `sudo`"
	ErrOnlyOneArgument = `Only one argument is allowed`

	HelpFlag       = "Show this message."
	ConnectFlag    = "Connect to the Tor network."
	ReconnectFlag  = "Reconnect to the Tor netork (changes IP address)."
	DisconnectFlag = "Disconnect from the Tor network."
	ShowIPFlag     = "Show your current IP address."
	FixDNSFlag     = "Use this if the website address can't be resolved."
	FixCFGFlag     = "Use this for restore configs if something went wrong."
	VersionFlag    = "Show the current version of the program."

	ConfigsRestored           = "Configs restored."
	ConfigsRestoringCancelled = "Configs restoring was cancelled."

	CurrentIP       = "Your current IP: %v"
	AndConnected    = "And you are CONNECTED to the Tor network."
	AndNotConnected = "And you are NOT CONNECTED to the Tor network."

	LocalDNSSet      = "Local DNS server (127.0.0.1) was set in %v"
	CloudflareDNSSet = "Cloudflare DNS server was set in %v"
	OpenDNSSet       = "OpenDNS server was set in %v"
	GoogleDNSSet     = "Google DNS server was set in %v"
	Quad9DNSSet      = "Quad9 DNS server was set in %v"

	StartConnection           = "Starting new Tor connection... Please, wait."
	AlreadyConnected          = "You are already connected to the Tor network!"
	SomethingWentWrongConnect = "Something went wrong. You haven't been connected to the Tor network."
	Connected                 = "Connected. Your current IP: %v"

	Disconnecting            = "Disconnecting from the Tor Network... Please, wait."
	AlreadyDisconnected      = "You are already not connected to Tor!"
	SomethingWrongDisconnect = "Something went wrong. You haven't been disconnected from the Tor network."
	Disconnected             = "Disconnected. Your current IP: %v"

	Version = "Current vesion: %v"
)

const (
	EnterNumber = `
	Please, enter a number between %d and %d

Your choise: `

	ChooseDNS = `
Choose DNS server:
	1 - Local DNS, 2 - CloudFlare, 3 - OpenDNS, 4 - Google, 5 - Quad9

Your choise: `

	FixCFGWarn = `
	You will be disconnected from the Tor network (if connected).
	Do you want to continue? (Y/N): `
)
