package gopwntools

import (
	"os"
	"time"
)

func Asm(shellcode string) []byte {
	f, err := os.Create("/tmp/1")

	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(".intel_syntax noprefix\n.p2align 0\n" + shellcode + "\n"))

	if err != nil {
		panic(err)
	}

	f.Close()

	SetContext(Context{LogLevel: CRITICAL})
	defer SetContext(Context{LogLevel: INFO})

	{
		defer Process("/usr/bin/as", "-o", "/tmp/2", "/tmp/1").Close()
		time.Sleep(time.Millisecond * 100)

		defer Process("/usr/bin/objcopy", "-O", "binary", "-j", ".text", "/tmp/2", "/tmp/3").Close()
		time.Sleep(time.Millisecond * 100)
	}

	sc, err := os.ReadFile("/tmp/3")

	if err != nil {
		panic(err)
	}

	os.Remove("/tmp/1")
	os.Remove("/tmp/2")
	os.Remove("/tmp/3")

	return sc
}
