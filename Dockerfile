FROM golang AS go-build

WORKDIR /build
COPY ./ ./
RUN go get -u
RUN go build -o ocr main.go

FROM golang AS go-run

WORKDIR /app
COPY --from=go-build /build/ocr .

EXPOSE 8080
CMD [ "./ocr" ]