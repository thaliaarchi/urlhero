// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tinytown

import (
	"fmt"

	trpc "github.com/hekmon/transmissionrpc"
)

func DownloadTransmission(c *trpc.Client, dir string) error {
	if err := checkVersion(c); err != nil {
		return err
	}
	ids, err := GetReleaseIDs()
	if err != nil {
		return err
	}
	for i, id := range ids {
		fmt.Printf("(%d/%d) Adding %s\n", i+1, len(ids), id)
		filename, err := saveTorrentFile(id, dir)
		if err != nil {
			return err
		}
		if _, err := c.TorrentAddFileDownloadDir(filename, dir); err != nil {
			return err
		}
	}
	return nil
}

func checkVersion(c *trpc.Client) error {
	ok, version, minVersion, err := c.RPCVersion()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("RPC remote v%d is incompatible with library v%d: remote needs at least v%d",
			version, trpc.RPCVersion, minVersion)
	}
	return nil
}
