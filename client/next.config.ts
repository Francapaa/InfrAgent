import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  distDir: 'build',
  images: {
    unoptimized: true,
  },
  experimental: {
  },
};

export default nextConfig;
