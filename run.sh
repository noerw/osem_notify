go fmt ./ ./cmd/ ./core/ && \
    go build ./ && \
    ./osem_notify check boxes \
    593bcd656ccf3b0011791f5a 5b26181b1fef04001b69093c 59b31b8dd67eb50011165a04 562bdcf3b3de1fe005e03d2a $@ \
    --log-level debug
