package transcache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type Cacher interface {
	Get(string) io.Reader
	Set(string) io.Writer
}

type CacherCtx interface {
	Cacher
	SetCtx(context.Context, string) io.WriteCloser
}

type Server struct {
	Addr      string
	Converter Converter
	Cache     CacherCtx

	mux *http.ServeMux
}

func handleConvert(s *Server) func(w http.ResponseWriter, r *http.Request) {
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
		wr := io.Writer(w)
		if s.Cache != nil {
			cr := s.Cache.Get(src)
			if cr != nil {
				fmt.Println("cache hit for:", src)

				_, err = io.Copy(w, cr)
				if err != nil {
					// TODO: Handle
				}
				return
			}

			cw := s.Cache.SetCtx(r.Context(), src)
			if cw != nil {
				wr = io.MultiWriter(w, cw)
			}
			defer cw.Close()
		}

		if err = s.Converter.ConvertCtx(r.Context(), res.Body, wr); err != nil {
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

	if s.Cache == nil {
		s.Cache = &MemoryCache{
			data: map[string]*bytes.Buffer{},
		}
	}

	s.mux.Handle("/convert/{src}", http.HandlerFunc(handleConvert(s)))
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
