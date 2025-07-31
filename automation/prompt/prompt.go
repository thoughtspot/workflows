package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"automation/logger"
)

type RW struct {
	Reader io.Reader
	Writer io.Writer
}

func NewRW(r io.Reader, w io.Writer) *RW {
	return &RW{
		Reader: r,
		Writer: w,
	}
}

func (rw *RW) Prompt(question, message string) string {
	l := logger.New()
	reader := bufio.NewReader(rw.Reader)
	fmt.Fprint(rw.Writer, "\n"+message+"\n")
	fmt.Fprint(rw.Writer, question+": ")

	answer, err := reader.ReadString('\n')
	if err != nil {
		l.Fatal(err)
	}
	return strings.TrimSpace(answer)
}
