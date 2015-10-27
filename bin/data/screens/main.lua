module("main", package.seeall)

-- local 
function OnScreenThreadBegin()
	LogDebug("场景线程"..ts:Get_thread_id().." 启动")

	ts:Add_screen("阿拉斯加2", 1)
end

-- local 
function OnScreenThreadEnd()
	LogDebug("场景线程"..ts:Get_thread_id().." 结束")
end
