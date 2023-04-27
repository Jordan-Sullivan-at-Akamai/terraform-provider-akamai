provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edge_hostname" "importedgehostname" {
  contract      = "ctr_1"
  group         = "grp_2"
  edge_hostname = "test.akamaized.net"
  ip_behavior   = "IPV4"
}
