# 最简单的创建资产并同步账号示例

terraform {
  required_providers {
    jumpserver = {
      source  = "fit2cloud/jumpserver"
      version = ">= 1.0.0"
    }
  }
}

provider "jumpserver" {
  endpoint   = "https://10.1.14.25"
  key_id     = "bb3cd120-585a-40ad-9a20-09991187ccc9"
  key_secret = "XPJ39dgHt6pPtlgmOsMsuU8FMPwBjRxAye0b"
  org_id     = "00000000-0000-0000-0000-000000000002"
  insecure_skip_verify = true
}

# 创建资产
resource "jumpserver_asset" "demo" {
  name     = "EC2测试服务器-test-12345"
  address  = "192.168.100.2"
  platform = "Linux"
  is_active = true
  comment  = "从AWS EC2同步的资产"
}

# 同步SSH账号
resource "jumpserver_account" "ec2_ssh" {
  asset       = jumpserver_asset.demo.id
  username    = "ec2-user"
  secret      = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAx6QY1AinZVJMCHFjKOufviWWkbpdNeUPQjByWkxbgtJsLyzA\nVwz9PGOo+4VrdH2zd7CLUvHLQ7AfRPhDXgrKDPoaknlXDqd3oqzK2chaClFp4IK\neozcThKWDq3SXCdhdYO6LrCvL3wFO84BV/M7wwFizPHu726h5J1f0OfqlvvtMSAK\nISS7Jc+rBnGcscsqht/5Gu9eJInomYuEqmsvykx6O4Zkb5k3xvFvfVKOv3f2Q4Fz\nHCWn271Q8UI9ZsB4mLINareOuzd2IieRiXMopyLYaSXe8hLnAoB5rnG5iZRP8U+0\n3ExQg25i4lpcQgeT9nzJc3p557h6D1M/bi5+YQIDAQABAoIBABCsCorCgkA65DCc\nT3yeWMPHXdCjsJ8Mlv6fDx2tXMMLEY/K+/EJG6jMZdNDbBrZWIB5VNlDagcoESRw\nWyfiXMdCp69txLBrmdkS9wnC6ooMDGpdGMTtOISolrF5IKUjgMcQjh7SEH81qzY4\nWPJgVLBPUFHvLlX+djSiU9sdUwDyuTTg2eZOKm+BDYYHFsOHorS7AVx4+G80W4lS\n/Es3a0vv5GpKdKRaN86F553TV8Zx7JiBMRx5oTh4oYSN085khq+Mnc9IMK+pSxyC\nBxrGm1l3FiMwh6OOH+9fuu2ipbhtmOeEsxqqETmB9A/oENJPGuwRg+qpzYU4aKvW\nOiveL5UCgYEA4hl6X4biB7Fcd4ggCOPFwqdiMJsuZa0kQ1fDzlaPrBTgI5AT8pgN\nf/tiyL1Ye57mVe+XhF5FkZtk7pxCRgjH1r0MBtGF+bv+0ipDlqgpe/BimUJYspcV\n4IDQtYCMCfsurY5tGw1KB1M4IF+ksivz7XD2KcwPxuTrCUHiY5z1t5MCgYEA4grf\n9WN730xvOVYoMPCpqzZFygUwBNqF7rVuM9zIFC4kDfkBlF73VweVSXI9E0cIMv6P\nKP4EHjLJlI9P70+8bdcw1mDIbuHY8wcKBrkrE/Fu7GS/SgHAqaJ919rAP8KhQcpS\nyPzXKhpQGYrLtUtzV2q0Cv6SZxgKrJ3II/ycwrsCgYEAwDiS7kX5QjsKduD3AzfK\nOKwfcV1s+6pQqyQhZvn2mYEB8ZobK2MUDxuEp086u5ajEqpoMXQIRztKewXD3lC2\nvRzp7Z4R/fhTMxAVeC8tXZ5H5S4fxG1ofv5k8foAlLfEvm7Y2WfZ6RJaJEPL/GIb\ntmEUFwLS4vBZ1fv6YV/fExsCgYEAzgqPnoQyM5beg2sPc5zLa68q6jzUSnhOQQrM\nCyYikpKEduAVGoN9/ayB3dLt7RaAWMtE/16brlMo/+uqNz99SLowYBkUWk4vjUdL\nUlmS9LjMHVqwKutyDK56+zkAqJ3mk6uyzlX6YvxdKwsjKHxABNzUhHkMRkDZ6gJg\nnrzyv2UCgYEAndlwbHlWq9gCFX6BeMo5XdxhZTCuOwc1q6Ja4/uE8IRsLmbJmKdt\n5F7Wj014dy3d3AIJE2WP5R2/OJ1c/jXAVCtNxB0DxFdaVu7v9LkSuFFZqqUfnJHf\nqaKd3HK8iMfzKXVW5GXEXIf/D/mtMboZg7BEPaA0kq127QbbzXMZs5I=\n-----END RSA PRIVATE KEY-----"
  secret_type = "ssh_key"
  comment  = "SSH Key for EC2测试服务器"
}
