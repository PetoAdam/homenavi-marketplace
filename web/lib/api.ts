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

const getApiBase = () => {
  const publicBase = (process.env.NEXT_PUBLIC_API_BASE || '').replace(/\/+$/, '');
  const internalBase = (process.env.INTERNAL_API_BASE || '').replace(/\/+$/, '');

  if (typeof window === 'undefined') {
    if (internalBase) {
      return internalBase;
    }
    if (publicBase.startsWith('http://') || publicBase.startsWith('https://')) {
      return publicBase;
    }
  }

  return publicBase;
};

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

  async getIntegration(id: string, version?: string): Promise<Integration | null> {
    try {
      const encodedId = encodeURIComponent(id);
      const query = version ? `?version=${encodeURIComponent(version)}` : '';
      const res = await fetch(this.buildUrl(`/integrations/${encodedId}${query}`), { cache: 'no-store' });
      if (!res.ok) {
        return null;
      }
      return await res.json();
    } catch {
      return null;
    }
  }
}

export const marketplaceApi = new MarketplaceApi(getApiBase());
