/*
Copyright © 2022 weak_ptr <weak_ptr@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nnnewb/nf/internal/pinger"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "send packet to ip:port",
	Long: "send packet to ip:port\n\n" +
		"Example:\n" +
		"    nf ping -p tcp 192.168.1.1 22\n",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		proto, err := cmd.PersistentFlags().GetString("protocol")
		cobra.CheckErr(err)

		interval, err := cmd.PersistentFlags().GetDuration("interval")
		cobra.CheckErr(err)

		switch proto {
		case "tcp", "udp", "icmp":
		default:
			cobra.CheckErr(fmt.Errorf("unknown protocol %s", proto))
		}

		var seq int
		for {
			switch proto {
			case "tcp":
				if len(args) != 2 {
					_ = cmd.Help()
					return
				}

				dpt, err := strconv.ParseInt(args[1], 10, 32)
				cobra.CheckErr(err)

				err = pinger.HandShakeTCP(net.ParseIP(args[0]), int(dpt), interval)
				if err != nil {
					log.Printf("error: %s", err)
					continue
				}

				log.Printf("%s:%s port is open\n", args[0], args[1])
			case "udp":
				cobra.CheckErr(errors.New("not implemented yet"))
			case "icmp":
				peer, rm, err := pinger.SendICMPEcho(os.Getpid()&0xffff, seq, net.UDPAddr{IP: net.ParseIP(args[0])}, interval)
				if err != nil {
					log.Printf("error: %s\n", err)
					continue
				}

				log.Printf("%s: %s\n", rm.Type, peer)
				seq++
			}

			time.Sleep(interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	pingCmd.PersistentFlags().StringP("protocol", "p", "", "protocol, tcp/udp/icmp")
	pingCmd.PersistentFlags().Duration("interval", time.Second, "wait time between two packet sent")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
