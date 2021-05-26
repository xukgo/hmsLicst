package hmsLicenter

import (
	"fmt"
	"github.com/flyaways/pool"
	"hmsList/rpcService"
	"sync/atomic"
	"time"
)

type Repo struct {
	licHost  string //ip eg:lic.vinck.cn
	connPool *pool.GRPCPool

	appName   string
	startUp   int64
	localUid  string
	buildTime time.Time
	version   string

	cloudEnable int32
	cloudQps    int32
	cloudExpire int64
	ipList      [][]byte

	timeout     time.Duration
	retry       int
	encryKey    []byte
	authUpdater func(response *rpcService.AuthResponse)
}

func NewRepo(licHost string, appName string, startTs int64, localid string, buildTime time.Time, version string) *Repo {
	repo := new(Repo)
	repo.licHost = licHost
	repo.appName = appName
	repo.startUp = startTs
	repo.localUid = localid
	repo.buildTime = buildTime
	repo.version = version

	repo.cloudEnable = 0
	repo.cloudQps = 0
	repo.cloudExpire = 0

	repo.timeout = time.Second * 3
	repo.retry = 1

	addr := repo.getGrpcAddr()
	repo.connPool, _ = initPool(addr)
	return repo
}

func (this *Repo) WithAuthUpdate(authUpdater func(response *rpcService.AuthResponse)) *Repo {
	this.authUpdater = authUpdater
	return this
}
func (this *Repo) WithEncryptKey(encryKey []byte) *Repo {
	this.encryKey = encryKey
	return this
}

func (this *Repo) WithTimeCtrl(timeout time.Duration, retry int) *Repo {
	this.timeout = timeout
	this.retry = retry
	return this
}

func (this *Repo) getGrpcAddr() string {
	return fmt.Sprintf("%s:9001", this.licHost)
}
func (this *Repo) getRemoteIPUrl() string {
	return fmt.Sprintf("http://%s:9002/remoteAddr", this.licHost)
}

func (this *Repo) GetAppName() string {
	return this.appName
}
func (this *Repo) GetLocalUid() string {
	return this.localUid
}
func (this *Repo) GetStartUp() int64 {
	return this.startUp
}

func (this *Repo) GetEnable() bool {
	res := atomic.LoadInt32(&this.cloudEnable)
	return res == 1
}
func (this *Repo) GetQps() int {
	res := atomic.LoadInt32(&this.cloudQps)
	return int(res)
}
func (this *Repo) GetExpireTime() time.Time {
	t := atomic.LoadInt64(&this.cloudExpire)
	return time.Unix(t, 0)
}

func (this *Repo) Start() error {
	var err error
	this.ipList = this.getIpList()
	err = this.register()
	go this.loopAuth()
	if err != nil {
		return err
	}
	return nil
}

func (this *Repo) loopAuth() {
	for {
		this.ipList = this.getIpList()
		err := this.register()
		if err != nil {
			atomic.StoreInt32(&this.cloudEnable, 0)
			time.Sleep(time.Second * 5)
			continue
		}

		for {
			err = this.auth()
			if err != nil {
				atomic.StoreInt32(&this.cloudEnable, 0)
				time.Sleep(time.Second * 5)
				break
			}
			time.Sleep(time.Second * 30)
		}
	}
}
