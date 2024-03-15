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

type VMRequestBody struct {
	// Type        string `json: "type"`
	OS      string `json:"os"`
	Region  string `json:"region"`
	Env     string `json:"env"`
	Product string `json:"product"`
	Name    string `json:"name"`
}
type VMResponseBody struct {
	// Type        string `json: "type"`
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	OS      string `json:"os"`
	Region  string `json:"region"`
	Env     string `json:"env"`
	Product string `json:"product"`
}

func resourceAzureVMs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureVMsCreate,
		ReadContext:   resourceAzureVMsRead,
		UpdateContext: resourceAzureVMsUpdate,
		DeleteContext: resourceAzureVMsDelete,
		Schema: map[string]*schema.Schema{
			"product": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"os": &schema.Schema{
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
func resourceAzureVMsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics
	requestBody := &VMRequestBody{
		// Type:        d.Get("type").(string),
		Product: d.Get("product").(string),
		OS:      d.Get("os").(string),
		Region:  d.Get("region").(string),
		Env:     d.Get("env").(string),
		Name:    d.Get("name").(string),
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := http.Post("http://localhost:8888/api/naming/vm/naming", "application/json", bytes.NewReader(requestBodyBytes))
	// req, err := http.NewRequest("POST", "http://example.com/my-resource", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to create resource: unexpected status code %d", resp.StatusCode)
	}

	var responseBody VMResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(responseBody.ID, 10))
	d.Set("name", responseBody.Name)
	d.Set("os", responseBody.OS)
	d.Set("region", responseBody.Region)
	d.Set("env", responseBody.Env)
	d.Set("product", responseBody.Product)

	// process response body as needed

	return resourceAzureVMsRead(ctx, d, m)
}

func resourceAzureVMsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Id()
	resp, err := http.Get(fmt.Sprintf("http://localhost:8888/api/naming/vm/naming/%s", id))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	var responseBody VMResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(responseBody.ID, 10))
	d.Set("name", responseBody.Name)
	d.Set("os", responseBody.OS)
	d.Set("region", responseBody.Region)
	d.Set("env", responseBody.Env)
	d.Set("product", responseBody.Product)
	return diags
}

func resourceAzureVMsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	requestBody := &VMRequestBody{
		Product: d.Get("product").(string),
		OS:      d.Get("os").(string),
		Region:  d.Get("region").(string),
		Env:     d.Get("env").(string),
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/naming/vm/naming/%s", "http://localhost:8888", id), bytes.NewReader(requestBodyBytes))
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
		return diag.Errorf("API returned status code %d", resp.StatusCode)
	}
	var responseBody VMResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Setando as propriedades atualizadas no ResourceData
	d.Set("os", responseBody.OS)
	d.Set("region", responseBody.Region)
	d.Set("env", responseBody.Env)
	d.Set("product", responseBody.Product)
	d.Set("name", responseBody.Name)
	return resourceAzureVMsRead(ctx, d, m)
}

func resourceAzureVMsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/naming/vm/naming/%s", "http://localhost:8888", id), nil)
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
