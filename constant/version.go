// Copyright 2017 The Ethermis Authors
// This file is part of Ethermis.
//
// Ethermis is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Ethermis is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Ethermis. If not, see <http://www.gnu.org/licenses/>.

package constant

import "fmt"

var (
	ClientIdentifier = "Ethermis" // Client identifier to advertise over the network
	VersionMajor     = 1          // Major version component of the current release
	VersionMinor     = 0          // Minor version component of the current release
	VersionPatch     = 0          // Patch version component of the current release
	VersionMeta      = "unstable" // Version metadata to append to the version string
	GitCommit        string
	BuildTime        string
)

func VersionString() string {
	if GitCommit != "" {
		return fmt.Sprintf("%v.%v.%v-%v-%v", VersionMajor, VersionMinor, VersionPatch, VersionMeta, GitCommit[:8])
	}

	return fmt.Sprintf("%v.%v.%v-%v", VersionMajor, VersionMinor, VersionPatch, VersionMeta)
}
