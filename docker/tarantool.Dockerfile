FROM tarantool/tarantool:2.10

WORKDIR /project

COPY scripts/setup_tarantool.lua .

CMD ["tarantool", "/project/setup_tarantool.lua"]