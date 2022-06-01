package profiles

import (
	cloudify_v5_0_5 "milkyway/pkg/tosca/profiles/cloudify/v5_0_5"
	common "milkyway/pkg/tosca/profiles/common/v1_0"
	hot_v1_0 "milkyway/pkg/tosca/profiles/hot/v1_0"
	implicit_v2_0 "milkyway/pkg/tosca/profiles/implicit/v2_0"
	simpleForNFV_v1_0 "milkyway/pkg/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_0 "milkyway/pkg/tosca/profiles/simple/v1_0"
	simple_v1_1 "milkyway/pkg/tosca/profiles/simple/v1_1"
	simple_v1_2 "milkyway/pkg/tosca/profiles/simple/v1_2"
	simple_v1_3 "milkyway/pkg/tosca/profiles/simple/v1_3"

	"github.com/tliron/kutil/url"
)

func init() {
	initProfile(common.Profile)
	initProfile(simple_v1_0.Profile)
	initProfile(simple_v1_1.Profile)
	initProfile(simple_v1_2.Profile)
	initProfile(simple_v1_3.Profile)
	initProfile(implicit_v2_0.Profile)
	initProfile(simpleForNFV_v1_0.Profile)
	initProfile(cloudify_v5_0_5.Profile)
	initProfile(hot_v1_0.Profile)
}

func initProfile(profile map[string]string) {
	for path, content := range profile {
		if err := url.RegisterInternalURL(path, content); err != nil {
			panic(err)
		}
	}
}
