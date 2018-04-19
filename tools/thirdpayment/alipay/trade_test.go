package alipay

import (
	"fmt"
	"testing"
	"net/http"
	"strings"
	"io/ioutil"
)

func TestAliPay_TradeQuery(t *testing.T) {
	fmt.Println("========== TradeQuery ==========")
	type arg struct {
		outTradeNo string
		wanted     error
		name       string
	}

	testCaes := []arg{
		{"trade_no_20170623022111", nil, "query success"},
		//TODO:add more test case
	}

	for _, tc := range testCaes {
		req := AliPayTradeQuery{
			OutTradeNo: tc.outTradeNo,
		}
		resp, err := client.TradeQuery(req)
		if err != tc.wanted {
			t.Errorf("%s input:%s wanted:%v get:%v", tc.name, tc.outTradeNo, tc.wanted, err)
		} else {
			t.Log(resp)
		}
	}
}

func TestAliPay_TradeAppPay(t *testing.T) {
	fmt.Println("========== TradeAppPay ==========")
	var p = AliPayTradeAppPay{}
	p.NotifyURL = "http://203.86.24.181:3000/alipay"
	p.Body = "body"
	p.Subject = "商品标题"
	p.OutTradeNo = "01010101"
	p.TotalAmount = "100.00"
	p.ProductCode = "p_1010101"
	fmt.Println(client.TradeAppPay(p))
}

func TestAliPay_TradePagePay(t *testing.T) {
	fmt.Println("========== TradePagePay ==========")
	var p = AliPayTradePagePay{}
	p.NotifyURL = "http://220.112.233.229:3000/alipay"
	p.ReturnURL = "http://220.112.233.229:3000"
	p.Subject = "修正了中文的 Bug"
	p.OutTradeNo = "trade_no_20170623011112"
	p.TotalAmount = "10.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	fmt.Println(client.TradePagePay(p))
}

func TestAliPay_TradePreCreate(t *testing.T) {
	fmt.Println("========== TradePreCreate ==========")
	var p = AliPayTradePreCreate{}
	p.OutTradeNo = "no_0001"
	p.Subject = "测试订单"
	p.TotalAmount = "10.10"

	fmt.Println(client.TradePreCreate(p))
}

func TestClient(t *testing.T) {
	//	str := `https://openapi.alipaydev.com/gateway.do?app_id=2088721993377622&biz_content={"body":"","subject":"超会玩支付","out_trade_no":"7RwRPmS1Qpj0","timeout_express:"15m","total_amount":"0.01","seller_id":"","product_code":""}&charset=utf-8&format=JSON&method=alipay.trade.app.pay&notify_url=http://sdk.abchwan.com/v1/aliPayNotify&sign=FPJCDOymS0a3CrF9uPgV6y8IOVAQ6BnJjOIkLibRcNewgRsa
	//3sLHKF/A0x2nXOa8AqgQnoQ1MVQnT6d7VTTdUT6A0SPM6rXtudMrWpKLyitPfNS+7mrHTWA06/PDUSuK26/43q1q816zrajwuZCfa3VtMAWr4BCx0mrfVnUjkSxw6D70KHOe1n9Ug0EyAVBsV4JGPrgv1I+/N2B21A5KpQH4RLGdsatkdxHqe7UshV2ePSkdYgCnFKo8R7YnjHJhkUOyGlMlN80JEI
	//WW0cVSDyLs2EI/gkZOzsCNCeyWJoDXgVFmyhLwxrtehpsYm56vFfw6fIOPLheQj4mAV0F8MA==&sign_type=RSA2&timestamp=2018-04-19 14:06:05&version=1.0`
	_str := `app_id=2018041602567981&biz_content={"notify_url":"http://sdk.abchwan.com/v1/aliPayNotify","body":"超会玩支付","subject":"游戏","out_tr
ade_no":"9njZ99E37vhC","timeout_express":"30m","total_amount":"0.10","seller_id":"","product_code":"QUICK_MSECURITY_PAY"}&charset=utf-8&format=JSON&method=alipay.trade.app.pay&
notify_url=http://sdk.abchwan.com/v1/aliPayNotify&sign=2rDGTjbY9hbvqbAAcZxeDr5HMdDXNE5fqwfGuAISE9hmy3uJDd4L7/DxlynGtb22gaqqI3rckf1tUTXbyyggdBB2oHwKm5jD+KN60HR4cvh4tEUCK4PqdqBo20bSINOMu8Y7T23mtcDFQ/LH+pVfEK9t75wjrSR9YVyH7QSpKL/bq4pG1gmN2uPeZ8Og0xOHGshXwq0YhGbx5A/wgUc/HXY0aWIGmdCgelaMFlWDEhEHpm9xRIvHsgiNdBGQgCMDoIGJjYtI7X7xMt/dCyets+vMwaXZodiFYc4TdwUREfZuoa4h+PSpw+r9grwgkB96qjo5KgT2/8rb/0nXsnvVyA==&sign_type=RSA2&timestamp=2018-04-19 16:49:07&version=1.0`
	buf := strings.NewReader(_str)
	req, err := http.NewRequest("POST", "https://openapi.alipaydev.com/gateway.do", buf)
	c := http.Client{}
	resp, err := c.Do(req)
	fmt.Println(resp.StatusCode, err)
	fmt.Println(resp, err)
	bts, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(bts), err)
	//data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("-----------", string(data), err)

}

func index(w http.ResponseWriter, r *http.Request) {
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "ok")
}

func main() {
	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/v1/aliPayNotify", index)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
