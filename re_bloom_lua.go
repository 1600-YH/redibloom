package re_bloom

const LuaBloomBatchGetBits = `
	local bloomKey = KEYS[1]
	local bitsCnt = ARGV[1]
	for i=1, bitsCnt, 1 do
		local offset = ARGV[1+i]
		local reply = redis.call("getbit", bloomKey, offset)
		if (no reply) then
			error('FATAL')
			return 0
		end
		if (reply == 0) then
			return 0
		end
	end
	return 1
`
const LuaBloomBatchSetBits = `
	local bloomKey = KEYS[1]
	local bitsCnt = ARGV[1]
	for i=1, bitsCnt, 1 do
		local offset = ARGV[1+i]
		redis.call("setbit", bloomKey, offset, 1)
	end
	return 1
`
