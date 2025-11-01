import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getAPIHost() {
  return process.env.NODE_ENV === 'production' ? document.location.host : 'localhost:8080';
}
