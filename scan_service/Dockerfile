FROM mkodockx/docker-clamav
RUN apt-get update && \
    apt-get install -y ca-certificates
COPY scan_service /app
EXPOSE 8000
CMD ./app