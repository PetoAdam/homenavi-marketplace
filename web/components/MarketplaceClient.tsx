"use client";

import { useMemo, useState } from "react";
import {
  FaArrowDown,
  FaCheckCircle,
  FaCompass,
  FaFire,
  FaSearch,
  FaSort,
  FaStar,
  FaThLarge,
  FaUsers
} from "react-icons/fa";
import type { Integration } from "../lib/api";
import IntegrationCard from "./IntegrationCard";

type Props = {
  integrations: Integration[];
};

export default function MarketplaceClient({ integrations }: Props) {
  const [query, setQuery] = useState("");
  const [mode, setMode] = useState<"discover" | "trending" | "downloads">("discover");
  const [filter, setFilter] = useState<"all" | "featured" | "verified" | "community">("all");
  const [sort, setSort] = useState<"name" | "version" | "downloads" | "trending">("trending");

  const featuredItems = useMemo(() => integrations.filter((item) => item.featured), [integrations]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    let items = integrations;

    if (filter === "featured") {
      items = items.filter((item) => item.featured);
    } else if (filter === "verified") {
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
    } else if (sort === "version") {
      items = [...items].sort((a, b) => b.version.localeCompare(a.version, undefined, { numeric: true }));
    } else if (sort === "downloads") {
      items = [...items].sort((a, b) => b.downloads - a.downloads);
    } else {
      items = [...items].sort((a, b) => b.trending_score - a.trending_score);
    }

    return items;
  }, [integrations, filter, query, sort]);

  return (
    <section className="mx-auto mt-10 max-w-6xl space-y-8">
      <div className="flex flex-wrap items-center gap-3 rounded-full border border-white/10 bg-panel/60 p-2 shadow-soft">
        {["discover", "trending", "downloads"].map((item) => (
          <button
            key={item}
            type="button"
            onClick={() => {
              setMode(item as "discover" | "trending" | "downloads");
              if (item === "downloads") {
                setSort("downloads");
              } else if (item === "trending") {
                setSort("trending");
              } else {
                setSort("name");
              }
            }}
            className={`rounded-full px-4 py-2 text-xs uppercase tracking-[0.25em] transition ${
              mode === item
                ? "bg-accent/30 text-accent"
                : "text-white/60 hover:bg-white/5 hover:text-white"
            }`}
          >
            <span className="flex items-center gap-2">
              {item === "discover" ? <FaCompass className="text-[0.7rem]" /> : null}
              {item === "trending" ? <FaFire className="text-[0.7rem]" /> : null}
              {item === "downloads" ? <FaArrowDown className="text-[0.7rem]" /> : null}
              {item === "downloads" ? "Downloads" : item}
            </span>
          </button>
        ))}
      </div>

      {featuredItems.length > 0 ? (
        <div className="rounded-3xl border border-white/10 bg-panel/50 p-6 shadow-soft">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-xs uppercase tracking-[0.3em] text-white/40">Featured</div>
              <div className="mt-2 text-lg font-semibold text-white">Spotlight integrations</div>
            </div>
            <button
              type="button"
              className="rounded-full border border-white/10 px-3 py-1 text-xs uppercase tracking-[0.2em] text-white/60"
              onClick={() => setFilter("featured")}
            >
              View all
            </button>
          </div>
          <div className="mt-5 grid gap-6 md:grid-cols-2">
            {featuredItems.map((integration, index) => (
              <IntegrationCard
                key={`${integration.id}-${integration.version}`}
                integration={integration}
                index={index}
                variant="featured"
              />
            ))}
          </div>
        </div>
      ) : null}

      <div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-white/10 bg-panel/60 p-4 shadow-soft">
        <div className="flex flex-1 flex-wrap items-center gap-3">
          <div className="flex min-w-[220px] flex-1 items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-white/70">
            <FaSearch className="text-sm text-white/50" />
            <input
              className="w-full bg-transparent text-sm text-white placeholder:text-white/40 focus:outline-none"
              placeholder="Search by name, publisher, or id"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>
          <div className="flex items-center gap-2">
            {["all", "featured", "verified", "community"].map((item) => (
              <button
                key={item}
                type="button"
                onClick={() => setFilter(item as "all" | "featured" | "verified" | "community")}
                className={`rounded-full border px-3 py-1 text-xs uppercase tracking-[0.2em] transition ${
                  filter === item
                    ? "border-accent bg-accent/20 text-accent"
                    : "border-white/10 text-white/60 hover:border-accent/50"
                }`}
              >
                <span className="flex items-center gap-2">
                  {item === "all" ? <FaThLarge className="text-[0.7rem]" /> : null}
                  {item === "featured" ? <FaStar className="text-[0.7rem]" /> : null}
                  {item === "verified" ? <FaCheckCircle className="text-[0.7rem]" /> : null}
                  {item === "community" ? <FaUsers className="text-[0.7rem]" /> : null}
                  {item}
                </span>
              </button>
            ))}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <FaSort className="text-xs text-white/50" />
          <span className="text-xs uppercase tracking-[0.2em] text-white/40">Sort</span>
          <select
            className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-white/80"
            value={sort}
            onChange={(event) => setSort(event.target.value as "name" | "version" | "downloads" | "trending")}
          >
            <option value="name">Name</option>
            <option value="version">Version</option>
            <option value="downloads">Downloads</option>
            <option value="trending">Trending</option>
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
