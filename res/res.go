package res

import (
	"bufio"
	"sync"
)

type Res struct {
	b *bufio.Writer
	l *sync.RWMutex
}

var PHP *Res
var MySQL *Res

func NewPHP(b *bufio.Writer) {
	PHP = &Res{
		b: b,
		l: &sync.RWMutex{},
	}
}

func NewMySQL(b *bufio.Writer) {
	MySQL = &Res{
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
