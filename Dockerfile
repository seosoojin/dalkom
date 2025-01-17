FROM public.ecr.aws/docker/library/golang:1.22.2-alpine3.18 AS build-env



RUN apk --no-cache add g++ ca-certificates tzdata build-base git

COPY . /project
WORKDIR /project
RUN go install -tags musl

FROM public.ecr.aws/docker/library/alpine:3.18
RUN apk --no-cache add g++ ca-certificates tzdata git
WORKDIR /workspace
COPY --from=build-env /go/bin/forge /go/bin/forge
ENTRYPOINT ["/go/bin/dalkom"]