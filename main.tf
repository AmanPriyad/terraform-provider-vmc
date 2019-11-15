provider "vmc" {
  refresh_token = "5aAHQLgiecL08W0L4qif0G9I6B0tge8Y3kapdgs6UTY8LRdL51dC5Z4QrVzAI5I3"

  # for staging environment only
  vmc_url       = "https://stg.skyscraper.vmware.com/vmc/api"
  csp_url       = "https://console-stg.cloud.vmware.com"
}

data "vmc_org" "my_org" {
  id = "05e0a625-3293-41bb-a01f-35e762781c2a:"
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}

data "vmc_customer_subnets" "my_subnets" {
  org_id               = "${data.vmc_org.my_org.id}"
  connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
  region               = "us-west-2"
}

resource "vmc_sddc" "sddc_1" {
  org_id = "${data.vmc_org.my_org.id}"

  # storage_capacity    = 100
  sddc_name           = ""
  vpc_cidr            = "10.2.0.0/16"
  num_host            = 1
  provider_type       = "AWS"
  region              = "${data.vmc_customer_subnets.my_subnets.region}"
  vxlan_subnet        = "192.168.1.0/24"
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  # sddc_template_id = ""
  deployment_type = "SingleAZ"

  account_link_sddc_config = [
    {
      customer_subnet_ids  = ["${data.vmc_customer_subnets.my_subnets.ids.0}"]
      connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
    },
  ]
  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}

resource "vmc_publicips" "IP1" {
  org_id = "${data.vmc_org.my_org.id}"
  sddc_id = "${vmc_sddc.sddc_1.id}"
  private_ip = "10.2.33.45"
  name = "vm1"
}

