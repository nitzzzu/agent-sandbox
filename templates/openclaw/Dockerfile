FROM node:24.13.0
LABEL org.opencontainers.image.source https://github.com/agent-sandbox/agent-sandbox

ENV NO_ONBOARD=1
RUN npm install -g openclaw@v2026.3.1

RUN npm install -g mcporter

RUN mkdir -p /root/.openclaw
COPY openclaw.json /root/.openclaw/openclaw.json
CMD openclaw gateway run