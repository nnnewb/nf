package cmd

import (
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	ifs, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatalf("pcap.FindAllDevs failed, error %+v", err)
	}

	if len(ifs) == 0 {
		log.Fatalf("no suitable network interface found")
	}

	rootCmd.PersistentFlags().StringP("interface", "i", ifs[0].Name, "network interface card")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nf.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
