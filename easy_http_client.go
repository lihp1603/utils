package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

//创建一个http client
func HttpClientCreate() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second, //限制建立TCP连接的时间
			KeepAlive: 30 * time.Second, //
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,  //限制tls握手时间
		ResponseHeaderTimeout: 10 * time.Second, //限制读取response header的时间
		ExpectContinueTimeout: 1 * time.Second,  //限制client在发送包含expect:100-continue的header到收到继续发送body的response之间的时间等待
	}
	//创建一个http的客服端
	var netClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: netTransport,
	}
	return netClient
}

//实现http client get方法的封装
func HttpClientGet(netClient *http.Client, url string) (bodyResp []byte, statusCode int, err error) {

	//创建一个http的客服端
	// 不应该每次get都创建一个client，这样的话，会出现问题，目前出现端口被消耗殆尽的问题，不知道是不是这个原因导致的
	// 根据go的手册上说的，出于效率考虑，应该一次建立，尽量重用
	// E0528 11:52:55.846561    7144 log_trace.go:26] http client get:http://10.205.1.81:8081/jobs/get,err:Get http://10.205.1.
	// 81:8081/jobs/get: dial tcp 10.205.1.81:8081: bind: An operation on a socket could not be performed because the system la
	// cked sufficient buffer space or because a queue was full.

	// netClient := HttpClientCreate()

	// get
	response, err := netClient.Get(url)
	if err != nil {
		LogTraceE("http client get:%s,err:%s", url, err.Error())
		return
	}
	//释放资源
	defer response.Body.Close()

	statusCode = response.StatusCode

	if statusCode == http.StatusOK { //如果状态码是成功的话，读取请求的body
		bodyResp, err = ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
	} else {
		LogTraceE("http get:%s,status:%d", url, statusCode)
	}

	return
}

//实现http client post方法的封装
func HttpClientPost(netClient *http.Client, url string, postBody string) (bodyResp []byte, statusCode int, err error) {
	//创建一个http的客服端
	// netClient := HttpClientCreate()

	//post
	response, err := netClient.Post(url, "application/json", strings.NewReader(postBody))
	if err != nil {
		LogTraceE("http client post:%s,body:%s,err:%s", url, postBody, err.Error())
		return
	}
	//释放资源
	defer response.Body.Close()

	statusCode = response.StatusCode

	if statusCode == http.StatusOK { //如果状态码是成功的话，读取请求的body
		bodyResp, err = ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
	} else {
		LogTraceE("http post:%s,status:%d", url, statusCode)
	}

	return
}
