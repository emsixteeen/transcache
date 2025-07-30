package transcache

import (
	"io"
	"os"
	"os/exec"
)

type Converter struct {
	Exec string
}

func (c *Converter) Convert(r io.Reader, w io.Writer) error {
	cmd := exec.Command(c.Exec,
		"-loglevel", "quiet",
		"-i", "pipe:",
		"-c:v", "libx265",
		"-an",
		"-f", "mp4",
		"-movflags", "empty_moov",
		"pipe:")

	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
