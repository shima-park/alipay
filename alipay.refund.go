package alipay

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RefundDetailData struct {
	AlipayTransID string
	Amount        float64
	RefundReason  string
}

func filterRefundReason(s string) string {
	s = strings.Replace(s, "^", "", -1)
	s = strings.Replace(s, "|", "", -1)
	s = strings.Replace(s, "$", "", -1)
	s = strings.Replace(s, "#", "", -1)
	return s
}

/*
service	接口名称	String	接口名称。	不可空	refund_fastpay_by_platform_pwd
partner	合作者身份ID	String(16)	签约的支付宝账号对应的支付宝唯一用户号。以2088开头的16位纯数字组成。	不可空	2088101008267254
_input_charset	参数编码字符集	String	商户网站使用的编码格式，如utf-8、gbk、gb2312等。	不可空	GBK
sign_type	签名方式	String	DSA、RSA、MD5三个值可选，必须大写。	不可空	MD5
sign	签名	String	请参见签名。	不可空	tphoyf4aoio5e6zxoaydjevem2c1s1zo
notify_url	服务器异步通知页面路径	String(200)	支付宝服务器主动通知商户网站里指定的页面http路径。	可空	http://api.test.alipay.net/atinterface/receive_notify.htm
seller_email	卖家支付宝账号	String	如果卖家Id已填，则此字段可为空。	不可空	Jier1105@alitest.com
seller_user_id	卖家用户ID	String	卖家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。登录时，seller_email和seller_user_id两者必填一个。如果两者都填，以seller_user_id为准。	不可空	2088101008267254
refund_date	退款请求时间	String	退款请求的当前时间。格式为：yyyy-MM-dd hh:mm:ss。	不可空	2011-01-12 11:21:00
batch_no	退款批次号	String	每进行一次即时到账批量退款，都需要提供一个批次号，通过该批次号可以查询这一批次的退款交易记录，对于每一个合作伙伴，传递的每一个批次号都必须保证唯一性。格式为：退款日期（8位）+流水号（3～24位）。不可重复，且退款日期必须是当天日期。流水号可以接受数字或英文字符，建议使用数字，但不可接受“000”。	不可空	201101120001
batch_num	总笔数	String	即参数detail_data的值中，“#”字符出现的数量加1，最大支持1000笔（即“#”字符出现的最大数量为999个）。	不可空	1
detail_data	单笔数据集	String	退款请求的明细数据。格式详情参见下面的“单笔数据集参数说明”。	不可空 2011011201037066^5.00^协商退款
*/
func (a *Alipay) Refund(batchNo string, detailDatas []RefundDetailData, notifyURL string) (refundURL string, err error) {
	if batchNo == "" {
		err = fmt.Errorf("%s batch_no:Required parameter missing", LogPrefix)
		return
	}

	if len(detailDatas) == 0 {
		err = fmt.Errorf("%s detail_data:Required parameter missing", LogPrefix)
		return
	}

	var datas []string
	for _, v := range detailDatas {
		reason := filterRefundReason(v.RefundReason)
		datas = append(datas, fmt.Sprintf("%s^%.2f^%s", v.AlipayTransID, v.Amount, reason))
	}

	kvs := KVpairs{}
	kvs = append(kvs, KVpair{K: "service", V: "refund_fastpay_by_platform_pwd"})
	kvs = append(kvs, KVpair{K: "partner", V: a.partner})
	kvs = append(kvs, KVpair{K: "_input_charset", V: "utf-8"})
	kvs = append(kvs, KVpair{K: "notify_url", V: notifyURL})
	kvs = append(kvs, KVpair{K: "seller_email", V: a.email})
	kvs = append(kvs, KVpair{K: "seller_user_id", V: a.partner})
	kvs = append(kvs, KVpair{K: "refund_date", V: time.Now().Format("2006-01-02 15:04:05")})
	kvs = append(kvs, KVpair{K: "batch_no", V: batchNo})
	kvs = append(kvs, KVpair{K: "batch_num", V: fmt.Sprint(len(detailDatas))})
	kvs = append(kvs, KVpair{K: "detail_data", V: strings.Join(datas, "#")})

	signStr := MD5(kvs.RemoveEmpty().Sort().Join("&"), a.key)

	kvs = append(kvs, KVpair{K: "sign", V: signStr})
	kvs = append(kvs, KVpair{K: "sign_type", V: "MD5"})

	vals := url.Values{}
	for _, v := range kvs {
		vals.Set(v.K, v.V)
	}

	refundURL = AlipayGateway + vals.Encode()
	return
}

type RefundNotifyResult struct {
	NotifyTime    string //  通知时间	Date	通知发送的时间。格式为：yyyy-MM-dd HH:mm:ss。	不可空	2009-08-12 11:08:32
	NotifyType    string //  通知类型	String	通知的类型。	不可空	batch_refund_notify
	NotifyID      string //  通知校验ID	String	通知校验ID。	不可空	70fec0c2730b27528665af4517c27b95
	SignType      string //  签名方式	String	DSA、RSA、MD5三个值可选，必须大写。	不可空	MD5
	Sign          string //  签名            String	请参见签名验证。	不可空	b7baf9af3c91b37bef4261849aa76281
	BatchNo       string //  退款批次号	String	原请求退款批次号。	不可空	20060702001
	SuccessNum    string //  退款成功总数	String	退交易成功的笔数。0<= success_num<= 总退款笔数。	不可空	2
	ResultDetails string //  退款结果明细	String	退款结果明细：退手续费结果返回格式：交易号^退款金额^处理结果\$退费账号^退费账户ID^退费金额^处理结果；不退手续费结果返回格式：交易号^退款金额^处理结果。若退款申请提交成功，处理结果会返回“SUCCESS”。若提交失败，退款的处理结果中会有报错码，参见即时到账批量退款业务错误码。	可空	2010031906272929^80^SUCCESS$jax_chuanhang@alipay.com^2088101003147483^0.01^SUCCESS
}

func (a *Alipay) RefundNotify(req *http.Request) (rnr *RefundNotifyResult, err error) {
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
		"batch_no",
		"success_num",
		"result_details",
	}

	err = a.verify(vals, fields)
	if err != nil {
		return
	}

	rnr = &RefundNotifyResult{
		NotifyTime:    vals.Get("notify_time"),
		NotifyType:    vals.Get("notify_type"),
		NotifyID:      vals.Get("notify_id"),
		SignType:      vals.Get("sign_type"),
		Sign:          vals.Get("sign"),
		BatchNo:       vals.Get("batch_no"),
		SuccessNum:    vals.Get("success_num"),
		ResultDetails: vals.Get("result_details"),
	}

	return
}
