# Production stage
FROM centos:7
LABEL org.opencontainers.image.source https://github.com/agent-sandbox/agent-sandbox

COPY ./agent-sandbox /app
COPY ./default-environments.json /default-environments.json
RUN chmod +x /app
CMD /app
