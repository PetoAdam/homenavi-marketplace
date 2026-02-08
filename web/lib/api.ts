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

const getApiBase = () => (process.env.NEXT_PUBLIC_API_BASE || '').replace(/\/+$/, '');

class MarketplaceApi {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private buildUrl(path: string) {
    if (this.baseUrl) {
      return `${this.baseUrl}${path}`;
    }
    return `/api${path}`;
  }

  async listIntegrations(): Promise<Integration[]> {
    try {
      const res = await fetch(this.buildUrl('/integrations'), { cache: 'no-store' });
      if (!res.ok) {
        return [];
      }
      const data = await res.json();
      return data.integrations || [];
    } catch {
      return [];
    }
  }
}

export const marketplaceApi = new MarketplaceApi(getApiBase());
