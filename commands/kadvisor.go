package commands

import (
	"time"

	log "github.com/Sirupsen/logrus"
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

var (
	Verbose      bool
	CfgFile      string
	PollDuration time.Duration
)

func Execute() {
	utils.StopOnErr(KadvisorCmd.Execute())
}

//Initializes flags
func init() {
	KadvisorCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file")
	KadvisorCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose logging")
	KadvisorCmd.PersistentFlags().DurationVarP(&PollDuration, "poll", "p", 10*time.Second, "poll duration")
	kadvisorCmdV = KadvisorCmd
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
	if len(CfgFile) > 0 {
		viper.SetConfigFile(CfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("$HOME/.kadvisor")
		viper.AddConfigPath("/etc/kadvisor")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("kadvisor")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug("Unable to locate Config file")
	}

	viper.SetDefault("verbose", false)
	viper.SetDefault("poll", 10*time.Second)

	if kadvisorCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("verbose", Verbose)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("poll").Changed {
		viper.Set("poll", PollDuration)
	}

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("config", viper.AllSettings()).Debug("Configured settings")
}
