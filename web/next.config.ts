import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone', // Enable standalone output for Docker
  eslint: {
    // Don't fail build on ESLint errors (warnings are still shown)
    ignoreDuringBuilds: true,
  },
  typescript: {
    // Don't fail build on TypeScript errors (warnings are still shown)
    ignoreBuildErrors: false, // Keep this false to catch real TS errors
  },
};

export default nextConfig;
