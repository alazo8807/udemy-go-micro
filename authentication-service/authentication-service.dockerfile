# Note: We are building the executable as a step in the Makefile, so in here we just need to copy
# the executable in the container and run it.
#build a tiny docker image
FROM alpine:latest

#Create a folder called app inside the docker container
RUN mkdir /app

#From our base container(above) copy the executable into the new container's app folder
# COPY --from=builder /app/brokerApp /app

COPY authApp /app

CMD [ "/app/authApp" ]