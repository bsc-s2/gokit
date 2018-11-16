# Min-Hash

## Introduction:
MinHash (or the min-wise independent permutations locality sensitive hashing scheme) is a technique for quickly estimating how similar two sets are

## Installation
To get the package use the standard:
```bash
go get github.com/bsc-s2/gokit/minhash
```

## Usage

calc similarity of two sets

To create two MinHashes with 128 buckets:
```
h0 := minhash.New(128)

h1 := minhash.New(128)
```

Add elements from set0 & set1:
```
for i := range set0 {
    h0.Add(set0[i])
}

for i := range set1 {
    h1.Add(set0[i])
}
```

Get similarity of two sets by MinHash(h0&h1):
```
similarity, err := h0.GetSimilarity(h1)
```

## Documentation
See the associated [GoDoc](http://godoc.org/github.com/bsc-s2/gokit/minhash)

## Example
See *func TestMinHashError(t `*`testing.T)* in minhash_test.go

## More Details
* [存储中的文件合并策略优化](http://drmingdrmer.github.io/tech/algorithm/2018/11/04/compact.html)
* [MinHash wiki](https://en.wikipedia.org/wiki/MinHash)