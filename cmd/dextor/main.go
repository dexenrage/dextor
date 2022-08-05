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

// The package main manages the user's connection to Tor.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	parseFlags()
}

func checkRoot() {
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	if (u.Gid != `0`) && (u.Uid != `0`) {
		log.Fatalln(viaSudo)
	}
}

func parseFlags() {
	var (
		helpFlag    bool
		connFlag    bool
		reconnFlag  bool
		disconnFlag bool
		ipFlag      bool
		dnsFlag     bool
		cfgFlag     bool
		verFlag     bool
	)

	pflag.BoolVarP(&helpFlag, "help", "h", false, msgFlagHelp)
	pflag.BoolVarP(&connFlag, "connect", "c", false, msgFlagConn)
	pflag.BoolVarP(&reconnFlag, "reconnect", "r", false, msgFlagReconn)
	pflag.BoolVarP(&disconnFlag, "disconnect", "d", false, msgFlagDisconn)
	pflag.BoolVarP(&ipFlag, "showip", "i", false, msgFlagIP)
	pflag.BoolVarP(&dnsFlag, "fixdns", "s", false, msgFlagDNS)
	pflag.BoolVarP(&cfgFlag, "fixcfg", "g", false, msgFlagCFG)
	pflag.BoolVarP(&verFlag, "version", "v", false, msgFlagVer)
	pflag.Parse()

	switch {
	case helpFlag:
		pflag.Usage()

	case connFlag:
		checkRoot()
		connect()

	case reconnFlag:
		checkRoot()
		disconnect()
		connect()

	case disconnFlag:
		checkRoot()
		disconnect()

	case ipFlag:
		showIP()

	case dnsFlag:
		checkRoot()
		fixDNS()

	case cfgFlag:
		checkRoot()
		fmt.Print(fixCFGWarn)

		var c string
		switch fmt.Scan(&c); c {
		case "Y", "y":
			fmt.Println(restCFGs)
			restoreCFG()
			fmt.Println(restored)
		default:
			fmt.Println(restCancel)
		}

	case verFlag:
		getVersion()

	default:
		pflag.Usage()
	}
}

func checkConn() (isTor bool, ip string) {
	var data struct {
		IsTor bool   `json:"IsTor"`
		IP    string `json:"IP"`
	}

	var attempt int
	const maxAttempts = 50

	resp, err := http.Get(apiURL)
	for err != nil {
		time.Sleep(time.Second) // To avoid request flood.

		resp, err = http.Get(apiURL)
		if err == nil {
			break
		}

		attempt++
		if attempt > maxAttempts {
			if err != nil {
				log.Fatalln(err)
			}
			log.Fatalln(errAttemptLimit)
		}
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	json.Unmarshal(b, &data)

	return data.IsTor, data.IP
}

func showIP() {
	fmt.Println(gettingInfo)
	isTor, ip := checkConn()
	fmt.Printf(currentIP_MSG, ip)

	switch isTor {
	case true:
		fmt.Println(andConn)
	default:
		fmt.Println(andNotConn)
	}
}

func promptNum(isFirstPrompt bool, min, max int) (num int) {
	if !isFirstPrompt {
		fmt.Printf(enterNum, min, max)
	}

	var input string
	fmt.Scanln(&input)

	for !valid.IsInt(input) {
		num = promptNum(false, min, max)
		input = strconv.Itoa(num)
	}

	num, _ = strconv.Atoi(input)
	for num < min || num > max {
		num = promptNum(false, min, max)
	}
	return num
}

func writeFile(dir, data string) {
	err := os.WriteFile(dir, []byte(data), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func fixDNS() {
	fmt.Print(chooseDNS)

	min, max := 1, 5 // DNS Server IDs range
	num := promptNum(true, min, max)

	err := os.Remove(resolvDIR)
	if err != nil {
		log.Fatalln(err)
	}

	switch num {
	case 1:
		writeFile(resolvDIR, localDNS)
		fmt.Printf(localDNSSet, resolvDIR)

	case 2:
		writeFile(resolvDIR, cloudflareDNS)
		fmt.Printf(cloudflareDNSSet, resolvDIR)

	case 3:
		writeFile(resolvDIR, openDNS)
		fmt.Printf(openDNSSet, resolvDIR)

	case 4:
		writeFile(resolvDIR, googleDNS)
		fmt.Printf(googleDNSSet, resolvDIR)

	case 5:
		writeFile(resolvDIR, quad9DNS)
		fmt.Printf(quad9DNSSet, resolvDIR)
	}
}

func restoreCFG() {
	err := exec.Command(cmdStopTor).Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command(cmdFuser).Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command(cmdResetIPTables).Run()
	if err != nil {
		log.Fatalln(err)
	}

	switch _, err = os.Stat(resolvBAK); err {
	case nil:
		err = os.Rename(resolvBAK, resolvDIR)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
	}

	switch _, err = os.Stat(sysctlBAK); err {
	case nil:
		err = os.Rename(sysctlBAK, sysctlDIR)
		if err != nil {
			log.Fatalln(err)
		}
		exec.Command(cmdSysctlReadValues).Run()
	default:
		if os.IsNotExist(err) {
			err = os.Remove(sysctlDIR)
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
	}

	switch _, err = os.Stat(torrcDIR); err {
	case nil:
		err = os.Remove(torrcDIR)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

	}

	err = exec.Command(cmdRestartNM).Run()
	if err != nil {
		log.Fatalln(err)
	}
}

// Gets actual Tor username and id based on Linux distro.
func getTorUser() (name string, uid string) {
	u, err := user.Lookup(`debian-tor`)
	if err != nil {
		u, err = user.Lookup(`tor`)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return u.Username, u.Uid
}

func connect() {
	fmt.Println(startConn)

	isTor, ip := checkConn()
	if isTor {
		log.Fatalf(alreadyConn, ip)
	}

	switch _, err := os.Stat(resolvBAK); err {
	case nil:
		// Backup of resolv.conf already exists.
		err = os.Remove(resolvDIR)
		if err != nil {
			log.Fatalln(err)
		}
		writeFile(resolvDIR, localDNS)
	default:
		if os.IsNotExist(err) {
			// User have their own resolv.conf, which must be safe and sound.
			err = os.Rename(resolvDIR, resolvBAK)
			if err != nil {
				log.Fatalln(err)
			}
			writeFile(resolvDIR, localDNS)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
	}

	switch _, err := os.Stat(torrcDIR); err {
	case nil:
		err = os.Remove(torrcDIR)
		if err != nil {
			log.Fatalln(err)
		}
		writeFile(torrcDIR, rcCFG)
	default:
		if os.IsNotExist(err) {
			writeFile(torrcDIR, rcCFG)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
	}

	switch _, err := os.Stat(sysctlDIR); err {
	case nil:
		// User can have their own sysctl.conf, which must be safe and sound.
		err = os.Rename(sysctlDIR, sysctlBAK)
		if err != nil {
			log.Fatalln(err)
		}
		writeFile(sysctlDIR, disableIPv6)
	default:
		if os.IsNotExist(err) {
			writeFile(sysctlDIR, disableIPv6)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := exec.Command(cmdSysctlReadValues).Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command(cmdIPTablesConfig()).Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command(cmdConnectTor()).Run()
	if err != nil {
		log.Fatalln(err)
	}

	switch isTor, ip = checkConn(); isTor {
	case true:
		fmt.Printf(connectedMSG, ip)
	default:
		restoreCFG()
		log.Fatalln(sthWrongConn)
	}
}

func disconnect() {
	fmt.Println(stopConn)

	isTor, ip := checkConn()
	if !isTor {
		log.Fatalf(alreadyDisconn, ip)
	}

	restoreCFG()

	switch isTor, ip = checkConn(); isTor {
	case true:
		log.Fatalln(sthWrongDisconn)
	default:
		fmt.Printf(disconnMSG, ip)
	}
}

func getVersion() {
	path, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	path = filepath.Dir(path)
	versionFile := filepath.Join(path, `version`)

	verBytes, err := os.ReadFile(versionFile)
	if err != nil {
		log.Fatalln(err)
	}
	version := strings.TrimSpace(string(verBytes))

	fmt.Printf(versionMSG, version)
}
