/**
 * Types for feed discovery
 */

export interface DiscoveredFeed {
  name: string;
  homepage: string;
  rss_feed: string;
  icon_url?: string;
  recent_articles?: Array<{
    title: string;
    date?: string;
  }>;
}

export interface ProgressCounts {
  current: number;
  total: number;
  found: number;
}

export interface StartResult {
  status: string;
  message?: string;
  total?: number;
}

export interface ProgressState {
  is_complete: boolean;
  error?: string;
  feeds?: DiscoveredFeed[];
  progress?: {
    stage: string;
    message?: string;
    detail?: string;
    current?: number;
    total?: number;
    found_count?: number;
    feed_name?: string;
  };
}
