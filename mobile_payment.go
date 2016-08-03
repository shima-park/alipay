package alipay

import (
	"fmt"
	"net/http"
	"strconv"

	"net/url"
)

var initMobileParamMap = map[string]bool{
	"service":        true,  //  service	接口名称	String	接口名称，固定值。	不可空	mobile.securitypay.pay
	"partner":        true,  //  partner	合作者身份ID	String(16)	签约的支付宝账号对应的支付宝唯一用户号。以2088开头的16位纯数字组成。	不可空	2088101568358171
	"_input_charset": true,  //  _input_charset	参数编码字符集	String	商户网站使用的编码格式，固定为utf-8。	不可空	utf-8
	"sign_type":      true,  //  sign_type	签名方式	String	签名类型，目前仅支持RSA。	不可空	RSA
	"sign":           true,  //  sign	签名	String	请参见签名。	不可空	lBBK%2F0w5LOajrMrji7DUgEqNjIhQbidR13GovA5r3TgIbNqv231yC1NksLdw%2Ba3JnfHXoXuet6XNNHtn7VE%2BeCoRO1O%2BR1KugLrQEZMtG5jmJI
	"notify_url":     true,  //  notify_url	服务器异步通知页面路径	String(200)	支付宝服务器主动通知商户网站里指定的页面http路径。	不可空	http://notify.msp.hk/notify.htm
	"app_id":         false, //  app_id	客户端号	String	标识客户端。	可空	external
	"appenv":         false, //  appenv	客户端来源	String	标识客户端来源。参数值内容约定如下：appenv=”system=客户端平台名^version=业务系统版本”	可空	appenv=”system=android^version=3.0.1.2”
	"out_trade_no":   true,  //  out_trade_no	商户网站唯一订单号	String(64)	支付宝合作商户网站唯一订单号。	不可空	0819145412-6177
	"subject":        true,  //  subject	商品名称	String(128)	商品的标题/交易标题/订单标题/订单关键字等。该参数最长为128个汉字。	不可空	测试
	"payment_type":   true,  //  payment_type	支付类型	String(4)	支付类型。默认值为：1（商品购买）。	不可空	1
	"seller_id":      true,  //  seller_id	卖家支付宝账号	String(16)	卖家支付宝账号（邮箱或手机号码格式）或其对应的支付宝唯一用户号（以2088开头的纯16位数字）。	不可空	xxx@alipay.com
	"total_fee":      true,  //  total_fee	总金额	Number	该笔订单的资金总额，单位为RMB-Yuan。取值范围为[0.01，100000000.00]，精确到小数点后两位。	不可空	0.01
	"body":           true,  //  body	商品详情	String(512)	对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body。	不可空	测试测试
	"goods_type":     false, //  goods_type	商品类型	String(1)	具体区分本地交易的商品类型。  1：实物交易；  0：虚拟交易。  默认为1（实物交易）。	可空	1
	"rn_check":       false, //  rn_check	是否发起实名校验	String(1)	T：发起实名校验；  F：不发起实名校验。	可空	T
	"it_b_pay":       false, //  it_b_pay	未付款交易的超时时间	String	设置未付款交易的超时时间，一旦超时，该笔交易就会自动被关闭。当用户输入支付密码、点击确认付款后（即创建支付宝交易后）开始计时。取值范围：1m～15d，或者使用绝对时间（示例格式：2014-06-13 16:00:00）。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。该参数数值不接受小数点，如1.5h，可转换为90m。	可空	30m
	"extern_token":   false, //  extern_token	授权令牌	String(32)	开放平台返回的包含账户信息的token（授权令牌，商户在一定时间内对支付宝某些服务的访问权限）。通过授权登录后获取的alipay_open_id，作为该参数的value，登录授权账户即会为支付账户。	可空	1b258b84ed2faf3e88b4d979ed9fd4db
	"out_context":    false, //  out_context	商户业务扩展参数	String(128)	业务扩展参数，支付宝特定的业务需要添加该字段，json格式。 商户接入时和支付宝协商确定。	可空	{“ccode”:“shanghai”,“no”:“2014052600006128”}

}

type MobilePaymentNotify struct {
	NotifyTime       string  //  notify_time	通知时间	Date	通知的发送时间。格式为yyyy-MM-dd HH:mm:ss。	不可空	2013-08-22 14:45:24
	NotifyType       string  //  notify_type	通知类型	String	通知的类型。	不可空	trade_status_sync
	NotifyID         string  //  notify_id	通知校验ID	String	通知校验ID。	不可空	64ce1b6ab92d00ede0ee56ade98fdf2f4c
	SignType         string  //  sign_type	签名方式	String	固定取值为RSA。	不可空	RSA
	Sign             string  //  sign	签名	String	请参见签名机制。	不可空	lBBK%2F0w5LOajrMrji7DUgEqNjIhQbidR13GovA5r3TgIbNqv231yC1NksLdw%2Ba3JnfHXoXuet6XNNHtn7VE%2BeCoRO1O%2BR1KugLrQEZMtG5jmJI
	OutTradeNo       string  //  out_trade_no	商户网站唯一订单号	String(64)	对应商户网站的订单系统中的唯一订单号，非支付宝交易号。需保证在商户网站中的唯一性。是请求时对应的参数，原样返回。	可空	082215222612710
	Subject          string  //  subject	商品名称	String(128)	商品的标题/交易标题/订单标题/订单关键字等。它在支付宝的交易明细中排在第一列，对于财务对账尤为重要。是请求时对应的参数，原样通知回来。	可空	测试
	PaymentType      string  //  payment_type	支付类型	String(4)	支付类型。默认值为：1（商品购买）。	可空	1
	TradeNo          string  //  trade_no	支付宝交易号	String(64)	该交易在支付宝系统中的交易流水号。最短16位，最长64位。	不可空	2013082244524842
	TradeStatus      string  //  trade_status	交易状态	String	交易状态，取值范围请参见“交易状态”。	不可空	TRADE_SUCCESS
	SellerID         string  //  seller_id	卖家支付宝用户号	String(30)	卖家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	不可空	2088501624816263
	SellerEmail      string  //  seller_email	卖家支付宝账号	String(100)	卖家支付宝账号，可以是email和手机号码。	不可空	xxx@alipay.com
	BuyerID          string  //  buyer_id	买家支付宝用户号	String(30)	买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	不可空	2088602315385429
	BuyerEmail       string  //  buyer_email	买家支付宝账号	String(100)	买家支付宝账号，可以是Email或手机号码。	不可空	dlwdgl@gmail.com
	TotalFee         float64 //  total_fee	交易金额	Number	该笔订单的总金额。请求时对应的参数，原样通知回来。	不可空	1.00
	Quantity         int64   //  quantity	购买数量	Number	购买数量，固定取值为1（请求时使用的是total_fee）。	可空	1
	Price            float64 //  price	商品单价	Number	price等于total_fee（请求时使用的是total_fee）。	可空	1.00
	Body             string  //  body	商品描述	String(512)	该笔订单的备注、描述、明细等。对应请求时的body参数，原样通知回来。	可空	测试测试
	GMTCreate        string  //  gmt_create	交易创建时间	Date	该笔交易创建的时间。格式为yyyy-MM-dd HH:mm:ss。	可空	2013-08-22 14:45:23
	GMTPayment       string  //  gmt_payment	交易付款时间	Date	该笔交易的买家付款时间。格式为yyyy-MM-dd HH:mm:ss。	可空	2013-08-22 14:45:24
	IsTotalFeeAdjust string  //  is_total_fee_adjust	是否调整总价	String(1)	该交易是否调整过价格。	可空	N
	UseCoupon        string  //  use_coupon	是否使用红包买家	String(1)	是否在交易过程中使用了红包。	可空	N
	Discount         float64 //  discount	折扣	String	支付宝系统会把discount的值加到交易金额上，如果有折扣，本参数为负数，单位为元。	可空	0.00
	RefundStatus     string  //  refund_status	退款状态	String	取值范围请参见“退款状态”。	可空	REFUND_SUCCESS
	GMTRefund        string  //  gmt_refund	退款时间	Date	卖家退款的时间，退款通知时会发送。格式为yyyy-MM-dd HH:mm:ss。	可空	2008-10-29 19:38:25
}

func (a *Alipay) MobilePayment(outTradeNo, subject string, totalFee float64, notifyURL string, extraParams map[string]string) (s string, err error) {
	if outTradeNo == "" {
		err = fmt.Errorf("%s out_trade_no : Required parameter missing", LogPrefix)
		return
	}

	if subject == "" {
		err = fmt.Errorf("%s subject is required parameter", LogPrefix)
		return
	}

	if notifyURL == "" {
		err = fmt.Errorf("%s notify_url is required parameter", LogPrefix)
		return
	}

	if totalFee == 0 {
		err = fmt.Errorf("%s total_fee is required parameter", LogPrefix)
		return
	}

	if a.privateKey == nil {
		err = fmt.Errorf("%s rsa private key is not init", LogPrefix)
		return
	}

	params := a.initParams(outTradeNo, subject, notifyURL, totalFee, extraParams)
	kvs, err := GenKVpairs(initMobileParamMap, params, "sign", "sign_type")
	if err != nil {
		return
	}

	for i, kv := range kvs {
		kvs[i] = KVpair{K: kv.K, V: fmt.Sprintf(`"%s"`, kv.V)}
	}

	var sig string
	sig, err = a.rsaSign(kvs)
	if err != nil {
		return
	}

	kvs = append(kvs, KVpair{K: "sign", V: fmt.Sprintf(`"%s"`, url.QueryEscape(sig))})
	kvs = append(kvs, KVpair{K: "sign_type", V: `"RSA"`})

	s = kvs.Join("&")
	return
}

func (a *Alipay) initParams(outTradeNo, subject, notifyURL string, totalFee float64, extraParams map[string]string) (params map[string]string) {
	params = make(map[string]string)

	params["service"] = "mobile.securitypay.pay"
	params["_input_charset"] = "utf-8"
	params["payment_type"] = "1"

	params["partner"] = a.partner
	params["seller_id"] = a.partner

	params["notify_url"] = notifyURL
	params["out_trade_no"] = outTradeNo
	params["total_fee"] = strconv.FormatFloat(totalFee, 'f', 2, 64)
	params["subject"] = subject
	params["body"] = subject

	if extraParams != nil {
		for k, v := range extraParams {
			_, ok := instantCreditParamMap[k]
			if ok {
				params[k] = v
			}
		}
	}
	return
}

func (a *Alipay) MobilePaymentNotify(req *http.Request) (result *MobilePaymentNotify, err error) {
	vals, err := parsePostData(req)
	if err != nil {
		return
	}

	if len(vals) == 0 {
		err = ErrNotifyDataIsEmpty
		return
	}

	var fields = []string{
		"notify_time",
		"notify_type",
		"notify_id",
		"sign_type",
		"sign",
		"out_trade_no",
		"subject",
		"payment_type",
		"trade_no",
		"trade_status",
		"seller_id",
		"seller_email",
		"buyer_id",
		"buyer_email",
		"total_fee",
		"quantity",
		"price",
		"body",
		"gmt_create",
		"gmt_payment",
		"is_total_fee_adjust",
		"use_coupon",
		"discount",
		"refund_status",
		"gmt_refund",
		"gmt_close",
	}

	err = a.rsaVerify(vals, fields)
	if err != nil {
		return
	}

	var price, totalFee, discount float64
	price, _ = strconv.ParseFloat(vals.Get("price"), 64)
	totalFee, _ = strconv.ParseFloat(vals.Get("total_fee"), 64)
	discount, _ = strconv.ParseFloat(vals.Get("discount"), 64)

	var quantity int64
	quantity, _ = strconv.ParseInt(vals.Get("quantity"), 10, 64)

	result = &MobilePaymentNotify{
		NotifyTime:       vals.Get("notify_time"),
		NotifyType:       vals.Get("notify_type"),
		NotifyID:         vals.Get("notify_id"),
		SignType:         vals.Get("sign_type"),
		Sign:             vals.Get("sign"),
		OutTradeNo:       vals.Get("out_trade_no"),
		Subject:          vals.Get("subject"),
		PaymentType:      vals.Get("payment_type"),
		TradeNo:          vals.Get("trade_no"),
		TradeStatus:      vals.Get("trade_status"),
		SellerID:         vals.Get("seller_id"),
		SellerEmail:      vals.Get("seller_email"),
		BuyerID:          vals.Get("buyer_id"),
		BuyerEmail:       vals.Get("buyer_email"),
		TotalFee:         totalFee,
		Quantity:         quantity,
		Price:            price,
		Body:             vals.Get("body"),
		GMTCreate:        vals.Get("gmt_create"),
		GMTPayment:       vals.Get("gmt_payment"),
		IsTotalFeeAdjust: vals.Get("is_total_fee_adjust"),
		UseCoupon:        vals.Get("use_coupon"),
		Discount:         discount,
		RefundStatus:     vals.Get("refund_status"),
		GMTRefund:        vals.Get("gmt_refund"),
	}

	return
}
