package transcache

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Converter struct {
	Exec    string
	Codec   string
	Options map[string]string
	MaxProc int

	sem  *semaphore.Weighted
	args []string
}

func (c *Converter) defaults() {
	if c.Codec == "" {
		c.Codec = "libx265"
	}

	if c.Options == nil {
		c.Options = map[string]string{}
	}

	if c.sem == nil {
		if c.MaxProc == 0 {
			c.MaxProc = runtime.NumCPU() / 4
		}

		c.sem = semaphore.NewWeighted(int64(c.MaxProc))
	}
}

func (c *Converter) buildArgs() []string {
	if c.args != nil {
		return c.args
	}

	c.defaults()
	args := []string{}

	args = append(args, "-loglevel", "quiet")
	//args = append(args, "-threads", "12")
	args = append(args, "-i", "pipe:")
	args = append(args, "-c:v", c.Codec)

	for k, v := range c.Options {
		k = fmt.Sprintf("-%s", strings.Replace(k, ";", ":", 1))
		v = strings.Replace(v, ";", ":", 1)
		args = append(args, k, v)
	}

	args = append(args, "-an")
	args = append(args, "-f", "mp4")
	args = append(args, "-movflags", "empty_moov")
	//args = append(args, "-threads", "12")
	args = append(args, "pipe:")

	c.args = args
	return c.args

}

func (c *Converter) Convert(r io.Reader, w io.Writer) error {
	return c.ConvertCtx(context.Background(), r, w)
}

func (c *Converter) ConvertCtx(ctx context.Context, r io.Reader, w io.Writer) error {
	args := c.buildArgs()
	cmd := exec.CommandContext(ctx, c.Exec, args...)

	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if !c.sem.TryAcquire(1) {
		return fmt.Errorf("no resources")
	}
	defer c.sem.Release(1)

	fmt.Println("running:", cmd.String())
	return cmd.Run()
}
