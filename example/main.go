package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"strconv"

	"net/http/httputil"

	"github.com/shima-park/alipay"
)

var (
	partner = "your pid"
	key     = "your key"
	email   = "your email"

	publicKeyPath  = "your rsa pubKey path" // "xxx/rsa_public_key.pem"
	privateKeyPath = "your rsa priKey path" // "xxx/rsa_private_key.pem"

	a = alipay.NewPayment(partner, key, email)
	// app 支付需要加入rsa公钥密钥
	// a.InitRSA(publicKeyPath, privateKeyPath)

	// 示例监听的端口
	port = ":9090"

	// 通过 lt --port 9090 获取的外网地址
	localTunnel = "http://eqfssupbgz.localtunnel.me"

	returnURL       = fmt.Sprintf("%s/%s", localTunnel, "alipay/return")
	notifyURL       = fmt.Sprintf("%s/%s", localTunnel, "alipay/notify")
	returnNotifyURL = fmt.Sprintf("%s/%s", localTunnel, "alipay/return-notify")
)

type MyServeMux struct {
	*http.ServeMux
}

func NewServeMux() *MyServeMux { return &MyServeMux{http.NewServeMux()} }

func (mux *MyServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	log.Println(string(dump))
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

func main() {
	mux := NewServeMux()
	mux.HandleFunc("/hello", HelloServer)
	mux.HandleFunc("/index", IndexServer)
	mux.HandleFunc("/alipay/payment-web", PaymentWebServer)     // 支付宝网页支付
	mux.HandleFunc("/alipay/payment-mobile", PaymentAPPServer)  // 支付宝app支付
	mux.HandleFunc("/alipay/return", ReturnWebServer)           // 支付宝网页支付返回处理
	mux.HandleFunc("/alipay/notify", NotifyWebServer)           // 支付宝支付通知
	mux.HandleFunc("/alipay/refund-notify", RefundNotifyServer) // 支付宝退款通知
	mux.HandleFunc("/alipay/refund", RefundServer)              // 支付宝退款

	log.Println("Listen", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func IndexServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	var html = `
<h1>Oh!!! It works!</h1>
<h2><a href="/alipay/payment-web">支付宝订单</a></h2>
<h2><a href="/alipay/payment-mobile">移动支付</a></h2>
`
	fmt.Fprintf(w, html)
	return
}

func PaymentWebServer(w http.ResponseWriter, req *http.Request) {
	var (
		outTradeNo  = time.Now().Format("20060102150405999")
		subject     = "test_subject"
		totalFee    = 0.01
		extraParams = map[string]string{
			"body":       "test body",
			"return_url": returnURL,
			"notify_url": notifyURL,
		}
	)
	checkourURL, err := a.InstantCredit(outTradeNo, subject, totalFee, extraParams)
	if err != nil {
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}
	http.Redirect(w, req, checkourURL, http.StatusFound)
	return
}

func PaymentAPPServer(w http.ResponseWriter, req *http.Request) {
	var (
		outTradeNo  = time.Now().Format("20060102150405999")
		subject     = "test_subject"
		totalFee    = 0.01
		extraParams = map[string]string{
			"body":       "test body",
			"return_url": returnURL,
			"notify_url": notifyURL,
		}
	)

	paymentParams, err := a.MobilePayment(outTradeNo, subject, totalFee, notifyURL, extraParams)
	if err != nil {
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", paymentParams)
	return
}

func ReturnWebServer(w http.ResponseWriter, req *http.Request) {
	r, err := a.InstantCreditReturn(req)
	if err != nil {
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}

	w.Header().Set("content-type", "text/html; charset=utf-8")

	var html = fmt.Sprintf(`
result:%+v<br>
<a href="/alipay/refund?trade_no=%s&refund_amount=%.2f&reason=不想买了">退款</a><br>
`, r, r.TradeNo, r.TotalFee)
	fmt.Fprintf(w, html)
	return
}

func RefundServer(w http.ResponseWriter, req *http.Request) {
	var (
		tradeNo         = req.URL.Query().Get("trade_no")
		reason          = req.URL.Query().Get("reason")
		refundAmount, _ = strconv.ParseFloat(req.URL.Query().Get("trade_amount"), 64)
		outRefundNo     = time.Now().Format("20060102150405000000")
		detailDatas     = []alipay.RefundDetailData{
			alipay.RefundDetailData{
				AlipayTransID: tradeNo,
				Amount:        refundAmount,
				RefundReason:  reason,
			},
		}
	)
	refundURL, err := a.Refund(outRefundNo, detailDatas, notifyURL)
	if err != nil {
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}
	http.Redirect(w, req, refundURL, http.StatusFound)
	return
}

func RefundNotifyServer(w http.ResponseWriter, req *http.Request) {
	r, err := a.RefundNotify(req)
	if err != nil {
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}

	var notifyResult = parseAlipayNotify(r.ResultDetails)
	for _, nr := range notifyResult {
		if nr.Result == "SUCCESS" {
			// do something...
		}
	}

	fmt.Fprintf(w, "%s", "success")
	return
}

type AlipayRefundNotifyResult struct {
	TradeNo string
	Amount  float64
	Result  string
}

func parseAlipayNotify(resultDetails string) (results []AlipayRefundNotifyResult) {
	var trades = strings.Split(resultDetails, "$")

	for _, trade := range trades {
		tradeInfo := strings.Split(trade, "^")
		if len(tradeInfo) == 3 {
			amount, _ := strconv.ParseFloat(tradeInfo[1], 64)
			results = append(results, AlipayRefundNotifyResult{
				TradeNo: tradeInfo[0],
				Amount:  amount,
				Result:  tradeInfo[2],
			})
		}
	}
	return
}

func NotifyWebServer(w http.ResponseWriter, req *http.Request) {
	r, err := a.InstantCreditNotify(req)
	if err != nil {
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}

	if r.TradeStatus != "TRADE_SUCCESS" && r.TradeStatus != "TRADE_FINISHED" {
		fmt.Fprintf(w, "%s", "fail")
		return
	}

	fmt.Fprintf(w, "%s", "success")
	return
}
