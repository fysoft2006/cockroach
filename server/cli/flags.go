// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Daniel Theophanes (kardianos@gmail.com)

package cli

import (
	"github.com/cockroachdb/cockroach/server"

	"github.com/spf13/cobra"
)

// initFlags sets the server.Context values to flag values.
// Keep in sync with "server/context.go". Values in Context should be
// settable here.
func initFlags(ctx *server.Context) {
	for _, cmd := range nodeCmds {
		f := cmd.Flags()

		// Server flags.

		f.StringVar(&ctx.Addr, "addr", ctx.Addr, "the host:port to bind for HTTP/RPC traffic")

		f.StringVar(&ctx.Stores, "stores", ctx.Stores, "specify a comma-separated list of stores, "+
			"specified by a colon-separated list of device attributes followed by '=' and "+
			"either a filepath for a persistent store or an integer size in bytes for an "+
			"in-memory store. Device attributes typically include whether the store is "+
			"flash (ssd), spinny disk (hdd), fusion-io (fio), in-memory (mem); device "+
			"attributes might also include speeds and other specs (7200rpm, 200kiops, etc.). "+
			"For example, --stores=hdd:7200rpm=/mnt/hda1,ssd=/mnt/ssd01,ssd=/mnt/ssd02,mem=1073741824.")

		f.StringVar(&ctx.Attrs, "attrs", ctx.Attrs, "specify an ordered, colon-separated list of node "+
			"attributes. Attributes are arbitrary strings specifying topography or "+
			"machine capabilities. Topography might include datacenter designation "+
			"(e.g. \"us-west-1a\", \"us-west-1b\", \"us-east-1c\"). Machine capabilities "+
			"might include specialized hardware or number of cores (e.g. \"gpu\", "+
			"\"x16c\"). "+
			"The relative geographic proximity of two nodes is inferred from the "+
			"common prefix of the attributes list, so topographic attributes should be "+
			"specified first and in the same order for all nodes. "+
			"For example: --attrs=us-west-1b,gpu.")

		f.DurationVar(&ctx.MaxOffset, "max-offset", ctx.MaxOffset, "specify "+
			"the maximum clock offset for the cluster. Clock offset is measured on all "+
			"node-to-node links and if any node notices it has clock offset in excess "+
			"of --max-offset, it will commit suicide. Setting this value too high may "+
			"decrease transaction performance in the presence of contention.")

		f.DurationVar(&ctx.MetricsFrequency, "metrics-frequency", ctx.MetricsFrequency, "specify "+
			"--metrics-frequency to adjust the frequency at which the server records "+
			"its own internal metrics.")

		// Gossip flags.

		f.StringVar(&ctx.GossipBootstrap, "gossip", ctx.GossipBootstrap, "specify a "+
			"comma-separated list of gossip addresses or resolvers for gossip bootstrap. "+
			"Each item in the list has an optional type: [type=]<address>. "+
			"Unspecified type means ip address or dns. Type can also be a load balancer (\"lb\"), "+
			"a unix socket (\"unix\") or, for single-node systems, \"self\".")

		f.DurationVar(&ctx.GossipInterval, "gossip-interval", ctx.GossipInterval,
			"approximate interval (time.Duration) for gossiping new information to peers.")

		// KV flags.

		f.BoolVar(&ctx.Linearizable, "linearizable", ctx.Linearizable, "enables linearizable behaviour "+
			"of operations on this node by making sure that no commit timestamp is reported "+
			"back to the client until all other node clocks have necessarily passed it.")

		// Engine flags.

		f.Int64Var(&ctx.CacheSize, "cache-size", ctx.CacheSize, "total size in bytes for "+
			"caches, shared evenly if there are multiple storage devices.")

		f.DurationVar(&ctx.ScanInterval, "scan-interval", ctx.ScanInterval, "specify "+
			"--scan-interval to adjust the target for the duration of a single scan "+
			"through a store's ranges. The scan is slowed as necessary to approximately "+
			"achieve this duration.")
	}

	var clientCmds []*cobra.Command
	clientCmds = append(clientCmds, kvCmds...)
	clientCmds = append(clientCmds, rangeCmds...)
	clientCmds = append(clientCmds, acctCmds...)
	clientCmds = append(clientCmds, permCmds...)
	clientCmds = append(clientCmds, zoneCmds...)

	for _, cmd := range clientCmds {
		cmd.Flags().StringVar(&ctx.Addr, "addr", ctx.Addr, "the address for connection to the cockroach cluster.")
	}

	for _, cmds := range [][]*cobra.Command{nodeCmds, clientCmds} {
		for _, cmd := range cmds {
			cmd.Flags().BoolVar(&ctx.Insecure, "insecure", ctx.Insecure, "run over plain HTTP. WARNING: "+
				"this is strongly discouraged.")
		}
	}

	for _, cmds := range [][]*cobra.Command{nodeCmds, clientCmds, certCmds} {
		for _, cmd := range cmds {
			cmd.Flags().StringVar(&ctx.Certs, "certs", ctx.Certs, "directory containing RSA key and x509 certs. "+
				"This flag is required if --insecure=false.")
		}
	}
}

func init() {
	initFlags(Context)
}
