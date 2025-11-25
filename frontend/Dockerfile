# Frontend Dockerfile for Homelab File Manager
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN npm install
RUN npm install @sveltejs/adapter-node

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage with Node adapter
FROM node:22-alpine AS runner

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 nodejs && \
    adduser -u 1001 -G nodejs -s /bin/sh -D sveltekit

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./

USER sveltekit

EXPOSE 3000

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

CMD ["node", "build"]
