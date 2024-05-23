package api

import (
	"net/http"

	"github.com/markel1974/webautoma/src/server/config"
	"github.com/markel1974/webautoma/src/server/handlers/httphandler"
)

type Server struct {
	router   *http.ServeMux
	basePath string
	cfg      *config.Config
	hh       *httphandler.HttpHandler
}

func New(hh *httphandler.HttpHandler, cfg *config.Config) *Server {
	return &Server{
		hh:     hh,
		cfg:    cfg,
		router: http.NewServeMux(),
	}
}

func (srv *Server) Setup(basePath string, addr string) error {
	srv.basePath = basePath
	srv.router.HandleFunc(basePath+"version", srv.hh.Version)
	srv.router.HandleFunc(basePath+"status", srv.hh.Status)
	srv.router.HandleFunc(basePath+"run", srv.hh.Run)
	srv.router.HandleFunc(basePath+"stop", srv.hh.Stop)
	srv.router.HandleFunc(basePath+"pause", srv.hh.Pause)
	srv.router.HandleFunc(basePath+"resume", srv.hh.Resume)
	return nil
}

func (srv *Server) GetPath() string {
	return srv.basePath
}

func (srv *Server) GetHandler() *http.ServeMux {
	return srv.router
}

/*
func (srv *Server) timeHandler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	_, _ = w.Write([]byte("The time is: " + tm))
}
*/

/*
func (srv *Server) SetNotFoundHandler(fn func(ctx *session.Context)) {
	srv.notFoundConnector = connectors.NewConnector("server", "", fn, nil)
	srv.router.NotFoundHandler = NewHandler("", srv.managers, srv.notFoundConnector, iaccounts.RoleUndefined, false)
}

func (srv *Server) GetVersion() connectors.IConnector {
	return connectors.NewConnector("Server", "current version", srv.getVersion, nil)
}

func (srv *Server) getVersion(ctx *session.Context) {
	v := map[string]interface{}{"version": version.AppVersion}
	helpers.RespondJsonOk(ctx.W, v)
}
*/
