package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Handler) GetDefault(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Works! %s", time.Now().Format(time.RFC850))
}