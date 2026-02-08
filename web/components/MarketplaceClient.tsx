"use client";

import { useMemo, useState } from "react";
import type { Integration } from "../lib/api";
import IntegrationCard from "./IntegrationCard";

type Props = {
  integrations: Integration[];
};

export default function MarketplaceClient({ integrations }: Props) {
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<"all" | "verified" | "community">("all");
  const [sort, setSort] = useState<"name" | "version">("name");

  const stats = useMemo(() => {
    const verified = integrations.filter((item) => item.verified).length;
    const community = integrations.length - verified;
    return { total: integrations.length, verified, community };
  }, [integrations]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    let items = integrations;

    if (filter === "verified") {
      items = items.filter((item) => item.verified);
    } else if (filter === "community") {
      items = items.filter((item) => !item.verified);
    }

    if (q) {
      items = items.filter((item) => {
        const haystack = [item.name, item.description, item.publisher, item.id]
          .filter(Boolean)
          .join(" ")
          .toLowerCase();
        return haystack.includes(q);
      });
    }

    if (sort === "name") {
      items = [...items].sort((a, b) => a.name.localeCompare(b.name));
    } else {
      items = [...items].sort((a, b) => b.version.localeCompare(a.version, undefined, { numeric: true }));
    }

    return items;
  }, [integrations, filter, query, sort]);

  return (
    <section className="mx-auto mt-10 max-w-6xl space-y-8">
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-2xl border border-white/10 bg-panel/60 p-4 shadow-soft">
          <div className="text-xs uppercase tracking-[0.3em] text-white/40">Total</div>
          <div className="mt-2 text-2xl font-semibold text-white">{stats.total}</div>
        </div>
        <div className="rounded-2xl border border-white/10 bg-panel/60 p-4 shadow-soft">
          <div className="text-xs uppercase tracking-[0.3em] text-white/40">Verified</div>
          <div className="mt-2 text-2xl font-semibold text-white">{stats.verified}</div>
        </div>
        <div className="rounded-2xl border border-white/10 bg-panel/60 p-4 shadow-soft">
          <div className="text-xs uppercase tracking-[0.3em] text-white/40">Community</div>
          <div className="mt-2 text-2xl font-semibold text-white">{stats.community}</div>
        </div>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-white/10 bg-panel/60 p-4 shadow-soft">
        <div className="flex flex-1 flex-wrap items-center gap-3">
          <div className="flex min-w-[220px] flex-1 items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-white/70">
            <span className="text-sm">Search</span>
            <input
              className="w-full bg-transparent text-sm text-white placeholder:text-white/40 focus:outline-none"
              placeholder="Search by name, publisher, or id"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>
          <div className="flex items-center gap-2">
            {["all", "verified", "community"].map((item) => (
              <button
                key={item}
                type="button"
                onClick={() => setFilter(item as "all" | "verified" | "community")}
                className={`rounded-full border px-3 py-1 text-xs uppercase tracking-[0.2em] transition ${
                  filter === item
                    ? "border-accent bg-accent/20 text-accent"
                    : "border-white/10 text-white/60 hover:border-accent/50"
                }`}
              >
                {item}
              </button>
            ))}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-xs uppercase tracking-[0.2em] text-white/40">Sort</span>
          <select
            className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-white/80"
            value={sort}
            onChange={(event) => setSort(event.target.value as "name" | "version")}
          >
            <option value="name">Name</option>
            <option value="version">Version</option>
          </select>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
        {filtered.length === 0 ? (
          <div className="col-span-full rounded-2xl border border-white/10 bg-panel/60 p-8 text-white/70">
            No integrations match this filter.
          </div>
        ) : (
          filtered.map((integration, index) => (
            <IntegrationCard key={`${integration.id}-${integration.version}`} integration={integration} index={index} />
          ))
        )}
      </div>
    </section>
  );
}
