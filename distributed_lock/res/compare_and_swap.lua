if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call("SET", KEYS[1], ARGV[2])
else
    return 0
end