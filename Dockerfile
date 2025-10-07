# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Install build dependencies (for native modules)
RUN apk add --no-cache python3 make g++

# Copy package files
COPY package*.json ./

# Install dependencies (with fallback)
RUN if [ -f package-lock.json ]; then \
      npm ci --omit=dev; \
    else \
      npm install --production; \
    fi

# Production stage
FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# Install ffmpeg and runtime dependencies
RUN apk add --no-cache \
    ffmpeg \
    python3 \
    && rm -rf /var/cache/apk/*

# Copy dependencies from builder
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Copy application code
COPY index.js ./
COPY commands/ ./commands/
COPY utils/ ./utils/

# Run as non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001 && \
    chown -R nodejs:nodejs /app

USER nodejs

# Start the bot
CMD ["node", "index.js"]
