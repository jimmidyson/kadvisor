package commands

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fabric8io/kadvisor/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	kube_client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
)

var KadvisorCmd = &cobra.Command{
	Use:   "kadvisor",
	Short: "KAdvisor is a metrics collector & publisher for Kubernetes",
	Long:  "A configurable metrics collector & publisher for Kubernetes",
}

var kadvisorCmdV *cobra.Command

var (
	Verbose                  bool
	CfgFile                  string
	PollDuration             time.Duration
	KubernetesMaster         string
	KubernetesApiVersion     string
	KubernetesInsecure       bool
	KubernetesCACertFile     string
	KubernetesClientCertFile string
	KubernetesClientKeyFile  string
	KubernetesCACertData     string
	KubernetesClientCertData string
	KubernetesClientKeyData  string
	InfluxdbSinkUrl          string
	InfluxdbServiceName      string
	InfluxdbSecure           bool
)

var KubernetesClient *kube_client.Client

func Execute() {
	AddCommands()
	utils.StopOnErr(KadvisorCmd.Execute())
}

func AddCommands() {
	KadvisorCmd.AddCommand(startCmd)
}

//Initializes flags
func init() {
	KadvisorCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file")
	KadvisorCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose logging")
	KadvisorCmd.PersistentFlags().DurationVarP(&PollDuration, "poll", "p", 10*time.Second, "poll duration")

	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesMaster, "kubernetes-master", "k", "", "Kubernetes master")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesApiVersion, "kubernetes-api-version", "", "v1beta2", "Kubernetes API version")
	KadvisorCmd.PersistentFlags().BoolVarP(&KubernetesInsecure, "kubernetes-skip-tls-verify", "", false, "Skip TLS verify of Kubernetes master certificate")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesCACertFile, "kubernetes-ca-cert-file", "", "", "Path to a certificate file for the certificate authority")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesClientCertFile, "kubernetes-client-cert-file", "", "", "Path to a client certificate file for TLS")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesClientKeyFile, "kubernetes-client-key-file", "", "", "Path to a client key file for TLS")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesCACertData, "kubernetes-ca-cert-data", "", "", "Base 64 encoded CA certificate data for the certificate authority")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesClientCertData, "kubernetes-client-cert-data", "", "", "Base 64 encoded client certificate data for TLS")
	KadvisorCmd.PersistentFlags().StringVarP(&KubernetesClientKeyData, "kubernetes-client-key-data", "", "", "Base 64 encoded client key data for TLS")

	KadvisorCmd.PersistentFlags().StringVarP(&InfluxdbSinkUrl, "influxdb", "i", "", "InfluxDB URL")
	KadvisorCmd.PersistentFlags().StringVarP(&InfluxdbServiceName, "influxdb-service", "", "INFLUXDB", "InfluxDB service name")
	KadvisorCmd.PersistentFlags().BoolVarP(&InfluxdbSecure, "influxdb-secure", "", false, "InfluxDB service name")

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
	viper.SetDefault("influxdbSecure", false)
	viper.SetDefault("kubernetesApiVersion", "v1beta2")
	viper.SetDefault("kubernetesInsecure", false)

	if kadvisorCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("verbose", Verbose)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("poll").Changed {
		viper.Set("poll", PollDuration)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-master").Changed {
		viper.Set("kubernetesMaster", KubernetesMaster)
	} else if len(os.Getenv("KUBERNETES_MASTER")) > 0 {
		viper.Set("kubernetesMaster", "${KUBERNETES_MASTER}")
	} else if len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0 && len(os.Getenv("KUBERNETES_SERVICE_PORT")) > 0 {
		viper.Set("kubernetesMaster", "https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT}")
	}
	viper.Set("kubernetesMaster", os.ExpandEnv(viper.GetString("kubernetesMaster")))
	if len(viper.GetString("kubernetesMaster")) == 0 {
		log.Fatal("Kubernetes master is not set")
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-api-version").Changed {
		viper.Set("kubernetesApiVersion", KubernetesApiVersion)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-skip-tls-verify").Changed {
		viper.Set("kubernetesInsecure", KubernetesInsecure)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-ca-cert-file").Changed {
		viper.Set("kubernetesCACertFile", KubernetesCACertFile)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-client-cert-file").Changed {
		viper.Set("kubernetesClientCertFile", KubernetesClientCertFile)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-client-key-file").Changed {
		viper.Set("kubernetesClientKeyFile", KubernetesClientKeyFile)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-ca-cert-data").Changed {
		viper.Set("kubernetesCACertData", KubernetesCACertData)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-client-cert-data").Changed {
		viper.Set("kubernetesClientCertData", KubernetesClientCertData)
	}
	if kadvisorCmdV.PersistentFlags().Lookup("kubernetes-client-key-data").Changed {
		viper.Set("kubernetesClientKeyData", KubernetesClientKeyData)
	}

	if kadvisorCmdV.PersistentFlags().Lookup("influxdb").Changed {
		viper.Set("influxdb", InfluxdbSinkUrl)
	} else {
		viper.Set("influxdb", fmt.Sprintf("${%#[1]s_SERVICE_HOST}:${%#[1]s_SERVICE_PORT}", InfluxdbServiceName))
	}
	viper.Set("influxdb", os.ExpandEnv(viper.GetString("influxdb")))
	if kadvisorCmdV.PersistentFlags().Lookup("influxdb-secure").Changed {
		viper.Set("influxdbSecure", InfluxdbSecure)
	}

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("config", viper.AllSettings()).Debug("Configured settings")
}

func InitializeKubeClient() {
	KubernetesClient = kube_client.NewOrDie(&kube_client.Config{
		Host:     viper.GetString("kubernetesMaster"),
		Version:  viper.GetString("kubernetesApiVersion"),
		Insecure: viper.GetBool("kubernetesInsecure"),
		TLSClientConfig: kube_client.TLSClientConfig{
			CAFile:   viper.GetString("kubernetesCACertFile"),
			CertFile: viper.GetString("kubernetesClientCertFile"),
			KeyFile:  viper.GetString("kubernetesClientKeyFile"),
			CAData:   []byte(viper.GetString("kubernetesCACertData")),
			CertData: []byte(viper.GetString("kubernetesClientCertData")),
			KeyData:  []byte(viper.GetString("kubernetesClientKeyData")),
		},
	})
	if _, err := KubernetesClient.ServerVersion(); err != nil {
		log.Error(err)
	}
}
