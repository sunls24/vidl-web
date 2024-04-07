export interface VideoInfo {
  id: string;
  title: string;
  description: string;
  uploader: string;
  upload_date: string;
  duration: number;
  view_count: number;
  thumbnail: string;
  extractor: string;
  webpage_url: string;
  formats: VideoFormat[];
}

export interface VideoFormat {
  id: string;
  ext: string;
  acodec: string;
  vcodec: string;
  format: string;
  size: number;
}

export const VERSION = "v1.0.0";
export const GITHUB_URL = "https://github.com/sunls24/vidlp";
