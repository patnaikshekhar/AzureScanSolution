FROM mkodockx/docker-clamav
RUN apt-get update && \
    apt-get install -y ca-certificates
COPY scan_service /app
COPY main.sh /
EXPOSE 8000
CMD ["sh", "main.sh"]