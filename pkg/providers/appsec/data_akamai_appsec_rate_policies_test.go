package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRatePolicies_data_basic(t *testing.T) {
	t.Run("match by RatePolicies ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetRatePoliciesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSRatePolicies/RatePolicies.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetRatePolicies",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRatePoliciesRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRatePolicies/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_rate_policies.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
