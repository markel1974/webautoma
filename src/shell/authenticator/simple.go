/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authenticator

import "github.com/markel1974/webautoma/src/shell/interfaces"

type SimpleAuthenticator struct {
	username string
	password string
}

func NewSimpleAuthenticator() *SimpleAuthenticator {
	a := &SimpleAuthenticator{}
	return a
}

func (a *SimpleAuthenticator) SetCredentials(username string, password string) {
	a.username = username
	a.password = password
}

func (a *SimpleAuthenticator) GetAuthenticationMode() interfaces.AuthMode {
	if len(a.username) > 0 {
		return interfaces.AuthModeFull
	}

	if len(a.password) > 0 {
		return interfaces.AuthModePassword
	}

	return interfaces.AuthModeNone
}

func (a *SimpleAuthenticator) IsAuthenticated(user string, password string) bool {
	switch a.GetAuthenticationMode() {
	case interfaces.AuthModeFull:
		return user == a.username && password == a.password

	case interfaces.AuthModePassword:
		return password == a.password

	case interfaces.AuthModeNone:
		return true

	default:
		return false
	}
}
