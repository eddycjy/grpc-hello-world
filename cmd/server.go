package cmd

import (
	"log"

	"github.com/spf13/cobra"
	
	"grpc-hello-world/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC hello-world server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : %v", err)
			}
		}()
		
		server.Serve()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.ServerPort, "port", "p", "50051", "server port")
	serverCmd.Flags().StringVarP(&server.CertPemPath, "cert-pem", "", "./certs/server.pem", "cert pem path")
	serverCmd.Flags().StringVarP(&server.CertKeyPath, "cert-key", "", "./certs/server.key", "cert key path")
	serverCmd.Flags().StringVarP(&server.CertName, "cert-name", "", "grpc server name", "server's hostname")
	rootCmd.AddCommand(serverCmd)
}