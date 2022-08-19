FROM registry.geoway.com/jdk/centos-rsmis-jdk:1.8

WORKDIR /var/spool/rsmis/service

COPY target/*.jar /var/spool/rsmis/service/webapp.jar
COPY target/classes/*.yaml /var/spool/rsmis/service/config/
# COPY target/*-encrypted.lic /var/spool/rsmis/service/webapp.lic

ENTRYPOINT java -jar webapp.jar

EXPOSE 8080

