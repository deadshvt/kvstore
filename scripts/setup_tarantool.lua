local user = 'user'
local password = 'password'

box.cfg{
    listen = 3301
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
    box.schema.user.grant(user, 'read,write', 'space', 'kv_store')
    box.schema.user.grant(user, 'read,write', 'space', 'user_store')
end

local admin_user = user_space:get('admin')
if not admin_user then
    user_space:insert{'admin', '6442daf21b7da3202fd27b8f76d79656c1f9f9e9e98b6bfee2b56c455d5660408533b081fd'}
    print('Admin user created with username: admin and password: presale')
else
    print('Admin user already exists')
end
