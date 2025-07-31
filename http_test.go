package transcache

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_handleConvert(t *testing.T) {
	// f5bb7a05e578b64d884a893dd0b14f1f373a1968.mp4
	// a61e1b399dc8bc26eed44aff46ce25c755ffeeec.mp4
	file := "f5bb7a05e578b64d884a893dd0b14f1f373a1968.mp4"
	fileServer := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	src := fmt.Sprintf("%s/%s", fileServer.URL, file)

	srvr := &Server{
		Converter: Converter{
			Exec: "ffmpeg",
		},
	}
	if err := srvr.Configure(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(srvr.mux)
	u := fmt.Sprintf("%s/convert/%s", server.URL, url.QueryEscape(src))

	res, err := http.Get(u)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	t.Log(u)
	t.Log(res.Status)

	w, err := io.Copy(io.Discard, res.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("converted:", w, "bytes")
}
