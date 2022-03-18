package property

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tj/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

// Alias of mock.Anything to use as a placeholder for any context.Context
var AnyCTX = mock.Anything

func TestResCPCode(t *testing.T) {
	// Helper to set up an expected call to mock papi.GetCPCodes with mock impl backed by the given slice
	expectGetCPCode := func(m *mockpapi, ContractID, GroupID string, CPCodes *[]papi.CPCode) *mock.Call {
		mockImpl := func(_ context.Context, req papi.GetCPCodesRequest) (*papi.GetCPCodesResponse, error) {
			res := &papi.GetCPCodesResponse{
				ContractID: req.ContractID,
				GroupID:    req.GroupID,
				CPCodes:    papi.CPCodeItems{Items: *CPCodes},
			}
			return res, nil
		}

		req := papi.GetCPCodesRequest{ContractID: ContractID, GroupID: GroupID}

		return m.OnGetCPCodes(mockImpl, AnyCTX, req)
	}

	// Helper to set up an expected call to mock papi.CreateCPCode with mock impl backed by the given slice
	expectCreateCPCode := func(m *mockpapi, CPCName, Product, Contract, Group string, CPCodes *[]papi.CPCode) *mock.Call {
		mockImpl := func(_ context.Context, req papi.CreateCPCodeRequest) (*papi.CreateCPCodeResponse, error) {
			cpc := papi.CPCode{
				ID:         fmt.Sprintf("cpc_%d", len(*CPCodes)),
				Name:       req.CPCode.CPCodeName,
				ProductIDs: []string{req.CPCode.ProductID},
			}

			*CPCodes = append(*CPCodes, cpc)
			res := &papi.CreateCPCodeResponse{CPCodeID: cpc.ID}

			return res, nil
		}

		req := papi.CreateCPCodeRequest{
			ContractID: Contract,
			GroupID:    Group,
			CPCode: papi.CreateCPCode{
				ProductID:  Product,
				CPCodeName: CPCName,
			},
		}

		return m.OnCreateCPCode(mockImpl, AnyCTX, req)
	}

	// Helper to set up an expected call to mock papi.UpdateCPCode with mock impl backed by the given slice
	expectUpdateCPCode := func(m *mockpapi, CPCodeID int, name string, CPCodes *[]papi.CPCode, err error) *mock.Call {
		mockImpl := func(_ context.Context, req papi.UpdateCPCodeRequest) (*papi.CPCodeDetailResponse, error) {
			if err != nil {
				return nil, err
			}
			(*CPCodes)[CPCodeID].Name = name

			res := &papi.CPCodeDetailResponse{
				ID:   req.ID,
				Name: req.Name,
			}

			return res, nil
		}

		f := false
		req := papi.UpdateCPCodeRequest{
			ID:               CPCodeID,
			Name:             name,
			Purgeable:        &f,
			OverrideTimeZone: &papi.CPCodeTimeZone{},
		}

		return m.OnUpdateCPCode(mockImpl, AnyCTX, req)
	}

	// Helper to set up an expected call to mock papi.GetCPCodeDetail
	expectGetCPCodeDetail := func(m *mockpapi, CPCodeID int, CPCodes *[]papi.CPCode, err error) *mock.Call {
		var call *mock.Call

		call = m.On("GetCPCodeDetail", AnyCTX, CPCodeID).Run(func(args mock.Arguments) {
			if err != nil {
				call.Return(nil, err)
			} else {
				res := &papi.CPCodeDetailResponse{
					ID:   CPCodeID,
					Name: (*CPCodes)[CPCodeID].Name,
				}
				call.Return(res, nil)
			}
		})
		return call
	}

	t.Run("create new CP Code", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_1", "grp_1", &CPCodes)
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/create_new_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
					),
				}},
			})
		})
	})

	t.Run("create new CP Code with deprecated attributes", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_1", "grp_1", &CPCodes)
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes)

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/create_new_cp_code_deprecated_attrs.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_1"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_1"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code with multiple products", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test", "prd_wrong", "another_wrong"}}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_test", "grp_test", &CPCodes)
		// No mock behavior for create because we're using an existing CP code

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_test"),
					),
				}},
			})
		})
	})

	t.Run("use existing CP Code", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs: []string{"prd_test"}}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_test", "grp_test", &CPCodes)
		// No mock behavior for create because we're using an existing CP code

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "group_id", "grp_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "contract_id", "ctr_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product", "prd_test"),
						resource.TestCheckResourceAttr("akamai_cp_code.test", "product_id", "prd_test"),
					),
				}},
			})
		})
	})

	t.Run("product missing from CP Code", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code", ProductIDs: []string{"prd_test"}},
			{ID: "cpc_test2", Name: "test cpcode"}, // Matches name from fixture
		}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_test", "grp_test", &CPCodes)
		// No mock behavior for create because we're using an existing CP code

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestResCPCode/use_existing_cp_code.tf"),
					ExpectError: regexp.MustCompile("Couldn't find product id on the CP Code"),
				}},
			})
		})
	})

	t.Run("change name", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_1", "grp_1", &CPCodes)
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()

		expectGetCPCodeDetail(client, 0, &CPCodes, nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", &CPCodes, nil).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "renamed cpcode"),
						),
					},
				},
			})
		})
	})

	t.Run("import existing cp code", func(t *testing.T) {
		client := &mockpapi{}
		id := "123,1,2"

		cpCodes := []papi.CPCode{{ID: "cpc_123", Name: "test cpcode", ProductIDs: []string{"prd_Web_Accel"}}}
		expectGetCPCode(client, "ctr_1", "grp_2", &cpCodes)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
					},
					{
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cp_code.test",
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "grp_2", rs.Attributes["group"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract"])
							assert.Equal(t, "prd_Web_Accel", rs.Attributes["product_id"])
							assert.Equal(t, "prd_Web_Accel", rs.Attributes["product"])
							assert.Equal(t, "cpc_123", rs.Attributes["id"])
							assert.Equal(t, "test cpcode", rs.Attributes["name"])
							return nil
						},
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("invalid import ID passed", func(t *testing.T) {
		client := &mockpapi{}
		id := "123"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:        loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cp_code.test",
						ExpectError:   regexp.MustCompile("comma-separated list of CP code ID, contract ID and group ID has to be supplied in import"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("empty CP code ID passed", func(t *testing.T) {
		client := &mockpapi{}
		id := ",ctr_1-1NC95D,grp_194665"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:        loadFixtureString("testdata/TestResCPCode/import_cp_code.tf"),
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cp_code.test",
						ExpectError:   regexp.MustCompile("CP Code is a mandatory parameter"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("immutable attributes updated", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_1", "grp_1", &CPCodes)
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_immutable.tf"),
						ExpectError: regexp.MustCompile(`cp code attribute 'contract' cannot be changed after creation \(immutable\)`),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_immutable.tf"),
						ExpectError: regexp.MustCompile(`cp code attribute 'product' cannot be changed after creation \(immutable\)`),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_immutable.tf"),
						ExpectError: regexp.MustCompile(`cp code attribute 'group' cannot be changed after creation \(immutable\)`),
					},
				},
			})
		})
	})

	t.Run("error fetching cpCode details", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_1", "grp_1", &CPCodes)
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()

		expectGetCPCodeDetail(client, 0, &CPCodes, fmt.Errorf("oops")).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
	})

	t.Run("error updating cpCode", func(t *testing.T) {
		client := &mockpapi{}
		defer client.AssertExpectations(t)

		// Contains CP Codes known to mock PAPI
		CPCodes := []papi.CPCode{}

		// Values are from fixture:
		expectGetCPCode(client, "ctr_1", "grp_1", &CPCodes)
		expectCreateCPCode(client, "test cpcode", "prd_1", "ctr_1", "grp_1", &CPCodes).Once()

		expectGetCPCodeDetail(client, 0, &CPCodes, nil).Once()
		expectUpdateCPCode(client, 0, "renamed cpcode", &CPCodes, fmt.Errorf("oops")).Once()

		// No mock behavior for delete because there is no delete operation for CP Codes

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCPCode/change_name_step0.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cp_code.test", "id", "cpc_0"),
							resource.TestCheckResourceAttr("akamai_cp_code.test", "name", "test cpcode"),
						),
					},
					{
						Config:      loadFixtureString("testdata/TestResCPCode/change_name_step1.tf"),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
	})
}
