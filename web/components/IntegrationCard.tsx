import type { Integration } from '../lib/api';

type Props = {
  integration: Integration;
  index?: number;
};

export default function IntegrationCard({ integration, index = 0 }: Props) {
  const cover = integration.images?.[0] || integration.assets?.icon || '';

  return (
    <div
      className="fade-up group relative overflow-hidden rounded-2xl border border-white/10 bg-panel/70 p-5 shadow-soft backdrop-blur transition-transform duration-300 hover:-translate-y-1"
      style={{ animationDelay: `${index * 60}ms` }}
    >
      <div className="absolute inset-0 bg-gradient-to-br from-white/5 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
      <div className="relative flex items-start justify-between gap-4">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-white/40">{integration.publisher || 'Community'}</p>
          <h3 className="text-lg font-semibold text-white">{integration.name}</h3>
          <p className="text-sm text-white/60">v{integration.version}</p>
        </div>
        <div className="h-12 w-12 overflow-hidden rounded-xl border border-white/10 bg-white/5">
          {cover ? (
            <img src={cover} alt="" className="h-full w-full object-cover" />
          ) : (
            <div className="h-full w-full text-center text-lg leading-[3rem] text-white/60">
              {integration.name.slice(0, 1).toUpperCase()}
            </div>
          )}
        </div>
      </div>
      <p className="relative mt-4 line-clamp-3 text-sm text-white/70">
        {integration.description || 'No description provided.'}
      </p>
      <div className="relative mt-5 flex flex-wrap items-center justify-between gap-3 text-xs text-white/50">
        <span className="rounded-full bg-white/5 px-2 py-1">{integration.listen_path}</span>
        {integration.verified ? (
          <span className="rounded-full bg-accent/20 px-2 py-1 text-accent">Verified</span>
        ) : (
          <span className="rounded-full bg-white/10 px-2 py-1 text-white/60">Community</span>
        )}
      </div>
      <div className="relative mt-4 flex flex-wrap items-center gap-3 text-xs">
        {integration.repo_url ? (
          <a
            className="rounded-full border border-white/10 px-3 py-1 text-white/70 transition hover:border-accent/60 hover:text-white"
            href={integration.repo_url}
            target="_blank"
            rel="noreferrer"
          >
            Repo
          </a>
        ) : null}
        <a
          className="rounded-full border border-white/10 px-3 py-1 text-white/70 transition hover:border-accent/60 hover:text-white"
          href={integration.manifest_url}
          target="_blank"
          rel="noreferrer"
        >
          Manifest
        </a>
      </div>
    </div>
  );
}
