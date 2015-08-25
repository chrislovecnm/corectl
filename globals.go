// Copyright 2015 - António Meireles  <antonio.meireles@reformi.st>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/codegangsta/cli"
)

type sessionInfo struct {
	configDir string
	pwd       string
	uid       string
	gid       string
	hasPowers bool
	debug     bool
	json      bool
	data      VMInfo
}

// SessionContext - running session
var SessionContext sessionInfo

// VMInfo - individual VM settings
type VMInfo struct {
	Channel     string
	Version     string
	Cpus        string
	Memory      string
	UUID        string
	Xhyve       string
	CloudConfig string
	CClocation  int
	SSHkey      string
	Extra       string
	Network     networkAssets
	Storage     storageAssets
}

// networkDevice types
const (
	_ = iota
	Raw
	Tap
)

// storageDevice types
const (
	_ = iota
	HDD
	CDROM
)

//
type NetworkInterface struct {
	Slot int
	Type int
}

//
type StorageDevice struct {
	Slot int
	Type int
	Path string
}

type storageAssets struct {
	CDDrives   map[string]StorageDevice
	HardDrives map[string]StorageDevice
}

type networkAssets struct {
	Tap map[string]NetworkInterface
	Raw map[string]NetworkInterface
}

// CloudConfig types
const (
	_ = iota
	Local
	Remote
)

func (session *sessionInfo) init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	log.SetPrefix("[coreos] ")

	var caller *user.User

	if uid := os.Geteuid(); uid == 0 {
		if usr := os.Getenv("SUDO_USER"); usr == "" {
			log.Fatalln("This shouldn't be called straight as 'root' user.",
				"It should be run either as a regular user or via 'sudo'.")
		} else {
			caller, _ = user.Lookup(usr)
			session.hasPowers = true
		}
	} else {
		session.hasPowers = false
		caller, _ = user.Current()
	}
	session.uid, session.gid = caller.Uid, caller.Gid
	session.configDir = fmt.Sprintf("%s/.coreos/", caller.HomeDir)
	if pwd, err := os.Getwd(); err != nil {
		log.Fatalln(err)
	} else {
		session.pwd = pwd
	}
	for _, i := range DefaultChannels {
		d := fmt.Sprintf("%s/images/%s", session.configDir, i)
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalln("unable to create", d)
		}
	}
	rundir := fmt.Sprintf("%s/running", session.configDir)
	if err := os.MkdirAll(rundir, 0755); err != nil {
		log.Fatalln("unable to create", rundir)
	}
	if session.hasPowers {
		if err := fixPerms(session.configDir); err != nil {
			log.Fatalln(err)
		}
	}
}
func (session *sessionInfo) canRun() {
	if !session.hasPowers {
		log.Fatalln("not enough permissions to run VMs. use 'sudo'.")
	}
}
func (vm *VMInfo) setChannel(channel string) {
	var b string
	for _, b = range DefaultChannels {
		if b == channel {
			vm.Channel = b
			return
		}
	}
	log.Printf("'%s' is not a recognized CoreOS image channel. %s",
		channel, "Using default ('alpha').")
	vm.Channel = "alpha"
	return
}
func (vm *VMInfo) setVersion(version string) {
	vm.Version = version
}

func (vm *VMInfo) setDefaultNIC() {
	vm.Network.Raw = make(map[string]NetworkInterface, 0)
	vm.Network.Raw[strconv.Itoa(0)] = NetworkInterface{
		Type: Raw,
		Slot: 0,
	}
}

// CoreOS public release streams
var DefaultChannels = []string{"alpha", "beta", "stable"}

// NOTyetImplementedCommand ...
func NOTyetImplementedCommand(c *cli.Context) {
	fmt.Println("placeholder. NOT YET IMPLEMENTED...")
}

// GPGLongId
const GPGLongID = "50E0885593D2DCB4"

// GPGKey
const GPGKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v2
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
Pg==
=LY4G
-----END PGP PUBLIC KEY BLOCK-----
`

//
const CoreOEMsetupEnv = `#!/bin/bash
UUID=$(cat /proc/cmdline | sed -e 's,.*uuid=,,' -e 's, .*,,')
CALLERID=$(cat /proc/cmdline | sed -e 's,.*localuser=,,' -e 's, .*,,')
STATUSDIR=/Users/${CALLERID}/.coreos/running/${UUID}
# wait for eth0 to get up...
while [ 1 ]; do
  COREOS_PUBLIC_IPV4=$(/bin/ifconfig eth0 | awk '/inet /{print $2}')
  if [ -n "${COREOS_PUBLIC_IPV4}" ]; then
     break
  fi
  sleep 1; done
COREOS_PRIVATE_IPV4=${COREOS_PUBLIC_IPV4}
( echo UUID=${UUID};
  echo CALLERID=${CALLERID};
  echo STATUSDIR=${STATUSDIR};
  echo COREOS_PUBLIC_IPV4=${COREOS_PUBLIC_IPV4};
  echo COREOS_PRIVATE_IPV4=${COREOS_PRIVATE_IPV4};
) > /etc/environment
`

//
const CoreOEMsetup = `#cloud-config

write-files:
  - path: /etc/conf.d/nfs
    permissions: '0644'
    content: |
      OPTS_RPC_MOUNTD=""
coreos:
  units:
    - name: rpc-statd.service
      command: start
      enable: true
    - name: Users.mount
      command: start
      content: |
        [Mount]
        ConditionPathExists=/etc/environment
        What=192.168.64.1:/Users
        Where=/Users
        Options=rw,async,nolock,noatime,rsize=32768,wsize=32768
        Type=nfs
    - name: local-cloud-config.service
      command: start
      enable: true
      content: |
        [Unit]
          Description=Load cloud-config from file
          Requires=coreos-setup-environment.service
          After=coreos-setup-environment.service
          Before=user-config.target
          Requires=Users.mount
          After=Users.mount

        [Service]
          Type=oneshot
          RemainAfterExit=yes
          EnvironmentFile=/etc/environment
          ExecStart=/bin/bash -c "[[ -f $STATUSDIR/cloud-config.local ]] && \
						/usr/bin/coreos-cloudinit \
							-from-file $STATUSDIR/cloud-config.local || true"
    - name: xhyve.service
      command: start
      content: |
        [Unit]
        Description=updates xhyve context
        Requires=Users.mount
        After=Users.mount
        ConditionPathExists=/etc/environment

        [Service]
        Type=oneshot
        RemainAfterExit=true
        EnvironmentFile=/etc/environment
        ExecStart=/bin/bash -c "echo ${COREOS_PUBLIC_IPV4} > ${STATUSDIR}/ip"
`