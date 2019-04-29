# store

## benchmark

- sqliteは過去の履歴にアクセスできるように実装したが，それにしても他の実装と比較して遅いことがわかる
- leveldb は平均的に read, write の性能ともに良い
    - 耐障害性がよくわからないため，第十日どうか判断できていないがとりあえずleveldbを利用することにする

### memory read

```
goos: linux
goarch: amd64
pkg: github.com/n0stack/n0stack/n0core/pkg/datastore/store/memory
BenchmarkMemoryStoreGet-8       20000000            59.2 ns/op           0 B/op           0 allocs/op
PASS
ok      github.com/n0stack/n0stack/n0core/pkg/datastore/store/memory    5.438s
Success: Benchmarks passed.
```

### memory write

```
goos: linux
goarch: amd64
pkg: github.com/n0stack/n0stack/n0core/pkg/datastore/store/memory
BenchmarkMemoryStoreApply-8       20000000            76.8 ns/op           7 B/op           1 allocs/op
PASS
ok      github.com/n0stack/n0stack/n0core/pkg/datastore/store/memory    3.630s
Success: Benchmarks passed.
```

### sqlite read

```
goos: linux
goarch: amd64
pkg: github.com/n0stack/n0stack/n0core/pkg/datastore/store/sqlite
BenchmarkSqliteStoreGet-8          30000         64989 ns/op        9757 B/op         198 allocs/op
PASS
ok      github.com/n0stack/n0stack/n0core/pkg/datastore/store/sqlite    81.512s
Success: Benchmarks passed.
```

### sqlite write


```
goos: linux
goarch: amd64
pkg: github.com/n0stack/n0stack/n0core/pkg/datastore/store/sqlite
BenchmarkSqliteStoreApply-8           1000       2184621 ns/op        7940 B/op         168 allocs/op
PASS
ok      github.com/n0stack/n0stack/n0core/pkg/datastore/store/sqlite    2.421s
Success: Benchmarks passed.
```

### leveldb read

```
goos: linux
goarch: amd64
pkg: github.com/n0stack/n0stack/n0core/pkg/datastore/store/leveldb
BenchmarkSqliteStoreApply-8        1000000          3615 ns/op         580 B/op          10 allocs/op
PASS
ok      github.com/n0stack/n0stack/n0core/pkg/datastore/store/leveldb    7.887s
Success: Benchmarks passed.
```

### leveldb write

```
goos: linux
goarch: amd64
pkg: github.com/n0stack/n0stack/n0core/pkg/datastore/store/leveldb
BenchmarkSqliteStoreApply-8         300000          4857 ns/op         287 B/op           4 allocs/op
PASS
ok      github.com/n0stack/n0stack/n0core/pkg/datastore/store/leveldb    3.239s
Success: Benchmarks passed.
```