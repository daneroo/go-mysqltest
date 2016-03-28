package jsonl

import (
	"log"
	"time"

	. "github.com/daneroo/go-ted1k/types"
	. "github.com/daneroo/go-ted1k/util"
	"github.com/daneroo/timewalker"
)

type Writer struct {
	Grain timewalker.Duration
	enc   FBJE
	intvl timewalker.Interval
}

// Consume the Entry (receive only) channel
// preforming batched writes (of size writeBatchSize)
// Also performs progress logging (and timing)
func (w *Writer) Write(src <-chan Entry) {
	start := time.Now()
	count := 0

	for entry := range src {
		count++

		w.openFor(entry)
		err := w.enc.Encode(&entry)
		Checkerr(err)

	}
	w.close()
	TimeTrack(start, "jsonl.Write", count)
}

func (w *Writer) close() {
	w.enc.Close()
}

// Does 4 things; open File/buffer/encoder/Interval
//
func (w *Writer) openFor(entry Entry) {
	// could test Start==End (not initialized)
	if !w.intvl.Start.IsZero() {
		// log.Printf("-I: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	} else {
		s := w.Grain.Floor(entry.Stamp)
		e := w.Grain.AddTo(s)
		w.intvl = timewalker.Interval{Start: s, End: e}
		log.Printf("+I: %s : %s %s", w.Grain, entry.Stamp, w.intvl)
	}

	if !entry.Stamp.Before(w.intvl.End) {
		if w.enc.isOpen {
			log.Printf("Should close: %s", w.intvl)

			// new interva: for loop
			s := w.Grain.Floor(entry.Stamp)
			e := w.Grain.AddTo(s)
			w.intvl = timewalker.Interval{Start: s, End: e}

			w.enc.Close()
		}
	}

	if !w.enc.isOpen {
		log.Printf("Should open: %s", w.intvl)

		// this make directories as well...
		file, err := pathFor(w.Grain, w.intvl)
		Checkerr(err)

		err = w.enc.Open(file)
		Checkerr(err)
	}

}
