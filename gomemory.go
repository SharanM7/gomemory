package gomemory

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"time"
)

type Memorycheck struct {
	Out      io.Writer
	Timer    int
	memstats runtime.MemStats
	Ctx      context.Context
}

func New(ctx context.Context, t int, out io.Writer) *Memorycheck {
	return &Memorycheck{
		Out:   out,
		Timer: t,
		Ctx:   ctx,
	}
}

func (m *Memorycheck) KeepChecking() {
	t := time.NewTicker(time.Duration(m.Timer) * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-m.Ctx.Done():
			return
		case <-t.C:
			m.print()
		}
	}
}

func (m *Memorycheck) print() {
	runtime.ReadMemStats(&m.memstats)

	// Output memory stats
	m.Out.Write([]byte(fmt.Sprintf("Alloc = %v MiB", bToMb(m.memstats.Alloc))))
	m.Out.Write([]byte(fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.memstats.TotalAlloc))))
	m.Out.Write([]byte(fmt.Sprintf("\tSys = %v MiB", bToMb(m.memstats.Sys))))
	m.Out.Write([]byte(fmt.Sprintf("\tNumGC = %v\n", m.memstats.NumGC)))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
