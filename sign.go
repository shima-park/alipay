package alipay

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type KVpair struct {
	K, V string
}

type KVpairs []KVpair

func (t KVpairs) Less(i, j int) bool {
	return t[i].K < t[j].K
}

func (t KVpairs) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t KVpairs) Len() int {
	return len(t)
}

func (t KVpairs) Sort() KVpairs {
	sort.Sort(t)
	return t
}

func (t KVpairs) RemoveEmpty() KVpairs {
	for i := 0; i < len(t); i++ {
		if t[i].V == "" {
			t = append(t[:i], t[i+1:]...)
			i--
		}
	}
	return t
}

func (t KVpairs) Join(sep string) string {
	var strs []string
	for _, kv := range t {
		strs = append(strs, kv.K+"="+kv.V)
	}
	return strings.Join(strs, sep)
}

func MD5(strs ...string) string {
	h := md5.New()
	for _, str := range strs {
		io.WriteString(h, str)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SHA1(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return h.Sum(nil)
}

func GenKVpairs(paramsKeyMap map[string]bool, initParams map[string]string, ignoreKey ...string) (kvs KVpairs, err error) {
	kvs = make(KVpairs, 0)
	for key, isMust := range paramsKeyMap {
		val, ok := initParams[key]
		if ok && val != "" {
			kvs = append(kvs, KVpair{K: key, V: val})
		} else {
			// sign 参数需要签名后才会生成这里跳过
			if isMust && !Contains(ignoreKey, key) {
				err = errors.New("must param is empty:" + key)
				return
			}
		}
	}
	return
}

func Contains(strs []string, key string) bool {
	for _, v := range strs {
		if v == key {
			return true
		}
	}
	return false
}
