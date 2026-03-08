package main

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"io"
	"log/slog"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	var count int
	count = 100

	randomURL := newURLGenerator(os.Stdout, 2)
	for range count {
		b := randomURL.generate()
		_, err := url.Parse(string(b))
		if err != nil {
			randomURL.Write([]byte("\n"))
			slog.Error(err.Error())
			randomURL.Write([]byte("\n"))
			randomURL.Write([]byte("\n"))
		}
		randomURL.Write(b)

		randomURL.Write([]byte("\n"))
	}
}

func newURLGenerator(w io.WriteCloser, urlPartLen int) *urlGenerator {
	return &urlGenerator{
		rBuf:              bytes.NewBuffer(make([]byte, urlPartLen)),
		defaultURLPartLen: 2,
		w:                 w,
	}
}

type urlGenerator struct {
	mu      sync.RWMutex
	builder strings.Builder
	rBuf    *bytes.Buffer

	defaultURLPartLen int

	w io.WriteCloser
}

func (ug *urlGenerator) Write(b []byte) (int, error) {
	return ug.w.Write(b)
}

func (ug *urlGenerator) generate() []byte {
	ug.mu.Lock()
	defer ug.mu.Unlock()
	defer ug.builder.Reset()

	// generate Authority
	ug.builder.WriteString(ug.randromAuthority())
	ug.builder.WriteString("://")

	// generate Host
	ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))
	ug.builder.WriteRune('.')
	ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))

	//generate a PathSegment?
	for flipCoin() {
		ug.builder.WriteRune('/')
		ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))
	}

	if genrateQuery := flipCoin(); genrateQuery {
		// generate the first Query
		ug.builder.WriteRune('?')
		//key
		ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))
		ug.builder.WriteRune('=')
		//value
		ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))

		// generate another Query?
		for flipCoin() {
			ug.builder.WriteRune('&')
			//key
			ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))
			ug.builder.WriteRune('=')
			//value
			ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))
		}
	}

	// generate Fragment?
	if flipCoin() {
		ug.builder.WriteRune('#')
		ug.builder.WriteString(ug.randomURLEncodedString(ug.defaultURLPartLen))
	}

	return []byte(ug.builder.String())
}

func (ug *urlGenerator) Close() error {
	ug.mu.Lock()
	defer ug.mu.Unlock()
	return ug.w.Close()
}

func (ug *urlGenerator) randromAuthority() string {
	switch rand.Intn(2) {
	case 0:
		return "http"
	case 1:
		return "https"
	}
	return ""
}

func (ug *urlGenerator) randomURLEncodedString(len int) string {
	for ug.rBuf.Len() < len {
		ug.rBuf.WriteByte(0) // inittialize with zeroes to fill up by random bytes
	}

	cryptoRand.Read(ug.rBuf.Bytes()) //revive:disable-line Fills up the buffer and never returns an error
	defer ug.rBuf.Reset()

	return base64.RawURLEncoding.EncodeToString(ug.rBuf.Bytes()[:len])
}

func flipCoin() bool {
	r := rand.Int()
	return r%2 == 0
}
