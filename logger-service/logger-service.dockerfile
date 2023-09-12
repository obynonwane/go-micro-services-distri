FROM alpine:latest

RUN mkdir /app
#compile and call it loggerServiceApp and copy it into the app directory of the docker container
COPY loggerServiceApp /app  

#Execute or run the compiled docker container
CMD [ "/app/loggerServiceApp"]