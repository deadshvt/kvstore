FROM tarantool/tarantool:2.10

WORKDIR /project

COPY scripts/setup_tarantool_replica.lua .
COPY .env .

CMD ["tarantool", "/project/setup_tarantool_replica.lua"]
