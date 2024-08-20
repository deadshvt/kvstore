FROM tarantool/tarantool:2.10

RUN apk add --no-cache cmake make gcc g++ musl-dev

RUN apk add --no-cache luarocks

RUN tarantoolctl rocks install vshard

WORKDIR /project

COPY scripts/setup_tarantool_router.lua .
COPY .env .

CMD ["tarantool", "/project/setup_tarantool_router.lua"]