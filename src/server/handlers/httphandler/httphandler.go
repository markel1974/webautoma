package httphandler

import (
	"github.com/markel1974/webautoma/src/server/helpers"
	"github.com/markel1974/webautoma/src/version"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

//curl 127.0.0.1:8080/v1/start -d '{"inputFile": "/Users/tinmr305/Development/go/src/markel/gm/build/verify.csv", "outputPath": "Result1", "csv": true, "json": true}'

type Options struct {
	InputFile  string `json:"inputFile"`
	OutputPath string `json:"outputPath"`
	CSV        bool   `json:"csv"`
	JSON       bool   `json:"json"`
}

type HttpHandler struct {
	lock     sync.RWMutex
	stubMode bool
}

func New() *HttpHandler {
	return &HttpHandler{}
}

func (fh *HttpHandler) Start() {
}

func (fh *HttpHandler) Version(w http.ResponseWriter, _ *http.Request) {
	m := map[string]string{
		"name":    version.AppName,
		"version": version.AppVersionShort,
		"date":    version.BuildDate,
	}
	helpers.RespondJsonOk(w, m)
}

func (fh *HttpHandler) Status(w http.ResponseWriter, r *http.Request) {
	/*
		name := urlQueryGetString(r.URL, "name")
		var ms *model.Status = nil
		fh.lock.RLock()
		if m, ok := fh.container[name]; ok {
			ms = m.Status(name)
		}
		fh.lock.RUnlock()

		if ms == nil {
			ms = model.NewStatus(fh.user, name, 0, "", nil, "unknown name")
		}
		helpers.RespondJsonOk(w, ms)
	*/
	helpers.RespondJsonOk(w, map[string]interface{}{"implement": "implement"})
}

func (fh *HttpHandler) Stop(w http.ResponseWriter, r *http.Request) {
	/*
		name := urlQueryGetString(r.URL, "name")
		var ms *model.Status = nil
		fh.lock.RLock()
		if m, ok := fh.container[name]; ok {
			ms = m.Stop(name)
		}
		fh.lock.RUnlock()
		if ms == nil {
			ms = model.NewStatus(fh.user, name, 0, "", nil, "unknown name")
		}
		helpers.RespondJsonOk(w, ms)

	*/
	helpers.RespondJsonOk(w, map[string]interface{}{"implement": "implement"})
}

func (fh *HttpHandler) Pause(w http.ResponseWriter, r *http.Request) {
	/*
		name := urlQueryGetString(r.URL, "name")
		var ms *model.Status = nil
		fh.lock.RLock()
		if m, ok := fh.container[name]; ok {
			ms = m.Pause(name)
		}
		fh.lock.RUnlock()
		if ms == nil {
			ms = model.NewStatus(fh.user, name, 0, "", nil, "unknown name")
		}
		helpers.RespondJsonOk(w, ms)
	*/
	helpers.RespondJsonOk(w, map[string]interface{}{"implement": "implement"})
}

func (fh *HttpHandler) Resume(w http.ResponseWriter, r *http.Request) {
	/*
		name := urlQueryGetString(r.URL, "name")
		var ms *model.Status = nil
		fh.lock.RLock()
		if m, ok := fh.container[name]; ok {
			ms = m.Resume(name)
		}
		fh.lock.RUnlock()
		if ms == nil {
			ms = model.NewStatus(fh.user, name, 0, "", nil, "unknown name")
		}
		helpers.RespondJsonOk(w, ms)
	*/
	helpers.RespondJsonOk(w, map[string]interface{}{"implement": "implement"})
}

func (fh *HttpHandler) Run(w http.ResponseWriter, r *http.Request) {
	/*
		if err := r.ParseMultipartForm(256 << 20); err != nil {
			helpers.RespondJsonErrorWithCode(w, http.StatusBadRequest, err)
			return
		}
		sideBuffer := getFormFile(r, "side")
		if len(sideBuffer) == 0 {
			helpers.RespondJsonErrorWithCode(w, http.StatusBadRequest, errors.New("empty side file"))
			return
		}
		varsBuffer := getFormFile(r, "vars")

		sideName := urlQueryGetString(r.URL, "sideName")
		if len(sideName) == 0 {
			helpers.RespondJsonErrorWithCode(w, http.StatusBadRequest, errors.New("empty side name"))
			return
		}
		workers := urlQueryGetInt(r.URL, "workers")
		vpnCounter := urlQueryGetInt(r.URL, "vpnCounter")
		vpnSeconds := urlQueryGetInt64(r.URL, "vpnSeconds")
		initialDelay := urlQueryGetInt(r.URL, "initialDelay")
		loop := urlQueryGetBool(r.URL, "loop")


		sideDriverId := fh.mDrivers.GetSideId()
		factories, err := fh.mDrivers.Compile([]string{sideDriverId})
		if err != nil {
			helpers.RespondJsonErrorWithCode(w, http.StatusBadRequest, errors.New("invalid side factory"))
			return
		}
		m := model.NewModel(fh.user, factories)
		if workers < 1 {
			workers = 1
		}
		opt := &idriver.CreateOptions{VpnSeconds: vpnSeconds, VpnCounter: vpnCounter, Workers: workers}
		if err = m.Setup(opt); err != nil {
			helpers.RespondJsonErrorWithCode(w, http.StatusBadRequest, err)
			return
		}

		fh.bindModel(sideName, m)

		defer fh.unbindModel(sideName, m)

		res := m.RunSide(sideDriverId, sideName, sideBuffer, varsBuffer, initialDelay, loop)
		helpers.RespondJsonOk(w, res)
	*/
	helpers.RespondJsonOk(w, map[string]interface{}{"implement": "implement"})
}

func (fh *HttpHandler) validPath(in string) string {
	var out []rune
	for _, v := range in {
		if v >= '0' && v <= '9' {
			out = append(out, v)
		} else if v >= 'a' && v <= 'z' {
			out = append(out, v)
		} else if v >= 'A' && v <= 'Z' {
			out = append(out, v)
		} else if v == '_' || v == '-' {
			out = append(out, v)
		} else {
			out = append(out, '_')
		}
	}
	return string(out)
}

func getFormFile(r *http.Request, id string) []byte {
	ff, _, err := r.FormFile(id)
	if err != nil {
		return nil
	}
	defer ff.Close()
	fb, err := io.ReadAll(ff)
	if err != nil {
		return nil
	}
	return fb
}

func urlQueryGetString(u *url.URL, id string) string {
	if u == nil {
		return ""
	}
	q := u.Query()
	if q == nil {
		return ""
	}
	return q.Get(id)
}

func urlQueryGetBool(u *url.URL, id string) bool {
	q := urlQueryGetString(u, id)
	q = strings.ToLower(q)
	if q == "true" || q == "1" {
		return true
	}
	return false
}

func urlQueryGetInt(u *url.URL, id string) int {
	q := urlQueryGetString(u, id)
	v, err := strconv.Atoi(q)
	if err != nil {
		return 0
	}
	return v
}

func urlQueryGetInt64(u *url.URL, id string) int64 {
	q := urlQueryGetString(u, id)
	v, err := strconv.ParseInt(q, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
