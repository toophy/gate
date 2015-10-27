#盘古网游引擎
@(我的第一个笔记本)[盘古|GO|Lua|网游|引擎]
**盘古网游引擎**使用**Go**语言开发，支持丰富的**事件**系统，简便的**线程间消息通信**，快捷的**Lua**脚本，使用本款引擎更容易开发出高效稳定功能强大的**MMORPG**，行走江湖必备利器。

----------
[TOC]
##Screen 游戏场景
场景, 玩家的活动场所, 游戏世界最重要组成部分.

----------
### Get_data
Screen:Get_data() table
*获取场景的自定义数据table*
>返回值 : data, 是一个table, 包含场景的自定义数据

例子
```
function OnInit(s)
	s:Get_data()["lolo"] = "lolo"
	print(s:Get_data()["lolo"])
end
```
----------
### Get_name
Screen:Get_name() string
*获取场景名称*
>返回值 : 是一个string, 场景名称

例子
```
function OnInit(s)
	print(s:Get_name())
end
```
----------
### Get_id
Screen:Get_id() int
*获取场景ID*
>返回值 : 是一个int, 场景编号

例子
```
function OnInit(s)
	print(s:Get_id())
end
```
----------
### Get_oid
Screen:Get_oid() int
*获取场景模板ID*
>返回值 : 是一个int, 场景模板编号

例子
```
function OnInit(s)
	print(s:Get_oid())
end
```
----------
### PostEvent
Screen:PostEvent(func string, touchtime string, param table) bool
*向场景投递一条事件*
>|参数|类型|功能|
>|-|-|
>|func|string|Lua函数名|
>|touchtime|int|事件触发时间(毫秒)|
>|param|table|回调参数,lua的table|

>返回值 : bool

例子
```
module("woLongShanZhuang", package.seeall)

function OnInit(s)
	s:PostEvent("OnDingdong",1000,{["log"]="叮咚!叮咚!"})
end

function OnDingdong( t )
	print(t["log"])
end
```
----------
##ScreenThread 场景线程
场景线程, 包含大量的场景, Lua脚本的最终宿主, 公共变量ts指的是Lua环境中当前线程.
### Add_screen
ScreenThread:Add_screen(name string,oid int) bool
*增加一张场景*
>|参数|类型|功能|
>|--||-|
>|name|string|场景新名称|
>|oid|int|场景原始ID|

>返回值 : bool

例子
```
module("main", package.seeall)

function OnScreenThreadBegin()
	LogDebug("场景线程"..ts:Get_thread_id().." 启动")
	ts:Add_screen("新手村", 1)
end
```
----------
###Del_screen
ScreenThread:Del_screen(id int) bool
*删除一张场景*
>|参数|类型|功能|
>|-|-|-|
>|id|int|场景ID|

>返回值 : bool

例子
```
module("main", package.seeall)

function OnScreenThreadBegin()
	LogDebug("场景线程"..ts:Get_thread_id().." 启动")
	ts:Add_screen("新手村", 1)
	ts:Del_screen(1)
end
```
----------
###Get_screen
ScreenThread:Get_screen(id int) *Screen
*获取一张场景的对象*
>|参数|类型|功能|
>|-|-|-|
>|id|int|场景ID|

>返回值 : bool

例子
```
module("main", package.seeall)

function OnScreenThreadBegin()
	LogDebug("场景线程"..ts:Get_thread_id().." 启动")
	ts:Add_screen("新手村", 1)
	local city = ts:Get_screen(1)
	print(city:Get_name())
end
```
----------
###Get_thread_id
ScreenThread:Get_thread_id() int
*获取线程ID*
>返回值 : int

例子
```
module("main", package.seeall)

function OnScreenThreadBegin()
	print(ts:Get_thread_id())
end
```
----------
###PostEventFromLua
ScreenThread:PostEventFromLua(mod string, func string, time int, param table) bool
*向场景线程投递一条事件*
>|参数|类型|功能|
>|-|-|
>|mod|string|Lua模块名|
>|func|string|Lua函数名|
>|touchtime|int|事件触发时间(毫秒)|
>|param|table|回调参数,lua的table|

>返回值 : bool


例子
```
module("main", package.seeall)

function OnScreenThreadBegin()
	ts:PostEventFromLua("main", "OnDingdong", 1000, {"log":"叮咚!叮咚!"})
end

function OnDingdong(t)
	print(t["log"])
end
```

## 反馈与建议
- 邮箱：<toophy@vip.qq.com>

---------
感谢阅读这份帮助文档。

[场景][1]

[1]: screen.html
