package model

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

type VideoFormat struct {
	Id     string `json:"id"`
	Ext    string `json:"ext"`
	ACodec string `json:"acodec"`
	VCodec string `json:"vcodec"`
	Format string `json:"format"`
	Size   int64  `json:"size"`
}

func NewVideoFormat(f gjson.Result, extractor string) (VideoFormat, bool) {
	var size = f.Get("filesize").Int()
	if size == 0 {
		size = f.Get("filesize_approx").Int()
	}
	if size == 0 {
		return VideoFormat{}, true
	}
	var vf = VideoFormat{
		Id:     f.Get("format_id").String(),
		ACodec: f.Get("acodec").String(),
		VCodec: f.Get("vcodec").String(),
		Format: f.Get("format").String(),
		Ext:    f.Get("ext").String(),
		Size:   size,
	}
	if vf.filterFormat(f, extractor) {
		return VideoFormat{}, true
	}
	return vf, false
}

func (vf *VideoFormat) filterFormat(f gjson.Result, extractor string) bool {
	switch extractor {
	case "BiliBili":
		// 1080P 高清 - hev1
		var resolution, more string
		if vf.VCodec == "none" {
			resolution = strings.TrimSpace(strings.Split(vf.Format, "-")[1])
			var tbr = f.Get("tbr").Int()
			more = strconv.FormatInt(tbr, 10) + "k"
		} else {
			quality := f.Get("quality").Int()
			if quality < 32 {
				// < 480p
				return true
			}
			resolution = vf.Format
			more = strings.Split(vf.VCodec, ".")[0]
		}
		vf.Format = strings.Join([]string{resolution, more}, " - ")
	case "youtube":
		// 3840x2160 (2160p) - vp09 - 4654k
		resolution := strings.TrimSpace(strings.Split(vf.Format, "-")[1])
		if resolution == "drc" || vf.Ext == "webm" {
			return true
		}
		var tbr = strconv.FormatInt(f.Get("tbr").Int(), 10) + "k"
		var codec = vf.VCodec
		if codec == "none" {
			vf.Format = strings.Join([]string{resolution, tbr}, " - ")
			return false
		}
		quality := f.Get("quality").Int()
		if quality < 8 {
			// < 720p
			return true
		}
		more := strings.Split(codec, ".")[0]
		vf.Format = strings.Join([]string{resolution, more, tbr}, " - ")
	case "Douyin":
		// 720p - h265 - Playback
		if strings.HasPrefix(vf.Id, "download_addr-") {
			return true
		}
		var index = vf.Id[len(vf.Id)-1:]
		if index != "0" && index != "2" {
			return true
		}
		resolution := strings.Split(vf.Id, "_")[1]
		var formatNote = f.Get("format_note").String()
		more := strings.ReplaceAll(formatNote, " video", "")
		vf.Format = strings.Join([]string{resolution, vf.VCodec, more}, " - ")
	}
	return false
}
