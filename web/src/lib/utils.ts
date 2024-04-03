import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function extractUrl(str: string) {
  return str.match(/(https?:\/\/\S+)/)?.[0];
}

export function convertSeconds(seconds: number): string {
  const format = (val: number, label: string) =>
    val > 0 ? `${val}${label}` : "";

  const hours = Math.floor(seconds / 3600);
  seconds %= 3600;
  const minutes = Math.floor(seconds / 60);
  seconds %= 60;
  return `${format(hours, "h")}${format(minutes, "m")}${format(
    Math.floor(seconds),
    "s",
  )}`;
}

export function formatCount(count: number) {
  if (count <= 1000) {
    return count.toString();
  }
  return `${Math.floor(count / 1000)}k`;
}

export function toMiB(bytes: number): number {
  return bytes / (1024 * 1024);
}

export function findAudio(formats: any[], index: number) {
  if (formats[index].acodec !== "none") {
    return null;
  }

  let audioCount = 0;
  for (let i = 0; i < formats.length; i++) {
    if (!formats[i].format.includes("audio only")) {
      audioCount = i;
      break;
    }
  }
  if (audioCount === 0 || index < audioCount) {
    return null;
  }
  const videoCount = formats.length - audioCount;
  const audioIndex = Math.floor(
    (audioCount / videoCount) * (index - audioCount),
  );
  return formats[audioIndex];
}
