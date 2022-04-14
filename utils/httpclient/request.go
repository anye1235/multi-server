package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"
)

type ResultHandle interface {
	HandleResult(body []byte, resObj interface{}) (interface{}, error)
}

//接口
type CommonResultHandle struct {
}

type Response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Status string      `json:"status"`
}

// HandleResult对接不同厂商，返回的格式不一样 这里需要制定返回
func (p *CommonResultHandle) HandleResult(body []byte, resObj interface{}) (interface{}, error) {
	res := Response{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}
	if res.Code != 200 {
		return nil, errors.New("返回值错误")
	}
	//不需要绑定对象直接返回rest结果
	if nil == resObj {
		return res.Data, nil
	}

	if err := DecodeToObj(res.Data, resObj); nil != err {
		return nil, err
	}
	return res.Data, nil
}

// DoGet Get发起请求
func DoGet(url string, query url.Values, ctx context.Context, res interface{}, rh ResultHandle, options ...Option) (interface{}, error) {
	options = AddDefaultoptions(ctx, options...)
	body, err := Get(url, query, options...)
	return result(body, err, rh, res)
}

// DoPost post json发起请求
func DoPost(url string, raw json.RawMessage, ctx context.Context, res interface{}, rh ResultHandle, options ...Option) (interface{}, error) {
	options = AddDefaultoptions(ctx, options...)
	body, err := PostJSON(url, raw, options...)
	return result(body, err, rh, res)
}

// DoPostForm post form发起请求
func DoPostForm(url string, form url.Values, ctx context.Context, res interface{}, rh ResultHandle, options ...Option) (interface{}, error) {
	options = AddDefaultoptions(ctx, options...)
	body, err := PostForm(url, form, options...)
	return result(body, err, rh, res)
}

// 待优化 options 不能覆盖外面传过来的值
func result(body []byte, err error, rh ResultHandle, res interface{}) (interface{}, error) {
	if err != nil {
		return body, err
	}
	if nil == rh {
		rh = &CommonResultHandle{}
	}
	return rh.HandleResult(body, res)
}

//HandleResult 对接不同厂商，返回的格式不一样 这里需要制定返回
func AddDefaultoptions(ctx context.Context, options ...Option) []Option {
	if nil != ctx {
		//options = append([]Option{WithLogger(ctx.)}, options...)
		//options = append([]Option{WithTrace(ctx.Trace())}, options...)
		//options = append([]Option{WithOriginRequest(ctx.Request())}, options...)
	}
	//增加至头部，业务传入的option覆盖默认的
	options = append([]Option{WithTTL(time.Second * 10)}, options...)
	options = append(options, WithOnFailedRetry(1, time.Millisecond*1, DemoGetRetryVerify))
	//options = append(options, WithOnFailedAlarm("接口告警", new(AlarmEmail), DemoGetAlarmVerify))
	return options
}

//DecodeToObj
//运用对象返回绑定的json 结果
//返回错误
func DecodeToObj(data interface{}, resObj interface{}) error {
	dataBy, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dataBy, resObj); err != nil {
		return err
	}
	return nil
}

// DemoGetRetryVerify 设置重试规则
func DemoGetRetryVerify(body []byte) (shouldRetry bool) {
	return false
	if len(body) == 0 {
		return true
	}
	return false
}

// DemoGetAlarmVerify 设置告警规则
func DemoGetAlarmVerify(body []byte) (shouldAlarm bool) {
	return false
	if len(body) == 0 {
		return true
	}
	return false
}
