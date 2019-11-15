package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/model"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs/publicips"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/tasks"
	"log"
	"time"
)

func resourcePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Update: resourcePublicIPUpdate,
		Delete: resourcePublicIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

<<<<<<< HEAD
=======
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

>>>>>>> master
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization identifier",
			},
			"sddc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sddc Identifier",
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP Allocation ID",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
<<<<<<< HEAD
				Description: "public IP allocated to the SDDC.",
=======
				Description: "Allocated Public IP",
>>>>>>> master
			},
			"private_ip": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Workload VM private IP",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workload VM name",
<<<<<<< HEAD
				Required:    true,
				Description: "private IP allocated to the SDDC.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name for the workload VM public IP assignment.",
=======
>>>>>>> master
			},
			"dnat_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
<<<<<<< HEAD
				Description: "DNAT rule identifier.",
=======
				Description: "DNAT rule ID",
>>>>>>> master
			},
			"snat_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
<<<<<<< HEAD
				Description: "SNAT rule identifier.",
=======
				Description: "SNAT rule ID",
>>>>>>> master
			},
		},
	}
}

func resourcePublicIPCreate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector

	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)

	privateIP := d.Get("private_ip").(string)
	workloadName := d.Get("name").(string)
	publicIPsClient := publicips.NewPublicipsClientImpl(connector)

	var sddcAllocatePublicIpSpec = &model.SddcAllocatePublicIpSpec{
		Count:      1,
		PrivateIps: []string{privateIP},
		Names:      []string{workloadName},
	}

	// Create Public IP
	task, err := publicIPsClient.Create(orgID, sddcID, *sddcAllocatePublicIpSpec)
	if err != nil {
		return fmt.Errorf("error while creating public IP : %v", err)
	}

	tasksClient := tasks.NewTasksClientImpl(connector)

<<<<<<< HEAD
	return resource.Retry(300*time.Minute, func() *resource.RetryError {
=======
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
>>>>>>> master
		task, err := tasksClient.Get(orgID, task.Id)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))
		}
		if *task.Status != "FINISHED" {
			log.Print("Task not finished yet")
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		} else {
			publicIPClient := publicips.NewPublicipsClientImpl(connector)
			publicIPs, err := publicIPClient.List(orgID, sddcID)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error while getting list of public IPs for SDDC %s: %v", d.Get("sddc_id").(string), err))
			}
			for i := 0; i < len(publicIPs); i++ {
				singleVal := publicIPs[i]
				if d.Get("private_ip").(string) == *(singleVal.AssociatedPrivateIp) {
					d.SetId(*(singleVal.AllocationId))
					break
				}
			}
			if d.Id() == "" {
				return resource.NonRetryableError(fmt.Errorf("error while getting the allocationID %v", err))
			}
			return resource.NonRetryableError(resourcePublicIPRead(d, m))
		}
	})
}

func resourcePublicIPRead(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector
	publicIPClient := publicips.NewPublicipsClientImpl(connector)

	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	allocationID := d.Id()
	publicIP, err := publicIPClient.Get(orgID, sddcID, allocationID)
	if err != nil {
		return fmt.Errorf("error while getting public IP details for %s: %v", allocationID, err)
	}

	d.SetId(*publicIP.AllocationId)
	d.Set("public_ip", publicIP.PublicIp)
	d.Set("name", publicIP.Name)
	d.Set("private_ip", publicIP.AssociatedPrivateIp)
	d.Set("dnat_rule_id", publicIP.DnatRuleId)
	d.Set("snat_rule_id", publicIP.SnatRuleId)
	return nil

}

func resourcePublicIPDelete(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector
	publicIPClient := publicips.NewPublicipsClientImpl(connector)

	allocationID := d.Id()
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	publicIP := d.Get("public_ip").(string)
	task, err := publicIPClient.Delete(orgID, sddcID, allocationID)
	if err != nil {
		return fmt.Errorf("Error while deleting public IP %s: %v", publicIP, err)
	}
<<<<<<< HEAD

	return resource.Retry(300*time.Minute, func() *resource.RetryError {
		tasksClient := tasks.NewTasksClientImpl(connector)
=======
	tasksClient := tasks.NewTasksClientImpl(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
>>>>>>> master
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error while deleting public IP %s: %v", publicIP, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}

func resourcePublicIPUpdate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	publicIPClient := publicips.NewPublicipsClientImpl(connector)
	allocationID := d.Id()
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	publicIPName := d.Get("name").(string)
<<<<<<< HEAD
=======
	associatedPrivateIP := d.Get("private_ip").(string)
	publicIP := d.Get("public_ip").(string)
>>>>>>> master

	if d.HasChange("private_ip") {

		if d.Get("private_ip") == "" {
			//detach privateIP case
<<<<<<< HEAD
			publicIP := d.Get("public_ip").(string)
=======
>>>>>>> master
			newSDDCPublicIP := model.SddcPublicIp{
				PublicIp: publicIP,
				Name:     &publicIPName,
			}
			task, err := publicIPClient.Update(orgID, sddcID, allocationID, "detach", newSDDCPublicIP)
			if err != nil {
				return fmt.Errorf("error while detaching the public ip: %v", err)
			}
<<<<<<< HEAD
			err = WaitForTask(connector, orgID, task.Id)
			if err != nil {
				return fmt.Errorf("Error while waiting for the detach task %s: %v", task.Id, err)
			}
		} else {
			//reattach privateIP case
			publicIP := d.Get("public_ip").(string)
			associatedPrivateIP := d.Get("private_ip").(string)
=======
			tasksClient := tasks.NewTasksClientImpl(connector)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				task, err := tasksClient.Get(orgID, task.Id)
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("Error while waiting for task sddc %s: %v", task.Id, err))
				}
				if *task.Status != "FINISHED" {
					return resource.RetryableError(fmt.Errorf("Expected IP to be detached but was in state %s", *task.Status))
				}
				return resource.NonRetryableError(resourcePublicIPRead(d, m))
			})
			if err != nil {
				return err
			}

		} else {
			//reattach privateIP case
>>>>>>> master
			newSDDCPublicIP := model.SddcPublicIp{
				PublicIp:            publicIP,
				AssociatedPrivateIp: &associatedPrivateIP,
				Name:                &publicIPName,
			}
			task, err := publicIPClient.Update(orgID, sddcID, allocationID, "reattach", newSDDCPublicIP)
			if err != nil {
				return fmt.Errorf("error while reattaching the public IP : %v", err)
			}
<<<<<<< HEAD
			err = WaitForTask(connector, orgID, task.Id)
			if err != nil {
				return fmt.Errorf("Error while waiting for the reattach task %s: %v", task.Id, err)
			}
		}
		d.Set("private_ip", d.Get("private_ip").(string))

	} else if d.HasChange("name") {

		newPublicIPName := d.Get("name").(string)
		associatedPrivateIP := d.Get("private_ip").(string)
		newSDDCPublicIP := model.SddcPublicIp{
			Name:                &newPublicIPName,
=======
			tasksClient := tasks.NewTasksClientImpl(connector)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				task, err := tasksClient.Get(orgID, task.Id)
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("error while waiting for task sddc %s: %v", task.Id, err))
				}
				if *task.Status != "FINISHED" {
					return resource.RetryableError(fmt.Errorf("expected IP to be reattached but was in state %s", *task.Status))
				}
				return resource.NonRetryableError(resourcePublicIPRead(d, m))
			})
			if err != nil {
				return err
			}
		}

	} else if d.HasChange("name") {
		//rename case
		newSDDCPublicIP := model.SddcPublicIp{
			Name:                &publicIPName,
>>>>>>> master
			AssociatedPrivateIp: &associatedPrivateIP,
		}
		task, err := publicIPClient.Update(orgID, sddcID, allocationID, "rename", newSDDCPublicIP)

		if err != nil {
			return fmt.Errorf("error while updating public IP for rename action type  : %v", err)
		}
<<<<<<< HEAD
		err = WaitForTask(connector, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for the rename task %s: %v", task.Id, err)
		}
		d.Set("name", d.Get("name").(string))
	}

	return resourcePublicIPRead(d, m)

=======

		tasksClient := tasks.NewTasksClientImpl(connector)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error while waiting for task sddc %s: %v", task.Id, err))
			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected IP to be renamed but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourcePublicIPRead(d, m))
		})
		if err != nil {
			return err
		}
	}
	return resourcePublicIPRead(d, m)
>>>>>>> master
}
