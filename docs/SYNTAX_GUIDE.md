# Terraform Provider 开发语法指南

本文档详细讲解开发 Terraform Provider for JumpServer 所需的核心 Go 语法和 Terraform Plugin Framework 概念。

---

## 目录

1. [Go 基础语法](#go-基础语法)
2. [JSON 处理](#json-处理)
3. [HTTP 客户端](#http-客户端)
4. [Terraform Plugin Framework](#terraform-plugin-framework)
5. [实战示例](#实战示例)

---

## Go 基础语法

### 1. 结构体（Struct）

#### 基本定义
```go
// 定义一个简单的结构体
type Task struct {
    ID   string
    Name string
}

// 使用结构体
task := Task{
    ID:   "123",
    Name: "My Task",
}
fmt.Println(task.Name)  // 输出: My Task
```

#### 结构体标签（Struct Tag）- **重要！**
结构体标签用于在序列化、反序列化时控制字段名映射。

```go
type Task struct {
    ID          string `json:"id"`           // JSON 字段名是小写 "id"
    Name        string `json:"name"`         // Go 变量名是大写 "Name"
    Description string `json:"description,omitempty"`  // omitempty: 空值时忽略
    Secret      string `json:"secret"`       // 不加 omitempty，空值也序列化
}

// 序列化: Go → JSON
task := Task{
    ID:          "123",
    Name:        "Test",
    Description: "",           // 空字符串
    Secret:      "",           // 空字符串
}

jsonData, _ := json.Marshal(task)
// 输出: {"id":"123","name":"Test","secret":""}
// 注意: Description 为空且带有 omitempty，所以被忽略
```

**标签语法格式**：
```
`key1:"value1" key2:"value2" key3:"value3"`
```

常见标签：
- `json:"field_name"` - JSON 字段名映射
- `json:"field_name,omitempty"` - 空值时忽略该字段
- `json:"-"` - 忽略该字段，不进行序列化
- `tfsdk:"field_name"` - Terraform 字段名映射

---

### 2. 指针（Pointer）

#### 指针基础
```go
// 值类型
var num int = 10
fmt.Println(num)  // 10

// 指针类型 - *int 表示指向 int 的指针
var ptr *int = &num
fmt.Println(ptr)    // 内存地址，如 0xc0000140a0
fmt.Println(*ptr)   // 解引用: 10

// 修改值
*ptr = 20
fmt.Println(num)    // 20 (原变量被修改)
```

#### 值类型 vs 指针类型作为函数参数
```go
// 值传递 - 函数内修改不影响外部
func modifyValue(t Task) {
    t.Name = "New Name"  // 只修改了副本
}

// 指针传递 - 函数内修改会影响外部
func modifyPointer(t *Task) {
    t.Name = "New Name"  // 修改了原变量
}

// 使用示例
task := Task{Name: "Old Name"}

modifyValue(task)
fmt.Println(task.Name)  // 输出: Old Name (未改变)

modifyPointer(&task)
fmt.Println(task.Name)  // 输出: New Name (已改变)
```

#### 在结构体方法中使用指针
```go
type Client struct {
    endpoint string
}

// 值接收者 - 不能修改结构体
func (c Client) GetValue() string {
    return c.endpoint
}

// 指针接收者 - 可以修改结构体
func (c *Client) SetValue(val string) {
    c.endpoint = val  // 修改原实例
}

// 使用
client := Client{endpoint: "api.example.com"}
client.SetValue("new.api.com")
```

**何时使用指针？**
- 结构体较大，避免拷贝
- 需要修改结构体内容
- 作为函数返回值，避免拷贝

---

### 3. 接口（Interface）

#### 接口定义和实现
```go
// 定义接口
type Resource interface {
    Create(ctx context.Context, req CreateRequest, resp *CreateResponse)
    Read(ctx context.Context, req ReadRequest, resp *ReadResponse)
}

// 实现接口 - 不需要显式声明 implements
type AssetResource struct {
    client *Client
}

// 只要方法签名匹配，就自动实现了接口
func (r *AssetResource) Create(ctx context.Context, req CreateRequest, resp *CreateResponse) {
    // 实现代码
}

func (r *AssetResource) Read(ctx context.Context, req ReadRequest, resp *ReadResponse) {
    // 实现代码
}
```

#### 类型断言
```go
// 检查是否实现了接口
var _ resource.Resource = &AssetResource{}
// 编译时检查：如果 AssetResource 没有实现 Resource 接口，编译报错
// _ 表示忽略结果，只用于类型检查

// 运行时类型断言
client, ok := req.ProviderData.(*jumpserver.Client)
if !ok {
    // 类型不匹配
    return
}
// 现在 client 是 *jumpserver.Client 类型
```

#### 空接口
```go
// 空接口 interface{} 可以持有任何类型的值
var data interface{} = "string"
data = 123
data = []string{"a", "b"}

// 类型断言获取具体类型
if str, ok := data.(string); ok {
    fmt.Println("是字符串:", str)
}
```

---

### 4. 错误处理（Error Handling）

#### 创建和返回错误
```go
import "fmt"

// 创建简单错误
err := fmt.Errorf("something went wrong")

// 使用 %w 包装错误，保留原始错误信息
func doTask() error {
    err := someOperation()
    if err != nil {
        return fmt.Errorf("failed to do task: %w", err)
    }
    return nil
}

// 使用 errors.New 创建错误
err := errors.New("invalid input")
```

#### 处理错误
```go
// 基本错误处理
task, err := client.CreateTask(req)
if err != nil {
    // 记录错误并返回
    return fmt.Errorf("could not create task: %w", err)
}

// 多个错误处理
var diags diag.Diagnostics
if err != nil {
    diags.AddError(
        "Error creating task",
        fmt.Sprintf("Could not create task: %s", err),
    )
}
if diags.HasError() {
    return
}
```

---

### 5. 切片（Slice）和 Map

#### 切片操作
```go
// 声明切片
var assets []string                    // nil 切片
assets := []string{}                   // 空切片
assets := []string{"server1", "server2"} // 有初始值

// 添加元素
assets = append(assets, "server3")

// 遍历切片
for i, asset := range assets {
    fmt.Printf("索引 %d: %s\n", i, asset)
}

// 切片长度和容量
fmt.Println(len(assets))  // 长度
fmt.Println(cap(assets))  // 容量
```

#### Map 操作
```go
// 声明 Map
params := make(map[string]string)

// 添加键值对
params["asset"] = "server1"
params["command"] = "ls -la"

// 读取值
command := params["command"]  // "ls -la"
value, exists := params["nonexistent"]  // 不存在时返回零值 + false

// 删除键
delete(params, "command")

// 遍历 Map
for key, value := range params {
    fmt.Printf("%s: %s\n", key, value)
}
```

#### 切片和 Map 结合使用
```go
// Map 的切片 - 用于批量操作
var paramList []map[string]string

paramList = append(paramList, map[string]string{
    "asset":   "server1",
    "command": "uptime",
})

paramList = append(paramList, map[string]string{
    "asset":   "server2",
    "command": "uptime",
})

// 遍历
for _, param := range paramList {
    fmt.Printf("在 %s 上执行 %s\n", param["asset"], param["command"])
}
```

---

## JSON 处理

### 序列化（Marshal）
```go
import "encoding/json"

type Task struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

task := Task{
    ID:   "123",
    Name: "Test Task",
}

// Go → JSON
jsonData, err := json.Marshal(task)
if err != nil {
    return err
}

fmt.Println(string(jsonData))
// 输出: {"id":"123","name":"Test Task"}
```

### 反序列化（Unmarshal）
```go
// JSON → Go
jsonStr := `{"id":"123","name":"Test Task"}`
var task Task

err := json.Unmarshal([]byte(jsonStr), &task)
if err != nil {
    return err
}

fmt.Println(task.ID, task.Name)  // 123 Test Task
```

### 处理嵌套结构
```go
type Task struct {
    ID     string    `json:"id"`
    Params []Param   `json:"params"`
}

type Param struct {
    Asset   string `json:"asset"`
    Command string `json:"command"`
}

jsonStr := `{
    "id": "123",
    "params": [
        {"asset": "server1", "command": "ls"},
        {"asset": "server2", "command": "uptime"}
    ]
}`

var task Task
json.Unmarshal([]byte(jsonStr), &task)
// task.Params[0].Asset = "server1"
```

---

## HTTP 客户端

### 基本请求
```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// GET 请求
func (c *Client) Get(path string, result interface{}) error {
    url := c.config.Endpoint + path
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // 检查状态码
    if resp.StatusCode != 200 {
        return fmt.Errorf("API error: %d: %s", resp.StatusCode, string(body))
    }

    // 解析 JSON 响应
    if err := json.Unmarshal(body, result); err != nil {
        return err
    }

    return nil
}

// POST 请求
func (c *Client) Post(path string, body interface{}, result interface{}) error {
    // 序列化请求体
    jsonData, err := json.Marshal(body)
    if err != nil {
        return err
    }

    // 创建请求
    url := c.config.Endpoint + path
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    // 执行请求
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 处理响应...
    return nil
}
```

---

## Terraform Plugin Framework

### 1. Schema 定义

#### String Attribute
```go
func (r *CommandTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "描述信息",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Required:    false,      // 是否必填
                Optional:    false,      // 是否可选
                Computed:    true,       // 是否计算得出
                Sensitive:   false,      // 是否敏感信息
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),  // Plan 修改器
                },
                Description: "字段描述",
                Validators: []validator.String{  // 验证器
                    stringvalidator.LengthAtLeast(1),
                },
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "任务名称",
            },
        },
    }
}
```

#### List Attribute
```go
"assets": schema.ListAttribute{
    ElementType: types.StringType,  // 列表元素类型
    Required:    true,
    Description: "资产 ID 列表",
}
```

#### Bool Attribute
```go
"is_active": schema.BoolAttribute{
    Optional:    true,
    Computed:    true,
    Description: "是否激活",
}
```

### 2. Resource Model 定义

```go
type CommandTaskResourceModel struct {
    ID      types.String `tfsdk:"id"`       // 对应 Terraform 中的 id
    Name    types.String `tfsdk:"name"`     // 对应 Terraform 中的 name
    Assets  types.List   `tfsdk:"assets"`   // 对应 Terraform 中的 assets
    Status  types.String `tfsdk:"status"`   // 对应 Terraform 中的 status
}

// types.String - Terraform 字符串类型
// types.Bool - Terraform 布尔类型
// types.List - Terraform 列表类型
// types.Int64 - Terraform 整数类型
```

### 3. 类型转换（重要！）

#### Terraform String ↔ Go String
```go
// Terraform String → Go String
planID := plan.ID.ValueString()
// planID 是 Go 的 string 类型

// Go String → Terraform String
plan.ID = types.StringValue("new-id")
// plan.ID 是 types.String 类型
```

#### Terraform List ↔ Go Slice
```go
// Terraform List → Go Slice
var assets []string
diags := plan.Assets.ElementsAs(ctx, &assets, false)
// assets 是 []string 类型
// false 表示不要严格类型检查

// Go Slice → Terraform List
var nodeIDs []any
for _, node := range asset.Nodes {
    nodeIDs = append(nodeIDs, node.ID)
}
diags := plan.Nodes.ElementsAs(ctx, nodeIDs, false)
```

#### Terraform Bool ↔ Go Bool
```go
// Terraform Bool → Go Bool
isActive := plan.IsActive.ValueBool()
// isActive 是 Go 的 bool 类型

// Go Bool → Terraform Bool
plan.IsActive = types.BoolValue(true)
```

### 4. Resource 接口实现

#### Create - 创建资源
```go
func (r *CommandTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // 1. 获取用户配置的 Plan
    var plan CommandTaskResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // 2. 类型转换: Terraform → Go
    var assets []string
    diags = plan.Assets.ElementsAs(ctx, &assets, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // 3. 构建 API 请求
    createReq := &jumpserver.CreateTaskRequest{
        Name:   plan.Name.ValueString(),
        Assets: assets,
    }

    // 4. 调用 API 创建资源
    task, err := r.client.CreateTask(createReq)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating task",
            fmt.Sprintf("Could not create task: %s", err),
        )
        return
    }

    // 5. 类型转换: Go → Terraform
    plan.ID = types.StringValue(task.ID)
    plan.Status = types.StringValue(task.Status)

    // 6. 保存到 State
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}
```

#### Read - 读取资源状态
```go
func (r *CommandTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // 1. 获取当前 State
    var state CommandTaskResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // 2. 调用 API 获取最新状态
    task, err := r.client.GetTask(state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading task",
            fmt.Sprintf("Could not read task: %s", err),
        )
        return
    }

    // 3. 更新 State
    state.Name = types.StringValue(task.Name)
    state.Status = types.StringValue(task.Status)

    // 4. 保存 State
    diags = resp.State.Set(ctx, state)
    resp.Diagnostics.Append(diags...)
}
```

#### Update - 更新资源
```go
func (r *CommandTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // 1. 获取新的 Plan
    var plan CommandTaskResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // 2. 获取旧的 State
    var state CommandTaskResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // 3. 调用 API 更新资源
    updateReq := &jumpserver.UpdateTaskRequest{
        Name: plan.Name.ValueString(),
    }

    task, err := r.client.UpdateTask(state.ID.ValueString(), updateReq)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error updating task",
            fmt.Sprintf("Could not update task: %s", err),
        )
        return
    }

    // 4. 更新 State
    plan.Name = types.StringValue(task.Name)
    plan.Status = types.StringValue(task.Status)

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
}
```

#### Delete - 删除资源
```go
func (r *CommandTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // 1. 获取 State
    var state CommandTaskResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // 2. 调用 API 删除资源
    err := r.client.DeleteTask(state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error deleting task",
            fmt.Sprintf("Could not delete task: %s", err),
        )
        return
    }

    // 3. 删除成功，不需要额外操作
}
```

---

## 学习路径建议

### 阶段 1: 理解基础概念
1. 通读本文档
2. 对照 `asset_resource.go` 理解每个概念
3. 手动抄写一遍代码，边抄边理解

### 阶段 2: 小练习
1. 给 `AssetResource` 添加一个新字段
2. 修改该字段的值
3. 测试编译和运行

### 阶段 3: 实战开发
1. 参考本文档的完整示例
2. 创建 `task.go` 定义 API 客户端
3. 创建 `command_task_resource.go` 实现 Resource
4. 在 `provider.go` 中注册
5. 编写测试用例

---

## 常见问题

### Q1: 什么时候使用指针？
A: 结构体较大、需要修改内容、作为函数返回值时使用指针

### Q2: 为什么要用 `tfsdk:"id"` 这种标签？
A: 告诉 Terraform 如何将 HCL 配置映射到 Go 结构体

### Q3: `ElementsAs(ctx, &assets, false)` 中的 false 是什么意思？
A: 表示不进行严格类型检查，允许隐式类型转换

### Q4: 为什么有些字段是 `Computed`？
A: `Computed` 字段由 Provider 计算得出，用户在配置中不需要指定

### Q5: Update 方法可以不实现吗？
A: 可以，但需要返回错误提示用户不支持更新

---

## 参考资源

- [Go by Example](https://gobyexample.com/) - 快速上手 Go
- [Go 语言圣经](https://gopl-zh.github.io/) - 深入理解 Go
- [Terraform Plugin Framework](https://www.terraform.io/plugin/framework) - 官方文档
- [Terraform Plugin Framework Examples](https://github.com/hashicorp/terraform-provider-aws) - AWS Provider 源码

---

祝学习顺利！如有问题，随时提问。
