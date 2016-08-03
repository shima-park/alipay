package alipay

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var instantCreditParamMap = map[string]bool{
	"service":             true,  //	接口名称	String	接口名称。	不可空	create_direct_pay_by_user
	"partner":             true,  //	   合作者身份ID	String(16)	签约的支付宝账号对应的支付宝唯一用户号。以2088开头的16位纯数字组成。	不可空	2088101011913539
	"_input_charset":      true,  //	参数编码字符集	String	商户网站使用的编码格式，如utf-8、gbk、gb2312等。	不可空	gbk
	"sign_type":           true,  //	签名方式	String	DSA、RSA、MD5三个值可选，必须大写。	不可空	MD5
	"sign":                true,  //	签名	String	请参见签名	不可空	7d314d22efba4f336fb187697793b9d2
	"notify_url":          false, //	服务器异步通知页面路径	String(190)	支付宝服务器主动通知商户网站里指定的页面http路径。	可空	http://api.test.alipay.net/atinterface/receive_return.htm
	"return_url":          false, //	页面跳转同步通知页面路径	String(200)	支付宝处理完请求后，当前页面自动跳转到商户网站里指定页面的http路径。	可空	http://api.test.alipay.net/atinterface/receive_return.htm
	"error_notify_url":    false, //	请求出错时的通知页面路径	String(200)	当商户通过该接口发起请求时，如果出现提示报错，支付宝会根据请求出错时的通知错误码通过异步的方式发送通知给商户。该功能需要联系支付宝开通。	可空	http://api.test.alipay.net/atinterface/receive_return.htm
	"out_trade_no":        true,  //	商户网站唯一订单号	String(64)	支付宝合作商户网站唯一订单号。	不可空	6843192280647118
	"subject":             true,  //	商品名称	String(256)	商品的标题/交易标题/订单标题/订单关键字等。该参数最长为128个汉字。	不可空	贝尔金护腕式
	"payment_type":        true,  //	支付类型	String(4)	取值范围请参见附录收款类型”。默认值为：1（商品购买）。注意：支付类型为“47”时，公共业务扩展参数（extend_param）中必须包含凭证号（evoucheprod_evouche_id）参数名和参数值。	不可空	1
	"total_fee":           true,  //	交易金额	Number	该笔订单的资金总额，单位为RMB-Yuan。取值范围为[0.01，100000000.00]，精确到小数点后两位。	不可空	100
	"seller_id":           true,  //	卖家支付宝用户号	String(16)	卖家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	不可空	2088002007018966
	"buyer_id":            false, //	买家支付宝用户号	String(16)	买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	可空	2088002007018955
	"seller_email":        false, //	卖家支付宝账号	String(100)	卖家支付宝账号，格式为邮箱或手机号。	可空	alipay-test01@alipay.com
	"buyer_email":         false, //	买家支付宝账号	String(100)	买家支付宝账号，格式为邮箱或手机号。	可空	tstable01@alipay.com
	"seller_account_name": false, //	卖家别名支付宝账号	String(100)	卖家别名支付宝账号。卖家信息优先级：seller_id>seller_account_name>seller_email。	可空	tstable02@alipay.com
	"buyer_account_name":  false, //	买家别名支付宝账号	String(100)	买家别名支付宝账号。买家信息优先级：buyer_id>buyer_account_name>buyer_email。	可空	tstable03@alipay.com
	"price":               false, //	商品单价	Number	单位为：RMB Yuan。取值范围为[0.01,100000000.00]，精确到小数点后两位。此参数为单价。规则：price、quantity能代替total_fee。即存在total_fee，就不能存在price和quantity；存在price、quantity，就不能存在total_fee。	可空	10.00
	"quantity":            false, //	购买数量	Number	price、quantity能代替total_fee。即存在total_fee，就不能存在price和quantity；存在price、quantity，就不能存在total_fee。	可空	1
	"body":                false, //	商品描述	String(1000)	对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body。	可空	美国专业护腕鼠标垫，舒缓式凝胶软垫模拟手腕的自然曲线和运动，创造和缓的GelFlex舒适地带！
	"show_url":            false, //	商品展示网址	String(400)	收银台页面上，商品展示的超链接。	可空	http://www.360buy.com/product/113714.html
	"paymethod":           false, //	默认支付方式	String	取值范围：creditPay（信用支付）；directPay（余额支付）。如果不设置，默认识别为余额支付。说明：必须注意区分大小写。	可空	directPay
	"enable_paymethod":    false, //	支付渠道	String	用于控制收银台支付渠道显示，该值的取值范围请参见支付渠道。可支持多种支付渠道显示，以“^”分隔。	可空	directPay^bankPay^cartoon^cash
	"need_ctu_check":      false, //	网银支付时是否做CTU校验	String	商户在配置了支持CTU（支付宝风险稽查系统）校验权限的前提下，可以选择本次交易是否需要经过CTU校验。Y：做CTU校验；N：不做CTU校验。	可空	Y
	"anti_phishing_key":   false, //	防钓鱼时间戳	String	通过时间戳查询接口获取的加密支付宝系统时间戳。如果已申请开通防钓鱼时间戳验证，则此字段必填。	可空	587FE3D2858E6B01E30104656E7805E2
	"exter_invoke_ip":     false, //	客户端IP	String(15)	用户在创建交易时，该用户当前所使用机器的IP。如果商户申请后台开通防钓鱼IP地址检查选项，此字段必填，校验用。	可空	128.214.222.111
	"extra_common_param":  false, //	公用回传参数	String(100)	如果用户请求时传递了该参数，则返回给商户时会回传该参数。	可空	你好，这是测试商户的广告。
	"extend_param":        false, //	公用业务扩展参数	String	用于商户的特定业务信息的传递，只有商户与支付宝约定了传递此参数且约定了参数含义，此参数才有效。参数格式：参数名1^参数值1|参数名2^参数值2|…… 多条数据用“|”间隔。支付类型（payment_type）为47（电子卡券）时，需要包含凭证号（evoucheprod_evouche_id）参数名和参数值。	可空	pnr^MFGXDW|start_ticket_no^123|end_ticket_no^234|b2b_login_name^abc
	"it_b_pay":            false, //	超时时间	String	设置未付款交易的超时时间，一旦超时，该笔交易就会自动被关闭。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（无论交易何时创建，都在0点关闭）。该参数数值不接受小数点，如1.5h，可转换为90m。该功能需要联系支付宝配置关闭时间。	可空	1h
	"default_login":       false, //	自动登录标识	String	用于标识商户是否使用自动登录的流程。如果和参数buyer_email一起使用时，就不会再让用户登录支付宝，即在收银台中不会出现登录页面。取值有以下情况：Y代表使用；N代表不使用。该功能需要联系支付宝配置。	可空	Y
	"product_type":        false, //	商户申请的产品类型	String(50)	用于针对不同的产品，采取不同的计费策略。如果开通了航旅垂直搜索平台产品，请填写CHANNEL_FAST_PAY；如果没有，则为空。	可空	CHANNEL_FAST_PAY
	"token":               false, //	快捷登录授权令牌	String(40)	如果开通了快捷登录产品，则需要填写；如果没有开通，则为空。	可空	201103290c9f9f2c03db4267a4c8e1bfe3adfd52
	"sign_id_ext":         false, //	商户买家签约号	String(50)	用于唯一标识商户买家。如果本参数不为空，则sign_name_ext不能为空。	可空	ZHANGSAN
	"sign_name_ext":       false, //	商户买家签约名	String(128)	商户买家唯一标识对应的名字。	可空	张三
	"qr_pay_mode":         false, //	扫码支付方式	String(1)	扫码支付的方式，支持前置模式和跳转模式。
	//前置模式是将二维码前置到商户的订单确认页的模式。需要商户在自己的页面中以iframe方式请求支付宝页面。具体分为以下3种：
	//0：订单码-简约前置模式，对应iframe宽度不能小于600px，高度不能小于300px；
	//1：订单码-前置模式，对应iframe宽度不能小于300px，高度不能小于600px；
	//3：订单码-迷你前置模式，对应iframe宽度不能小于75px，高度不能小于75px。
	//跳转模式下，用户的扫码界面是由支付宝生成的，不在商户的域名下。
	//2：订单码-跳转模式	可空	1
}

type InstantCreditReturn struct {
	IsSuccess        string  //  成功标识	String(1)	表示接口调用是否成功，并不表明业务处理结果。	不可空	T
	SignType         string  //  签名方式	String	DSA、RSA、MD5三个值可选，必须大写。	不可空	MD5
	Sign             string  //  签名	String(32)	请参见签名验证	不可空	b1af584504b8e845ebe40b8e0e733729
	OutTradeNo       string  //  商户网站唯一订单号	String(64)	对应商户网站的订单系统中的唯一订单号，非支付宝交易号。需保证在商户网站中的唯一性。是请求时对应的参数，原样返回。	可空	6402757654153618
	Subject          string  //  商品名称	String(256)	商品的标题/交易标题/订单标题/订单关键字等。	可空	手套
	PaymentType      string  //  支付类型	String(4)	对应请求时的payment_type参数，原样返回。	可空	1
	Exterface        string  //  接口名称	String	标志调用哪个接口返回的链接。	可空	create_direct_pay_by_user
	TradeNo          string  //  支付宝交易号	String(64)	该交易在支付宝系统中的交易流水号。最长64位。	可空	2014040311001004370000361525
	TradeStatus      string  //  交易状态	String	交易目前所处的状态。成功状态的值只有两个：TRADE_FINISHED（普通即时到账的交易成功状态）；TRADE_SUCCESS（开通了高级即时到账或机票分销产品后的交易成功状态）	可空	TRADE_FINISHED
	NotifyID         string  //  通知校验ID	String	支付宝通知校验ID，商户可以用这个流水号询问支付宝该条通知的合法性。	可空	RqPnCoPT3K9%2Fvwbh3I%2BODmZS9o4qChHwPWbaS7UMBJpUnBJlzg42y9A8gQlzU6m3fOhG
	NotifyTime       string  //  通知时间	Date	通知时间（支付宝时间）。格式为yyyy-MM-dd HH:mm:ss。	可空	2008-10-23 13:17:39
	NotifyType       string  //  通知类型	String	返回通知类型。	可空	trade_status_sync
	SellerEmail      string  //  卖家支付宝账号	String(100)	卖家支付宝账号，可以是Email或手机号码。	可空	chao.chenc1@alipay.com
	BuyerEmail       string  //  买家支付宝账号	String(100)	买家支付宝账号，可以是Email或手机号码。	可空	tstable01@alipay.com
	SellerID         string  //  卖家支付宝账户号	String(30)	卖家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	可空	2088002007018916
	BuyerID          string  //  买家支付宝账户号	String(30)	买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	可空	2088101000082594
	TotalFee         float64 //  交易金额	Number	该笔订单的资金总额，单位为RMB-Yuan。取值范围为[0.01,100000000.00]，精确到小数点后两位。	可空	10.00
	Body             string  //  商品描述	String(400)	对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body。	可空	Hello
	ExtraCommonParam string  //  	公用回传参数	String	用于商户回传参数，该值不能包含“=”、“&”等特殊字符。如果用户请求时传递了该参数，则返回给商户时会回传该参数。	可空	你好，这是测试商户的广告。
	AgentUserID      string  //  信用支付购票员的代理人ID	String	本参数用于信用支付。它代表执行支付操作的操作员账号所属的代理人的支付宝唯一用户号。以2088开头的纯16位数字。	可空	2088101000071628
}

type InstantCreditNotify struct {
	NotifyTime       string  //  通知时间	Date	通知的发送时间。格式为yyyy-MM-dd HH:mm:ss。	不可空	2009-08-12 11:08:32
	NotifyType       string  //  通知类型	String	通知的类型。	不可空	trade_status_sync
	NotifyID         string  //  通知校验ID	String	通知校验ID。	不可空	70fec0c2730b27528665af4517c27b95
	SignType         string  //  签名方式	String	DSA、RSA、MD5三个值可选，必须大写。	不可空	DSA
	Sign             string  //  签名	String	请参见签名验证。	不可空	_p_w_l_h_j0b_gd_aejia7n_ko4_m%2Fu_w_jd3_nx_s_k_mxus9_hoxg_y_r_lunli_pmma29_t_q%3D
	OutTradeNo       string  //  商户网站唯一订单号	String(64)	对应商户网站的订单系统中的唯一订单号，非支付宝交易号。需保证在商户网站中的唯一性。是请求时对应的参数，原样返回。	可空	3618810634349901
	Subject          string  //  商品名称	String(256)	商品的标题/交易标题/订单标题/订单关键字等。它在支付宝的交易明细中排在第一列，对于财务对账尤为重要。是请求时对应的参数，原样通知回来。	可空	phone手机
	PaymentType      string  //  支付类型	String(4)	取值范围请参见收款类型。	可空	1
	TradeNo          string  //  支付宝交易号	String(64)	该交易在支付宝系统中的交易流水号。最长64位。	可空	2014040311001004370000361525
	TradeStatus      string  //  交易状态	String	取值范围请参见交易状态。	可空	TRADE_FINISHED
	GMTCreate        string  //  交易创建时间	Date	该笔交易创建的时间。格式为yyyy-MM-dd HH:mm:ss。	可空	2008-10-22 20:49:31
	GMTPayment       string  //  交易付款时间	Date	该笔交易的买家付款时间。格式为yyyy-MM-dd HH:mm:ss。	可空	2008-10-22 20:49:50
	GMTClose         string  //  交易关闭时间	Date	交易关闭时间。格式为yyyy-MM-dd HH:mm:ss。	可空	2008-10-22 20:49:46
	RefundStatus     string  //  退款状态	String	取值范围请参见退款状态。	可空	REFUND_SUCCESS
	GMTRefund        string  //  退款时间	Date	卖家退款的时间，退款通知时会发送。格式为yyyy-MM-dd HH:mm:ss。	可空	2008-10-29 19:38:25
	SellerEmail      string  //  卖家支付宝账号	String(100)	卖家支付宝账号，可以是email和手机号码。	可空	chao.chenc1@alipay.com
	BuyerEmail       string  //  买家支付宝账号	String(100)	买家支付宝账号，可以是Email或手机号码。	可空	13758698870
	SellerID         string  //  卖家支付宝账户号	String(30)	卖家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	可空	2088002007018916
	BuyerID          string  //  买家支付宝账户号	String(30)	买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。	可空	2088002007013600
	Price            float64 //  商品单价	Number	如果请求时使用的是total_fee，那么price等于total_fee；如果请求时使用的是price，那么对应请求时的price参数，原样通知回来。	可空	10.00
	TotalFee         float64 //  交易金额	Number	该笔订单的总金额。请求时对应的参数，原样通知回来。	可空	10.00
	Quantity         uint    //  购买数量	Number	如果请求时使用的是total_fee，那么quantity等于1；如果请求时使用的是quantity，那么对应请求时的quantity参数，原样通知回来。	可空	1
	Body             string  //  商品描述	String(400)	该笔订单的备注、描述、明细等。对应请求时的body参数，原样通知回来。	可空	Hello
	Discount         float64 //  折扣	Number	支付宝系统会把discount的值加到交易金额上，如果需要折扣，本参数为负数。	可空	-5
	IsTotalFeeAdjust string  //  是否调整总价	String(1)	该交易是否调整过价格。	可空	N
	UseCoupon        string  //  是否使用红包买家	String(1)	是否在交易过程中使用了红包。	可空	N
	ExtraCommonParam string  //  公用回传参数	String	用于商户回传参数，该值不能包含“=”、“&”等特殊字符。如果用户请求时传递了该参数，则返回给商户时会回传该参数。	可空	你好，这是测试商户的广告。
	OutChannelType   string  //  支付渠道组合信息	String	该笔交易所使用的支付渠道。格式为：渠道1|渠道2|…，如果有多个渠道，用“|”隔开。取值范围请参见渠道类型说明与币种列表。	可空	OPTIMIZED_MOTO|BALANCE
	OutChannelAmount string  //  支付金额组合信息	String	该笔交易通过使用各支付渠道所支付的金额。格式为：金额1|金额2|…，如果有多个支付渠道，各渠道所支付金额用“|”隔开。	可空	90.00|10.00
	OutChannelInst   string  //  实际支付渠道	String	该交易支付时实际使用的银行渠道。格式为：支付渠道1|支付渠道2|…，如果有多个支付渠道，用“|”隔开。取值范围请参见实际支付渠道列表。该参数需要联系支付宝开通。	可空	ICBC
	BusinessScene    string  //  是否扫码支付	String	回传给商户此标识为qrpay时，表示对应交易为扫码支付。目前只有qrpay一种回传值。非扫码支付方式下，目前不会返回该参数。	可空	qrpay
}

func (a *Alipay) InstantCredit(outTradeNo, subject string, totalFee float64, extraParams map[string]string) (checkoutURL string, err error) {
	if outTradeNo == "" {
		err = fmt.Errorf("%s out_trade_no : Required parameter missing", LogPrefix)
		return
	}

	if subject == "" {
		err = fmt.Errorf("%s subject is required parameter", LogPrefix)
		return
	}

	if totalFee == 0 {
		err = fmt.Errorf("%s total_fee is required parameter", LogPrefix)
		return
	}

	params := a.initInstantCreditParams(outTradeNo, subject, totalFee, extraParams)
	kvs, err := GenKVpairs(instantCreditParamMap, params, "sign", "sign_type")
	if err != nil {
		return
	}

	signStr := MD5(kvs.RemoveEmpty().Sort().Join("&"), a.key)

	vals := url.Values{}
	for _, v := range kvs {
		vals.Set(v.K, v.V)
	}
	vals.Set("sign", signStr)
	vals.Set("sign_type", "MD5")

	checkoutURL = AlipayGateway + vals.Encode()
	return
}

func (a *Alipay) initInstantCreditParams(outTradeNo, subject string, totalFee float64, extraParams map[string]string) (params map[string]string) {
	params = make(map[string]string)

	params["service"] = "create_direct_pay_by_user"
	params["_input_charset"] = "utf-8"
	params["payment_type"] = "1"

	params["partner"] = a.partner
	params["key"] = a.key
	params["seller_id"] = a.partner
	params["seller_email"] = a.email

	params["out_trade_no"] = outTradeNo
	params["total_fee"] = strconv.FormatFloat(totalFee, 'f', 2, 64)
	params["subject"] = subject

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

/*
商户需要验证该通知数据中的out_trade_no是否为商户系统中创建的订单号，
并判断total_fee是否确实为该订单的实际金额（即商户订单创建时的金额），
同时需要校验通知中的seller_id（或者seller_email)
是否为out_trade_no这笔单据的对应的操作方（有的时候，
一个商户可能有多个seller_id/seller_email），
上述有任何一个验证不通过，则表明本次通知是异常通知，务必忽略。
在上述验证通过后商户必须根据支付宝不同类型的业务通知，正确的进行不同的业务处理，并且过滤重复的通知结果数据。
在支付宝的业务通知中，只有交易通知状态为TRADE_SUCCESS或TRADE_FINISHED时，支付宝才会认定为买家付款成功。
如果商户需要对同步返回的数据做验签，必须通过服务端的签名验签代码逻辑来实现。如果商户未正确处理业务通知，存在潜在的风险，商户自行承担因此而产生的所有损失。

交易状态TRADE_SUCCESS的通知触发条件是商户签约的产品支持退款功能的前提下，买家付款成功；
交易状态TRADE_FINISHED的通知触发条件是商户签约的产品不支持退款功能的前提下，买家付款成功；或者，商户签约的产品支持退款功能的前提下，交易已经成功并且已经超过可退款期限；
交易成功之后，商户（高级即时到账或机票平台商）可调用批量退款接口，系统会发送退款通知给商户，具体内容请参见批量退款接口文档；
当商户使用站内退款时，系统会发送包含refund_status和gmt_refund字段的通知给商户。
*/
func (a *Alipay) InstantCreditReturn(req *http.Request) (result *InstantCreditReturn, err error) {
	vals := req.URL.Query()
	if len(vals) == 0 {
		err = ErrReturnDataIsEmpty
		return
	}

	var fields = []string{
		"is_success",
		"sign_type",
		"sign",
		"out_trade_no",
		"subject",
		"payment_type",
		"exterface",
		"trade_no",
		"trade_status",
		"notify_id",
		"notify_time",
		"notify_type",
		"seller_email",
		"buyer_email",
		"seller_id",
		"buyer_id",
		"total_fee",
		"body",
		"extra_common_param",
		"agent_user_id",
	}

	err = a.verify(vals, fields)
	if err != nil {
		return
	}

	totalFee, _ := strconv.ParseFloat(vals.Get("total_fee"), 64)

	result = &InstantCreditReturn{
		IsSuccess:        vals.Get("is_success"),
		SignType:         vals.Get("sign_type"),
		Sign:             vals.Get("sign"),
		OutTradeNo:       vals.Get("out_trade_no"),
		Subject:          vals.Get("subject"),
		PaymentType:      vals.Get("payment_type"),
		Exterface:        vals.Get("exterface"),
		TradeNo:          vals.Get("trade_no"),
		TradeStatus:      vals.Get("trade_status"),
		NotifyID:         vals.Get("notify_id"),
		NotifyTime:       vals.Get("notify_time"),
		NotifyType:       vals.Get("notify_type"),
		SellerEmail:      vals.Get("seller_email"),
		BuyerEmail:       vals.Get("buyer_email"),
		SellerID:         vals.Get("seller_id"),
		BuyerID:          vals.Get("buyer_id"),
		TotalFee:         totalFee,
		Body:             vals.Get("body"),
		ExtraCommonParam: vals.Get("extra_common_param"),
		AgentUserID:      vals.Get("agent_user_id"),
	}
	return
}

func (a *Alipay) InstantCreditNotify(req *http.Request) (result *InstantCreditNotify, err error) {
	vals, err := parsePostData(req)
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
		"gmt_create",
		"gmt_payment",
		"gmt_close",
		"refund_status",
		"gmt_refund",
		"seller_email",
		"buyer_email",
		"seller_id",
		"buyer_id",
		"price",
		"total_fee",
		"quantity",
		"body",
		"discount",
		"is_total_fee_adjust",
		"use_coupon",
		"extra_common_param",
		"out_channel_type",
		"out_channel_amount",
		"out_channel_inst",
		"business_scene",
	}

	err = a.verify(vals, fields)
	if err != nil {
		return
	}

	var price, totalFee, discount float64
	price, _ = strconv.ParseFloat(vals.Get("price"), 64)
	totalFee, _ = strconv.ParseFloat(vals.Get("total_fee"), 64)
	discount, _ = strconv.ParseFloat(vals.Get("discount"), 64)

	var quantity uint64
	quantity, _ = strconv.ParseUint(vals.Get("quantity"), 10, 64)

	result = &InstantCreditNotify{
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
		GMTCreate:        vals.Get("gmt_create"),
		GMTPayment:       vals.Get("gmt_payment"),
		GMTClose:         vals.Get("gmt_close"),
		RefundStatus:     vals.Get("refund_status"),
		GMTRefund:        vals.Get("gmt_refund"),
		SellerEmail:      vals.Get("seller_email"),
		BuyerEmail:       vals.Get("buyer_email"),
		SellerID:         vals.Get("seller_id"),
		BuyerID:          vals.Get("buyer_id"),
		Price:            price,
		TotalFee:         totalFee,
		Quantity:         uint(quantity),
		Body:             vals.Get("body"),
		Discount:         discount,
		IsTotalFeeAdjust: vals.Get("is_total_fee_adjust"),
		UseCoupon:        vals.Get("use_coupon"),
		ExtraCommonParam: vals.Get("extra_common_param"),
		OutChannelType:   vals.Get("out_channel_type"),
		OutChannelAmount: vals.Get("out_channel_amount"),
		OutChannelInst:   vals.Get("out_channel_inst"),
		BusinessScene:    vals.Get("business_scene"),
	}

	return
}
