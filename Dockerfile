# build stage
FROM golang:1.22-alpine AS build

# set working directory
WORKDIR /app

# copy source code
COPY . .

RUN go mod download

RUN go build -o ./bin/instagram_recipe_parser ./main.go


FROM alpine:latest AS final

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ffmpeg

COPY --from=build /app/bin/instagram_recipe_parser ./
COPY .env ./

EXPOSE 8080

ENTRYPOINT [ "./instagram_recipe_parser" ]
