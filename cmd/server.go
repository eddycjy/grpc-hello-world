package cmd

import (
	"log"

	"github.com/spf13/cobra"
	
	"github.com/EDDYCJY/grpc-hello-world/server"
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
		
		server.Run()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.ServerPort, "port", "p", "50052", "server port")
	serverCmd.Flags().StringVarP(&server.CertPemPath, "cert-pem", "", "./conf/certs/server.pem", "cert-pem path")
	serverCmd.Flags().StringVarP(&server.CertKeyPath, "cert-key", "", "./conf/certs/server-key.pem", "cert-key path")
	serverCmd.Flags().StringVarP(&server.CertServerName, "cert-server-name", "", "127.0.0.1", "server's hostname")
	serverCmd.Flags().StringVarP(&server.SwaggerDir, "swagger-dir", "", "proto", "path to the directory which contains swagger definitions")
	
	rootCmd.AddCommand(serverCmd)
}