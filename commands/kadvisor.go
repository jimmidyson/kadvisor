package commands

import (
	"github.com/fabric8io/kadvisor/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var KadvisorCmd = &cobra.Command{
	Use:   "kadvisor",
	Short: "KAdvisor is a metrics collector & publisher for Kubernetes",
	Long:  "A configurable metrics collector & publisher for Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
	},
}

var kadvisorCmdV *cobra.Command

var VerboseLog bool
var CfgFile string

func Execute() {
	utils.StopOnErr(KadvisorCmd.Execute())
}

//Initializes flags
func init() {
	KadvisorCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml)")
	KadvisorCmd.PersistentFlags().BoolVar(&VerboseLog, "verboseLog", false, "verbose logging")
	kadvisorCmdV = KadvisorCmd
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
	viper.SetConfigFile(CfgFile)
}
