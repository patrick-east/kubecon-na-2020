package api

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/util"
	"google.golang.org/grpc"
)

const PluginName = "grpc_api"

type Config struct {
	Listen string `json:"listen"`
}

type Server struct {
	UnimplementedAuthorizerServer
	manager *plugins.Manager
	mtx     sync.Mutex
	config  Config
	grpcServer  *grpc.Server
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.config.Listen)
	if err != nil {
		s.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	RegisterAuthorizerServer(s.grpcServer, s)

	s.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateOK})

	go func() {
		err = s.grpcServer.Serve(lis)
		if err != nil {
			s.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
		} else {
			s.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	s.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})
	s.grpcServer.GracefulStop()
}

func (s *Server) Reconfigure(ctx context.Context, config interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.config = config.(Config)
	s.Stop(ctx)
	s.Start(ctx)
}

func (s *Server) Authz(ctx context.Context, req *AuthzRequest) (*AuthzResponse, error) {
	resp := &AuthzResponse{
		Allow: false,
	}

	input := map[string]interface{}{
		"jwt": req.Jwt,
	}

	pq, err := rego.New(
		rego.Query("data.authz.allow"),

		// Policies?? External Data???

	).PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	rs, err := pq.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	if len(rs) > 0 {
		resp.Allow = rs[0].Expressions[0].Value.(bool)
	}

	// Decision logs??

	return resp, nil
}

type Factory struct{}

func (Factory) New(m *plugins.Manager, config interface{}) plugins.Plugin {

	m.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})

	return &Server{
		manager: m,
		config:  config.(Config),
	}
}

func (Factory) Validate(_ *plugins.Manager, config []byte) (interface{}, error) {
	parsedConfig := Config{}
	err := util.Unmarshal(config, &parsedConfig)
	return parsedConfig, err
}