package cmd

import "github.com/spf13/cobra"

var rootCmd = cobra.Command{
    Use: "auction",
    Short: "Auction CLI Application",
    Long: "Auction is a real-time bidding tool to auction and bid for items from your terminal",
}

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}

func init() {
    rootCmd.AddCommand(joinCmd)
}
