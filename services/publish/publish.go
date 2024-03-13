package publish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/utils"
)

type PublishMsg struct {
	Timestamp              string
	AzCode                 string
	ResourceId             string
	ResourceType           string
	ResourceState          string
	ErrorMsg               interface{}

	VolumeAttachment       []string      //volume attachment
	VolumeSize             float64       //volume size
	FixedIps               []string      //nova instance fixed_ips
	Host                   string        //nova instance host
	Node                   string        //nova hypervisor hostname
	ProvisionState         string        //baremetal node provision state
	Maintenance            bool          //baremetal node provision state
	FlavorId               string        //flavor id
}

type Publisher struct {
	Url            string
}

func (p *Publisher) Publish(ctx context.Context, body []byte) error {
	global.LOG.Debug(fmt.Sprintf("[req-%s]Request body %s", utils.GetRequestID(ctx), string(body)))
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(p.Url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		global.LOG.Error(fmt.Sprintf("[req-%s]Publish msg err: %v", utils.GetRequestID(ctx), err))
		return err
	}
	if resp.StatusCode() > 500 {
		global.LOG.Info(fmt.Sprintf("[req-%s]Publish msg failed: %s", utils.GetRequestID(ctx), resp.Body()))
		return errors.New(string(resp.Body()))
	}
	content := resp.Body()
	global.LOG.Info(fmt.Sprintf("[req-%s]Response body %s", utils.GetRequestID(ctx), string(content)))
	return nil
}

func (p *Publisher) Call(ctx context.Context, publishMsg *PublishMsg, buffer chan *PublishMsg) {
	body, _ := json.Marshal(publishMsg)
	if err := RetryDo(func(ctx context.Context, body []byte) error {
		return p.Publish(ctx, body)
	}, ctx, body, DefaultRetries, DefaultSleep); err != nil {
		buffer <- publishMsg
		global.LOG.Info(fmt.Sprintf("[req-%s]Push to buffer", utils.GetRequestID(ctx)))
	}
}
