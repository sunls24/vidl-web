package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vidl-web/common"
	"vidl-web/common/logger"
)

var usefulKeys = []string{"title", "description", "duration", "view_count", "uploader", "upload_date", "webpage_url", "id"}

func Analyze(c *gin.Context) {
	link, exit := common.Query(c, "link")
	if exit {
		return
	}
	output, err := common.Exec(common.YTBin, "-j", "--cookies", "cookies.txt", link)
	if err != nil {
		logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, common.NewErrorMsg("Video analyze failed"))
		return
	}

	var resp = map[string]any{}
	for _, key := range usefulKeys {
		resp[key] = gjson.GetBytes(output, key).Value()
	}

	var thumbnail = gjson.GetBytes(output, "thumbnail").String()
	var extractor = gjson.GetBytes(output, "extractor").String()
	if extractor == "BiliBili" {
		thumbnail = "/api/proxy?url=" + thumbnail
	}
	resp["thumbnail"] = thumbnail
	resp["extractor"] = extractor

	var formatArray = gjson.GetBytes(output, "formats").Array()
	var formats = make([]map[string]any, 0, len(formatArray))
	for _, f := range formatArray {
		var size = f.Get("filesize").Float()
		if size == 0 {
			size = f.Get("filesize_approx").Float()
		}
		if size == 0 {
			continue
		}
		var format = f.Get("format").String()
		var formatSp = strings.Split(format, "-")
		if len(formatSp) > 1 {
			format = strings.TrimSpace(formatSp[1])
		}
		if format == "drc" {
			continue
		}
		formats = append(formats, map[string]any{
			"id":     f.Get("format_id").Value(),
			"tbr":    f.Get("tbr").Value(),
			"ext":    f.Get("ext").Value(),
			"acodec": f.Get("acodec").Value(),
			"vcodec": f.Get("vcodec").Value(),
			"format": format,
			"size":   size,
		})
	}
	resp["formats"] = formats

	c.JSON(http.StatusOK, resp)
}

func Download(c *gin.Context) {
	link, exit := common.Query(c, "link")
	if exit {
		return
	}
	formatId, exit := common.Query(c, "formatId")
	if exit {
		return
	}
	filename, exit := common.Query(c, "filename")
	if exit {
		return
	}
	if !strings.Contains(formatId, "+") {
		steamDownload(c, link, formatId, filename)
		return
	}
	filename = common.CacheDir + filename
	output, err := common.Exec(common.YTBin, "--ffmpeg-location", common.FfmpegBin, "--cookies", "cookies.txt", "--force-overwrites", "-f", formatId, "--merge-output-format", "mp4/mkv", "-o", filename, link)
	if err != nil {
		logger.Error().Err(err).Msg(string(output))
		c.JSON(http.StatusInternalServerError, common.NewErrorMsg("Video download failed"))
		return
	}
	if strings.Contains(string(output), filename+common.ExtMp4) {
		filename += common.ExtMp4
	} else {
		filename += common.ExtMkv
	}
	f, err := os.Open(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorMsg(err.Error()))
		return
	}
	defer f.Close()
	stat, _ := f.Stat()
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename[len(common.CacheDir):]))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(stat.Size(), 10))
	_, _ = io.Copy(c.Writer, f)
	go func() {
		_ = os.Remove(filename)
	}()
}

func steamDownload(c *gin.Context, link, formatId, filename string) {
	ext, exit := common.Query(c, "ext")
	if exit {
		return
	}
	cmd := exec.Command(common.YTBin, "--ffmpeg-location", common.FfmpegBin, "-f", formatId, "-o", "-", link)
	stdout, _ := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_ = cmd.Start()

	filename += "." + ext
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))
	c.Header("Content-Type", "application/octet-stream")
	size := c.Query("size")
	if size != "" {
		c.Header("Content-Length", size)
	}
	_, _ = io.Copy(c.Writer, stdout)

	if err := cmd.Wait(); err != nil {
		logger.Error().Err(err).Msg(stderr.String())
		c.JSON(http.StatusInternalServerError, common.NewErrorMsg("Video download failed"))
		return
	}
}
