FROM golang as builder

COPY . /src
RUN cd /src && go build -o /org-contributors
RUN chmod +x /org-contributors
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
