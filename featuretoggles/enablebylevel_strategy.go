package featuretoggles

import (
	unleashcontext "github.com/Unleash/unleash-client-go/context"
	"github.com/fabric8-services/fabric8-auth/log"
)

const (
	EnableByLevel string = "enableByLevel"
	Level         string = "level"
)

// EnableByLevelStrategy the strategy to roll out a feature if the user opted-in for a compatible level of features
type EnableByLevelStrategy struct {
}

// Name the name of the stragegy. Must match the name on the Unleash server.
func (s *EnableByLevelStrategy) Name() string {
	return EnableByLevel
}

// IsEnabled returns `true` if the given context is compatible with the settings configured on the Unleash server
func (s *EnableByLevelStrategy) IsEnabled(settings map[string]interface{}, ctx *unleashcontext.Context) bool {
	log.Debug(nil, map[string]interface{}{"settings_group_id": settings[Level], "properties_group_id": ctx.Properties[Level]}, "checking if feature is enabled for user, based on his/her group...")
	return settings[Level] == ctx.Properties[Level]
}