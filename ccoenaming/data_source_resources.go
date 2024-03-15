package ccoenaming

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAzureResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAzureResourcesRead,
		Schema: map[string]*schema.Schema{
			"resources": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"function": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"env": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"application": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAzureResourcesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	var diags diag.Diagnostics
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/naming/resource/naming", "http://localhost:8888"), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()
	resources := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&resources)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("resources", resources); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
