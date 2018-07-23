package term

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/wzshiming/winseq"
	"golang.org/x/crypto/ssh/terminal"
)

type CmdFunc func(cmd ...string) (string, int, error)

type Terminal struct {
	Reader  io.Reader
	Writer  io.Writer
	Prompt  string
	CmdFunc CmdFunc
}

func NewTerminal(prompt string, cmd CmdFunc) *Terminal {
	return &Terminal{
		Reader:  os.Stdin,
		Writer:  os.Stdout,
		Prompt:  prompt,
		CmdFunc: cmd,
	}
}

func (c *Terminal) Run() error {

	ter := terminal.NewTerminal(struct {
		io.Reader
		io.Writer
	}{
		c.Reader,
		c.Writer,
	}, "")

	for {
		line, err := ter.ReadPassword(c.Prompt)
		if err != nil {
			if err == io.EOF {
				continue
			}
			return err
		}

		read := csv.NewReader(bytes.NewBufferString(line))
		read.Comma = ' '
		read.TrimLeadingSpace = true
		da, err := read.ReadAll()
		if err != nil {
			return err
		}
		for _, v := range da {
			beg := time.Now()
			result, size, err := c.CmdFunc(v...)
			if err != nil {
				return err
			}
			ter.Write([]byte(result))

			sub := time.Now().Sub(beg).Truncate(time.Second / 1000)
			if size == 0 {
				ter.Write([]byte(fmt.Sprintf("\n(%s)", sub)))
			} else {
				ter.Write([]byte(fmt.Sprintf("\n%d result(s) (%s)", size, sub)))
			}

		}
	}
	return nil
}