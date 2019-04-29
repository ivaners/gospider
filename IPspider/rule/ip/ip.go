package ip

import (
	"fmt"
	"github.com/nange/gospider/spider"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	"github.com/ivaners/gospider/IPspider/model"
	"github.com/nange/gospider/web/core"
)

func init() {
	core.Register(&model.Ip{})
	spider.Register(rule)
}

var rule = &spider.TaskRule{
	Name:           "IP 池(xicidaili)",
	Description:    "抓取代理IP",
	Namespace:      "ip",
	DisableCookies: true,
	OutputFields:   []string{"ip", "port", "type", "addr"},
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error { 
			var err error
			for page := 1; page <= 2; page++ {
				if err = ctx.VisitForNext("http://www.xicidaili.com/nt/" + strconv.Itoa(page)); err != nil {
					break
				}
			}
			return err
		},
		Nodes: map[int]*spider.Node{
			0: fetachIp, 
		},
	},
}

var fetachIp = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnError: func(ctx *spider.Context, res *spider.Response, err error) error {
		log.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
		// 出错时重试三次
		return Retry(ctx, 3)
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`#ip_list tbody tr:not(:has(th))`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			ctx.PutReqContextValue("ip", el.ChildText("td:nth-child(2)"))
			ctx.PutReqContextValue("port", el.ChildText("td:nth-child(3)"))
			ctx.PutReqContextValue("type", el.ChildText("td:nth-child(4)"))
			ctx.PutReqContextValue("addr", el.ChildText("td:nth-child(6)"))
			return ctx.Output(map[int]interface{}{
				0: ctx.GetReqContextValue("ip"),
				1: ctx.GetReqContextValue("port"),
				2: ctx.GetReqContextValue("addr"),
				3: ctx.GetReqContextValue("type"),
			})
		},
	},
}

func Retry(ctx *spider.Context, count int) error {
	req := ctx.GetRequest()
	key := fmt.Sprintf("err_req_%s", req.URL.String())

	var et int
	if errCount := ctx.GetAnyReqContextValue(key); errCount != nil {
		et = errCount.(int)
		if et >= count {
			return fmt.Errorf("exceed %d counts", count)
		}
	}
	log.Infof("errCount:%d, we wil retry url:%s, after 1 second", et+1, req.URL.String())
	time.Sleep(time.Second)
	ctx.PutReqContextValue(key, et+1)
	ctx.Retry()

	return nil
}
