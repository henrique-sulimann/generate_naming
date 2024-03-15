# Criando um Terraform Provider

# Pré-requisitos
Para criarmos um `Terraform Provider` são necessários alguns pré-requisitos:
- terraform cli versão >= 0.14
- golang versão >= 1.18

# Baixando as Dependências

Assim que você clonar este repositório, na raiz do repositório execute o seguinte comando para instalar as dependências do projeto
```
go mod download
```

# Criando um resource no Terraform Provider

Para criarmos um resource no terraform provider, nós temos que criar o arquivo que siga a seguinte convensão `resource_<nome do resource>.go` e dentro deste arquivo nós iremos criar o `Schema` que este resource irá precisar receber para ser criado.
Abaixo um Exemplo de Schema de Resource
```
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
```
Como podemos ver no exemplo acima, nós criamos o Schema do resource e passamos as quatro operações de um CRUD de uma API, que é o Create, Read, Update e Delete, pois são essas quatro operações que o nosso provider de terraform precisa para conseguir executar o `terraform plan`, `terraform apply` e `terraform destroy` corretamente.

Portanto, vamos entender como criar as funções para essas 4 operações.

# Criando a função para `CreateContext`

Agora nós iremos aprender a criar um função que irá realizar a criação do nosso recurso, para que isso seja possível, nós precisamos saber quais os parâmetros são requeridos pela nossa API, após descobrir, podemos criar a função de criação da seguinte maneira

```
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
func resourceAzureResourcesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics
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
```

Como podemos perceber no exemplo acima, basicamente o que nós fizemos foi pegar os valores que foram recebidos pelo Schema e armazenar na nossa `Struct`, após isso, realizamos o `POST` para a API e, com base no `RESPONSE` da API, nós armazenamos os dados no nosso `STATE`

# Criando a função para o `ReadContext`

Para criar a funcão para o `ReadContext` nós apenas precisar realizar um GET para a nossa API e armazenar o `response` para o nosso `ResourceData` do terraform

```
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
```
Percebamos que a primeira coisa que é realizado o `ReadContext` é tentar resgatar o `id` do nosso `state`, pois é assim que o terraform saberá se ele irá criar ou não o recurso.
E Percebamos também que o `GET` é realizado pelo `ID` que nós pegamos do `ResourceData` e, com base na `ResponseBody` desse GET, é armazenado estes valores no `ResourceData` para que possa ser comparado com o valor do `state` e identificar o que será ou não alterado caso rode o `terraform apply`

# Criando a função para o `UpdateContext`

Agora nós iremos criar a funcão que irá realizar o `PUT` por baixo dos panos na nossa API e que irá atualizar os dados do nosso Resource.
Para isso, vamos criar o seguinte código

```
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
```
Como podemos ver, o Update é muito similar com o `CreateContext`, porém, muda apenas o método que é utilizado na requisição HTTP.

# Criando a função para o `DeleteContext`

Agora iremos criar uma função que irá realizar o `Delete` do nosso recurso, ou seja, quando executarmos o `terraform destroy`, será esta função que será chamada.

```
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
```
E como podemos ver, a funcão de `Delete` é bem simples, basicamente nós precisamos resgatar o `ID` do nosso `state`, manda a requisição do `Delete` e, ao final, zerar o nosso `state`.

Após criar todas as funções para o nosso `Resource` nós precisamos referente ela no nosso arquivo de configuração do nosso provider chamado `provider.go`

```
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

```

E como podemos ver, na configuração do nosso provider temos os `Resources` e os `Data Sources`, e para cada um deles, nós passamos a função que contém todos os métodos e schemas de cada resource.

# Observação
Um ponto muito importante para se atentar na hora de criar cada resource é com a nomenclatura que é utilizada para cada resource.
Ficando da seguinte forma:
```
<nome do provider>_<nome do resource>
```
## Exemplo
Se o nosso provider será chamado de `ccoe-naming` e o nosso `Resource` será chamado de `resources`, na configuração do nosso provider nós temos que referenciar o resource da seguinte forma

```
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ccoe-naming_resources": resourceAzureResources(),
		},
	}
}
```

# Buildando o provider e movendo para a pasta de plugins do Terraform
Agora que nós já temos um provider com um resource criado, podemos buildar este código e testar no terraform.
Para fazer isso, fazemos o seguinte:

## Build
```
go clean -cache
go build -o terraform-provider-ccoe-naming
```

## Criação do diretório de plugins do terraform localmente e cópia do binário
```
mkdir -p ~/.terraform.d/plugins/gms.corp/dev/ccoe-naming/0.0.1/darwin_amd64
mv terraform-provider-ccoe-naming ~/.terraform.d/plugins/gms.corp/dev/ccoe-naming/0.0.1/darwin_amd64
```