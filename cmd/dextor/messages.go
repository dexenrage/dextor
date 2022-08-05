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

const (
	viaSudo   = "Please run program via `sudo`"
	errGetUID = "Error getting user ID.\nReinstall `tor` package and try again."

	msgFlagHelp    = "Show this message."
	msgFlagConn    = "Connect to the Tor network."
	msgFlagReconn  = "Reconnect to the Tor netork (changes IP address)."
	msgFlagDisconn = "Disconnect from the Tor network."
	msgFlagIP      = "Show your current IP address."
	msgFlagDNS     = "Use this if the website address can't be resolved."
	msgFlagCFG     = "Use this for restore configs if something went wrong."
	msgFlagUPD     = "Check the program for updates."
	msgFlagVer     = "Show the current version of the program."

	fixCFGWarn = `
	You will be disconnected from the Tor network (if connected).
	Do you want to continue? (Y/N): `
	restCFGs   = "\nRestoring configs..."
	restored   = "\nConfigs restored."
	restCancel = "\nConfigs restoring was cancelled."

	gettingInfo     = "Getting information..."
	currentIP_MSG   = "Your current IP: %v\n"
	andConn         = "And you are CONNECTED to the Tor network."
	andNotConn      = "And you are NOT CONNECTED to the Tor network."
	errAttemptLimit = "Maximum number of connection attempts exceeded."

	enterNum = `
	Please, enter a number between %d and %d

Your choise: `
	chooseDNS = `
Choose DNS server:
	1 - Local DNS, 2 - CloudFlare, 3 - OpenDNS, 4 - Google, 5 - Quad9

Your choise: `
	localDNSSet      = "Local DNS server (127.0.0.1) was set in %v\n"
	cloudflareDNSSet = "Cloudflare DNS server was set in %v\n"
	openDNSSet       = "OpenDNS server was set in %v\n"
	googleDNSSet     = "Google DNS server was set in %v\n"
	quad9DNSSet      = "Quad9 DNS server was set in %v\n"

	startConn    = "Starting new Tor connection... Please, wait."
	alreadyConn  = "You are already connected to the Tor network!\nYour current IP: %v\n"
	sthWrongConn = "Something went wrong. You haven't been connected to the Tor network."
	connectedMSG = "Connected. Your current IP: %v\n"

	stopConn        = "Disconnecting from the Tor Network... Please, wait."
	alreadyDisconn  = "You are already not connected to Tor!\nYour current IP: %v\n"
	sthWrongDisconn = "Something went wrong. You haven't been disconnected from the Tor network."
	disconnMSG      = "Disconnected. Your current IP: %v\n"

	versionMSG = "Current vesion: %v\n"
)
