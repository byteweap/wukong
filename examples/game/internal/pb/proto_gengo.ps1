Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# 使用 Docker 的 protobuf 镜像生成 Go 代码
& docker run --rm -v "${PSScriptRoot}:/defs" rvolosatovs/protoc --proto_path=/defs --go_out=/defs /defs/*.proto

# 将生成的代码移动到当前目录
Move-Item -Force (Join-Path $PSScriptRoot "pb\*.go") $PSScriptRoot
# 删除临时目录
Remove-Item -Recurse -Force (Join-Path $PSScriptRoot "pb")
