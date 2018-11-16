package minhash

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestMinHashError(t *testing.T) {

	for i := 0.2; i < 1.0; i += 0.1 {
		testMinHashError(t, 128, 1024, i)
	}
}

const tolerance = 0.1

func testMinHashError(t *testing.T, bucketCnt, n int, sameRatio float64) {
	h0 := New(bucketCnt)
	h1 := New(bucketCnt)
	sameCnt := int(2 * float64(n) * sameRatio / (sameRatio+1))
	for i := 0; i < n; i++ {
		h0.AddString(strconv.Itoa(i))
	}
	for i := n - sameCnt; i < 2*n-sameCnt; i++ {
		h1.AddString(strconv.Itoa(i))
	}

	ret, err := h0.GetSimilarity(h1)
	if err != nil {
		t.Fatal(err)
	}
	exp := float64(sameCnt) / float64(n*2-sameCnt)
	fmt.Printf("exp:%f, minHash estimated:%f\n", exp, ret)
	if math.Abs(ret-exp) > tolerance {
		t.Fatal("error rate too high")
	}
}

func TestMinHashUnmarshal(t *testing.T) {
	bucketCnt := 128
	exp := New(bucketCnt)
	elementCnt := 1024
	for i := 0; i < elementCnt; i++ {
		exp.AddString(strconv.Itoa(i))
	}

	b, err := exp.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	ret := new(MinHash)
	n, err := ret.Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(bucketCnt*8+8) {
		t.Fatal("bytes len mismatch")
	}

	for i := 0; i < bucketCnt; i++ {
		if exp.Get(i) != ret.Get(i) {
			t.Fatal("decode mismatch", exp.Get(i), ret.Get(i))
		}
	}
}
