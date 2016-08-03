package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// http://doc.open.alipay.com/doc2/detail?spm=0.0.0.0.1LPUGt&treeId=63&articleId=103758&docType=1

var AlipayGateway = "https://mapi.alipay.com/gateway.do?"
var LogPrefix = "[Alipay]"

var ErrNotFoundNotifyID = errors.New("not found notify_id")
var ErrReturnDataIsEmpty = errors.New("return data is empty")
var ErrNotifyDataIsEmpty = errors.New("notify data is empty")

type Alipay struct {
	partner string
	key     string
	email   string

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewPayment(partner, key, email string) *Alipay {
	a := &Alipay{
		partner: partner,
		key:     key,
		email:   email,
	}

	return a
}

func (a *Alipay) InitRSA(pubPath, priPath string) *Alipay {
	var (
		err        error
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
	)

	publicKey, err = newPublicKey(pubPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err = newPrivateKey(priPath)
	if err != nil {
		log.Fatal(err)
	}

	a.publicKey = publicKey
	a.privateKey = privateKey

	return a
}

func newPublicKey(path string) (pub *rsa.PublicKey, err error) {
	// Read the verify sign certification key
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		err = fmt.Errorf("%s bad key data: %s", LogPrefix, "not PEM-encoded")
		return
	}
	//	if got, want := block.Type, "CERTIFICATE"; got != want {
	//		err = fmt.Errorf("%s unknown key type %q, want %q", LogPrefix, got, want)
	//		return
	//	}

	// Decode the certification
	//cert, err = x509.ParseCertificate(block.Bytes)

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		err = fmt.Errorf("%s bad public key: %s", LogPrefix, err)
		return
	}
	pub = pubInterface.(*rsa.PublicKey)
	return
}

func newPrivateKey(path string) (priKey *rsa.PrivateKey, err error) {
	// Read the private key
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("%s read key file: %s", LogPrefix, err)
		return
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		err = fmt.Errorf("%s bad key data: %s", LogPrefix, "not PEM-encoded")
		return
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		err = fmt.Errorf("%s unknown key type %q, want %q", LogPrefix, got, want)
		return
	}

	// Decode the RSA private key
	priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		err = fmt.Errorf("%sbad private key: %s", LogPrefix, err)
		return
	}

	return
}

func (a *Alipay) rsaSign(kvs KVpairs) (sig string, err error) {
	h := sha1.New()
	io.WriteString(h, kvs.RemoveEmpty().Sort().Join("&"))
	hashed := h.Sum(nil)

	rsaSign, err := rsa.SignPKCS1v15(rand.Reader, a.privateKey, crypto.SHA1, hashed)
	if err != nil {
		return
	}

	sig = base64.StdEncoding.EncodeToString(rsaSign)
	return
}

func (a *Alipay) rsaVerify(vals url.Values, fields []string) (err error) {
	var signature, notifyID string
	kvs := KVpairs{}
	for key := range vals {
		if len(fields) > 0 && !Contains(fields, key) {
			continue
		}

		if key == "sign" {
			signature, _ = url.QueryUnescape(vals.Get(key))
			continue
		}

		if key == "sign_type" {
			continue
		}

		if key == "notify_id" {
			notifyID, _ = url.QueryUnescape(vals.Get(key))
		}

		var k, v string
		k, err = url.QueryUnescape(key)
		if err != nil {
			return
		}

		v, err = url.QueryUnescape(vals.Get(key))
		if err != nil {
			return
		}
		kvs = append(kvs, KVpair{K: k, V: v})
	}

	hashed := SHA1([]byte(kvs.RemoveEmpty().Sort().Join("&")))

	var inSign []byte
	inSign, err = base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return
	}

	err = rsa.VerifyPKCS1v15(a.publicKey, crypto.SHA1, hashed, inSign)
	if err != nil {
		return
	}

	err = a.checkNotify(notifyID)
	if err != nil {
		return
	}
	return
}

// verify 判断过来的参数是否有效
// TODO 只把支付非空参数加入验证 它把我传的参数也带过来了!!!!!!!!!!!
func (a *Alipay) verify(vals url.Values, fields []string) (err error) {
	var signature, notifyID string
	kvs := KVpairs{}
	for key := range vals {
		if len(fields) > 0 && !Contains(fields, key) {
			continue
		}

		val := vals.Get(key)

		if key == "sign" {
			signature = val
			continue
		}

		if key == "sign_type" {
			continue
		}

		if key == "notify_id" {
			notifyID = val
		}
		/*
			var v string
			v, err = url.QueryUnescape(val)
			if err != nil {
				return
			}
		*/kvs = append(kvs, KVpair{K: key, V: val})
	}

	signStr := MD5(kvs.RemoveEmpty().Sort().Join("&"), a.key)
	if signStr != signature {
		err = fmt.Errorf("%s illegal signature, want %s, got %s", LogPrefix, signature, signStr)
		return
	}
	_ = notifyID
	/*
		err = a.checkNotify(notifyID)
		if err != nil {
			return
		}
	*/
	return
}

/*
checkNotify 直接访问支付宝借口判断请求是否有效

得到的处理结果有两种：
成功时：true
不成功时：报对应错误
*/
func (a *Alipay) checkNotify(notifyID string) (err error) {
	if notifyID == "" {
		err = ErrNotFoundNotifyID
		log.Println("error:", err)
		return
	}

	vals := url.Values{}
	vals.Set("service", "notify_verify")
	vals.Set("partner", a.partner)
	vals.Set("notify_id", notifyID)

	r, err := http.NewRequest("GET", AlipayGateway+vals.Encode(), nil)
	if err != nil {
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		},
	}

	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if string(body) != "true" {
		err = fmt.Errorf("%s illegal notify.%s got %s", LogPrefix, AlipayGateway+vals.Encode(), string(body))
		return
	}
	return
}

func parsePostData(req *http.Request) (vals url.Values, err error) {
	var formStr = []byte(req.Form.Encode())
	if len(formStr) > 0 {
		var fields []string
		fields = strings.Split(string(formStr), "&")

		vals = url.Values{}
		data := map[string]string{}
		for _, field := range fields {
			f := strings.SplitN(field, "=", 2)
			key, val := f[0], f[1]
			data[key] = val
			vals.Set(key, val)
		}
	}
	return
}
