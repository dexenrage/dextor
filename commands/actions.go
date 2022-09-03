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
	"dextor/messages"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"time"
)

func checkRoot() (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	u, err := user.Current()
	catcherr.HandleError(err)

	if (u.Gid != `0`) && (u.Uid != `0`) {
		catcherr.HandleErrorString(messages.RunViaSudo)
	}
	return err
}

func requestToApi(ctx context.Context, api string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
}

func checkConnection(ctx context.Context) (data connectionData, err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	const api = `https://check.torproject.org/api/ip`
	resp, err := requestToApi(ctx, api)
	for err != nil {
		if ctx.Err() != nil {
			catcherr.HandleError(err)
		}
		resp, err = requestToApi(ctx, api)
		time.Sleep(1 * time.Second)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	catcherr.HandleError(err)

	err = json.Unmarshal(b, &data)
	catcherr.HandleError(err)
	return data, err
}

func promptInputNumber(isFirstPrompt bool, min, max int) (number int) {
	if !isFirstPrompt {
		fmt.Printf(messages.EnterNumber, min, max)
	}

	var input string
	fmt.Scanln(input)

	for len(input) == 0 {
		input = strconv.Itoa(promptInputNumber(false, min, max))
	}

	number, err := strconv.Atoi(input)
	for err != nil {
		input = strconv.Itoa(promptInputNumber(false, min, max))
		number, err = strconv.Atoi(input)
	}

	for number < min || number > max {
		number = promptInputNumber(false, min, max)
	}
	return number
}

func statFile(name string) error {
	_, err := os.Stat(name)
	return err
}

func saveBackupFile(oldpath, newpath string) error {
	err := statFile(newpath)
	if errors.Is(err, fs.ErrNotExist) {
		return os.Rename(oldpath, newpath)
	}
	return err
}

func restoreFile(backupPath, newPath string) error {
	err := removeFile(newPath)
	if err == nil {
		return os.Rename(backupPath, newPath)
	}
	return err
}

func writeFile(name, data string) error {
	return os.WriteFile(name, []byte(data), os.ModePerm)
}

func saveBackupAndWriteFile(oldpath, newpath, data string) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	err = saveBackupFile(oldpath, newpath)
	catcherr.HandleError(err)

	return writeFile(oldpath, data)
}

func truncateFile(name, data string) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	err = statFile(name)
	if err == nil {
		catcherr.HandleError(os.Truncate(name, 0))
	}
	if !errors.Is(err, fs.ErrNotExist) {
		catcherr.HandleError(err)
	}

	return writeFile(name, data)
}

func removeFile(name string) error {
	err := os.Remove(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}
