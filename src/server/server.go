package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/markel1974/webautoma/src/server/api"
	"github.com/markel1974/webautoma/src/server/asset"
	"github.com/markel1974/webautoma/src/server/config"
	"github.com/markel1974/webautoma/src/server/handlers/httphandler"
	"github.com/markel1974/webautoma/src/server/helpers"
	"github.com/markel1974/webautoma/src/server/web"
)

const (
	//readTimeout = 30
	//writeTimeout = 15
	baseApiPath = "/v1/"
	//baseWebsocketPath = "/websocket/"
	baseUIPath  = "/ui/"
	baseUIIndex = "index.html"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	resources  *asset.Resources
	//adapter    *core.Adapter
}

func New() *Server {
	fe := &Server{
		resources: asset.New(),
	}
	return fe
}

func (fe *Server) Setup(hh *httphandler.HttpHandler, cfg *config.Config) error {
	//fe.adapter = adapter
	fe.cfg = cfg

	if err := fe.resources.SetupBuffer(helpers.DefinitionUI, web.Content()); err != nil {
		return err
	}

	apiServer := api.New(hh, cfg)
	if err := apiServer.Setup(baseApiPath, helpers.CompileUrl(cfg.Listen, cfg.Domain, cfg.TLS)); err != nil {
		return err
	}

	uiPath := baseUIPath
	if len(cfg.UIPath) > 0 {
		uiPath = cfg.UIPath
	}

	uiServer := web.New(fe.resources)
	if err := uiServer.Setup(uiPath, "", baseUIIndex); err != nil {
		return err
	}

	mux := apiServer.GetHandler()

	mux.Handle(uiServer.GetPath(), http.StripPrefix(uiServer.GetPath(), uiServer))

	mux.Handle("/", http.StripPrefix("/", uiServer))

	fe.httpServer = &http.Server{
		Addr:    cfg.Listen,
		Handler: mux,
		//WriteTimeout: time.Duration(writeTimeout) * time.Second,
		//ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	return nil
}

func (fe *Server) getTLSConfig(host string, caCertFile string, certOpt tls.ClientAuthType) (*tls.Config, error) {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	//if certOpt > tls.RequestClientCert {
	caCert, err = os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	//}

	return &tls.Config{
		ServerName: host,
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		// TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
		MinVersion: tls.VersionTLS12,
	}, nil
}

func (fe *Server) Start() error {
	var err error

	log.Printf("Starting Web on %s", fe.cfg.Listen)

	if fe.cfg.TLS {
		var tlsConfig *tls.Config
		if len(fe.cfg.CA) > 0 {
			tlsConfig, err = fe.getTLSConfig(fe.cfg.Domain, fe.cfg.CA, tls.NoClientCert)
		}
		if err == nil {
			fe.httpServer.TLSConfig = tlsConfig
			err = fe.httpServer.ListenAndServeTLS(fe.cfg.Cert, fe.cfg.Key)
		}
	} else {
		err = fe.httpServer.ListenAndServe()
	}

	return err
}

func (fe *Server) serveHTTP(h http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		var r = recover()
		if r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}
			//logger.Error(nil, "Internal Server Error %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	h.ServeHTTP(w, r)
}
