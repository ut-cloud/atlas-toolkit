package utils

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosAkConf struct {
	Address     string `json:"address,omitempty"`
	Port        uint64 `json:"port,omitempty"`
	AccessKey   string `json:"accessKey,omitempty"`
	SecretKey   string `json:"secretKey,omitempty"`
	NamespaceId string `json:"namespaceId,omitempty"`
	GroupName   string `json:"groupName,omitempty"`
}

func InitAkNacos(conf *NacosAkConf) (config_client.IConfigClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(conf.Address, conf.Port),
	}
	cc := &constant.ClientConfig{
		NamespaceId:         conf.NamespaceId,
		AccessKey:           conf.AccessKey,
		SecretKey:           conf.SecretKey,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	return client, err
}
