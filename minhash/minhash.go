package minhash

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/cespare/xxhash"
	"io"
	"math"
)

// MinHash stores a Set of min hashes(per bucket has one min hash)
type MinHash struct {
	hashes    []uint64
	bucketCnt int
}

// New creates a new MinHash with bucketCnt
func New(bucketCnt int) *MinHash {
	h := make([]uint64, bucketCnt)
	for i := range h {
		h[i] = math.MaxUint64
	}
	return &MinHash{
		h, bucketCnt,
	}
}

// Add a element in MinHash
func (h *MinHash) Add(element []byte) {
	digest := xxhash.Sum64(element)

	pos := int(digest % uint64(h.bucketCnt))
	if digest < h.Get(pos) {
		h.Set(pos, digest)
	}
}

func (h *MinHash) AddString(element string) {
	h.Add([]byte(element))
}

// Get value from MinHash by position
func (h *MinHash) Get(pos int) uint64 {
	return h.hashes[pos]
}

// Set value in specified position
func (h *MinHash) Set(pos int, value uint64) {
	h.hashes[pos] = value
}

// WriteTo writes a binary representation of the MinHash to an i/o stream.
// It returns the number of bytes written.
func (h *MinHash) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.LittleEndian, uint64(h.bucketCnt))
	if err != nil {
		return
	}

	err = binary.Write(w, binary.LittleEndian, h.hashes)
	if err != nil {
		return
	}
	return int64(h.bucketCnt*8 + 8), nil
}

// ReadFrom reads a binary representation of the MinHash (such as might
// have been written by WriteTo()) from an i/o stream.
// It returns the number of bytes read.
func (h *MinHash) ReadFrom(r io.Reader) (n int64, err error) {
	var bucketCnt uint64
	err = binary.Read(r, binary.LittleEndian, &bucketCnt)
	if err != nil {
		return
	}
	h.bucketCnt = int(bucketCnt)

	h.hashes = make([]uint64, h.bucketCnt)
	err = binary.Read(r, binary.LittleEndian, h.hashes)
	if err != nil {
		return
	}
	return int64(h.bucketCnt*8 + 8), nil
}

// Marshal transfer MinHash to bytes
func (h *MinHash) Marshal() (b []byte, err error) {
	buf := bytes.NewBuffer(nil)
	_, err = h.WriteTo(buf)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

// Unmarshal transfer bytes to MinHash
func (h *MinHash) Unmarshal(b []byte) (n int64, err error) {
	buf := bytes.NewBuffer(b)
	return h.ReadFrom(buf)
}

// GetSimilarity of two MinHash
func (h *MinHash) GetSimilarity(h1 *MinHash) (float64, error) {
	if h.bucketCnt != h1.bucketCnt {
		return 0, errors.New("bucketCnt mismatch")
	}

	hits, estimated := 0, 0
	for i := range h.hashes {
		if h.Get(i) != math.MaxUint64 || h1.Get(i) != math.MaxUint64 {
		hits++
		if h.Get(i) == h1.Get(i) {
			estimated++
		}
		}
	}

	return float64(estimated) / float64(hits), nil
}

// GetSignature returns a signature for the set.
func (h *MinHash) GetSignature() []uint64 {
	return h.hashes
}