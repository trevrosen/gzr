FROM node:7

ADD . /usr/src/app/public
WORKDIR /usr/src/app
COPY start.sh /usr/local/bin/start.sh
COPY docker-config/kube.conf /root/.kube/config
COPY docker-config/gzr.json /root/.gzr.json

# Add key
RUN mkdir -p /root/.ssh && \
    ssh-keyscan -H github.com >> ~/.ssh/known_hosts
COPY id_rsa /root/.ssh/id_rsa
RUN chmod -R go-rwx /root/.ssh

RUN npm install -g webpack && \
    cd public && \
    npm install && \
    npm run build

# NOTE: Go and Glide install shouldn't be needed after this project is public, same with key

# Go
RUN wget -q https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz && \
    tar xvf go1.7.4.linux-amd64.tar.gz && \
    /bin/rm go1.7.4.linux-amd64.tar.gz && \
    mv ./go /usr/local && \
    cp /usr/local/go/bin/* /usr/local/bin && \
    mkdir -p /go/bin /go/src
ENV GOPATH "/go"
ENV PATH "$PATH:/usr/local/go/bin:$GOPATH/bin"

# Glide
# NOTE: at the time of writing glide.sh does not resolve https://github.com/Masterminds/glide/issues/784
RUN wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz && \
    tar xzvf glide-v0.12.3-linux-amd64.tar.gz && \
    cp linux-amd64/glide /go/bin/glide && \
    rm -r linux-amd64 && rm glide-v0.12.3-linux-amd64.tar.gz

# Kubectl
RUN curl -sLO https://storage.googleapis.com/kubernetes-release/release/v1.4.7/bin/linux/amd64/kubectl && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin

# Gzr
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts && \
    curl gzr.io/get | bash && \
    mv gzr /usr/local/bin

# Remove key
RUN rm -r ~/.ssh/id_rsa

ENTRYPOINT ["start.sh"]
