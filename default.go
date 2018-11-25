package motan

import (
	"net/http"
	"sync"

	motan "github.com/weibocom/motan-go/core"
	"github.com/weibocom/motan-go/endpoint"
	"github.com/weibocom/motan-go/filter"
	"github.com/weibocom/motan-go/ha"
	"github.com/weibocom/motan-go/lb"
	"github.com/weibocom/motan-go/provider"
	"github.com/weibocom/motan-go/registry"
	"github.com/weibocom/motan-go/serialize"
	"github.com/weibocom/motan-go/server"
)

var (
	extOnce           sync.Once
	handlerOnce       sync.Once
	defaultExtFactory *motan.DefaultExtensionFactory
	// all default manage handlers
	defaultManageHandlers map[string]http.Handler
	// PermissionCheck is default permission check for manage request
	PermissionCheck = NoPermissionCheck
)

type PermissionCheckFunc func(r *http.Request) bool

func NoPermissionCheck(r *http.Request) bool {
	return true
}

func GetDefaultManageHandlers() map[string]http.Handler {
	handlerOnce.Do(func() {
		defaultManageHandlers = make(map[string]http.Handler, 16)

		status := &StatusHandler{}
		defaultManageHandlers["/"] = status
		defaultManageHandlers["/200"] = status
		defaultManageHandlers["/503"] = status
		defaultManageHandlers["/version"] = status

		info := &InfoHandler{}
		defaultManageHandlers["/getConfig"] = info
		defaultManageHandlers["/getReferService"] = info

		debug := &DebugHandler{}
		defaultManageHandlers["/debug/pprof/"] = debug
		defaultManageHandlers["/debug/pprof/cmdline"] = debug
		defaultManageHandlers["/debug/pprof/profile"] = debug
		defaultManageHandlers["/debug/pprof/symbol"] = debug
		defaultManageHandlers["/debug/pprof/trace"] = debug
		defaultManageHandlers["/debug/mesh/trace"] = debug
		defaultManageHandlers["/debug/pprof/sw"] = debug

		switcher := &SwitcherHandle{}
		defaultManageHandlers["/switcher/set"] = switcher
		defaultManageHandlers["/switcher/get"] = switcher
		defaultManageHandlers["/switcher/getAll"] = switcher
	})
	return defaultManageHandlers
}

func GetDefaultExtFactory() motan.ExtensionFactory {
	// 初始化默认的扩展
	extOnce.Do(func() {
		// 初始化对象
		defaultExtFactory = &motan.DefaultExtensionFactory{}
		defaultExtFactory.Initialize()

		// 为扩展对象添加扩展
		AddDefaultExt(defaultExtFactory)
	})
	return defaultExtFactory
}

func AddDefaultExt(d motan.ExtensionFactory) {

	// all default extension
	// 注册过滤器
	filter.RegistDefaultFilters(d)
	// 注册ha
	ha.RegistDefaultHa(d)
	// 注册负载均衡
	lb.RegistDefaultLb(d)
	// 注册end point
	endpoint.RegistDefaultEndpoint(d)
	// 注册provider
	provider.RegistDefaultProvider(d)
	// 注册服务注册服务
	registry.RegistDefaultRegistry(d)
	// 注册服务
	server.RegistDefaultServers(d)
	// TODO
	server.RegistDefaultMessageHandlers(d)
	// 注册序列化方式
	serialize.RegistDefaultSerializations(d)
}
