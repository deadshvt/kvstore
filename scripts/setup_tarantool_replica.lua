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

local instance_id = os.getenv('INSTANCE_ID')
local listen_port = os.getenv('TARANTOOL_LISTEN_PORT')
local user = os.getenv('TARANTOOL_USER')
local password = os.getenv('TARANTOOL_USER_PASSWORD')

box.cfg{
    listen = listen_port
}

local kv_space = box.schema.space.create('kv_store', {if_not_exists=true})
kv_space:format({
    {name='key', type='string'},
    {name='value', type='any'}
})
kv_space:create_index('primary', {parts={'key'}, if_not_exists=true})

local user_space = box.schema.space.create('user_store', {if_not_exists=true})
user_space:format({
    {name='username', type='string'},
    {name='password', type='string'}
})
user_space:create_index('primary', {parts={'username'}, if_not_exists=true})

if not box.schema.user.exists(user) then
    box.schema.user.create(user, {password = password})
    box.schema.user.grant(user, 'read,write', 'space.kv_store')
    box.schema.user.grant(user, 'read,write', 'space.user_store')
end

local admin_user = user_space:get('admin')
if not admin_user then
    user_space:insert{'admin', 'presale'}
    print('Admin user created with username: admin and password: presale')
else
    print('Admin user already exists')
end

print("Tarantool replica " .. instance_id .. " configured and started successfully on port " .. listen_port)
