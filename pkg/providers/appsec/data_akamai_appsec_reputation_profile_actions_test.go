package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationProfileActions_data_basic(t *testing.T) {
	t.Run("match by ReputationProfileActions ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getReputationProfileActionsResponse := appsec.GetReputationProfileActionsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSReputationProfileActions/ReputationProfileActions.json"), &getReputationProfileActionsResponse)

		client.On("GetReputationProfileActions",
			mock.Anything,
			appsec.GetReputationProfileActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ReputationProfileID: 321456, Action: ""},
		).Return(&getReputationProfileActionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSReputationProfileActions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_reputation_profile_actions.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
