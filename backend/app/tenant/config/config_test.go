package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"

	adminconfig "zeebee-backend/app/admin/config"
	borrowerconfig "zeebee-backend/app/borrower/config"
	lenderconfig "zeebee-backend/app/lender/config"
)

func TestDeploymentJwtConfigMatchesClients(t *testing.T) {
	envs := []string{"staging", "prod"}
	if _, err := os.Stat(filepath.Join("..", "etc", "config-local.yaml")); err == nil {
		envs = append(envs, "local")
	}

	for _, env := range envs {
		t.Run(env, func(t *testing.T) {
			var tenant Config
			var admin adminconfig.Config
			var lender lenderconfig.Config
			var borrower borrowerconfig.Config

			conf.MustLoad(filepath.Join("..", "etc", "config-"+env+".yaml"), &tenant)
			conf.MustLoad(filepath.Join("..", "..", "admin", "etc", "config-"+env+".yaml"), &admin)
			conf.MustLoad(filepath.Join("..", "..", "lender", "etc", "config-"+env+".yaml"), &lender)
			conf.MustLoad(filepath.Join("..", "..", "borrower", "etc", "config-"+env+".yaml"), &borrower)

			require.Equal(t, tenant.RpcConf.JwtConf, admin.RpcConf.JwtConf)
			require.Equal(t, tenant.RpcConf.JwtConf, lender.RpcConf.JwtConf)
			require.Equal(t, tenant.RpcConf.JwtConf, borrower.RpcConf.JwtConf)
		})
	}
}
