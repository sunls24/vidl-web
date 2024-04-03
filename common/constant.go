package common

type ErrorMsg struct {
	Error string `json:"error"`
}

func NewErrorMsg(msg string) ErrorMsg {
	return ErrorMsg{Error: msg}
}

var (
	YTBin     = GetEnv("YTBin", "yt-dlp")
	FfmpegBin = GetEnv("Ffmpeg", "ffmpeg")
)

const (
	ExtMp4   = ".mp4"
	ExtMkv   = ".mkv"
	CacheDir = "cache/"
)
