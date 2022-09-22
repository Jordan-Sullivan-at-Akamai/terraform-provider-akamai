package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiWAPSelectedHostnames_data_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfigurationWAP.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSWAPSelectedHostnames/WAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&getWAPSelectedHostnamesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiWAPSelectedHostnames_data_error_retrieving_hostnames(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfigurationWAP.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSWAPSelectedHostnames/WAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetWAPSelectedHostnames failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
						ExpectError: regexp.MustCompile(`GetWAPSelectedHostnames failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiWAPSelectedHostnames_NonWAP_data_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSSelectedHostnames/SelectedHostnames.json"), &getSelectedHostnamesResponse)

		getMatchTargetsResponse := appsec.GetMatchTargetsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSMatchTargets/MatchTargets.json"), &getMatchTargetsResponse)

		client.On("GetSelectedHostnames",
			mock.Anything,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.On("GetMatchTargets",
			mock.Anything,
			appsec.GetMatchTargetsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getMatchTargetsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
