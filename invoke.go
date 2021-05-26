package hmsLicenter

import (
	"bytes"
	"context"
	"hmsList/httpClient"
	"hmsList/netUtil"
	"hmsList/rpcService"
	"net"
	"sync/atomic"
	"time"
)

func (this *Repo) register() error{
	var err error
	for i := 0; i <= this.retry; i++ {
		err = this.innerRegister()
		if err != nil {
			continue
		}
		return nil
	}
	return err
}

func (this *Repo) auth() error{
	var err error
	for i := 0; i <= this.retry; i++ {
		err = this.innerAuth()
		if err != nil {
			continue
		}
		return nil
	}
	return err
}

type StringPair struct {
	Key string
	Value string
}

func (this *Repo) Report(action string, arr []*StringPair) (*rpcService.BasicResponse,error){
	var resp *rpcService.BasicResponse
	var err error

	barr := make([]*rpcService.BytesPair,0,len(arr))
	for _,v := range arr{
		barr = append(barr, &rpcService.BytesPair{
			Key: []byte(v.Key),
			Value: []byte(v.Value),
		})
	}
	for i := 0; i <= this.retry; i++ {
		resp, err = this.innerReport(action, barr)
		if err == nil {
			return resp,err
		}
		continue
	}
	return nil,err
}

func (this *Repo) innerRegister() error {
	conn, err := this.connPool.Get()
	if err != nil {
		return err
	}
	defer this.connPool.Put(conn)

	client := rpcService.NewServiceClient(conn)
	rootCtx := context.Background()

	req := new(rpcService.RegisterRequest)
	req.Id = []byte(this.appName)
	req.Pid = []byte(this.localUid)
	req.Timestamp = time.Now().UnixNano() / 1000000
	req.BuildTime = this.buildTime.UnixNano()
	req.Version = []byte(this.version)
	req.LocalIPs = this.ipList
	req.IsEncrypt = len(this.encryKey) > 0
	req.Encrypt(this.encryKey)

	timeCtx, _ := context.WithTimeout(rootCtx, this.timeout)
	resp, err := client.Register(timeCtx, req)
	if err != nil {
		return err
	}

	resp.Decrypt(this.encryKey)
	this.parseAuthResponse(resp)
	return nil
}

func (this *Repo) innerAuth() error {
	conn, err := this.connPool.Get()
	if err != nil {
		//atomic.StoreInt32(&this.cloudEnable, 0)
		return err
	}
	defer this.connPool.Put(conn)

	client := rpcService.NewServiceClient(conn)
	rootCtx := context.Background()

	request := new(rpcService.AuthRequest)
	request.Id = []byte(this.appName)
	request.Pid = []byte(this.localUid)
	request.Timestamp = time.Now().UnixNano() / 1000000
	request.IsEncrypt = len(this.encryKey) > 0
	request.Encrypt(this.encryKey)

	tctx, _ := context.WithTimeout(rootCtx, this.timeout)
	resp, err := client.GetAuth(tctx, request)
	if err != nil {
		return err
	}

	resp.Decrypt(this.encryKey)
	this.parseAuthResponse(resp)
	return nil
}

func (this *Repo) parseAuthResponse(resp *rpcService.AuthResponse) {
	if resp.Enable {
		atomic.StoreInt32(&this.cloudEnable, 1)
	} else {
		atomic.StoreInt32(&this.cloudEnable, 0)
	}

	atomic.StoreInt32(&this.cloudQps, resp.Qps)
	atomic.StoreInt64(&this.cloudExpire, resp.ExpireTime)
	if this.authUpdater != nil{
		this.authUpdater(resp)
	}
}

func (this *Repo) innerReport(action string, arr []*rpcService.BytesPair) (*rpcService.BasicResponse,error){
	conn, err := this.connPool.Get()
	if err != nil {
		return nil,err
	}
	defer this.connPool.Put(conn)

	client := rpcService.NewServiceClient(conn)
	rootCtx := context.Background()

	req := new(rpcService.ReportRequest)
	req.IsEncrypt = len(this.encryKey) > 0
	req.Id = []byte(this.appName)
	req.Pid = []byte(this.localUid)
	req.Timestamp = time.Now().UnixNano() / 1000000
	req.Action = []byte(action)
	req.Attrs =arr
	req.Encrypt(this.encryKey)

	timeCtx, _ := context.WithTimeout(rootCtx, this.timeout)
	resp, err := client.Report(timeCtx, req)
	if err != nil {
		return nil,err
	}

	return resp,nil
}

func (this *Repo) getIpList() [][]byte {
	privateIps, _ := netUtil.GetPrivateIPList()
	publicIps, _ := netUtil.GetPublicIPList()
	list := make([][]byte, 0, len(privateIps)+len(publicIps)+1)

	httpcli := httpClient.NewHttpClient(this.getRemoteIPUrl(), this.timeout, false)
	httpResp := httpcli.Get()
	if httpResp.Error == nil && httpResp.StatusCode == 200 {
		index := bytes.Index(httpResp.Data, []byte(":"))
		if index > 0 {
			ip := string(httpResp.Data[:index])
			exist := false
			for idx := range publicIps {
				if publicIps[idx].String() == ip {
					exist = true
					break
				}
			}
			if !exist {
				list = append(list, net.ParseIP(ip))
			}
		}
	}

	for idx := range publicIps {
		list = append(list, publicIps[idx])
	}
	for idx := range privateIps {
		list = append(list, privateIps[idx])
	}
	return list
}

//func parseBuildTime(str string) time.Time {
//	str = strings.ReplaceAll(str, ".", "")
//	loc, _ := time.LoadLocation("Asia/Shanghai")
//	t, _ := time.ParseInLocation("20060102150405", str, loc)
//	return t
//}

//func getIP(host string) string {
//	addr, err := net.ResolveIPAddr("ip", host)
//	if err != nil {
//		return ""
//	}
//
//	return addr.String()
//}
