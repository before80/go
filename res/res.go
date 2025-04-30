package res

import (
	"bufio"
	"sync"
)

type Res struct {
	b *bufio.Writer
	l *sync.RWMutex
}

var R *Res

func NewRes(b *bufio.Writer) {
	R = &Res{
		b: b,
		l: &sync.RWMutex{},
	}
}

func (res *Res) WriteStringAndFlush(str string) {
	res.l.Lock()
	defer res.l.Unlock()
	_, _ = res.b.WriteString(str)
	_ = res.b.Flush()
}
