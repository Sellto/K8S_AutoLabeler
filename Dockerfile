FROM alpine
ADD AutoLabeler.go AutoLabeler.go
ADD Install.sh Install.sh
RUN chmod +x Install.sh && ./Install.sh
VOLUME /dev /dev
CMD udevd --daemon --debug && udevadm control --reload-rules && udevadm trigger && ./AutoLabeler
