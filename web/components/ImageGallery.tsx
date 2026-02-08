"use client";

import { useEffect, useState } from "react";
import { FaChevronLeft, FaChevronRight, FaImages, FaTimes } from "react-icons/fa";

type Props = {
  images: string[];
};

export default function ImageGallery({ images }: Props) {
  const [activeIndex, setActiveIndex] = useState<number | null>(null);

  const hasImages = images.length > 0;

  useEffect(() => {
    if (activeIndex === null) {
      return undefined;
    }

    const handleKey = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setActiveIndex(null);
      } else if (event.key === "ArrowRight") {
        setActiveIndex((current) => (current === null ? 0 : (current + 1) % images.length));
      } else if (event.key === "ArrowLeft") {
        setActiveIndex((current) =>
          current === null ? images.length - 1 : (current - 1 + images.length) % images.length
        );
      }
    };

    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [activeIndex, images.length]);

  if (!hasImages) {
    return null;
  }

  return (
    <div className="rounded-2xl border border-white/10 bg-panel/60 p-6 shadow-soft">
      <h2 className="flex items-center gap-2 text-lg font-semibold text-white">
        <FaImages className="text-white/60" />
        Gallery
      </h2>
      <div className="mt-4 grid gap-4 sm:grid-cols-2">
        {images.map((src, index) => (
          <button
            key={src}
            type="button"
            onClick={() => setActiveIndex(index)}
            className="group relative overflow-hidden rounded-2xl border border-white/10 bg-white/5 text-left"
          >
            <img src={src} alt="" className="h-full w-full object-cover transition duration-300 group-hover:scale-105" />
          </button>
        ))}
      </div>

      {activeIndex !== null ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-6">
          <div className="relative w-full max-w-5xl">
            <button
              type="button"
              onClick={() => setActiveIndex(null)}
              aria-label="Close"
              className="absolute right-4 top-4 rounded-full border border-white/20 bg-black/60 p-2 text-white/80"
            >
              <FaTimes />
            </button>
            <div className="flex items-center justify-between gap-4">
              <button
                type="button"
                onClick={() =>
                  setActiveIndex((current) => (current === null ? 0 : (current - 1 + images.length) % images.length))
                }
                aria-label="Previous image"
                className="hidden rounded-full border border-white/20 bg-black/40 p-3 text-white/70 sm:block"
              >
                <FaChevronLeft />
              </button>
              <div className="flex-1 overflow-hidden rounded-3xl border border-white/10 bg-black/40">
                <img src={images[activeIndex]} alt="" className="h-full w-full object-contain" />
              </div>
              <button
                type="button"
                onClick={() =>
                  setActiveIndex((current) => (current === null ? 0 : (current + 1) % images.length))
                }
                aria-label="Next image"
                className="hidden rounded-full border border-white/20 bg-black/40 p-3 text-white/70 sm:block"
              >
                <FaChevronRight />
              </button>
            </div>
            <div className="mt-4 flex items-center justify-center gap-2 text-xs text-white/70">
              {images.map((_, index) => (
                <button
                  key={index}
                  type="button"
                  onClick={() => setActiveIndex(index)}
                  className={`h-2 w-2 rounded-full ${
                    activeIndex === index ? "bg-white" : "bg-white/30"
                  }`}
                  aria-label={`Go to image ${index + 1}`}
                />
              ))}
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
