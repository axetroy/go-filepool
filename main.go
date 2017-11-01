package main

import (
	"github.com/suboat/go-filepool/config"
	"github.com/suboat/go-filepool/upload"
	"github.com/suboat/sorm/log"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"path/filepath"
	"github.com/Sirupsen/logrus"
	"net/http"
	"strings"
)

var (
	CfgMap *mainCfg = nil
)

type mainCfg struct {
	FilePath    string `ini:"-"`
	Address     string
	UploadUrl   string
	DownloadUrl string
}

func newMainCfg(s *mainCfg) (d *mainCfg, err error) {
	d = &mainCfg{
		Address:     "0.0.0.0:9091",
		UploadUrl:   "/upload",
		DownloadUrl: "/download/:size/:file",
	}
	d.FilePath, _ = filepath.Abs(path.Join(path.Dir(os.Args[0]), "./"))
	d.FilePath = path.Join(d.FilePath, "config.ini")
	return
}

func main() {

	var (
		h = &upload.UploadHandler{}
	)

	CfgMap, _ = newMainCfg(nil)

	// set log level
	log.SetLevel(logrus.InfoLevel)

	h.FormName = "file"          // 文件名
	h.RequireImage = true        // 要求是图片
	h.MaxSize = 20 * 1024 * 1024 // 大小限制 20MB

	r := gin.Default()

	/**
	下载图片
	 */
	r.GET(CfgMap.DownloadUrl, func(c *gin.Context) {

		size := c.Param("size")
		file := c.Param("file")

		// 如果路径是以/结尾的话，那么这个图片不存在
		if strings.HasSuffix(c.Request.URL.Path, "/") {
			http.NotFound(c.Writer, c.Request)
			return
		}

		// 尺寸不对
		if size != "origin" && size != "thumbnail" {
			http.NotFound(c.Writer, c.Request)
			return
		}

		fileType := size

		absFilePath := path.Join("./upload", fileType, file)

		if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
			// if the path not found
			http.NotFound(c.Writer, c.Request)
			return
		}

		http.ServeFile(c.Writer, c.Request, absFilePath)
	})

	/**
	上传图片
	 */
	r.POST(CfgMap.UploadUrl, func(c *gin.Context) {
		log.Info("FileDir:", config.DownloadDir, " DownloadUrl:", CfgMap.DownloadUrl)
		h.ServeHTTP(c.Writer, c.Request)
	})

	// 检查缺失的缩略图
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		if err := upload.ToolFixThumbnail(); err != nil {
			log.Error(err)
		}
	}()

	r.Run(CfgMap.Address) // listen and serve on 0.0.0.0:8080

	// run
	log.Info("ListenAndServe ", CfgMap.Address)
}
