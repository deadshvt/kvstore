local function load_env(filename)
    local file = io.open(filename, "r")
    if not file then return end

    for line in file:lines() do
        for key, value in string.gmatch(line, "([%w_]+)=([%w%p]+)") do
            os.setenv(key, value)
        end
    end

    file:close()
end

load_env('/project/.env')

local vshard = require('vshard')

local function create_shard_config()
    local sharding = {}
    local shard_count = tonumber(os.getenv('TARANTOOL_SHARD_COUNT'))

    for i = 1, shard_count do
        local shard_name = 'shard' .. i
        local replica_uri = os.getenv('TARANTOOL_REPLICA' .. i .. '_URI')
        local replica_name = 'replica' .. i

        sharding[shard_name] = {
            replicas = {
                [replica_name] = {
                    uri = replica_uri,
                    name = replica_name
                }
            }
        }
    end

    return sharding
end

vshard.router.cfg({
    bucket_count = 3000,
    sharding = create_shard_config()
})

local user = os.getenv('TARANTOOL_USER')
local password = os.getenv('TARANTOOL_USER_PASSWORD')

if not box.schema.user.exists(user) then
    box.schema.user.create(user, {password = password})
    box.schema.user.grant(user, 'read,write,execute', 'universe')
end

print("Tarantool router configured and started successfully")

