package ccoenaming

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceRequestBody struct {
	// Type        string `json: "type"`
	Product     string `json: "product"`
	Function    string `json: "function"`
	Application string `json: "application"`
	Region      string `json: "region"`
	Env         string `json: "env"`
	Name        string `json: "name"`
}
type ResourceResponseBody struct {
	// Type        string `json: "type"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Product     string `json: "product"`
	Function    string `json: "function"`
	Application string `json: "application"`
	Region      string `json: "region"`
	Env         string `json: "env"`
}

func resourceAzureResources() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureResourcesCreate,
		ReadContext:   resourceAzureResourcesRead,
		UpdateContext: resourceAzureResourcesUpdate,
		DeleteContext: resourceAzureResourcesDelete,
		Schema: map[string]*schema.Schema{
			// "type": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			"product": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"function": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"application": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"env": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type: schema.TypeString,
				// Optional: true,
				Computed: true,
			},
		},
	}
}
func resourceAzureResourcesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	requestBody := &ResourceRequestBody{
		Product:     d.Get("product").(string),
		Function:    d.Get("function").(string),
		Application: d.Get("application").(string),
		Region:      d.Get("region").(string),
		Env:         d.Get("env").(string),
		Name:        d.Get("name").(string),
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := http.Post("http://localhost:8888/api/naming/resource/naming", "application/json", bytes.NewReader(requestBodyBytes))
	// req, err := http.NewRequest("POST", "http://example.com/my-resource", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to create resource: unexpected status code %d of follow error: %s", resp.StatusCode, resp.Body)
	}

	// responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	var responseBody ResourceResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	id := int64(responseBody.ID)
	d.SetId(fmt.Sprintf("%d", id))
	d.Set("name", responseBody.Name)
	d.Set("function", responseBody.Function)
	d.Set("region", responseBody.Region)
	d.Set("env", responseBody.Env)
	d.Set("product", responseBody.Product)
	d.Set("application", responseBody.Application)

	return resourceAzureResourcesRead(ctx, d, m)
}

func resourceAzureResourcesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Id()
	resp, err := http.Get(fmt.Sprintf("http://localhost:8888/api/naming/resource/naming/%s", id))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	var responseBody ResourceResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(responseBody.ID, 10))
	d.Set("name", responseBody.Name)
	d.Set("function", responseBody.Function)
	d.Set("region", responseBody.Region)
	d.Set("env", responseBody.Env)
	d.Set("product", responseBody.Product)
	d.Set("application", responseBody.Application)
	return diags
}

func resourceAzureResourcesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	requestBody := &ResourceRequestBody{
		Product:     d.Get("product").(string),
		Function:    d.Get("function").(string),
		Application: d.Get("application").(string),
		Region:      d.Get("region").(string),
		Env:         d.Get("env").(string),
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/naming/resource/naming/%s", "http://localhost:8888", id), bytes.NewReader(requestBodyBytes))
	// req, err := http.NewRequest("POST", "http://example.com/my-resource", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("API returned status code %d with follow error: %s", resp.StatusCode, resp.Body)
	}
	var responseBody ResourceResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Setando as propriedades atualizadas no ResourceData
	d.Set("function", responseBody.Function)
	d.Set("region", responseBody.Region)
	d.Set("env", responseBody.Env)
	d.Set("product", responseBody.Product)
	d.Set("application", responseBody.Application)
	d.Set("name", responseBody.Name)
	return resourceAzureResourcesRead(ctx, d, m)
}

func resourceAzureResourcesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/naming/resource/naming/%s", "http://localhost:8888", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("API returned status code %d", resp.StatusCode)
	}
	d.SetId("")
	return diags
}
