import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import React, { useRef, useState } from "react";
import { ArrowDownToLine, Github, RotateCw, TextSearch } from "lucide-react";
import { toast } from "sonner";
import Badge from "@/components/badge.tsx";
import { GITHUB_URL, VERSION, VideoInfo } from "@/lib/constant.ts";
import {
  convertSeconds,
  extractUrl,
  findAudio,
  formatCount,
  toMiB,
} from "@/lib/utils.ts";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import Download from "@/components/download.tsx";

function App() {
  const [alert, setAlert] = useState(false);
  const [progress, setProgress] = useState(0);
  const [loading, setLoading] = useState(false);
  const [formatId, setFormatId] = useState<string | undefined>(undefined);
  const [videoInfo, setVideoInfo] = useState<VideoInfo | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  async function onAnalyze() {
    if (!inputRef.current || !inputRef.current.value) {
      toast.warning("No link detected");
      return;
    }
    const link = extractUrl(inputRef.current.value) ?? "";
    if (!link) {
      toast.warning("No available link found");
      return;
    }
    setLoading(true);
    setVideoInfo(null);
    setFormatId(undefined);
    try {
      const res = await fetch(`/api/analyze?link=${encodeURIComponent(link)}`);
      const info = await res.json();
      if (info.error) {
        toast.error(info.error);
        return;
      }
      console.debug("video:", info);
      setVideoInfo(info);
      let findBest = false;
      for (let i = info.formats.length - 1; i >= 0; i--) {
        const fmt = info.formats[i];
        if (fmt.vcodec !== "none" && fmt.acodec !== "none") {
          setFormatId(fmt.id);
          findBest = true;
          break;
        }
      }
      if (!findBest) {
        setFormatId(info.formats[0].id);
      }
    } catch (err: any) {
      toast.error(err);
    } finally {
      setLoading(false);
    }
  }

  async function onDownload() {
    if (!videoInfo) {
      return;
    }
    let formatIndex = 0;
    for (let i = 0; i < videoInfo.formats.length; i++) {
      if (videoInfo.formats[i].id === formatId) {
        formatIndex = i;
        break;
      }
    }
    const format = videoInfo.formats[formatIndex];
    const filename = `${videoInfo.extractor}-${videoInfo.id}-${format.format.replace(/ /g, "")}`;
    let size = format.size;
    let complexFormatId = formatId;
    const audio = findAudio(videoInfo.formats, formatIndex);
    if (audio) {
      size += audio.size;
      complexFormatId = encodeURIComponent(complexFormatId + "+" + audio.id);
    }

    console.debug("-> download");
    console.debug("format:", format);
    console.debug("audio:", audio);

    let query = `link=${encodeURIComponent(videoInfo.webpage_url)}&formatId=${complexFormatId}&filename=${encodeURIComponent(filename)}`;
    if (!audio) {
      // not merge
      query += `&ext=${format.ext}`;
    }
    try {
      const resp = await fetch(`/api/download?${query}`);
      if (!resp.ok) {
        toast.error((await resp.json()).error);
        return;
      }
      const progress = new Response(
        new ReadableStream({
          async start(controller) {
            const reader = resp.body!.getReader();
            let downLen = 0;
            let changeLen = 0;
            for (;;) {
              const { done, value } = await reader.read();
              if (done) break;
              controller.enqueue(value);
              downLen += value.byteLength;
              changeLen += value.byteLength;
              if (changeLen > 102400) {
                changeLen = 0;
                setProgress((downLen / size) * 100);
              }
            }
            controller.close();
          },
        }),
      );
      const url = window.URL.createObjectURL(await progress.blob());
      const a = document.createElement("a");
      a.href = url;
      a.download = decodeURIComponent(
        resp.headers.get("Content-Disposition")!.split("filename=")[1],
      );
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (err: any) {
      toast.error(err);
    } finally {
      setAlert(false);
      setProgress(0);
    }
  }

  async function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key !== "Enter" || e.nativeEvent.isComposing) {
      return;
    }
    if (loading) {
      return;
    }
    e.preventDefault();
    await onAnalyze();
  }

  return (
    <div className="mt-8 flex flex-col items-center gap-6 px-4 md:mt-16 md:gap-8">
      <h1 className="text-center text-3xl font-bold sm:text-4xl md:text-5xl">
        Download videos from anywhere
      </h1>
      <p className="max-w-[500px] text-center text-xl text-muted-foreground">
        Enter the link to the video you want to download, and we'll make it
        happen.
      </p>
      <div className="flex w-full max-w-[600px] gap-2 md:py-4">
        <Input
          ref={inputRef}
          disabled={loading}
          placeholder="Enter video link"
          onKeyDown={onKeyDown}
        />
        <Button disabled={loading} onClick={onAnalyze}>
          {loading ? <RotateCw className="animate-spin" /> : <TextSearch />}
        </Button>
      </div>
      {videoInfo && (
        <div className="flex w-full max-w-[600px] flex-col gap-3 md:max-w-[730px] md:flex-row">
          <div className="flex basis-1/2 flex-col gap-3">
            <img src={videoInfo.thumbnail} />
            <div className="flex gap-2">
              <Select value={formatId} onValueChange={setFormatId}>
                <SelectTrigger>
                  <SelectValue placeholder="Select format" />
                </SelectTrigger>
                <SelectContent>
                  {videoInfo.formats.map((v, i) => (
                    <SelectItem key={v.id} value={v.id}>
                      {`${v.format} ≈${toMiB(
                        v.size + (findAudio(videoInfo.formats, i)?.size ?? 0),
                      ).toFixed(2)}Mib`}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Download
                open={alert}
                onChange={setAlert}
                progress={progress}
                trigger={
                  <Button onClick={onDownload}>
                    <ArrowDownToLine />
                  </Button>
                }
              />
            </div>
          </div>
          <div className="flex max-w-[96%] basis-1/2 flex-col gap-2 overflow-auto">
            <h3 className="line-clamp-2 text-xl font-medium underline underline-offset-4">
              {videoInfo.title}
            </h3>
            {videoInfo.description && (
              <p className="line-clamp-3 text-muted-foreground">
                {videoInfo.description}
              </p>
            )}
            {videoInfo.uploader && (
              <Badge title="Author" value={videoInfo.uploader} />
            )}
            {videoInfo.upload_date && (
              <Badge title="Date" value={videoInfo.upload_date} />
            )}
            {videoInfo.duration && (
              <Badge
                title="Duration"
                value={convertSeconds(videoInfo.duration)}
              />
            )}
            {videoInfo.view_count > 0 ? (
              <Badge title="View" value={formatCount(videoInfo.view_count)} />
            ) : null}
          </div>
        </div>
      )}
      <div className="flex gap-2">
        <div className="flex gap-1 rounded-md bg-secondary px-2 py-0.5">
          <span className="text-muted-foreground">Support</span>
          <img className="h-6 w-6" src="/img/douyin.svg" />
          <img className="h-6 w-6" src="/img/youtube.svg" />
          <img className="h-6 w-6" src="/img/bilibili.svg" />
        </div>
        <a
          href={GITHUB_URL}
          target="_blank"
          className="flex cursor-pointer items-center gap-1 rounded-md bg-secondary px-2 py-0.5 text-muted-foreground"
        >
          <Github size={18} />
          <span className="underline underline-offset-2">{VERSION}</span>
        </a>
      </div>
    </div>
  );
}

export default App;
