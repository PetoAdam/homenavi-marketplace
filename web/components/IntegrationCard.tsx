"use client";

import { useRouter } from 'next/navigation';
import { FaCheckCircle, FaDownload, FaStar, FaUsers } from 'react-icons/fa';
import type { Integration } from '../lib/api';
import IntegrationIcon from './IntegrationIcon';

type Props = {
  integration: Integration;
  index?: number;
  variant?: "default" | "featured";
};

export default function IntegrationCard({ integration, index = 0, variant = "default" }: Props) {
  const cover = integration.assets?.icon || integration.images?.[0] || '';
  const router = useRouter();
  const isFeatured = variant === "featured";

  const versionLabel = integration.version.startsWith('v')
    ? integration.version.slice(1)
    : integration.version;

  const handleNavigate = () => {
    router.push(`/integrations/${encodeURIComponent(integration.id)}`);
  };

  return (
    <div
      className={`fade-up group relative cursor-pointer overflow-hidden rounded-2xl border border-white/10 bg-panel/70 p-5 shadow-soft backdrop-blur transition-transform duration-300 hover:-translate-y-1 ${
        isFeatured ? "md:p-6" : ""
      }`}
      style={{ animationDelay: `${index * 60}ms` }}
      role="button"
      tabIndex={0}
      onClick={handleNavigate}
      onKeyDown={(event) => {
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault();
          handleNavigate();
        }
      }}
    >
      <div className="absolute inset-0 bg-gradient-to-br from-white/5 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
      <div className={`relative flex items-start justify-between gap-4 ${isFeatured ? "md:gap-6" : ""}`}>
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-white/40">{integration.publisher || 'Community'}</p>
          <h3 className={`${isFeatured ? "text-xl" : "text-lg"} font-semibold text-white`}>{integration.name}</h3>
          <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-white/60">
            <span className="rounded-full bg-white/10 px-2 py-1">Version {versionLabel}</span>
          </div>
        </div>
        <div
          className={`flex items-center justify-center overflow-hidden rounded-xl border border-white/10 bg-white/5 ${
            isFeatured ? "h-16 w-16" : "h-12 w-12"
          }`}
        >
          <IntegrationIcon
            icon={cover}
            fallback={integration.name}
            className={isFeatured ? "h-8 w-8 text-white/80" : undefined}
          />
        </div>
      </div>
      <p className={`relative mt-4 ${isFeatured ? "line-clamp-4" : "line-clamp-3"} text-sm text-white/70`}>
        {integration.description || 'No description provided.'}
      </p>
      <div className="relative mt-5 flex flex-wrap items-center justify-between gap-3 text-xs text-white/50">
        <span className="rounded-full bg-white/5 px-2 py-1">{integration.listen_path}</span>
        {integration.verified ? (
          <span className="inline-flex items-center gap-2 rounded-full bg-accent/20 px-2 py-1 text-accent">
            <FaCheckCircle className="text-[0.7rem]" />
            Verified
          </span>
        ) : (
          <span className="inline-flex items-center gap-2 rounded-full bg-white/10 px-2 py-1 text-white/60">
            <FaUsers className="text-[0.7rem]" />
            Community
          </span>
        )}
      </div>
      <div className="relative mt-3 flex flex-wrap items-center gap-2 text-xs text-white/60">
        <span className="inline-flex items-center gap-2 rounded-full bg-white/5 px-2 py-1">
          <FaDownload className="text-[0.7rem]" />
          {integration.downloads} downloads
        </span>
        {integration.featured ? (
          <span className="inline-flex items-center gap-2 rounded-full bg-accent/20 px-2 py-1 text-accent">
            <FaStar className="text-[0.7rem]" />
            Featured
          </span>
        ) : null}
      </div>
      <div className="relative mt-4 flex flex-wrap items-center gap-3 text-xs">
        {integration.repo_url ? (
          <a
            className="rounded-full border border-white/10 px-3 py-1 text-white/70 transition hover:border-accent/60 hover:text-white"
            href={integration.repo_url}
            target="_blank"
            rel="noreferrer"
            onClick={(event) => event.stopPropagation()}
          >
            Repo
          </a>
        ) : null}
        <a
          className="rounded-full border border-white/10 px-3 py-1 text-white/70 transition hover:border-accent/60 hover:text-white"
          href={integration.manifest_url}
          target="_blank"
          rel="noreferrer"
          onClick={(event) => event.stopPropagation()}
        >
          Manifest
        </a>
      </div>
    </div>
  );
}
