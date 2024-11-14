package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/gomodule/redigo/redis"
	"github.com/ut-cloud/atlas-toolkit/casbin/redis_adapter"
)

func InitCasbin(rPool *redis.Pool, perms []*RoleMenuPerm, roleKeys []string, operaPolicies []string) (model.Model, *redis_adapter.Adapter) {
	m, _ := model.NewModelFromString(ModelConf)
	_, err := rPool.Get().Do("DEL", CacheCasbin)
	if err != nil {
		panic(fmt.Sprintf("[middleware] redis pool err: %v", err))
	}
	a, err := redis_adapter.NewAdapterWithPoolAndOptions(rPool, redis_adapter.WithKey(CacheCasbin))
	policies := make([][]string, len(operaPolicies))
	for i, opera := range operaPolicies {
		policies[i] = []string{"admin", opera, "*"}
	}
	for i := range perms {
		policies = append(policies, []string{perms[i].RoleKey, fmt.Sprintf("/%s", perms[i].Perms), "*"})
	}
	err = a.AddPolicies("", "p", policies)
	//err = a.AddPolicies("p", "p",[][]string{{"api_admin", "/api.system.v1.*", "*"}})
	rolePolicies := make([][]string, len(roleKeys))
	for i := range roleKeys {
		rolePolicies[i] = []string{roleKeys[i], roleKeys[i]}
	}
	err = a.AddPolicies("", "g", rolePolicies)
	if err != nil {
		panic(fmt.Sprintf("[middleware] new redis adapter err: %s", err))
	}
	return m, a
}
