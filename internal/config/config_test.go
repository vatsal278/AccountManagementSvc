package config

import (
	"github.com/PereRohit/util/config"
	"github.com/PereRohit/util/testutil"
	"testing"
)

func TestInitSvcConfig(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args func() args
		want func(args) *SvcConfig
	}{
		{
			name: "Success",
			args: func() args {
				return args{
					cfg: Config{
						ServiceRouteVersion: "v2",
						ServerConfig:        config.ServerConfig{},
						DataBase: DbCfg{
							Driver: "mysql",
						},
					},
				}
			},
			want: func(arg args) *SvcConfig {
				required := &SvcConfig{
					Cfg: &Config{
						ServiceRouteVersion: "v2",
						ServerConfig:        config.ServerConfig{},
						DataBase: DbCfg{
							Driver: "mysql",
						},
					},
					ServiceRouteVersion: "v2",
					SvrCfg:              config.ServerConfig{},
				}
				return required
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.args()
			got := InitSvcConfig(s.cfg)

			diff := testutil.Diff(got, tt.want(s))
			if diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		})
	}
}
