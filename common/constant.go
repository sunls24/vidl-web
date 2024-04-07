package common

import (
	"vidlp/common/utils"
)

var (
	YtDlpBin  = utils.GetEnv("YT_DLP", "yt-dlp")
	FfmpegBin = utils.GetEnv("FFMPEG", "ffmpeg")
	Cookies   = utils.GetEnv("COOKIES", "cookies.txt")
)

const (
	ExtMp4   = ".mp4"
	ExtMkv   = ".mkv"
	CacheDir = "cache/"
)
