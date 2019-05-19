package server

import (
    "crypto/tls"
    "net"
    "net/http"
    "log"
    "strings"
    "path"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/elazarl/go-bindata-assetfs"
    
    pb "github.com/EDDYCJY/grpc-hello-world/proto"
    "github.com/EDDYCJY/grpc-hello-world/pkg/util"
    "github.com/EDDYCJY/grpc-hello-world/pkg/ui/data/swagger"
)

var (
    ServerPort string
    CertServerName string
    CertPemPath string
    CertKeyPath string
    SwaggerDir string
    EndPoint string

    tlsConfig *tls.Config
)

func Run() (err error) {
    EndPoint = ":" + ServerPort
    tlsConfig = util.GetTLSConfig(CertPemPath, CertKeyPath)

    conn, err := net.Listen("tcp", EndPoint)
    if err != nil {
        log.Printf("TCP Listen err:%v\n", err)
    }

    srv := newServer(conn)

    log.Printf("gRPC and https listen on: %s\n", ServerPort)

    if err = srv.Serve(util.NewTLSListener(conn, tlsConfig)); err != nil {
        log.Printf("ListenAndServe: %v\n", err)
    }

    return err
}
 
func newServer(conn net.Listener) (*http.Server) {
    grpcServer := newGrpc()
    gwmux, err := newGateway()
    if err != nil {
        panic(err)
    }

    mux := http.NewServeMux()
    mux.Handle("/", gwmux)
    mux.HandleFunc("/swagger/", serveSwaggerFile)
    serveSwaggerUI(mux)

    return &http.Server{
        Addr:      EndPoint,
        Handler:   util.GrpcHandlerFunc(grpcServer, mux),
        TLSConfig: tlsConfig,
    }
}

func newGrpc() *grpc.Server {
    creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)
    if err != nil {
        panic(err)
    }

    opts := []grpc.ServerOption{
        grpc.Creds(creds),
    }

    server := grpc.NewServer(opts...)

    pb.RegisterHelloWorldServer(server, NewHelloService())

    return server
}

func newGateway() (http.Handler, error) {
    ctx := context.Background()
    dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertServerName)
    if err != nil {
        return nil, err
    }
    dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
    
    gwmux := runtime.NewServeMux()
    if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
        return nil, err
    }

    return gwmux, nil
}

func serveSwaggerFile(w http.ResponseWriter, r *http.Request) {
      if ! strings.HasSuffix(r.URL.Path, "swagger.json") {
        log.Printf("Not Found: %s", r.URL.Path)
        http.NotFound(w, r)
        return
    }

    p := strings.TrimPrefix(r.URL.Path, "/swagger/")
    p = path.Join(SwaggerDir, p)

    log.Printf("Serving swagger-file: %s", p)

    http.ServeFile(w, r, p)
}

func serveSwaggerUI(mux *http.ServeMux) {
    fileServer := http.FileServer(&assetfs.AssetFS{
        Asset:    swagger.Asset,
        AssetDir: swagger.AssetDir,
        Prefix:   "third_party/swagger-ui",
    })
    prefix := "/swagger-ui/"
    mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}