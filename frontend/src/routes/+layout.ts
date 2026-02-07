// Disable SSR for static SPA deployment
// All rendering happens client-side, API calls go to /api/v1/*
export const prerender = false;
export const ssr = false;
export const csr = true;
