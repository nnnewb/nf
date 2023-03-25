package cmd

import (
	"fmt"
	"github.com/nnnewb/nf/internal/elevation"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/nnnewb/nf/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nf",
	Short: "A libpcap based traffic capture tool.",
	Long: "A libpcap based traffic capture tool.\n\n" +
		"Example:\n" +
		"    nf -i any -- udp and dst port 1701\n",
	Version: fmt.Sprintf("%s BUILD: %s %s %s", constants.VERSION, constants.BUILD_COMMIT, constants.BUILD_TIME, constants.BUILD_LINKAGE),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if !elevation.IsElevated() {
			err := elevation.RunElevated()
			if err != nil {
				log.Fatalf("run elevated failed, error %+v", err)
			}
			os.Exit(0)
		}

		err := cmd.ParseFlags(os.Args)
		if err != nil {
			log.Fatalf("parse flags failed, error %+v", err)
		}

		i, err := cmd.PersistentFlags().GetString("interface")
		if err != nil {
			log.Fatalf("get interface parameter failed, error %+v", err)
		}

		log.Printf("start capture packets on NIC %s\n", i)
		handle, err := pcap.OpenLive(i, 0, false, time.Second)
		if err != nil {
			log.Fatalf("pcap.OpenLive failed, error %+v", err)
		}
		defer handle.Close()

		log.Printf("used BPF: %s\n", strings.Join(args, " "))
		err = handle.SetBPFFilter(strings.Join(args, " "))
		if err != nil {
			log.Fatalf("handle.SetBPFFilter failed, error %+v", err)
		}

		pktSrc := gopacket.NewPacketSource(handle, handle.LinkType())
		for pkt := range pktSrc.Packets() {
			var (
				proto       string
				src         net.IP
				dst         net.IP
				spt         int
				dpt         int
				description string
			)

			if ipLayer := pkt.Layer(layers.LayerTypeIPv4); ipLayer != nil {
				ip := ipLayer.(*layers.IPv4)
				src = ip.SrcIP
				dst = ip.DstIP
			}

			if icmp4Layer := pkt.Layer(layers.LayerTypeICMPv4); icmp4Layer != nil {
				icmp4 := icmp4Layer.(*layers.ICMPv4)
				proto = "icmp4"
				description = fmt.Sprintf("ICMP %s, id %d, seq %d", icmp4.TypeCode, icmp4.Id, icmp4.Seq)
			}

			if tcpLayer := pkt.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp := tcpLayer.(*layers.TCP)
				proto = "tcp"
				spt = int(tcp.SrcPort)
				dpt = int(tcp.DstPort)
			}

			if udpLayer := pkt.Layer(layers.LayerTypeUDP); udpLayer != nil {
				udp := udpLayer.(*layers.UDP)
				proto = "udp"
				spt = int(udp.SrcPort)
				dpt = int(udp.DstPort)
			}

			switch proto {
			case "tcp", "udp":
				log.Printf("%s %s:%d > %s:%d %s", proto, src, spt, dst, dpt, description)
			case "icmp4", "icmp6":
				log.Printf("%s %s > %s %s", proto, src, dst, description)
			default:
				log.Printf("unknown pkt: %s", pkt.String())
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".nf" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".nf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
