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

package catcherr

import (
	"errors"
	"fmt"
)

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleErrorString(str string) { HandleError(errors.New(str)) }

func Recover(sender string) {
	if r := recover(); r != nil {
		LogError(sender, r)
	}
}

func LogError(sender string, a ...any) {
	const template = `[ %s ]: `

	_, err := fmt.Printf(template, a...)
	if err != nil {
		panic(err)
	}
}

func RecoverAndReturnError() error {
	if r := recover(); r != nil {
		return errors.New(fmt.Sprint(r))
	}
	return nil
}
