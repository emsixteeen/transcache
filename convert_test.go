package transcache

import (
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	c := Converter{
		Exec: "ffmpeg",
	}

	// f5bb7a05e578b64d884a893dd0b14f1f373a1968.mp4
	// a61e1b399dc8bc26eed44aff46ce25c755ffeeec.mp4
	i, err := os.Open("testdata/f5bb7a05e578b64d884a893dd0b14f1f373a1968.mp4")
	if err != nil {
		t.Fatal(err)
	}
	defer i.Close()

	o, err := os.CreateTemp("/tmp", "h265-")
	if err != nil {
		panic(err)
	}
	defer o.Close()

	if err = c.Convert(i, o); err != nil {
		panic(err)
	}

	fii, _ := i.Stat()
	fio, _ := o.Stat()

	t.Logf("input: %s, size: %d", i.Name(), fii.Size())
	t.Logf("output: %s, size: %d", o.Name(), fio.Size())
	t.Logf("ratio: 1:%f", float64(fii.Size())/float64(fio.Size()))
}
