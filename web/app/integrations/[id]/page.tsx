import { notFound } from 'next/navigation';

import ImageGallery from '../../../components/ImageGallery';
import IntegrationIcon from '../../../components/IntegrationIcon';
import { marketplaceApi } from '../../../lib/api';

export const dynamic = 'force-dynamic';

type PageProps = {
  params: {
    id: string;
  };
};

const formatVersion = (version: string) => (version.startsWith('v') ? version.slice(1) : version);

export default async function IntegrationPage({ params }: PageProps) {
  const integration = await marketplaceApi.getIntegration(params.id);

  if (!integration) {
    notFound();
  }

  const cover = integration.assets?.icon || integration.images?.[0] || '';
  const versionLabel = formatVersion(integration.version);

  return (
    <main className="min-h-screen px-8 py-12">
      <section className="mx-auto max-w-5xl space-y-10">
        <header className="flex flex-wrap items-start justify-between gap-6">
          <div className="space-y-3">
            <p className="text-xs uppercase tracking-[0.3em] text-white/40">
              {integration.publisher || 'Community'}
            </p>
            <h1 className="text-4xl font-semibold text-white">{integration.name}</h1>
            <div className="flex flex-wrap items-center gap-3 text-sm text-white/70">
              <span className="rounded-full bg-white/10 px-3 py-1">Version {versionLabel}</span>
              <span className="rounded-full bg-white/10 px-3 py-1">
                {integration.downloads} downloads
              </span>
              {integration.featured ? (
                <span className="rounded-full bg-accent/20 px-3 py-1 text-accent">Featured</span>
              ) : null}
              {integration.verified ? (
                <span className="rounded-full bg-accent/20 px-3 py-1 text-accent">Verified</span>
              ) : (
                <span className="rounded-full bg-white/10 px-3 py-1 text-white/60">Community</span>
              )}
            </div>
            <p className="max-w-2xl text-white/70">
              {integration.description || 'No description provided.'}
            </p>
          </div>
          <div className="flex h-24 w-24 items-center justify-center overflow-hidden rounded-3xl border border-white/10 bg-white/5">
            <IntegrationIcon
              icon={cover}
              fallback={integration.name}
              className="h-10 w-10 text-white/80"
              imageClassName="h-full w-full object-cover"
            />
          </div>
        </header>

        <section className="grid gap-6 md:grid-cols-2">
          <div className="rounded-2xl border border-white/10 bg-panel/60 p-6 shadow-soft">
            <h2 className="text-lg font-semibold text-white">Details</h2>
            <dl className="mt-4 space-y-3 text-sm text-white/70">
              <div>
                <dt className="text-xs uppercase tracking-[0.2em] text-white/40">Integration ID</dt>
                <dd className="mt-1 rounded-xl bg-white/5 px-3 py-2">{integration.id}</dd>
              </div>
              <div>
                <dt className="text-xs uppercase tracking-[0.2em] text-white/40">Listen Path</dt>
                <dd className="mt-1 rounded-xl bg-white/5 px-3 py-2">{integration.listen_path}</dd>
              </div>
              <div>
                <dt className="text-xs uppercase tracking-[0.2em] text-white/40">Release Tag</dt>
                <dd className="mt-1 rounded-xl bg-white/5 px-3 py-2">{integration.release_tag || 'N/A'}</dd>
              </div>
            </dl>
          </div>

          <div className="rounded-2xl border border-white/10 bg-panel/60 p-6 shadow-soft">
            <h2 className="text-lg font-semibold text-white">Links</h2>
            <div className="mt-4 flex flex-wrap gap-3 text-sm">
              {integration.repo_url ? (
                <a
                  className="rounded-full border border-white/10 px-4 py-2 text-white/70 transition hover:border-accent/60 hover:text-white"
                  href={integration.repo_url}
                  target="_blank"
                  rel="noreferrer"
                >
                  Repository
                </a>
              ) : null}
              <a
                className="rounded-full border border-white/10 px-4 py-2 text-white/70 transition hover:border-accent/60 hover:text-white"
                href={integration.manifest_url}
                target="_blank"
                rel="noreferrer"
              >
                Manifest
              </a>
            </div>
          </div>
        </section>

        {integration.images?.length ? <ImageGallery images={integration.images} /> : null}
      </section>
    </main>
  );
}
