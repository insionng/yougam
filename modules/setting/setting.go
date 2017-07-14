package setting

import (
	"github.com/insionng/yougam/libraries/karlseguin/ccache"
)

const (
	FlashPolicyService = true

	CreateNodesOfGoldCoins = -1000
)

var Cache = ccache.New(ccache.Configure())
