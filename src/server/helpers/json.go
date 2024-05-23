package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondJsonBufferWithCode(w http.ResponseWriter, httpCode int, buffer []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, _ = fmt.Fprint(w, string(buffer))
}

func RespondJsonWithCode(w http.ResponseWriter, httpCode int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		RespondJsonErrorWithCode(w, http.StatusInternalServerError, err)
		return
	}
	RespondJsonBufferWithCode(w, httpCode, b)
}

func RespondJsonErrorWithCode(w http.ResponseWriter, httpCode int, err error) {
	var e string
	if err != nil {
		e = err.Error()
	} else {
		e = "unknown error"
	}
	data := map[string]interface{}{"error": e}
	RespondJsonWithCode(w, httpCode, data)
}

func RespondJsonOk(w http.ResponseWriter, data interface{}) {
	RespondJsonWithCode(w, http.StatusOK, data)
}
