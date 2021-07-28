FROM golang:1.16.6-buster
RUN go get -u github.com/kejith/ginbar-backend
ENV APP_USER app
ENV APP_HOME /go/src/ginbar-backend
ARG GROUP_ID
ARG USER_ID
RUN groupadd --gid $GROUP_ID app && useradd -m -l --uid $USER_ID --gid  $GROUP_ID $APP_USER
RUN mkdir $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME
USER $APP_USER
WORKDIR $APP_HOME
EXPOSE 8080
RUN make build && make RUN