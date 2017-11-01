## GO 文件上传库实例

### 使用

```bash
git clone https://github.com/axetroy/go-filepool.git
cd ./go-filepool
go run main.go
```

### /upload  POST

上传文件

### /download/:size/:file   [GET]

获取上传的文件

- size:
  - origin: 获取原始图片
  - thumbnail: 获取缩略图

- file: 文件哈希+文件后缀名

