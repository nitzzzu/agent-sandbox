# Production stage
FROM centos:7
LABEL org.opencontainers.image.source https://github.com/agent-sandbox/agent-sandbox

COPY ./config/ /config/

COPY ./agent-sandbox /app
RUN chmod +x /app

CMD /app
