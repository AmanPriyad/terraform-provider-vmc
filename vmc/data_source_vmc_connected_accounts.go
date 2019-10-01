package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/account_link/connectedAccounts"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/runtime/protocol/client"
	"log"
)

func dataSourceVmcConnectedAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcConnectedAccountsRead,

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Description: "Organization identifier.",
				Required:    true,
			},
			"provider_type": {
				Type:        schema.TypeString,
				Description: "The cloud provider of the SDDC (AWS or ZeroCloud).",
				Optional:    true,
				Default:     "AWS",
			},
			"ids": {
				Type:        schema.TypeList,
				Description: "The corresponding connected (customer) account UUID this connection is attached to.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVmcConnectedAccountsRead(d *schema.ResourceData, m interface{}) error {

	orgID := d.Get("org_id").(string)
	providerType := d.Get("provider_type").(string)

	connectedAccountsClient := connectedAccounts.NewConnectedAccountsClientImpl(m.(client.Connector))
	accounts, err := connectedAccountsClient.Get(orgID, &providerType)

	ids := []string{}
	for _, account := range accounts {
		ids = append(ids, account.Id)
	}

	log.Printf("[DEBUG] Connected accounts are %v\n", accounts)

	if err != nil {
		return fmt.Errorf("Error while reading accounts from org %q: %v", orgID, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", orgID, providerType))
	d.Set("ids", ids)
	return nil
}
