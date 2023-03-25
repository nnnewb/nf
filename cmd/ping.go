package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
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

				useSynScan, err := cmd.PersistentFlags().GetBool("syn")
				cobra.CheckErr(err)

				dpt, err := strconv.ParseInt(args[1], 10, 32)
				cobra.CheckErr(err)

				if useSynScan {
					err = pinger.SynScan(net.ParseIP(args[0]), int(dpt), interval)
					if err != nil {
						log.Printf("error: %+v", err)
						continue
					}
				} else {
					err = pinger.HandShakeTCP(net.ParseIP(args[0]), int(dpt), interval)
					if err != nil {
						log.Printf("error: %+v", err)
						continue
					}
				}

				log.Printf("%s:%s port is open\n", args[0], args[1])
			case "udp":
				if len(args) != 2 {
					_ = cmd.Help()
					return
				}

				dst := net.ParseIP(args[0])
				if dst == nil {
					cobra.CheckErr(fmt.Errorf("%s is not valid IP address", args[1]))
				}

				dpt, err := strconv.ParseInt(args[1], 10, 32)
				cobra.CheckErr(err)

				err = pinger.SendUDPPacket(dst, int(dpt))
				if err != nil {
					log.Printf("error: %s\n", err)
					continue
				}

				log.Printf("%s:%s packet sent\n", args[0], args[1])
			}

			if proto == "icmp" {
				if runtime.GOOS == "windows" {
					log.Fatalf("ICMP echo not implemented on windows.")
				}

				peer, rm, err := pinger.SendICMPEcho(os.Getpid()&0xffff, seq, net.UDPAddr{IP: net.ParseIP(args[0])}, interval*time.Second)
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
	pingCmd.PersistentFlags().Bool("syn", false, "try syn scan")
	pingCmd.PersistentFlags().Duration("interval", time.Second, "wait time between two packet sent")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
