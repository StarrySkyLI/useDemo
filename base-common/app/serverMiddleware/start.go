package serverMiddleware

import (
	"github.com/zeromicro/go-zero/rest"
	"gitlab.coolgame.world/go-template/base-common/middleware"
)

type SMOption func(s *ServerMiddleware)

func WithWhiteHeaderPathSMOption(whiteHeader map[string]int) SMOption {
	return func(s *ServerMiddleware) {
		s.whiteHeader = whiteHeader
	}
}

func WithPlatformBusinessFuncSMOption(fun middleware.GetPlatformBusinessFunc) SMOption {
	return func(s *ServerMiddleware) {
		s.getPlatformBusinessFunc = fun
	}
}

func WithDebugOption() SMOption {
	return func(s *ServerMiddleware) {
		s.isDebug = true
	}
}

func WithTestOption() SMOption {
	return func(s *ServerMiddleware) {
		s.isTest = true
	}
}

func WithCheckTokenHandleSMOption(fun middleware.CheckRequestTokenFunc) SMOption {
	return func(s *ServerMiddleware) {
		s.checkTokenHandle = fun
	}
}

type ServerMiddleware struct {
	whiteHeader             map[string]int
	getPlatformBusinessFunc middleware.GetPlatformBusinessFunc
	checkTokenHandle        middleware.CheckRequestTokenFunc

	Server *rest.Server

	isDebug bool
	isTest  bool
}

func NewServerMiddleware(s *rest.Server, opt ...SMOption) *ServerMiddleware {
	res := &ServerMiddleware{
		Server: s,
	}

	for _, item := range opt {
		item(res)
	}

	return res
}

func (s *ServerMiddleware) ApiUseMiddleware() {
	// ------------- cant edit sort -----------------------
	s.useApiHeaderMiddleware()
	s.platformBusinessMiddleware()
	// ------------- cant edit sort end -----------------------

	s.signVerifyMiddleware()
	s.mustUserAgentMiddleware()
	s.apiRequestDecryptMiddleware()
}

func (s *ServerMiddleware) useApiHeaderMiddleware() {
	var apiHeaderOption = []middleware.ApiHeadOption{
		middleware.CloseVerifyOption(s.whiteHeader),
	}
	if s.isDebug {
		apiHeaderOption = append(apiHeaderOption, middleware.WithDebugOption())
	}
	s.Server.Use(middleware.NewApiHeaderMiddleware(
		apiHeaderOption...,
	).Handle)
}

func (s *ServerMiddleware) platformBusinessMiddleware() {
	var option = []middleware.PlatformBusinessMiddlewareOption{
		middleware.WithGetPlatformBusinessFuncOption(
			s.getPlatformBusinessFunc,
		),
	}
	if s.isDebug {
		option = append(option, middleware.WithPlatformBusinessDebugOption())
	}

	s.Server.Use(middleware.NewPlatformBusinessMiddleware(
		option...,
	).Handle)
}

func (s *ServerMiddleware) signVerifyMiddleware() {
	var key string = middleware.ProdSignKey

	if s.isTest {
		key = middleware.TestSignKey
	}

	s.Server.Use(middleware.NewSignVerifyMiddleware(key, s.isDebug == false, s.whiteHeader).Handle)

}

func (s *ServerMiddleware) mustUserAgentMiddleware() {
	if s.checkTokenHandle == nil {
		panic("must use CheckTokenHandleSMOption.")
	}

	s.Server.Use(middleware.NewUserAgentMiddleware(
		middleware.WithCheckOption(s.checkTokenHandle),
	).Handle)
}

func (s *ServerMiddleware) apiRequestDecryptMiddleware() {
	var key []byte = []byte(middleware.ProdDecryptKey)
	if s.isTest {
		key = []byte(middleware.TestDecryptKey)
	}

	s.Server.Use(middleware.NewApiRequestDecryptMiddleware(key, s.isDebug == false).Handle)
}
