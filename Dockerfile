FROM harbor.voneyun.com/k8s-dev/alpine:3.4-v1
WORKDIR /app
RUN export LANG=zh_CN.UTF-8
ENV env sbx
ENV project_name  nonick
ENV service_name nonick-notifier-service
RUN mkdir /app/conf
COPY nonick-notifier-service /app/
COPY start.sh /app/
RUN chmod u+x ./start.sh
ENTRYPOINT ["sh","start.sh"]
EXPOSE 80
