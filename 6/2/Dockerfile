FROM nginx:alpine

COPY index.html /usr/share/nginx/html/index.html

ARG GREETING="Hello, Docker!"

ENV GREETING $GREETING

RUN sed -i "s/<h1>.*<\/h1>/<h1>$GREETING<\/h1>/" /usr/share/nginx/html/index.html

CMD ["nginx", "-g", "daemon off;"]
