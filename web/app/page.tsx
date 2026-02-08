import MarketplaceClient from '../components/MarketplaceClient';
import { marketplaceApi } from '../lib/api';

export const dynamic = 'force-dynamic';

export default async function Home() {
  const integrations = await marketplaceApi.listIntegrations();

  return (
    <main className="min-h-screen px-8 py-12">
      <header className="mx-auto max-w-6xl">
        <div className="flex flex-wrap items-center justify-between gap-6">
          <div>
            <div className="flex items-center gap-4">
              <div className="h-14 w-14 rounded-2xl border border-white/10 bg-panel/70 p-3 shadow-soft">
                <img src="/icon.png" alt="Homenavi" className="h-full w-full" />
              </div>
              <div>
                <p className="text-xs uppercase tracking-[0.35em] text-white/40">Homenavi Marketplace</p>
                <h1 className="mt-3 text-4xl font-semibold text-white">Discover integrations</h1>
              </div>
            </div>
            <p className="mt-3 max-w-2xl text-white/60">
              Curated releases published through secure pipelines. Browse what is available and keep your
              deployment fresh.
            </p>
          </div>
          <div className="hidden rounded-3xl border border-white/10 bg-panel/60 px-6 py-4 text-sm text-white/70 shadow-soft md:block">
            <div className="text-xs uppercase tracking-[0.3em] text-white/40">Status</div>
            <div className="mt-2 text-lg text-white">{integrations.length} integrations</div>
          </div>
        </div>
      </header>

      <MarketplaceClient integrations={integrations} />
    </main>
  );
}
