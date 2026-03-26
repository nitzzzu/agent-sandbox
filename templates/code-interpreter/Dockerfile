FROM e2bdev/code-interpreter:latest
 
USER root
 
COPY envd /workspace/envd

COPY warmup.sh /root/.server/warmup.sh
RUN chmod +x /root/.server/warmup.sh

COPY startup.sh /root/.server/startup.sh
RUN chmod +x /root/.server/startup.sh

COPY start-up.sh /root/.server/start-up.sh
RUN chmod +x /root/.server/start-up.sh

COPY main.py /root/.server/main.py


ENTRYPOINT ["/root/.server/start-up.sh"]
