package initiator

import (
	"log"

	"github.com/casbin/casbin/v2"
)

func CasbinEnforcer(modelPath, policyPath string) casbin.IEnforcer {

	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatalf("failed to create casbin policy enforcer %v", err)
	}
	return e
}
