package kubeproxy

import (
	"github.com/go-redis/redis/v8"
)

type Manager interface{
	// Add a rule to iptable
	AddRule(msg *redis.Message) error

	// Delete a rule from iptable
	DelRule(msg *redis.Message) error
}

type KubeproxyManager struct {

}

func New() KubeproxyManager {
	return KubeproxyManager{}
}

func (km *KubeproxyManager) AddRule(msg *redis.Message) error {
	
}

func (km *KubeproxyManager) DelRule(msg *redis.Message) error {

}


