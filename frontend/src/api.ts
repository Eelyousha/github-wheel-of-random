const API_BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8000";

export interface Project {
  id: number;
  full_name: string;
  name: string;
  owner: string;
  description: string | null;
  stars: number;
  language: string | null;
  topics: string[];
  html_url: string;
  avatar_url: string | null;
  last_updated: string;
}

export interface Filters {
  languages: string[];
  topics: string[];
}

export interface Status {
  last_crawl: string | null;
  project_count: number;
  is_crawling: boolean;
}

export interface ProjectFilter {
  page?: number;
  limit?: number;
  lang?: string;
  topic?: string;
  min_stars?: number;
}

async function fetchJSON<T>(path: string, params?: Record<string, string>): Promise<T> {
  const url = new URL(`${API_BASE}${path}`);
  if (params) {
    Object.entries(params).forEach(([k, v]) => {
      if (v) url.searchParams.set(k, v);
    });
  }
  const res = await fetch(url.toString());
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}

export function getProjects(filter: ProjectFilter): Promise<Project[]> {
  return fetchJSON<Project[]>("/api/v1/projects", {
    page: String(filter.page || 1),
    limit: String(filter.limit || 50),
    lang: filter.lang || "",
    topic: filter.topic || "",
    min_stars: String(filter.min_stars || 0),
  });
}

export function getRandomProjects(limit: number, filter?: { lang?: string; topic?: string; min_stars?: number }): Promise<Project[]> {
  return fetchJSON<Project[]>("/api/v1/projects/random", {
    limit: String(limit),
    lang: filter?.lang || "",
    topic: filter?.topic || "",
    min_stars: String(filter?.min_stars || 0),
  });
}

export function getFilters(): Promise<Filters> {
  return fetchJSON<Filters>("/api/v1/filters");
}

export function getStatus(): Promise<Status> {
  return fetchJSON<Status>("/api/v1/status");
}

export async function triggerCrawl(): Promise<void> {
  const res = await fetch(`${API_BASE}/api/v1/crawl`, { method: "POST" });
  if (!res.ok) throw new Error(`crawl trigger failed: ${res.status}`);
}
