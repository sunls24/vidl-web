package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"vidlp/common"
	"vidlp/common/logger"
	"vidlp/common/utils"
	"vidlp/model"
)

var usefulKeys = []string{"title", "description", "duration", "view_count", "uploader", "upload_date", "webpage_url", "id", "extractor", "thumbnail"}

func Analyze(c *gin.Context) {
	link, exit := utils.MustQuery(c, "link")
	if exit {
		return
	}
	output, err := common.Exec(common.YtDlpBin, "--cookies", common.Cookies, "-j", link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewMsg("Video analyze failed"))
		logger.Error().Err(err).Send()
		return
	}

	var resp = make(map[string]any, len(usefulKeys)+1)
	for _, key := range usefulKeys {
		resp[key] = gjson.GetBytes(output, key).Value()
	}
	resp["thumbnail"] = "/api/proxy?url=" + url.QueryEscape(resp["thumbnail"].(string))
	var channel = gjson.GetBytes(output, "channel")
	if channel.Exists() {
		resp["uploader"] = channel.String()
	}

	var extractor = resp["extractor"].(string)
	var formatArray = gjson.GetBytes(output, "formats").Array()
	var formats = make([]model.VideoFormat, 0, len(formatArray))
	for _, f := range formatArray {
		vf, skip := model.NewVideoFormat(f, extractor)
		if skip {
			continue
		}
		formats = append(formats, vf)
	}
	resp["formats"] = formats

	c.JSON(http.StatusOK, resp)
}

func Download(c *gin.Context) {
	link, exit := utils.MustQuery(c, "link")
	if exit {
		return
	}
	formatId, exit := utils.MustQuery(c, "formatId")
	if exit {
		return
	}
	filename, exit := utils.MustQuery(c, "filename")
	if exit {
		return
	}
	if !strings.Contains(formatId, "+") {
		streamDownload(c, link, formatId, filename)
		return
	}
	filename = common.CacheDir + filename
	output, err := common.Exec(common.YtDlpBin, "--ffmpeg-location", common.FfmpegBin, "--cookies", common.Cookies, "--force-overwrites", "-f", formatId, "--merge-output-format", "mp4/mkv", "-o", filename, link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewMsg("Video download failed"))
		logger.Error().Err(err).Send()
		return
	}
	if strings.Contains(string(output), filename+common.ExtMp4) {
		filename += common.ExtMp4
	} else {
		filename += common.ExtMkv
	}
	f, err := os.Open(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewMsg(err.Error()))
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

func streamDownload(c *gin.Context, link, formatId, filename string) {
	ext, exit := utils.MustQuery(c, "ext")
	if exit {
		return
	}
	cmd := exec.Command(common.YtDlpBin, "--ffmpeg-location", common.FfmpegBin, "--cookies", common.Cookies, "-f", formatId, "-o", "-", link)
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
		c.JSON(http.StatusInternalServerError, model.NewMsg("Video stream download failed"))
		logger.Error().Err(fmt.Errorf("%w: %s", err, stderr.String())).Send()
	}
}
