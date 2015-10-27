module("main", package.seeall)

-- local 
function OnWorldBegin()
	LogDebug("世界线程启动")
	ts:CreateScreenThread(Tid_screen_1, "场景线程1", 100, Evt_lay1_time, 60000)
	-- ts:CreateScreenThread(Tid_screen_2, "场景线程2", 100, Evt_lay1_time, 60000)
end

-- local 
function OnWorldEnd()
	LogDebug("世界线程结束")
end
