package scf

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/chroblert/jlog"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

var (
	ScfApiProxyUrl string
)

func HandlerHttp(w http.ResponseWriter, r *http.Request) {
	dumpReq, err := httputil.DumpRequest(r, true)
	if err != nil {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := &DefineEvent{
		URL:     r.URL.String(),
		Content: base64.StdEncoding.EncodeToString(dumpReq),
	}
	//jlog.Println(event.URL, event.Content)
	bytejson, err := json.Marshal(event)
	if err != nil {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	req, err := http.NewRequest("POST", ScfApiProxyUrl, bytes.NewReader(bytejson))
	if err != nil {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	bytersp, _ := ioutil.ReadAll(resp.Body)
	var respevent RespEvent
	if err := json.Unmarshal(bytersp, &respevent); err != nil {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if resp.StatusCode > 0 && resp.StatusCode != 200 {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	//处理头+内容
	resp1 := strings.Split(respevent.Data, "^")
	respHeaders, err := base64.StdEncoding.DecodeString(resp1[0])
	respBody, err := base64.StdEncoding.DecodeString(resp1[1])
	//retByte, err := base64.StdEncoding.DecodeString(respevent.Data)
	if err != nil {
		jlog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	resp.Body.Close()
	respHeadersMap := make(map[string][]string)
	err = json.Unmarshal(respHeaders, &respHeadersMap)
	for k, v := range respHeadersMap {
		var s []string
		for _, val := range v {
			s = append(s, val)
		}
		//jlog.Printf("%s:%s\n", k, s[0])
		w.Header().Set(k, s[0])
	}
	w.Write(respBody)

	//w.Write(retByte)
	return
}

//}
