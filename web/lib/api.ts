export type Integration = {
  id: string;
  name: string;
  version: string;
  description: string;
  manifest_url: string;
  image: string;
  images: string[];
  assets: Record<string, string>;
  listen_path: string;
  repo_url?: string;
  release_tag?: string;
  publisher?: string;
  verified: boolean;
  latest: boolean;
};

const apiBase = process.env.NEXT_PUBLIC_API_BASE || '/api';

export async function fetchIntegrations(): Promise<Integration[]> {
  try {
    const res = await fetch(`${apiBase}/api/integrations`, { cache: 'no-store' });
    if (!res.ok) {
      return [];
    }
    const data = await res.json();
    return data.integrations || [];
  } catch {
    return [];
  }
}
