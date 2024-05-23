package web

import (
	"errors"
	"github.com/markel1974/webautoma/src/server/asset"
	"github.com/markel1974/webautoma/src/server/helpers"
	"net/http"
	"path"
	"strings"
)

const htmlNotFound = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>404 page not found</title>
	</head>
	<body>
		The page you are looking for is not found. You will redirected after <span id="s">5</span> seconds.
	</body>
	<script type="text/javascript">
	var s = 5;
	setInterval(function(){
		s--;
		document.getElementById('s').innerText = s;
		if (s === 0) location.href = '/ui/';
	}, 1000);
	</script>
	</html>
`

type Server struct {
	basePath  string
	prefix    string
	indexFile string
	redirect  string
	res       *asset.Resources
}

func New(res *asset.Resources) *Server {
	return &Server{
		res: res,
	}
}

func (as *Server) Setup(basePath string, prefix string, index string) error {
	as.prefix = prefix
	as.indexFile = strings.TrimLeft(index, "/")
	as.basePath = basePath
	if !strings.HasPrefix(as.basePath, "/") {
		as.basePath = "/" + as.basePath
	}
	if !strings.HasSuffix(as.basePath, "/") {
		as.basePath = as.basePath + "/"
	}
	as.redirect = strings.TrimSuffix(as.basePath, "/")
	return nil
}

func (as *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fp := path.Clean(as.prefix + r.URL.Path)
	if fp == "." {
		fp = ""
	} else {
		fp = strings.TrimLeft(fp, "/")
	}
	if err := as.assetHandler(w, r, fp); err == nil {
		return
	}
	if len(fp) > 0 {
		fp += "/"
	}
	fp += as.indexFile
	if err := as.assetHandler(w, r, fp); err == nil {
		return
	}
	notFoundHandler(w, r)
}

func (as *Server) assetHandler(w http.ResponseWriter, r *http.Request, fp string) error {
	content, hash, err := as.res.GetHash(helpers.DefinitionUI, fp)
	if err != nil {
		return err
	}
	noneMatchHeader := r.Header.Get("If-None-Match")
	noneMatch := strings.Trim(noneMatchHeader, "\"")
	if noneMatch == hash {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", hash)
		_, _ = w.Write(content)
	}
	return nil
}

/*
func (as *Server) NotFoundHandler(c *session.Context) {
	notFoundHandler(c.W, c.R)
}
*/

func (as *Server) GetPath() string {
	return as.basePath
}

//func (as *Server) RedirectHandler(w http.ResponseWriter, r *http.Request) {
//	http.Redirect(w, r, as.redirect, http.StatusSeeOther)
//}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strings.TrimLeft(r.URL.Path, "/"), "v1") {
		helpers.RespondJsonErrorWithCode(w, http.StatusNotFound, errors.New("api not found"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(htmlNotFound))
	}
}
