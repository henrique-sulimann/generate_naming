// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ccoenaming

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ccoe-naming_resources": resourceAzureResources(),
			"ccoe-naming_vms":       resourceAzureVMs(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ccoe-naming_resources": dataSourceAzureResources(),
		},
	}
}
