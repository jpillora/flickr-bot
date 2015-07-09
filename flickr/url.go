package flickr

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
)

type Args map[string]string

func URL(apiKey, secret, endpoint string, args Args) string {
	//set defaults
	vals := url.Values{}
	if apiKey != "" {
		vals.Set("api_key", apiKey)
	}
	if args != nil {
		//set args
		for k, v := range args {
			vals.Set(k, v)
		}
	}
	//set signature
	if secret != "" {
		i := 0
		keys := make([]string, len(vals))
		for k, _ := range vals {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		hash := md5.New()
		hash.Write([]byte(secret))
		for _, k := range keys {
			hash.Write([]byte(fmt.Sprintf("%s%s", k, vals.Get(k))))
		}
		vals.Set("api_sig", hex.EncodeToString(hash.Sum(nil)))
	}
	//final url!
	return endpoint + "?" + vals.Encode()
}
