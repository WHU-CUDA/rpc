// 服务发现
package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int
// 负载均衡模式
const (
	RandomSelect SelectMode = iota
	RoundRobinSelect
)
// 服务发现接口
type Discovery interface {
	// 刷新服务发现
	Refresh() error
	// 更新服务列表
	Update(servers []string) error
	// 获取服务
	Get(mode SelectMode) (string, error)
	// 获取全部服务
	GetAll() ([]string, error)
}

// 服务发现结构体(类)
type MultiServerDiscovery struct {
	r       *rand.Rand
	// 读写锁
	mu      sync.RWMutex
	// 服务列表
	servers []string
	// 索引
	index   int
}

func (d *MultiServerDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Lock()
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}
	log.Println("rpc registry: refresh servers from registry", d.registry)
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry: refresh err: ", err)
		return err
	}
	servers := strings.Split(resp.Header.Get("X-Geerpc-Servers"), ",")
	d.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		if strings.TrimpSpace(server) != "" {
			d.servers = append(d.servers, strings.TrimpSpace(server))
		}
	}
	d.lastUpdate = time.Now()
	return nil
}

func (d *MultiServerDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d *MultiServerDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// 服务列表长度
	n := len(d.servers)
	// 没有服务，即服务列表为0
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	// 负载均衡
	switch mode {
		// 随机选择
	case RandomSelect:
		return d.servers[d.r.Intn(n)], nil
		// 轮询
	case RoundRobinSelect:
		s := d.servers[d.index%n]
		d.index = (d.index + 1) % n
		return s, nil
		// .... 更多负载均衡方式
	default:
		// 不支持的负载均衡方式，返回错误
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

func (d *MultiServerDiscovery) GetAll() ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

func NewMultiServerDiscovery(servers []string) *MultiServerDiscovery {
	d := &MultiServerDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

var _ Discovery = (*MultiServerDiscovery)(nil)
