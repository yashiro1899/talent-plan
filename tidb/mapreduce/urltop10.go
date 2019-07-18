package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// URLTop10 .
func URLTop10(nWorkers int) RoundsArgs {
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs
	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

func URLCountMap(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	m := make(map[string]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		m[line]++
	}

	kvs := make([]KeyValue, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, KeyValue{k, strconv.Itoa(v)})
	}
	return kvs
}

func URLCountReduce(key string, values []string) string {
	var total int
	for _, v := range values {
		i, _ := strconv.Atoi(v)
		total += i
	}
	return fmt.Sprintf("%s %d\n", key, total)
}

func URLTop10Map(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	m := make(map[string]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		kv := strings.Split(line, " ")
		v, _ := strconv.Atoi(kv[1])
		m[kv[0]] += v
	}

	urls, cnts := TopN(m, 10)
	kvs := make([]KeyValue, len(urls))
	for i, url := range urls {
		kvs[i] = KeyValue{"", fmt.Sprintf("%s %d", url, cnts[i])}
	}

	return kvs
}

func URLTop10Reduce(key string, values []string) string {
	m := make(map[string]int)
	for _, line := range values {
		kv := strings.Split(line, " ")
		v, _ := strconv.Atoi(kv[1])
		m[kv[0]] += v
	}

	urls, cnts := TopN(m, 10)
	buf := new(bytes.Buffer)
	for i, url := range urls {
		fmt.Fprintf(buf, "%s: %d\n", url, cnts[i])
	}
	return buf.String()
}
