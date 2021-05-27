FROM scratch

WORKDIR /app

COPY ./main /app/switcheroo

EXPOSE 9543

ENTRYPOINT ["./switcheroo"]
