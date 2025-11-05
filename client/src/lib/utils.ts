import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function buildApiUrl(path: string, { websocket=false }: { websocket?: boolean; } = {}) {
  let protocol = location.protocol
  if (websocket) {
    protocol = protocol.replace("http", "ws")
  }
  let hostname = location.host
  if (import.meta.env.MODE !== "production") {
    hostname = "localhost:8080"
  }
  if (path.charAt(0) === "/") {
    path = path.replace("/", "")
  }
  let url = `${protocol}//${hostname}/${path}`
  return url
}
