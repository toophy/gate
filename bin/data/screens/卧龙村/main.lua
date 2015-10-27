module("woLongShanZhuang", package.seeall)

-- local 
function OnInit(s)
	print("欢迎来到卧龙山庄.")
	-- s:Get_data()["lolo"] = "lolo"
	-- print(s:Get_data()["lolo"])
	s:PostEvent("OnHeartBeat",5000,{})
	s:PostEvent("Eon_Qiguan",10000,{["log"]="咕咕鸟在鸣叫!"})

	LogInfo(s:Get_name())
	LogInfo("id="..s:Get_id()..",oid="..s:Get_oid())
end

-- local 
function OnHeartBeat(s,t)
	LogInfo(s:Get_name().."心跳 "..os.time())
	s:PostEvent("OnHeartBeat",5000,{})
	-- print(s:Get_data()["lolo"])
	-- s:Get_data()["lolo"] = "lolo"..os.time()
end

-- local 
function Eon_Qiguan(s,t)
	LogInfo(t["log"])
end
