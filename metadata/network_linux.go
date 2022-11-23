//go:build linux

/*
 * Copyright (c) 2021 Michael Morris. All Rights Reserved.
 *
 * Licensed under the MIT license (the "License"). You may not use this file except in compliance
 * with the License. A copy of the License is located at
 *
 * https://github.com/mmmorris1975/aws-runas/blob/master/LICENSE
 *
 * or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License
 * for the specific language governing permissions and limitations under the License.
 */

package metadata

import "net"

func addAddress(iface *net.Interface, cidrAddr string) error {
	cmd := []string{"ip", "address", "add", cidrAddr, "dev", iface.Name}
	return doCommand(cmd)
}

func removeAddress() error {
	mu.Lock()
	defer mu.Unlock()
	iface, err := findInterfaceByAddress(DefaultEc2ImdsAddr)
	if err != nil {
		return err
	}

	cmd := []string{"ip", "address", "del", DefaultEc2ImdsCidr, "dev", iface.Name}
	return doCommand(cmd)
}
