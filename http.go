package transcache

import (
	"fmt"
	"net/http"
)

type Server struct {
	Addr      string
	Converter Converter

	mux *http.ServeMux
}

func handleConvert(c Converter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		src := r.PathValue("src")
		if src == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: copy all request values from request, validate src is valid URL, is allowed, etc.
		res, err := http.Get(src)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "upstream error: %s", err)
			return
		}
		defer res.Body.Close()

		// TODO: maybe validate content-type, etc.
		// TODO: we can't really write an error after we start writing already ...
		if err = c.ConvertCtx(r.Context(), res.Body, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "convert error: %s", err)
			return
		}
	}
}

func (s *Server) Configure() error {
	if s.Addr == "" {
		s.Addr = ":8080"
	}

	if s.mux == nil {
		s.mux = new(http.ServeMux)
	}

	if s.Converter.Exec == "" {
		return fmt.Errorf("missing converter")
	}

	s.mux.Handle("/convert/{src}", http.HandlerFunc(handleConvert(s.Converter)))
	return nil
}

func (s *Server) Run() error {
	if err := s.Configure(); err != nil {
		return err
	}

	srvr := http.Server{
		Addr:    s.Addr,
		Handler: s.mux,
	}

	return srvr.ListenAndServe()
}
