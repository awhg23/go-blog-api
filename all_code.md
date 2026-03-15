# go-blog-api — 全项目代码汇总

> 本文件由脚本自动生成，包含项目中所有源码文件，方便整体上传给 AI 工具参考。

## `.env.example`

```
DB_PASSWORD=your_database_password_here
JWT_SECRET=your_jwt_secret_key_here
```

## `.gitignore`

```
.idea/
.vscode/
*.exe
vendor/
*.txt
config/config.yaml
config/config.docker.yaml
.env
```

## `Dockerfile`

```dockerfile
# ==========================================
# 第一阶段：构建阶段（Builder）
# 使用带有 Go 环境的镜像来编译我们的代码
# ==========================================
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装 git（如果需要从私有仓库拉取依赖）和其他必要的工具
RUN apk add --no-cache git
# 配置 Go 代理（国内加速下载依赖）
ENV GOPROXY=https://proxy.golang.org,direct

# 先复制 go.mod 和 go.sum 并下载依赖（利用 Docker 缓存加速）
COPY go.mod go.sum ./
RUN go mod download

# 把项目所有源代码复制进容器
COPY . .

# 编译 Go 语言程序，生成一个名为 main 的可执行文件
RUN go build -o main ./cmd


# ==========================================
# 第二阶段：运行阶段（Runner）
# 使用一个极小的 alpine 镜像来运行编译好的文件，减小最终体积
# ==========================================
FROM alpine:latest

WORKDIR /app

# 从构建阶段把编译好的二进制文件复制过来
COPY --from=builder /app/main .
# 复制配置文件和 Swagger 文档（运行程序必须的静态文件）
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs ./docs

# 声明容器内部监听的端口
EXPOSE 8080

# 启动容器时执行的命令
CMD ["./main"]
```

## `docker-compose.yml`

```yaml
services:
  # 1. 数据库服务
  mysql-db: # 【这里就是房间名，对应 config 里的 host】
    image: mysql:8.0
    container_name: blog-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD}"
      MYSQL_DATABASE: "go_blog_api"
    ports:
      - "3307:3306"
  # 2. Redis 缓存服务
  redis-cache: # 【这里就是房间名，对应 config 里的 addr】
    image: redis:alpine
    container_name: blog-redis
    restart: always
    ports:
      - "6379:6379"
    
    # 3. 我们的 Go 博客后端服务
  blog-app:
    build: .
    container_name: blog-api-server
    restart: always
    ports:
      - "8080:8080"
    volumes:
      # 【魔法在这里】：启动时，把我们刚才写的 config.docker.yaml 替换掉容器里的 config.yaml
      - ./config/config.docker.yaml:/app/config/config.yaml
    depends_on:
        - mysql-db
        - redis-cache
```

## `go.mod`

```go
module go-blog-api

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/redis/go-redis/v9 v9.18.0
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.11.1
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.1
	github.com/swaggo/swag v1.16.6
	golang.org/x/crypto v0.48.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.15.0 // indirect
	github.com/bytedance/sonic/loader v0.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.22.5 // indirect
	github.com/go-openapi/jsonreference v0.21.5 // indirect
	github.com/go-openapi/spec v0.22.4 // indirect
	github.com/go-openapi/swag/conv v0.25.5 // indirect
	github.com/go-openapi/swag/jsonname v0.25.5 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.5 // indirect
	github.com/go-openapi/swag/loading v0.25.5 // indirect
	github.com/go-openapi/swag/stringutils v0.25.5 // indirect
	github.com/go-openapi/swag/typeutils v0.25.5 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.5 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.0 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	go.mongodb.org/mongo-driver/v2 v2.5.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.25.0 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

```

## `go.sum`

```
filippo.io/edwards25519 v1.2.0 h1:crnVqOiS4jqYleHd9vaKZ+HKtHfllngJIiOpNpoJsjo=
filippo.io/edwards25519 v1.2.0/go.mod h1:xzAOLCNug/yB62zG1bQ8uziwrIqIuxhctzJT18Q77mc=
github.com/KyleBanks/depth v1.2.1 h1:5h8fQADFrWtarTdtDudMmGsC7GPbOAu6RVB3ffsVFHc=
github.com/KyleBanks/depth v1.2.1/go.mod h1:jzSb9d0L43HxTQfT+oSA1EEp2q+ne2uh6XgeJcm8brE=
github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
github.com/bsm/gomega v1.27.10/go.mod h1:JyEr/xRbxbtgWNi8tIEVPUYZ5Dzef52k01W3YH0H+O0=
github.com/bytedance/gopkg v0.1.3 h1:TPBSwH8RsouGCBcMBktLt1AymVo2TVsBVCY4b6TnZ/M=
github.com/bytedance/gopkg v0.1.3/go.mod h1:576VvJ+eJgyCzdjS+c4+77QF3p7ubbtiKARP3TxducM=
github.com/bytedance/sonic v1.15.0 h1:/PXeWFaR5ElNcVE84U0dOHjiMHQOwNIx3K4ymzh/uSE=
github.com/bytedance/sonic v1.15.0/go.mod h1:tFkWrPz0/CUCLEF4ri4UkHekCIcdnkqXw9VduqpJh0k=
github.com/bytedance/sonic/loader v0.5.0 h1:gXH3KVnatgY7loH5/TkeVyXPfESoqSBSBEiDd5VjlgE=
github.com/bytedance/sonic/loader v0.5.0/go.mod h1:AR4NYCk5DdzZizZ5djGqQ92eEhCCcdf5x77udYiSJRo=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/cloudwego/base64x v0.1.6 h1:t11wG9AECkCDk5fMSoxmufanudBtJ+/HemLstXDLI2M=
github.com/cloudwego/base64x v0.1.6/go.mod h1:OFcloc187FXDaYHvrNIjxSe8ncn0OOM8gEHfghB2IPU=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f/go.mod h1:cuUVRXasLTGF7a8hSLbxyZXjz+1KgoB3wDUb6vlszIc=
github.com/frankban/quicktest v1.14.6 h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=
github.com/frankban/quicktest v1.14.6/go.mod h1:4ptaffx2x8+WTWXmUCuVU6aPUX1/Mz7zb5vbUoiM6w0=
github.com/fsnotify/fsnotify v1.9.0 h1:2Ml+OJNzbYCTzsxtv8vKSFD9PbJjmhYF14k/jKC7S9k=
github.com/fsnotify/fsnotify v1.9.0/go.mod h1:8jBTzvmWwFyi3Pb8djgCCO5IBqzKJ/Jwo8TRcHyHii0=
github.com/gabriel-vasile/mimetype v1.4.13 h1:46nXokslUBsAJE/wMsp5gtO500a4F3Nkz9Ufpk2AcUM=
github.com/gabriel-vasile/mimetype v1.4.13/go.mod h1:d+9Oxyo1wTzWdyVUPMmXFvp4F9tea18J8ufA774AB3s=
github.com/gin-contrib/gzip v0.0.6 h1:NjcunTcGAj5CO1gn4N8jHOSIeRFHIbn51z6K+xaN4d4=
github.com/gin-contrib/gzip v0.0.6/go.mod h1:QOJlmV2xmayAjkNS2Y8NQsMneuRShOU/kjovCXNuzzk=
github.com/gin-contrib/sse v1.1.0 h1:n0w2GMuUpWDVp7qSpvze6fAu9iRxJY4Hmj6AmBOU05w=
github.com/gin-contrib/sse v1.1.0/go.mod h1:hxRZ5gVpWMT7Z0B0gSNYqqsSCNIJMjzvm6fqCz9vjwM=
github.com/gin-gonic/gin v1.12.0 h1:b3YAbrZtnf8N//yjKeU2+MQsh2mY5htkZidOM7O0wG8=
github.com/gin-gonic/gin v1.12.0/go.mod h1:VxccKfsSllpKshkBWgVgRniFFAzFb9csfngsqANjnLc=
github.com/go-openapi/jsonpointer v0.22.5 h1:8on/0Yp4uTb9f4XvTrM2+1CPrV05QPZXu+rvu2o9jcA=
github.com/go-openapi/jsonpointer v0.22.5/go.mod h1:gyUR3sCvGSWchA2sUBJGluYMbe1zazrYWIkWPjjMUY0=
github.com/go-openapi/jsonreference v0.21.5 h1:6uCGVXU/aNF13AQNggxfysJ+5ZcU4nEAe+pJyVWRdiE=
github.com/go-openapi/jsonreference v0.21.5/go.mod h1:u25Bw85sX4E2jzFodh1FOKMTZLcfifd1Q+iKKOUxExw=
github.com/go-openapi/spec v0.22.4 h1:4pxGjipMKu0FzFiu/DPwN3CTBRlVM2yLf/YTWorYfDQ=
github.com/go-openapi/spec v0.22.4/go.mod h1:WQ6Ai0VPWMZgMT4XySjlRIE6GP1bGQOtEThn3gcWLtQ=
github.com/go-openapi/swag v0.19.15 h1:D2NRCBzS9/pEY3gP9Nl8aDqGUcPFrwG2p+CNFrLyrCM=
github.com/go-openapi/swag/conv v0.25.5 h1:wAXBYEXJjoKwE5+vc9YHhpQOFj2JYBMF2DUi+tGu97g=
github.com/go-openapi/swag/conv v0.25.5/go.mod h1:CuJ1eWvh1c4ORKx7unQnFGyvBbNlRKbnRyAvDvzWA4k=
github.com/go-openapi/swag/jsonname v0.25.5 h1:8p150i44rv/Drip4vWI3kGi9+4W9TdI3US3uUYSFhSo=
github.com/go-openapi/swag/jsonname v0.25.5/go.mod h1:jNqqikyiAK56uS7n8sLkdaNY/uq6+D2m2LANat09pKU=
github.com/go-openapi/swag/jsonutils v0.25.5 h1:XUZF8awQr75MXeC+/iaw5usY/iM7nXPDwdG3Jbl9vYo=
github.com/go-openapi/swag/jsonutils v0.25.5/go.mod h1:48FXUaz8YsDAA9s5AnaUvAmry1UcLcNVWUjY42XkrN4=
github.com/go-openapi/swag/jsonutils/fixtures_test v0.25.5 h1:SX6sE4FrGb4sEnnxbFL/25yZBb5Hcg1inLeErd86Y1U=
github.com/go-openapi/swag/jsonutils/fixtures_test v0.25.5/go.mod h1:/2KvOTrKWjVA5Xli3DZWdMCZDzz3uV/T7bXwrKWPquo=
github.com/go-openapi/swag/loading v0.25.5 h1:odQ/umlIZ1ZVRteI6ckSrvP6e2w9UTF5qgNdemJHjuU=
github.com/go-openapi/swag/loading v0.25.5/go.mod h1:I8A8RaaQ4DApxhPSWLNYWh9NvmX2YKMoB9nwvv6oW6g=
github.com/go-openapi/swag/stringutils v0.25.5 h1:NVkoDOA8YBgtAR/zvCx5rhJKtZF3IzXcDdwOsYzrB6M=
github.com/go-openapi/swag/stringutils v0.25.5/go.mod h1:PKK8EZdu4QJq8iezt17HM8RXnLAzY7gW0O1KKarrZII=
github.com/go-openapi/swag/typeutils v0.25.5 h1:EFJ+PCga2HfHGdo8s8VJXEVbeXRCYwzzr9u4rJk7L7E=
github.com/go-openapi/swag/typeutils v0.25.5/go.mod h1:itmFmScAYE1bSD8C4rS0W+0InZUBrB2xSPbWt6DLGuc=
github.com/go-openapi/swag/yamlutils v0.25.5 h1:kASCIS+oIeoc55j28T4o8KwlV2S4ZLPT6G0iq2SSbVQ=
github.com/go-openapi/swag/yamlutils v0.25.5/go.mod h1:Gek1/SjjfbYvM+Iq4QGwa/2lEXde9n2j4a3wI3pNuOQ=
github.com/go-openapi/testify/enable/yaml/v2 v2.4.0 h1:7SgOMTvJkM8yWrQlU8Jm18VeDPuAvB/xWrdxFJkoFag=
github.com/go-openapi/testify/enable/yaml/v2 v2.4.0/go.mod h1:14iV8jyyQlinc9StD7w1xVPW3CO3q1Gj04Jy//Kw4VM=
github.com/go-openapi/testify/v2 v2.4.0 h1:8nsPrHVCWkQ4p8h1EsRVymA2XABB4OT40gcvAu+voFM=
github.com/go-openapi/testify/v2 v2.4.0/go.mod h1:HCPmvFFnheKK2BuwSA0TbbdxJ3I16pjwMkYkP4Ywn54=
github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
github.com/go-playground/validator/v10 v10.30.1 h1:f3zDSN/zOma+w6+1Wswgd9fLkdwy06ntQJp0BBvFG0w=
github.com/go-playground/validator/v10 v10.30.1/go.mod h1:oSuBIQzuJxL//3MelwSLD5hc2Tu889bF0Idm9Dg26cM=
github.com/go-sql-driver/mysql v1.9.3 h1:U/N249h2WzJ3Ukj8SowVFjdtZKfu9vlLZxjPXV1aweo=
github.com/go-sql-driver/mysql v1.9.3/go.mod h1:qn46aNg1333BRMNU69Lq93t8du/dwxI64Gl8i5p1WMU=
github.com/go-viper/mapstructure/v2 v2.5.0 h1:vM5IJoUAy3d7zRSVtIwQgBj7BiWtMPfmPEgAXnvj1Ro=
github.com/go-viper/mapstructure/v2 v2.5.0/go.mod h1:oJDH3BJKyqBA2TXFhDsKDGDTlndYOZ6rGS0BRZIxGhM=
github.com/goccy/go-json v0.10.5 h1:Fq85nIqj+gXn/S5ahsiTlK3TmC85qgirsdTP/+DeaC4=
github.com/goccy/go-json v0.10.5/go.mod h1:oq7eo15ShAhp70Anwd5lgX2pLfOS3QCiwU/PULtXL6M=
github.com/goccy/go-yaml v1.19.2 h1:PmFC1S6h8ljIz6gMRBopkjP1TVT7xuwrButHID66PoM=
github.com/goccy/go-yaml v1.19.2/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
github.com/golang-jwt/jwt/v5 v5.3.1 h1:kYf81DTWFe7t+1VvL7eS+jKFVWaUnK9cB1qbwn63YCY=
github.com/golang-jwt/jwt/v5 v5.3.1/go.mod h1:fxCRLWMO43lRc8nhHWY6LGqRcf+1gQWArsqaEUEa5bE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/jinzhu/inflection v1.0.0 h1:K317FqzuhWc8YvSVlFMCCUb36O/S9MCKRDI7QkRKD/E=
github.com/jinzhu/inflection v1.0.0/go.mod h1:h+uFLlag+Qp1Va5pdKtLDYj+kHp5pxUVkryuEj+Srlc=
github.com/jinzhu/now v1.1.5 h1:/o9tlHleP7gOFmsnYNz3RGnqzefHA47wQpKrrdTIwXQ=
github.com/jinzhu/now v1.1.5/go.mod h1:d3SSVoowX0Lcu0IBviAWJpolVfI5UJVZZ7cO71lE/z8=
github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
github.com/klauspost/cpuid/v2 v2.3.0 h1:S4CRMLnYUhGeDFDqkGriYKdfoFlDnMtqTiI/sFzhA9Y=
github.com/klauspost/cpuid/v2 v2.3.0/go.mod h1:hqwkgyIinND0mEev00jJYCxPNVRVXFQeu1XKlok6oO0=
github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
github.com/pelletier/go-toml/v2 v2.2.4 h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=
github.com/pelletier/go-toml/v2 v2.2.4/go.mod h1:2gIqNv+qfxSVS7cM2xJQKtLSTLUE9V8t9Stt+h56mCY=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/quic-go/qpack v0.6.0 h1:g7W+BMYynC1LbYLSqRt8PBg5Tgwxn214ZZR34VIOjz8=
github.com/quic-go/qpack v0.6.0/go.mod h1:lUpLKChi8njB4ty2bFLX2x4gzDqXwUpaO1DP9qMDZII=
github.com/quic-go/quic-go v0.59.0 h1:OLJkp1Mlm/aS7dpKgTc6cnpynnD2Xg7C1pwL6vy/SAw=
github.com/quic-go/quic-go v0.59.0/go.mod h1:upnsH4Ju1YkqpLXC305eW3yDZ4NfnNbmQRCMWS58IKU=
github.com/redis/go-redis/v9 v9.18.0 h1:pMkxYPkEbMPwRdenAzUNyFNrDgHx9U+DrBabWNfSRQs=
github.com/redis/go-redis/v9 v9.18.0/go.mod h1:k3ufPphLU5YXwNTUcCRXGxUoF1fqxnhFQmscfkCoDA0=
github.com/rogpeppe/go-internal v1.10.0 h1:TMyTOH3F/DB16zRVcYyreMH6GnZZrwQVAoYjRBZyWFQ=
github.com/rogpeppe/go-internal v1.10.0/go.mod h1:UQnix2H7Ngw/k4C5ijL5+65zddjncjaFoBhdsK/akog=
github.com/sagikazarmark/locafero v0.12.0 h1:/NQhBAkUb4+fH1jivKHWusDYFjMOOKU88eegjfxfHb4=
github.com/sagikazarmark/locafero v0.12.0/go.mod h1:sZh36u/YSZ918v0Io+U9ogLYQJ9tLLBmM4eneO6WwsI=
github.com/spf13/afero v1.15.0 h1:b/YBCLWAJdFWJTN9cLhiXXcD7mzKn9Dm86dNnfyQw1I=
github.com/spf13/afero v1.15.0/go.mod h1:NC2ByUVxtQs4b3sIUphxK0NioZnmxgyCrfzeuq8lxMg=
github.com/spf13/cast v1.10.0 h1:h2x0u2shc1QuLHfxi+cTJvs30+ZAHOGRic8uyGTDWxY=
github.com/spf13/cast v1.10.0/go.mod h1:jNfB8QC9IA6ZuY2ZjDp0KtFO2LZZlg4S/7bzP6qqeHo=
github.com/spf13/pflag v1.0.10 h1:4EBh2KAYBwaONj6b2Ye1GiHfwjqyROoF4RwYO+vPwFk=
github.com/spf13/pflag v1.0.10/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
github.com/spf13/viper v1.21.0 h1:x5S+0EU27Lbphp4UKm1C+1oQO+rKx36vfCoaVebLFSU=
github.com/spf13/viper v1.21.0/go.mod h1:P0lhsswPGWD/1lZJ9ny3fYnVqxiegrlNrEmgLjbTCAY=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
github.com/stretchr/objx v0.5.2/go.mod h1:FRsXN1f5AsAjCGJKqEizvkpNtU+EGNCLh3NxZ/8L+MA=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
github.com/stretchr/testify v1.8.4/go.mod h1:sz/lmYIOXD/1dqDmKjjqLyZ2RngseejIcXlSw2iwfAo=
github.com/stretchr/testify v1.10.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/subosito/gotenv v1.6.0 h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=
github.com/subosito/gotenv v1.6.0/go.mod h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=
github.com/swaggo/files v1.0.1 h1:J1bVJ4XHZNq0I46UU90611i9/YzdrF7x92oX1ig5IdE=
github.com/swaggo/files v1.0.1/go.mod h1:0qXmMNH6sXNf+73t65aKeB+ApmgxdnkQzVTAj2uaMUg=
github.com/swaggo/gin-swagger v1.6.1 h1:Ri06G4gc9N4t4k8hekMigJ9zKTFSlqj/9paAQCQs7cY=
github.com/swaggo/gin-swagger v1.6.1/go.mod h1:LQ+hJStHakCWRiK/YNYtJOu4mR2FP+pxLnILT/qNiTw=
github.com/swaggo/swag v1.16.6 h1:qBNcx53ZaX+M5dxVyTrgQ0PJ/ACK+NzhwcbieTt+9yI=
github.com/swaggo/swag v1.16.6/go.mod h1:ngP2etMK5a0P3QBizic5MEwpRmluJZPHjXcMoj4Xesg=
github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
github.com/ugorji/go/codec v1.3.1 h1:waO7eEiFDwidsBN6agj1vJQ4AG7lh2yqXyOXqhgQuyY=
github.com/ugorji/go/codec v1.3.1/go.mod h1:pRBVtBSKl77K30Bv8R2P+cLSGaTtex6fsA2Wjqmfxj4=
github.com/yuin/goldmark v1.4.13/go.mod h1:6yULJ656Px+3vBD8DxQVa3kxgyrAnzto9xy5taEt/CY=
github.com/zeebo/xxh3 v1.0.2 h1:xZmwmqxHZA8AI603jOQ0tMqmBr9lPeFwGg6d+xy9DC0=
github.com/zeebo/xxh3 v1.0.2/go.mod h1:5NWz9Sef7zIDm2JHfFlcQvNekmcEl9ekUZQQKCYaDcA=
go.mongodb.org/mongo-driver/v2 v2.5.0 h1:yXUhImUjjAInNcpTcAlPHiT7bIXhshCTL3jVBkF3xaE=
go.mongodb.org/mongo-driver/v2 v2.5.0/go.mod h1:yOI9kBsufol30iFsl1slpdq1I0eHPzybRWdyYUs8K/0=
go.uber.org/atomic v1.11.0 h1:ZvwS0R+56ePWxUNi+Atn9dWONBPp/AUETXlHW0DxSjE=
go.uber.org/atomic v1.11.0/go.mod h1:LUxbIzbOniOlMKjJjyPfpl4v+PKK2cNJn91OQbhoJI0=
go.uber.org/mock v0.6.0 h1:hyF9dfmbgIX5EfOdasqLsWD6xqpNZlXblLB/Dbnwv3Y=
go.uber.org/mock v0.6.0/go.mod h1:KiVJ4BqZJaMj4svdfmHM0AUx4NJYO8ZNpPnZn1Z+BBU=
go.yaml.in/yaml/v3 v3.0.4 h1:tfq32ie2Jv2UxXFdLJdh3jXuOzWiL1fo0bu/FbuKpbc=
go.yaml.in/yaml/v3 v3.0.4/go.mod h1:DhzuOOF2ATzADvBadXxruRBLzYTpT36CKvDb3+aBEFg=
golang.org/x/arch v0.25.0 h1:qnk6Ksugpi5Bz32947rkUgDt9/s5qvqDPl/gBKdMJLE=
golang.org/x/arch v0.25.0/go.mod h1:0X+GdSIP+kL5wPmpK7sdkEVTt2XoYP0cSjQSbZBwOi8=
golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
golang.org/x/crypto v0.0.0-20210921155107-089bfa567519/go.mod h1:GvvjBRRGRdwPK5ydBHafDWAxML/pGHZbMvKqRZ5+Abc=
golang.org/x/crypto v0.48.0 h1:/VRzVqiRSggnhY7gNRxPauEQ5Drw9haKdM0jqfcCFts=
golang.org/x/crypto v0.48.0/go.mod h1:r0kV5h3qnFPlQnBSrULhlsRfryS2pmewsg+XfMgkVos=
golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4/go.mod h1:jJ57K6gSWd91VN4djpZkiMVwK6gcyfeH4XE8wZrZaV4=
golang.org/x/mod v0.33.0 h1:tHFzIWbBifEmbwtGz65eaWyGiGZatSrT9prnU8DbVL8=
golang.org/x/mod v0.33.0/go.mod h1:swjeQEj+6r7fODbD2cqrnje9PnziFuw4bmLbBZFrQ5w=
golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20220722155237-a158d28d115b/go.mod h1:XRhObCWvk6IyKnWLug+ECip1KBveYUHfp+8e9klMJ9c=
golang.org/x/net v0.7.0/go.mod h1:2Tu9+aMcznHK/AK1HMvgo6xiTLG5rD5rZLDS+rp2Bjs=
golang.org/x/net v0.51.0 h1:94R/GTO7mt3/4wIKpcR5gkGmRLOuE/2hNGeWq/GBIFo=
golang.org/x/net v0.51.0/go.mod h1:aamm+2QF5ogm02fjy5Bb7CQ0WMt1/WVM7FtyaTLlA9Y=
golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.20.0 h1:e0PTpb7pjO8GAtTs2dQ6jYa5BWYlMuX047Dco/pItO4=
golang.org/x/sync v0.20.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.42.0 h1:omrd2nAlyT5ESRdCLYdm3+fMfNFE/+Rf4bDIQImRJeo=
golang.org/x/sys v0.42.0/go.mod h1:4GL1E5IUh+htKOUEOaiffhrAeqysfVGipDYzABqnCmw=
golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
golang.org/x/term v0.0.0-20210927222741-03fcf44c2211/go.mod h1:jbD1KX2456YbFQfuXm/mYQcufACuNUgVhRMnK/tPxf8=
golang.org/x/term v0.5.0/go.mod h1:jMB1sMXY+tzblOD4FWmEbocvup2/aLOaQEp7JmGp78k=
golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/text v0.3.7/go.mod h1:u+2+/6zg+i71rQMx5EYifcz6MCKuco9NR6JIITiCfzQ=
golang.org/x/text v0.7.0/go.mod h1:mrYo+phRRbMaCq/xk9113O4dZlRixOauAjOtrjsXDZ8=
golang.org/x/text v0.34.0 h1:oL/Qq0Kdaqxa1KbNeMKwQq0reLCCaFtqu2eNuSeNHbk=
golang.org/x/text v0.34.0/go.mod h1:homfLqTYRFyVYemLBFl5GgL/DWEiH5wcsQ5gSh1yziA=
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.1.12/go.mod h1:hNGJHUnrk76NpqgfD5Aqm5Crs+Hm0VOH/i9J2+nxYbc=
golang.org/x/tools v0.42.0 h1:uNgphsn75Tdz5Ji2q36v/nsFSfR/9BRFvqhGBaJGd5k=
golang.org/x/tools v0.42.0/go.mod h1:Ma6lCIwGZvHK6XtgbswSoWroEkhugApmsXyrUmBhfr0=
golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
google.golang.org/protobuf v1.36.11 h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=
google.golang.org/protobuf v1.36.11/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=
gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gorm.io/driver/mysql v1.6.0 h1:eNbLmNTpPpTOVZi8MMxCi2aaIm0ZpInbORNXDwyLGvg=
gorm.io/driver/mysql v1.6.0/go.mod h1:D/oCC2GWK3M/dqoLxnOlaNKmXz8WNTfcS9y5ovaSqKo=
gorm.io/gorm v1.31.1 h1:7CA8FTFz/gRfgqgpeKIBcervUn3xSyPUmr6B2WXJ7kg=
gorm.io/gorm v1.31.1/go.mod h1:XyQVbO2k6YkOis7C2437jSit3SsDK72s7n7rsSHd+Gs=

```

## `cmd/main.go`

```go
package main

import (
	"log"

	"go-blog-api/config"
	_ "go-blog-api/docs"
	"go-blog-api/internal/cache"
	"go-blog-api/internal/db"
	"go-blog-api/internal/router"
)

// @title Go 博客后端 API
// @version 1.0
// @description 这是一个基于 Go + Gin + GORM 开发的博客项目后端 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email awhg23@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// [新增] 1.初始化配置
	config.InitConfig()

	// 2.初始化数据库连接（此时 db.go 中才能安全读取到 config.App.Database)
	db.InitDB()

	// [新增] 3.初始化 Redis 连接
	cache.InitRedis()

	// 4.初始化路由
	r := router.SetupRouter()

	// 5.从配置中读取端口并启动服务器
	port := ":" + config.App.Server.Port
	log.Printf("服务器启动中，监听端口 %s...", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

```

## `config/config.go`

```go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var App *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	App = &Config{}
	if err := viper.Unmarshal(App); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	log.Println("配置文件加载成功!")
}

```

## `config/config.exmaple.yaml`

```yaml
server:
  port: "8080"

database:
  host: "127.0.0.1"
  port: "3306"
  user: "root"
  password: "your_mysql_password_here"
  dbname: "go_blog_api"

JWT:
  secret: "my_super_secret_key_123456"



```

## `config/config.docker.yaml.example`

```
server:
  port: "8080"

database:
  host: "mysql-db" # 【核心改变】将 127.0.0.1 改成 MySQL 容器的名字
  port: "3306"
  user: "root"
  password: "your_database_password_here"
  dbname: "go_blog_api"

jwt:
  secret: "my_super_secret_key_123456"

redis:
  addr: "redis-cache:6379" # 【核心改变】将 127.0.0.1 改成 Redis 容器的名字
  password: ""
  db: 0
```

## `internal/cache/redis.go`

```go
package cache

import (
	"context"
	"log"

	"go-blog-api/config"

	"github.com/redis/go-redis/v9"
)

// RDB 是全局的 Redis 客户端实例
var RDB *redis.Client

// InitReids 初始化 Redis 连接
func InitRedis() {
	cfg := config.App.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 连接成功!")
}

```

## `internal/db/db.go`

```go
package db

import (
	"fmt"
	"go-blog-api/internal/model"
	"log"

	"go-blog-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 从全局配置中获取数据库信息
	dbCfg := config.App.Database

	// DSN (Data Source Name): 账号:密码@tcp(地址:端口)/数据库名?参数
	// 动态拼接DNS
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	// 自动迁移（根据模型创建表）
	DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	log.Println("数据库连接成功")
}

```

## `internal/handler/user.go`

```go
package handler

import (
	"net/http"

	"go-blog-api/internal/db"
	"go-blog-api/internal/model"
	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// Register 处理用户注册请求
// @Summary 用户注册
// @Description 注册新用户（账号需3-20位，密码需6-20位)
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "账号和密码"
// @Success 200 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "参数错误，账号需3-20位，密码需6-20位"
// @Failure 409 {object} map[string]interface{} "用户名已存在"
// @Failure 500 {object} map[string]interface{} "密码加密失败"
// @Failure 500 {object} map[string]interface{} "用户创建失败"
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误（账号需3-20位，密码需6-20位）"})
		return
	}
	// 检查用户名是否已存在
	var count int64
	db.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user := model.User{
		Username:       req.Username,
		PasswordDigest: string(hashedPassword),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 处理用户登录请求
// @Summary 用户登录
// @Description 校验账号密码并返回 JWT Token
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body LoginRequest true "账号和密码"
// @Success 200 {object} map[string]interface{} "登录成功，返回 JWT Token"
// @Failure 400 {object} map[string]interface{} "参数错误（账号和密码不能为空）"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Failure 500 {object} map[string]interface{} "系统异常，生成 Token 失败"
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误，请提供账号和密码"})
		return
	}

	var user model.User
	// 根据用户名查询用户
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	// 密码校验
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统异常，生成Token失败"})
		return
	}

	//把Token返回给前端
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

```

## `internal/handler/post.go`

```go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-blog-api/internal/cache"
	"go-blog-api/internal/db"
	"go-blog-api/internal/model"
	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required"`
}

// CreatePost 处理发布文章请求
// @Summary 发布新文章
// @Description 发布一篇新的博客文章（需登录验证）
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer // 这个标签说明该接口需要右上角的 Token 鉴权
// @Param post body CreatePostRequest true "文章标题和内容"
// @Success 200 {object} map[string]interface{}	"文章发布成功"
// @Failure 400 {object} map[string]interface{} "参数错误（标题和内容不为空）"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 500 {object} map[string]interface{} "发布文章失败"
// @Router /posts [post]
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误（标题和内容不为空）"})
		return
	}
	//核心逻辑：从上下文中获取刚刚中间件塞进去的userID
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//构建要插入数据库的文章模型
	post := model.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}
	//存入数据库
	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章发布成功",
		"data": gin.H{
			"post_id": post.ID,
			"title":   post.Title,
		},
	})
}

// GetPosts 获取文章列表
// @Summary 获取文章列表（分页）
// @Description 获取文章列表，支持分页和缓存（无需登录，默认每页10条）
// @Tags 文章模块
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}	"文章列表获取成功"
// @Failure 500 {object} map[string]interface{} "获取文章列表失败"
// @Router /posts [get]
func GetPosts(c *gin.Context) {
	//1.获取分页参数（默认 page=1，size=10)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	//防御性校验：page和size不能小于1
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	} //防止单次请求过多导致崩溃

	// ====================[新增：1. 缓存拦截层] ====================
	ctx := context.Background()
	// 构建这个分页专属的 Redis Key，比如“posts:page:1:size:10"
	cacheKey := fmt.Sprintf("posts:page:%d:size:%d", page, size)

	// 尝试从 Redis 中获取数据
	cachedData, err := cache.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		// [缓存命中 Cache Hit]！
		// 极致性能优化：因为存在 Redis 里的直接就是 JSON 字符串
		// 所以不需要反序列化，直接指定 Header 返回给前端就行
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedData)
		return
	}
	// ==============================================================

	// ==================== [原有：2. 数据库查询层] ====================
	// 如果代码走到这，说明【缓存未命中 Cache Miss】（或者是第一次访问，或者缓存过期了）
	offset := (page - 1) * size

	var posts []model.Post
	var total int64

	//2.获取文章总数（用于前端分页组件计算总页数)
	db.DB.Model(&model.Post{}).Count(&total)

	//3.分页查询并预加载作者信息
	if err := db.DB.Preload("User").
		Order("created_at desc").
		Limit(size).
		Offset(offset).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	//4.返回结果
	response := gin.H{
		"data": posts,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
		},
	}
	// ==============================================================

	// ====================[新增：3. 缓存回写层] ====================
	// 将组装好的 response 序列化为 JSON 字符串，存入 Redis
	if jsonData, err := json.Marshal(response); err == nil {
		// 设置缓存过期时间（如 60 秒）
		// 这样既能抵挡这 60 秒内的高并发，又保证了 60 秒后能拉取到别人发的新文章
		cache.RDB.Set(ctx, cacheKey, jsonData, 60*time.Second)
	}
	// ==============================================================

	// 最后把本次查到的数据返回给前端（一般只有触发查库的第一名用户才会走到这）
	c.JSON(http.StatusOK, response)
}

// UpdatePost 修改文章（需要 JWT鉴权）
// @Summary 修改文章
// @Description 修改指定 ID 的文章（需登录，只能修改自己的文章）
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "文章ID"
// @Param post body CreatePostRequest true "修改后的文章标题和内容"
// @Success 200 {object} map[string]interface{} "文章更新成功"
// @Failure 400 {object} map[string]interface{} "参数错误（标题和内容不为空）"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 403 {object} map[string]interface{} "越权操作：只能修改自己的文章"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "更新文章失败"
// @Router /posts/:id [put]
func UpdatePost(c *gin.Context) {
	postID := c.Param("id")

	//1.获取当前登录的 userID （从 JWT 中间件中取得）
	currentUserID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//2.绑定请求体参数
	var req struct {
		Title   string `json:"title" binding:"required,min=1,max=100"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不合法"})
		return
	}

	//3.校验越权：查询文章是否存在，并判断 user_id 是否匹配
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	//【核心安全校验】：防止A修改B的文章
	if post.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "越权操作：只能修改自己的文章"})
		return
	}

	//4.执行更新
	//使用 map 更新可以避免 GORM 忽略零值的问题，或者直接传入结构体
	if err := db.DB.Model(&post).Updates(map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    post,
	})
}

// DeletePost 删除文章 （需要 JWT 鉴权）
// @Summary 删除文章
// @Description 删除指定 ID 的文章（需登录，只能删除自己的文章）
// @Tags 文章模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "文章删除成功"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 403 {object} map[string]interface{} "越权操作：只能删除自己的文章"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "删除文章失败"
// @Router /posts/:id [delete]
func DeletePost(c *gin.Context) {
	postID := c.Param("id")

	//1. 获取当前登录的 userID （从 JWT 中间件中取得）
	currentUserID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//2.校验越权：查询文章并检查归属
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if post.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "越权操作，只能删除自己的文章"})
		return
	}

	//3.执行删除（这里是物理删除，想要软删除可以在 Model 中引入 gormDeletedAt）
	if err := db.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

```

## `internal/handler/comment.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"go-blog-api/internal/db"
	"go-blog-api/internal/model"
	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}

// ==================== 1. 发表评论 (需要 JWT 鉴权) ====================
// CreateComment 发表评论
// @Summary 发表评论
// @Description 发表评论（需要登录）
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateCommentRequest true "评论内容"
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "评论成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "发表评论失败"
// @Router /posts/:id/comments [post]
func CreateComment(c *gin.Context) {
	// 1. 获取 URL 路径中的文章 ID，例如 /api/posts/1/comments
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	//2.检查文章是否存在（防御性编程）
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	//3.从上下文获取当前登录的用户 ID （复用之前的 switch 断言逻辑）
	currentUserID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//4.解析前端传来的评论内容
	var req struct {
		Content string `json:"content" binding:"required,max=500"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空且不能超过500字"})
		return
	}

	//5.组装并入库
	comment := model.Comment{
		PostID:  postID,
		UserID:  currentUserID,
		Content: req.Content,
	}
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发表评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论成功",
		"data":    comment,
	})
}

// ==================== 2. 获取某篇文章的评论列表 (公开接口) ====================
// GetPostComments 获取某篇文章的评论列表
// @Summary 获取某篇文章的评论列表（分页）
// @Description 获取某篇文章的评论列表（公开接口）
// @Tags 评论模块
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Success 200 {object} map[string]interface{} "返回评论列表和分页信息"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "获取评论列表失败"
// @Router /posts/:id/comments [get]
func GetPostComments(c *gin.Context) {
	postID := c.Param("id")

	//检查文章是否存在（防御性编程）
	var post model.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 获取分页参数（默认 page=1，size=10）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	//防御性校验： page和size不能小于1
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	offset := (page - 1) * size

	var total int64
	var comments []model.Comment

	//获取评论数量（用于前端分页组件计算总页数）
	db.DB.Model(&model.Comment{}).
		Where("post_id = ?", postID).
		Count(&total)

	// 按照创建时间倒序排列（最新的评论在最上面），并且预加载评论的作者信息
	if err := db.DB.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at desc").
		Limit(size).
		Offset(offset).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"data": comments,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// ==================== 3. 删除评论 (需要 JWT 鉴权) ====================
// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除评论（需要登录，只能删除自己的评论）
// @Tags 评论模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "评论ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 401 {object} map[string]interface{} "未登录或Token无效"
// @Failure 403 {object} map[string]interface{} "越权操作：只能删除自己的评论"
// @Failure 404 {object} map[string]interface{} "评论不存在"
// @Failure 500 {object} map[string]interface{} "删除失败"
// @Router /comments/:id [delete]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	//1.获取当前登录用户ID
	currentUserID, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return
	}

	//2.查找评论
	var comment model.Comment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	//3.越权校验：只能删除自己的评论
	if comment.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "越权操作：只能删除自己的评论"})
		return
	}

	//4.执行删除
	if err := db.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}

```

## `internal/middleware/jwt_auth.go`

```go
package middleware

import (
	"net/http"
	"strings"

	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 JWT 的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.从 HTTP 请求头中获取 Authorization 字段
		//标准的 Token 格式是放在 Header 里：Authorization：Bearer xxxxx.yyyyy.zzzzz
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头中缺失 Authorization 字段，请先登录"})
			c.Abort()
			return
		}
		//2. 按空格分割，提取出真正的 Token 字符串
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 格式错误，应为 Bearer <Token>"})
			c.Abort()
			return
		}
		//parts[1] 就是我们要的 jwt Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期，请重新登录"})
			c.Abort()
			return
		}
		//3. 验证通过！把提取出的 userID 存入 Gin 的上下文（Context）中
		//这样后续的业务接口（比如发文章）就能直接从Context里知道当前操作的是谁
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		//4.放行，继续执行后面的业务逻辑
		c.Next()

	}
}

```

## `internal/middleware/jwt_auth_test.go`

```go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-blog-api/config"
	"go-blog-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	// 设置 Gin 为测试模式，减少多余日志
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 挂载要测试的中间件
	r.Use(JWTAuthMiddleware())

	// 写一个假的受保护的接口
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"message": "success", "userID": userID})
	})
	return r
}

func TestJWTAuthMiddleware(t *testing.T) {
	// 初始化测试配置
	config.App = &config.Config{
		JWT: config.JWTConfig{Secret: "test_secret"},
	}

	router := setupTestRouter()

	t.Run("不带 Token 应该返回 401", func(t *testing.T) {
		// 1. 创建一个模拟的 HTTP 请求 （GET /protected）
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		// 2. 创建一个模拟的相应记录器（用于接收返回值）
		w := httptest.NewRecorder()

		// 3. 将请求发给 Gin 路由
		router.ServeHTTP(w, req)

		// 4. 断言：预期被中间件拦截，返回401
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("携带无效 Token 应该返回 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer fake_invalid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("携带有效的 Token 应该返回 200 并放行", func(t *testing.T) {
		// 生成一个合法的 Token
		validToken, _ := utils.GenerateToken(1024, "test_user")

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 断言：预期中间件放行，成功到达接口，返回 200
		assert.Equal(t, http.StatusOK, w.Code)
		// 还可以进一步断言返回的 JSON 中包含正确的 userID
		assert.Contains(t, w.Body.String(), "1024")
	})
}

```

## `internal/model/user.go`

```go
package model

import "time"

// User 对应数据库里的 users 表
type User struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	Username       string    `gorm:"type:varchar(50);not null;unique"`
	PasswordDigest string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt      time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

```

## `internal/model/post.go`

```go
package model

import "time"

// Post 对应数据库里的 posts 表
type Post struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"` //加上普通索引加速查询
	User      User      `gorm:"foreignKey:UserID"`
	Title     string    `gorm:"type:varchar(100);not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

```

## `internal/model/comment.go`

```go
package model

import "time"

type Comment struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int64     `gorm:"not null;index" json:"post_id"` // 关联的文章ID(索引加速查询)
	UserID    int64     `gorm:"not null;index" json:"user_id"` // 评论的作者ID(索引加速查询)
	User      User      `gorm:"foreignKey:UserID" json:"user"` //预加载关联：评论者信息
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

```

## `internal/router/router.go`

```go
package router

import (
	"go-blog-api/internal/handler"
	"go-blog-api/internal/middleware"

	_ "go-blog-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// [新增] 注册 Swagger 的路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api")
	{
		// ========= 公开接口（所有人可访问） =========
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		api.GET("/posts", handler.GetPosts)

		api.GET("/posts/:id/comments", handler.GetPostComments)
		// ========= 私有接口（需要 JWT Token 鉴权） =========
		authApi := api.Group("")
		authApi.Use(middleware.JWTAuthMiddleware())
		{
			authApi.POST("/posts", handler.CreatePost)
			authApi.PUT("/posts/:id", handler.UpdatePost)
			authApi.DELETE("/posts/:id", handler.DeletePost)
			authApi.POST("/posts/:id/comments", handler.CreateComment)
			authApi.DELETE("/comments/:id", handler.DeleteComment)
		}
	}

	return r
}

```

## `internal/utils/context.go`

```go
package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(c *gin.Context) (int64, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或Token无效"})
		return 0, errors.New("未登录或Token无效")
	}

	switch v := userIDValue.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "上下文中无效的userID类型"})
		return 0, errors.New("上下文中无效的userID类型")
	}
}

```

## `internal/utils/jwt.go`

```go
package utils

import (
	"go-blog-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义我们要在 Token 里携带的信息（载荷）
type Claims struct {
	UserID               int64  `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        //JWT官方内置的字段 ？？
}

func GenerateToken(userID int64, username string) (string, error) {
	//1.设置Token过期时间
	expirationTime := time.Now().Add(24 * time.Hour)
	//组装要装进Token的数据
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), //过期时间
			Issuer:    "go-blog-api",                      //签发人
		},
	}

	//使用 HS256 算法生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置里的 Secret
	secret := []byte(config.App.JWT.Secret)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否正确
		return []byte(config.App.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	// 验证Token是否有效，并提取CLaims的数据
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

```

## `internal/utils/jwt_test.go`

```go
package utils

import (
	"go-blog-api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateAndParseToken 测试 JWT 的生成和解析逻辑是否闭环
func TestGenerateAndParseToken(t *testing.T) {
	// 1.准备工作：自定义一个配置，防止空指针报错，因为测试环境未执行 main.go 中的 InitConfig 函数
	config.App = &config.Config{
		JWT: config.JWTConfig{
			Secret: "test_super_secret_key",
		},
	}

	testUserID := int64(888)
	testUsername := "golang_tester"

	// 2.执行 Token 生成
	token, err := GenerateToken(testUserID, testUsername)

	// 断言（Assert）：预期 err 为 nil，且 token 不为空
	assert.NoError(t, err, "生成 Token 时不应返回错误")
	assert.NotEmpty(t, token, "生成的 Token 不应为空")

	// 3.执行解析 Token 测试
	claims, err := ParseToken(token)

	// 断言：解析不能报错，且解析出来的数据必须和存进去的一模一样！
	assert.NoError(t, err, "解析 Token 时不应返回错误")
	assert.Equal(t, testUserID, claims.UserID, "解析出的 UserID 不匹配")
	assert.Equal(t, testUsername, claims.Username, "解析出的 Username 不匹配")
}

```

## `docs/docs.go`

```go
// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "awhg23@outlook.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/comments/:id": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除评论（需要登录，只能删除自己的评论）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "评论模块"
                ],
                "summary": "删除评论",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "评论ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "删除成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "403": {
                        "description": "越权操作：只能删除自己的评论",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "评论不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "删除失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "校验账号密码并返回 JWT Token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "账号和密码",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功，返回 JWT Token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误（账号和密码不能为空）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "用户名或密码错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "系统异常，生成 Token 失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/posts": {
            "get": {
                "description": "获取文章列表，支持分页和缓存（无需登录，默认每页10条）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "获取文章列表（分页）",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章列表获取成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "获取文章列表失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer // 这个标签说明该接口需要右上角的 Token 鉴权": []
                    }
                ],
                "description": "发布一篇新的博客文章（需登录验证）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "发布新文章",
                "parameters": [
                    {
                        "description": "文章标题和内容",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreatePostRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章发布成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误（标题和内容不为空）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "发布文章失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/posts/:id": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "修改指定 ID 的文章（需登录，只能修改自己的文章）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "修改文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "修改后的文章标题和内容",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreatePostRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章更新成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误（标题和内容不为空）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "403": {
                        "description": "越权操作：只能修改自己的文章",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "更新文章失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除指定 ID 的文章（需登录，只能删除自己的文章）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "删除文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章删除成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "403": {
                        "description": "越权操作：只能删除自己的文章",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "删除文章失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/posts/:id/comments": {
            "get": {
                "description": "获取某篇文章的评论列表（公开接口）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "评论模块"
                ],
                "summary": "获取某篇文章的评论列表（分页）",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "页码，默认1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数量，默认10",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回评论列表和分页信息",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "获取评论列表失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "发表评论（需要登录）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "评论模块"
                ],
                "summary": "发表评论",
                "parameters": [
                    {
                        "description": "评论内容",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateCommentRequest"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "评论成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "发表评论失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "注册新用户（账号需3-20位，密码需6-20位)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "账号和密码",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "注册成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误，账号需3-20位，密码需6-20位",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "409": {
                        "description": "用户名已存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "用户创建失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.CreateCommentRequest": {
            "type": "object",
            "required": [
                "content"
            ],
            "properties": {
                "content": {
                    "type": "string",
                    "maxLength": 500
                }
            }
        },
        "handler.CreatePostRequest": {
            "type": "object",
            "required": [
                "content",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                }
            }
        },
        "handler.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handler.RegisterRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                },
                "username": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 3
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Go 博客后端 API",
	Description:      "这是一个基于 Go + Gin + GORM 开发的博客项目后端 API 文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

```

## `docs/swagger.json`

```json
{
    "swagger": "2.0",
    "info": {
        "description": "这是一个基于 Go + Gin + GORM 开发的博客项目后端 API 文档",
        "title": "Go 博客后端 API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "awhg23@outlook.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api",
    "paths": {
        "/comments/:id": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除评论（需要登录，只能删除自己的评论）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "评论模块"
                ],
                "summary": "删除评论",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "评论ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "删除成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "403": {
                        "description": "越权操作：只能删除自己的评论",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "评论不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "删除失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "校验账号密码并返回 JWT Token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "账号和密码",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功，返回 JWT Token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误（账号和密码不能为空）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "用户名或密码错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "系统异常，生成 Token 失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/posts": {
            "get": {
                "description": "获取文章列表，支持分页和缓存（无需登录，默认每页10条）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "获取文章列表（分页）",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章列表获取成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "获取文章列表失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer // 这个标签说明该接口需要右上角的 Token 鉴权": []
                    }
                ],
                "description": "发布一篇新的博客文章（需登录验证）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "发布新文章",
                "parameters": [
                    {
                        "description": "文章标题和内容",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreatePostRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章发布成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误（标题和内容不为空）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "发布文章失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/posts/:id": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "修改指定 ID 的文章（需登录，只能修改自己的文章）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "修改文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "修改后的文章标题和内容",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreatePostRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章更新成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误（标题和内容不为空）",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "403": {
                        "description": "越权操作：只能修改自己的文章",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "更新文章失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除指定 ID 的文章（需登录，只能删除自己的文章）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文章模块"
                ],
                "summary": "删除文章",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "文章删除成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "403": {
                        "description": "越权操作：只能删除自己的文章",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "删除文章失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/posts/:id/comments": {
            "get": {
                "description": "获取某篇文章的评论列表（公开接口）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "评论模块"
                ],
                "summary": "获取某篇文章的评论列表（分页）",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "页码，默认1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数量，默认10",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回评论列表和分页信息",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "获取评论列表失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "发表评论（需要登录）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "评论模块"
                ],
                "summary": "发表评论",
                "parameters": [
                    {
                        "description": "评论内容",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateCommentRequest"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "文章ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "评论成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "未登录或Token无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "文章不存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "发表评论失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "注册新用户（账号需3-20位，密码需6-20位)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "账号和密码",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "注册成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "参数错误，账号需3-20位，密码需6-20位",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "409": {
                        "description": "用户名已存在",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "用户创建失败",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.CreateCommentRequest": {
            "type": "object",
            "required": [
                "content"
            ],
            "properties": {
                "content": {
                    "type": "string",
                    "maxLength": 500
                }
            }
        },
        "handler.CreatePostRequest": {
            "type": "object",
            "required": [
                "content",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 1
                }
            }
        },
        "handler.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handler.RegisterRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                },
                "username": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 3
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}
```

## `docs/swagger.yaml`

```yaml
basePath: /api
definitions:
  handler.CreateCommentRequest:
    properties:
      content:
        maxLength: 500
        type: string
    required:
    - content
    type: object
  handler.CreatePostRequest:
    properties:
      content:
        type: string
      title:
        maxLength: 100
        minLength: 1
        type: string
    required:
    - content
    - title
    type: object
  handler.LoginRequest:
    properties:
      password:
        type: string
      username:
        type: string
    required:
    - password
    - username
    type: object
  handler.RegisterRequest:
    properties:
      password:
        maxLength: 20
        minLength: 6
        type: string
      username:
        maxLength: 20
        minLength: 3
        type: string
    required:
    - password
    - username
    type: object
host: localhost:8080
info:
  contact:
    email: awhg23@outlook.com
    name: API Support
  description: 这是一个基于 Go + Gin + GORM 开发的博客项目后端 API 文档
  license:
    name: Apache 2.0
    url: http://www.apache.org/licenses/LICENSE-2.0.html
  termsOfService: http://swagger.io/terms/
  title: Go 博客后端 API
  version: "1.0"
paths:
  /comments/:id:
    delete:
      consumes:
      - application/json
      description: 删除评论（需要登录，只能删除自己的评论）
      parameters:
      - description: 评论ID
        in: path
        name: id
        required: true
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: 删除成功
          schema:
            additionalProperties: true
            type: object
        "401":
          description: 未登录或Token无效
          schema:
            additionalProperties: true
            type: object
        "403":
          description: 越权操作：只能删除自己的评论
          schema:
            additionalProperties: true
            type: object
        "404":
          description: 评论不存在
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 删除失败
          schema:
            additionalProperties: true
            type: object
      security:
      - Bearer: []
      summary: 删除评论
      tags:
      - 评论模块
  /login:
    post:
      consumes:
      - application/json
      description: 校验账号密码并返回 JWT Token
      parameters:
      - description: 账号和密码
        in: body
        name: request
        required: true
        schema:
          $ref: '#/definitions/handler.LoginRequest'
      produces:
      - application/json
      responses:
        "200":
          description: 登录成功，返回 JWT Token
          schema:
            additionalProperties: true
            type: object
        "400":
          description: 参数错误（账号和密码不能为空）
          schema:
            additionalProperties: true
            type: object
        "401":
          description: 用户名或密码错误
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 系统异常，生成 Token 失败
          schema:
            additionalProperties: true
            type: object
      summary: 用户登录
      tags:
      - 用户模块
  /posts:
    get:
      consumes:
      - application/json
      description: 获取文章列表，支持分页和缓存（无需登录，默认每页10条）
      parameters:
      - default: 1
        description: 页码
        in: query
        name: page
        type: integer
      - default: 10
        description: 每页数量
        in: query
        name: size
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: 文章列表获取成功
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 获取文章列表失败
          schema:
            additionalProperties: true
            type: object
      summary: 获取文章列表（分页）
      tags:
      - 文章模块
    post:
      consumes:
      - application/json
      description: 发布一篇新的博客文章（需登录验证）
      parameters:
      - description: 文章标题和内容
        in: body
        name: post
        required: true
        schema:
          $ref: '#/definitions/handler.CreatePostRequest'
      produces:
      - application/json
      responses:
        "200":
          description: 文章发布成功
          schema:
            additionalProperties: true
            type: object
        "400":
          description: 参数错误（标题和内容不为空）
          schema:
            additionalProperties: true
            type: object
        "401":
          description: 未登录或Token无效
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 发布文章失败
          schema:
            additionalProperties: true
            type: object
      security:
      - Bearer // 这个标签说明该接口需要右上角的 Token 鉴权: []
      summary: 发布新文章
      tags:
      - 文章模块
  /posts/:id:
    delete:
      consumes:
      - application/json
      description: 删除指定 ID 的文章（需登录，只能删除自己的文章）
      parameters:
      - description: 文章ID
        in: path
        name: id
        required: true
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: 文章删除成功
          schema:
            additionalProperties: true
            type: object
        "401":
          description: 未登录或Token无效
          schema:
            additionalProperties: true
            type: object
        "403":
          description: 越权操作：只能删除自己的文章
          schema:
            additionalProperties: true
            type: object
        "404":
          description: 文章不存在
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 删除文章失败
          schema:
            additionalProperties: true
            type: object
      security:
      - Bearer: []
      summary: 删除文章
      tags:
      - 文章模块
    put:
      consumes:
      - application/json
      description: 修改指定 ID 的文章（需登录，只能修改自己的文章）
      parameters:
      - description: 文章ID
        in: path
        name: id
        required: true
        type: integer
      - description: 修改后的文章标题和内容
        in: body
        name: post
        required: true
        schema:
          $ref: '#/definitions/handler.CreatePostRequest'
      produces:
      - application/json
      responses:
        "200":
          description: 文章更新成功
          schema:
            additionalProperties: true
            type: object
        "400":
          description: 参数错误（标题和内容不为空）
          schema:
            additionalProperties: true
            type: object
        "401":
          description: 未登录或Token无效
          schema:
            additionalProperties: true
            type: object
        "403":
          description: 越权操作：只能修改自己的文章
          schema:
            additionalProperties: true
            type: object
        "404":
          description: 文章不存在
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 更新文章失败
          schema:
            additionalProperties: true
            type: object
      security:
      - Bearer: []
      summary: 修改文章
      tags:
      - 文章模块
  /posts/:id/comments:
    get:
      consumes:
      - application/json
      description: 获取某篇文章的评论列表（公开接口）
      parameters:
      - description: 文章ID
        in: path
        name: id
        required: true
        type: integer
      - description: 页码，默认1
        in: query
        name: page
        type: integer
      - description: 每页数量，默认10
        in: query
        name: size
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: 返回评论列表和分页信息
          schema:
            additionalProperties: true
            type: object
        "404":
          description: 文章不存在
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 获取评论列表失败
          schema:
            additionalProperties: true
            type: object
      summary: 获取某篇文章的评论列表（分页）
      tags:
      - 评论模块
    post:
      consumes:
      - application/json
      description: 发表评论（需要登录）
      parameters:
      - description: 评论内容
        in: body
        name: request
        required: true
        schema:
          $ref: '#/definitions/handler.CreateCommentRequest'
      - description: 文章ID
        in: path
        name: id
        required: true
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: 评论成功
          schema:
            additionalProperties: true
            type: object
        "400":
          description: 参数错误
          schema:
            additionalProperties: true
            type: object
        "401":
          description: 未登录或Token无效
          schema:
            additionalProperties: true
            type: object
        "404":
          description: 文章不存在
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 发表评论失败
          schema:
            additionalProperties: true
            type: object
      security:
      - Bearer: []
      summary: 发表评论
      tags:
      - 评论模块
  /register:
    post:
      consumes:
      - application/json
      description: 注册新用户（账号需3-20位，密码需6-20位)
      parameters:
      - description: 账号和密码
        in: body
        name: request
        required: true
        schema:
          $ref: '#/definitions/handler.RegisterRequest'
      produces:
      - application/json
      responses:
        "200":
          description: 注册成功
          schema:
            additionalProperties: true
            type: object
        "400":
          description: 参数错误，账号需3-20位，密码需6-20位
          schema:
            additionalProperties: true
            type: object
        "409":
          description: 用户名已存在
          schema:
            additionalProperties: true
            type: object
        "500":
          description: 用户创建失败
          schema:
            additionalProperties: true
            type: object
      summary: 用户注册
      tags:
      - 用户模块
securityDefinitions:
  Bearer:
    description: Type "Bearer" followed by a space and JWT token.
    in: header
    name: Authorization
    type: apiKey
swagger: "2.0"

```

