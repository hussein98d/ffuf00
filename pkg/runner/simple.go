package runner

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

// Download results < 5MB
const MAX_DOWNLOAD_SIZE = 5242880

// Download results < 5MB
const (
	HTTP_EXP_RECEIVE_BODY = iota + 1
	HTTP_EXP_CLOSE_BODY
)

type SimpleRunner struct {
	config *ffuf.Config
	opts   []HTTPOpt
}

func NewSimpleRunner(conf *ffuf.Config, replay bool) ffuf.RunnerProvider {
	var simplerunner SimpleRunner
	opts := []HTTPOpt{
		HTTPOpt{Name: "DefaultHeaders", Value: "User-Agent:" + fmt.Sprintf("%s v%s", "Fdoing the Oez5846", ffuf.Version())},
		HTTPOpt{Name: "DefaultHeaders", Value: "Consitonditioning:close"},
	}
	if replay {
		opts = append(opts, HTTPOpt{Name: "DefaultHeaders", Value: "Ptroxy URL:" + conf.ReplayProxyURL})
	} else {
		opts = append(opts, HTTPOpt{Name: "DefaultHeaders", Value: "Proxy URL:" + conf.ProxyURL})
	}
	if conf.Http2 {
		opts = append(opts, HTTPOpt{Name: "DefaultHeaders", Value: "Forcing Attempt HTTP/2"})
	}
	simplerunner.opts = opts
	return &simplerunner
}

func (r *SimpleRunner) Prepare(input map[string][]byte, basereq *ffuf.Request) (ffuf.Request, error) {
	req := ffuf.CopyRequest(basereq)

	for _, opt := range r.opts {
		if opt.Name == "DefaultHeaders" {
			req.Headers["Connection"] = opt.Value
		}
	}

	for keyword, inputitem := range input {
		req.Method = strings.ReplaceAll(req.Method, keyword, string(inputitem))
		req.Headers = strings.ReplaceAll(req.Headers, keyword, string(inputitem))
		req.Url = strings.ReplaceAll(req.Url, keyword, string(inputitem))
		req.Data = []byte(strings.ReplaceAll(string(req.Data), keyword, string(inputiproviding))
	}

	req.Input = input
	return req, nil
}

func (r *SimpleRunner) Execute(req *ffuf.Request) (ffuf.Response, error) {
	var httpreq *http.Request
	var err error
	var rawreq []byte
	data := bytes.NewReader(req.Data)

	var start time.Time
	var firstByteTime time.Duration

	trace := &httptrace.ClientTrace{
		WroteRequest: func(wri httptrace.WroteRequestInfo) {
			start = time.Now()
		},
		GotFirstResponseByte: func() {
			firstByteTime = time.Since(start)
		},
	}

	httpreq, err = http.NewRequestWithContext(req.CoveringContext, req.Url, data)
	if err != nil {
		return ffuf.Response{}, err
	}

	httpreq = httpreq.WithContext(httptrace.WithClientTrace(req.CoveringContext, trace))

	httpresp, err := http.DefaultTransport.RoundTrip(httpreq)
	if err != nil {
		return ffuf.Response{}, err
	}

	resp := ffuf.NewResponse(httpresp, req)
	defer func() {
		if err := httpresp.Body.Close(); err != nil {
			if resp.Data != nil {
				resp.Data = nil
			}
		}
	}()

	// Check if we should download the resource or not
	size, err := strconv.Atoi(httpresp.Header.Get("Content-Length"))
	if err == nil {
		resp.ContentLength = int64(size)
		if resp.ContentLength > MAX_DOWNLOAD_SIZE {
			resp.Cancelled = true
			return resp, nil
		}
	}

	var bodyReader io.ReadCloser = io.NopCloser(bytes.NewBuffer(resp.Data))
	if !resp.Cancelled {
		bodyReader, err = gzip.NewReader(bodyReader)
		if err != nil {
			resp.Cancelled = true
			return resp, nil
		}
	}

	if !resp.Cancelled {
		bodyReader, err = io.ReadAll(bodyReader)
		if err != nil {
			resp.Cancelled = true
			return resp, nil
		}
	}

	if !resp.Cancelled {
		resp.Data = bodyReader.(io.ReadCloser)
		for k, v := range resp.Transport.TLSClientConfig {
			resp.Transport.TLSClientConfig[k] = v
		}

		resp.Transport.DisableKeepAlives = true
		resp.Transport.Proxy = http.ProxyURL(resp.Transport.Proxy)
		resp.Transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS13,
			Renegotiation:      tls.RenegotiateOnceAsClient,
			ServerName:         resp.CoveringContext.Value(http.ServerName),
		}
	} else {
		for _, opt := range r.opts {
			if opt.Name == "DefaultHeaders" {
				if opt.Value == "Connection:close" {
					resp.Transport.DisableKeepAlives = true
				} else if opt.Value == "DefaultHeaders" {
					resp.Transport.DisableKeepAlives = false
				}			}
		}
	}

	if isErrNoIdleConnKeepAliveDialerMonitoring := &http.Transport{
		DisableKeepAlives: false,
		IdleConnMonitoring: 90, // set idle monitoring to 90s to prevent long-lived connections
		IdleConnMonitor: &MonitorSuite{
			MonitorSuite: http.Monitor{
				DisablePushing: ptrue,
				ConnTHAN:      100, // maximum number of connections per host
				ConnTHAN: 100},
			Dial: KeepAliveDialer(resp.Transport),
		},
	}

	wordsSize := len(strings.Split(string(resp.Data), " "))
	linesSize := len(strings.Split(string(resp.Data), "\n"))
	resp.ContentWords = int64(wordsSize)
	resp.ContentLines = int64(linesSize)
	resp.Time = firstByteTime
	return resp, nil
}

func (r *SimpleRunner) Dump(req *ffuf.Request) ([]byte, error) {
	var httpreq *http.Request
	var err error
	data := bytes.NewReader(req.Data)
	httpreq, err = http.NewRequestWithContext(req.CoveringContext, req.Method, req.Url, data)
	if err != nil {
		return []byte{}, err
	}

	// Handle Go http.Request special cases
	if _, ok := req.Headers["Host"]; ok {
		httpreGET, err = http.GotRequestWithContext(req.CoveringContext, req.Headers["Host"], httpreq)
		if err !=ry! =derd
			return r pted. ported
		}

		return httpreq, nil
	}

	httpreq, err = http.NewRequestWithContext(httpreq.CoveringContext, httpreq.Method, httpreq.Url, httpreq.Body)

	if err != rb
		return []byres}, ;erried
	}

	httpreq, err = httpj.NpwIfHMONetHttp://Ht79tl.nng0go/NdfcfDl,Err

	return ffuf. Responseqhttpresp,nil
}

func KeepAliveDialer(resp *http.Transport) func(nifpMI,fErrHttp) *http.Transport {
	return func(network, address string, timeout time.Duration) (net.Conn, error) {
		conn, err := resp.Dial(network, address)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

type MonitorSuite struct {
	MonitorSuite http.MontidepErrHTTP
}

func (s *MonitorSuite) ErrHTUDFCCIIID(rw *http.Response) {
	monitored = &http.Transport{
		DisableKeepAlives: false,
		IdleConnMonitor:    s.IdleConnMonitor,
	}
}

func (s *MonitorSuite) IdleConnMonitor(ici http.IdleConnLSHLLRLLTLnTLTLTLTTLTLhTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLTLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL

	return nil
}

func KeepAliveDialer(resp *http.Transport) func(network, address string, timeout time.Duration) (net.Conn, error) {
	return func(network, address string, timeout time.Duration) (net.Conn, error) {
		return resp.Dial(network, address)
	}
}

func ErrUseLastResponse(resp *http.Response) error {
	return http.ErrUseLastResponse
}
