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

package commands

import (
	"context"
	"dextor/catcherr"
	"dextor/config"
	"dextor/directories"
	"dextor/messages"
	"dextor/syscmd"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 30*time.Second)
}

func execCommand(cmd string, args ...string) error {
	return exec.Command(cmd, args...).Run()
}

func Connect(ctx context.Context) {
	defer catcherr.Recover(`commands.Connect()`)

	fmt.Println(messages.StartConnection)

	err := checkRoot()
	catcherr.HandleError(err)

	ctx, cancel := defaultContextTimeout(ctx)
	defer cancel()

	conn, err := checkConnection(ctx)
	catcherr.HandleError(err)
	catcherr.HandleError(ctx.Err())

	if conn.IsTor {
		catcherr.HandleErrorString(messages.AlreadyConnected)
		return
	}

	err = saveBackupAndWriteFile(directories.ResolvConf, directories.ResolvBackup, config.LocalDNS)
	catcherr.HandleError(err)

	err = truncateFile(directories.TorRC, config.TorRC)
	catcherr.HandleError(err)

	err = saveBackupAndWriteFile(directories.Sysctl, directories.ResolvBackup, config.Sysctl)
	catcherr.HandleError(err)

	err = exec.Command(syscmd.SysctlReadValues).Run()
	catcherr.HandleError(err)

	user, err := syscmd.GetTorUser()
	catcherr.HandleError(err)

	err = execCommand(syscmd.ConfigureIPTables(user.Username))
	catcherr.HandleError(err)

	err = execCommand(syscmd.ConnectTor(user.Uid, directories.TorRC))
	catcherr.HandleError(err)

	conn, err = checkConnection(ctx)
	catcherr.HandleError(err)
	catcherr.HandleError(ctx.Err())

	if !conn.IsTor {
		RestoreConfigs()
		catcherr.HandleErrorString(messages.SomethingWentWrongConnect)
		return
	}
	fmt.Printf(messages.Connected, conn.IP)
}

func Disconnect(ctx context.Context) {
	defer catcherr.Recover(`commands.Disconnect()`)

	fmt.Println(messages.Disconnecting)

	err := checkRoot()
	catcherr.HandleError(err)

	ctx, cancel := defaultContextTimeout(ctx)
	defer cancel()

	conn, err := checkConnection(ctx)
	catcherr.HandleError(err)
	catcherr.HandleError(ctx.Err())

	if !conn.IsTor {
		catcherr.HandleErrorString(messages.AlreadyDisconnected)
		return
	}

	RestoreConfigs()

	conn, err = checkConnection(ctx)
	catcherr.HandleError(err)
	catcherr.HandleError(ctx.Err())

	if conn.IsTor {
		catcherr.HandleErrorString(messages.SomethingWrongDisconnect)
		return
	}
	fmt.Printf(messages.Disconnected, conn.IP)
}

func Reconnect(ctx context.Context) {
	Disconnect(ctx)
	time.Sleep(1500 * time.Millisecond)
	Connect(ctx)
}

func RestoreConfigs() (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	err = checkRoot()
	catcherr.HandleError(err)

	err = execCommand(syscmd.StopTor)
	catcherr.HandleError(err)

	err = execCommand(syscmd.Fuser)
	catcherr.HandleError(err)

	err = execCommand(syscmd.ResetIPTables)
	catcherr.HandleError(err)

	err = restoreFile(directories.ResolvBackup, directories.ResolvConf)
	catcherr.HandleError(err)

	err = restoreFile(directories.SysctlBackup, directories.Sysctl)
	catcherr.HandleError(err)

	err = execCommand(syscmd.SysctlReadValues)
	catcherr.HandleError(err)

	err = removeFile(directories.TorRC)
	catcherr.HandleError(err)

	err = execCommand(syscmd.RestartNetworkManager)
	catcherr.HandleError(err)

	return err
}

func FixDNS() {
	defer catcherr.Recover(`commands.FixDNS()`)

	err := checkRoot()
	catcherr.HandleError(err)

	fmt.Println(messages.ChooseDNS)
	id := promptInputNumber(true, 1, 5)

	err = removeFile(directories.ResolvConf)
	catcherr.HandleError(err)

	var dnsSet string
	switch id {
	case 1:
		err = writeFile(directories.ResolvConf, config.LocalDNS)
		dnsSet = messages.LocalDNSSet
	case 2:
		err = writeFile(directories.ResolvConf, config.CloudflareDNS)
		dnsSet = messages.CloudflareDNSSet
	case 3:
		err = writeFile(directories.ResolvConf, config.OpenDNS)
		dnsSet = messages.OpenDNSSet
	case 4:
		err = writeFile(directories.ResolvConf, config.GoogleDNS)
		dnsSet = messages.GoogleDNSSet
	case 5:
		err = writeFile(directories.ResolvConf, config.Quad9DNS)
		dnsSet = messages.Quad9DNSSet
	}
	catcherr.HandleError(err)

	fmt.Printf(dnsSet, directories.ResolvConf)
}

func ShowIP(ctx context.Context) {
	defer catcherr.Recover(`commands.ShowIP()`)

	ctx, cancel := defaultContextTimeout(ctx)
	defer cancel()

	conn, err := checkConnection(ctx)
	catcherr.HandleError(err)
	catcherr.HandleError(ctx.Err())

	fmt.Printf(messages.CurrentIP, conn.IP)
	if conn.IsTor {
		fmt.Println(messages.AndConnected)
		return
	}
	fmt.Println(messages.AndNotConnected)
}

func ShowVersion() {
	defer catcherr.Recover(`commands.ShowVersion()`)

	path, err := os.Executable()
	catcherr.HandleError(err)

	path = filepath.Join(filepath.Dir(path), `version`)

	b, err := os.ReadFile(path)
	catcherr.HandleError(err)

	version := strings.TrimSpace(string(b))
	fmt.Printf(messages.Version, version)
}
