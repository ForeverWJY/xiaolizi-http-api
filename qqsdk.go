package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//测试方法 启动时会调用一次
func testFunc() {
	//getnickname(0, 1066231345, true)
	//getgroupnamefromcache(62827649)
	//getfriendlist(LoginQQ)
	//getgrouplist(LoginQQ)
	//getgroupmemberlist(LoginQQ, 62827649)
	//getgroupmgrlist(LoginQQ, 62827649)
	//getgroupcard(LoginQQ, 62827649, 1066231345)
	//getsignat(LoginQQ, 1066231345)
	//setnickname(LoginQQ, "bubucom")
	//setsignat(LoginQQ, "bubucom")
	//getmutetime(LoginQQ, 62827649)
	//getphotourl("[pic,hash=3FBD19CBF124ECF646333B19069BD51C,guid=\\/1066231345-2421157924-3FBD19CBF124ECF646333B19069BD51C]", LoginQQ, 0)
	//sendfreepackage(LoginQQ, 62827649, 1066231345, 299)
	//getqqonlinestate(LoginQQ, 1066231345)
	//sharemusic(LoginQQ, 0, 1066231345, "Uptown Funk", "Mark Ronson/Bruno Mars",
	//	"https://y.music.163.com/m/song/29722263/", "https://cpic.url.cn/v1/vk04ujo7pbgh88vh8ppt96pqmg139518ae10nth1gku65sivp3ft79b7psfekit5s53uf7gig7ddr8ia5genc1p627p657tdlbo6bfshtq053dghadu10b0u84hq8dak/7ggqo6k576kp89pk1cei1d8ea6khc37k2bph6p19fj6e2cp9fr10",
	//	"http://music.163.com/song/media/outer/url?id=29722263", 4)
	//getskey(LoginQQ)
	//getpskey(LoginQQ, "openmobile.qq.com")
	//getclientkey(LoginQQ)
}

func reply(reqMsg *EventResponse, msg string) *ReplyReturnRet {
	switch reqMsg.MessageType {
	case "private":
		return sendPrivateMsg(reqMsg.Sender.UserID, msg)
	case "group":
		return sendGroupMsg(reqMsg.GroupID, msg)
	}
	return nil
}

type ReplyReturnRet struct {
	MessageID int32 `json:"message_id"`
}

type Ret struct {
	Ret string `json:"ret"`
}

func getLoginQQ() int {
	qq := get(fmt.Sprintf("http://%v/getlogonqq", addr))
	if qq == nil {
		logApi.Debug("获取QQ号出错")
		return 0
	}
	var ret = new(Ret)
	_ = json.Unmarshal(qq, &ret)
	logApi.Debug(ret.Ret)

	var event map[string]interface{}
	_ = json.Unmarshal([]byte(ret.Ret), &event)
	for qq := range event["QQlist"].(map[string]interface{}) {
		logApi.Debug("当前登录QQ：", qq)
		atoi, _ := strconv.Atoi(qq)
		return atoi
	}
	return 0
}

//发送好友消息
func sendPrivateMsg(toqq int, text string) *ReplyReturnRet {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("user_id", fmt.Sprintf("%v", toqq))
	r.Form.Add("message", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/send_private_msg", httpAddr), bodystr)
	if data != nil {
		var ret = new(ReplyReturnRet)
		_ = json.Unmarshal(data, &ret)
		return ret
	}
	return nil
}

//发送群消息
func sendGroupMsg(togroup int, text string) *ReplyReturnRet {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("group_id", fmt.Sprintf("%v", togroup))
	r.Form.Add("message", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/send_group_msg", addr), bodystr)
	if data != nil {
		var ret = new(ReplyReturnRet)
		_ = json.Unmarshal(data, &ret)
		return ret
	}
	return nil
}

//发送群临时消息
func sendgrouptempmsg(fromqq int, togroup int, toqq int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("text", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgrouptempmsg", addr), bodystr)
}

//添加好友
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [可选]text 指定附言
func addfriend(fromqq int, toqq int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("text", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/addfriend", addr), bodystr)
}

//添加群
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [可选]text 指定附言
func addgroup(fromqq int, togroup int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("text", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/addgroup", addr), bodystr)
}

//删除好友
func deletefriend(fromqq int, toqq int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/deletefriend", addr), bodystr)
}

//置屏蔽好友
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [必须]ignore 指定是否屏蔽(true,false)
func setfriendignmsg(fromqq int, toqq int, ignore bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	if ignore == true {
		r.Form.Add("ignore", "true")
	} else {
		r.Form.Add("ignore", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/setfriendignmsg", addr), bodystr)
}

//置特别关心好友
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [必须]care 指定是否关心(true,false)
func setfriendcare(fromqq int, toqq int, care bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	if care == true {
		r.Form.Add("care", "true")
	} else {
		r.Form.Add("care", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/setfriendcare", addr), bodystr)
}

//发送好友XML消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [必须]xml 指定消息内容(存在特殊字符请使用URL编码)
func sendprivatexmlmsg(fromqq int, toqq int, xml string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("xml", xml)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendprivatexmlmsg", addr), bodystr)
}

//发送群XML消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]xml 指定消息内容(存在特殊字符请使用URL编码) [可选]anonymous 指定是否匿名(true,false)
func sendgroupxmlmsg(fromqq int, togroup int, xml string, anonymous bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("xml", xml)
	if anonymous == true {
		r.Form.Add("anonymous", "true")
	} else {
		r.Form.Add("anonymous", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgroupxmlmsg", addr), bodystr)
}

//发送好友JSON消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [必须]json 指定消息内容(存在特殊字符请使用URL编码)
func sendprivatejsonmsg(fromqq int, toqq int, json string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("json", json)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendprivatejsonmsg", addr), bodystr)
}

//发送群JSON消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]json 指定消息内容(存在特殊字符请使用URL编码) [可选]anonymous 指定是否匿名(true,false)
func sendgroupjsonlmsg(fromqq int, togroup int, json string, anonymous bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("json", json)
	if anonymous == true {
		r.Form.Add("anonymous", "true")
	} else {
		r.Form.Add("anonymous", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgroupjsonlmsg", addr), bodystr)
}

//上传好友图片，返回值可用于发送图片
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]toqq 指定好友QQ
//[可选]fromtype 指定图片来源类型(0:pic参数,1:本地文件,2:网络文件 默认为0)
//[fromtype=0时必须]pic 指定数据(请使用BASE64+URL编码:url_encode(base64_encode(src)))
//[fromtype=1时必须]path 指定文件路径(请使用绝对路径,存在特殊字符请使用URL编码)
//[fromtype=2时必须]url 指定文件url(存在特殊字符请使用URL编码)
//[可选]flashpic 指定是否闪照(true,false)
func sendprivatepic(fromqq int, toqq int, fromtype int, pic string, path string, url string, flashpic bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	if fromtype == 0 || fromtype == 1 || fromtype == 2 {
		r.Form.Add("fromtype", fmt.Sprintf("%v", fromtype))
	}
	if fromtype == 0 {
		r.Form.Add("pic", pic)
	} else if fromtype == 1 {
		r.Form.Add("path", path)
	} else if fromtype == 2 {
		r.Form.Add("url", url)
	}
	if flashpic == true {
		r.Form.Add("flashpic", "true")
	} else {
		r.Form.Add("flashpic", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendprivatepic", addr), bodystr)
}

//上传群图片，返回值可用于发送图片
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]togroup 指定群号
//[可选]fromtype 指定图片来源类型(0:pic参数,1:本地文件,2:网络文件 默认为0)
//[fromtype=0时必须]pic 指定数据(请使用BASE64+URL编码:url_encode(base64_encode(src)))
//[fromtype=1时必须]path 指定文件路径(请使用绝对路径,存在特殊字符请使用URL编码)
//[fromtype=2时必须]url 指定文件url(存在特殊字符请使用URL编码)
//[可选]flashpic 指定是否闪照(true,false)
func sendgrouppic(fromqq int, togroup int, fromtype int, pic string, path string, url string, flashpic bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	if fromtype == 0 || fromtype == 1 || fromtype == 2 {
		r.Form.Add("fromtype", fmt.Sprintf("%v", fromtype))
	}
	if fromtype == 0 {
		r.Form.Add("pic", pic)
	} else if fromtype == 1 {
		r.Form.Add("path", path)
	} else if fromtype == 2 {
		r.Form.Add("url", url)
	}
	if flashpic == true {
		r.Form.Add("flashpic", "true")
	} else {
		r.Form.Add("flashpic", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgrouppic", addr), bodystr)
}

//向好友发送语音
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]toqq 指定好友QQ
//[可选]audiotype 指定语音类型(0普通语音,1变声语音,2文字语音,3红包匹配语音)
//[可选]text 指定语音文字
//[可选]fromtype 指定语音来源类型(0:pic参数,1:本地文件,2:网络文件 默认为0)
//[fromtype=0时必须]audio 指定数据(请使用BASE64+URL编码:url_encode(base64_encode(src)))
//[fromtype=1时必须]path 指定文件路径(请使用绝对路径,存在特殊字符请使用URL编码)
//[fromtype=2时必须]url 指定文件url(存在特殊字符请使用URL编码)
func sendprivateaudio(fromqq int, toqq int, audiotype int, text string, fromtype int, audio string, path string, url string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	if audiotype == 0 || audiotype == 1 || audiotype == 2 || audiotype == 3 {
		r.Form.Add("type", fmt.Sprintf("%v", audiotype))
	}
	if fromtype == 0 || fromtype == 1 || fromtype == 2 {
		r.Form.Add("fromtype", fmt.Sprintf("%v", fromtype))
	}
	if text != "" {
		r.Form.Add("text", text)
	}
	if fromtype == 0 {
		r.Form.Add("audio", audio)
	} else if fromtype == 1 {
		r.Form.Add("path", path)
	} else if fromtype == 2 {
		r.Form.Add("url", url)
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendprivateaudio", addr), bodystr)
}

//向群发送语音
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]togroup 指定群号
//[可选]type 指定语音类型(0普通语音,1变声语音,2文字语音,3红包匹配语音)
//[可选]text 指定语音文字
//[可选]fromtype 指定语音来源类型(0:pic参数,1:本地文件,2:网络文件 默认为0)
//[fromtype=0时必须]audio 指定数据(请使用BASE64+URL编码:url_encode(base64_encode(src)))
//[fromtype=1时必须]path 指定文件路径(请使用绝对路径,存在特殊字符请使用URL编码)
//[fromtype=2时必须]url 指定文件url(存在特殊字符请使用URL编码)
func sendgroupaudio(fromqq int, togroup int, audiotype int, text string, fromtype int, audio string, path string, url string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	if audiotype == 0 || audiotype == 1 || audiotype == 2 || audiotype == 3 {
		r.Form.Add("type", fmt.Sprintf("%v", audiotype))
	}
	if fromtype == 0 || fromtype == 1 || fromtype == 2 {
		r.Form.Add("fromtype", fmt.Sprintf("%v", fromtype))
	}
	if text != "" {
		r.Form.Add("text", text)
	}
	if fromtype == 0 {
		r.Form.Add("audio", audio)
	} else if fromtype == 1 {
		r.Form.Add("path", path)
	} else if fromtype == 2 {
		r.Form.Add("url", url)
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgroupaudio", addr), bodystr)
}

//上传头像
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[可选]fromtype 指定图片来源类型(0:pic参数,1:本地文件,2:网络文件 默认为0)
//[fromtype=0时必须]pic 指定数据(请使用BASE64+URL编码:url_encode(base64_encode(src)))
//[fromtype=1时必须]path 指定文件路径(请使用绝对路径,存在特殊字符请使用URL编码)
//[fromtype=2时必须]url 指定文件url(存在特殊字符请使用URL编码)
func uploadfacepic(fromqq int, fromtype int, pic string, path string, url string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	if fromtype == 0 || fromtype == 1 || fromtype == 2 {
		r.Form.Add("fromtype", fmt.Sprintf("%v", fromtype))
	}
	if fromtype == 0 {
		r.Form.Add("pic", pic)
	} else if fromtype == 1 {
		r.Form.Add("path", path)
	} else if fromtype == 2 {
		r.Form.Add("url", url)
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/uploadfacepic", addr), bodystr)
}

//上传群头像
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]group 指定群号
//[可选]fromtype 指定图片来源类型(0:pic参数,1:本地文件,2:网络文件 默认为0)
//[fromtype=0时必须]pic 指定数据(请使用BASE64+URL编码:url_encode(base64_encode(src)))
//[fromtype=1时必须]path 指定文件路径(请使用绝对路径,存在特殊字符请使用URL编码)
//[fromtype=2时必须]url 指定文件url(存在特殊字符请使用URL编码)
func uploadgroupfacepic(fromqq int, group int, fromtype int, pic string, path string, url string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if fromtype == 0 || fromtype == 1 || fromtype == 2 {
		r.Form.Add("fromtype", fmt.Sprintf("%v", fromtype))
	}
	if fromtype == 0 {
		r.Form.Add("pic", pic)
	} else if fromtype == 1 {
		r.Form.Add("path", path)
	} else if fromtype == 2 {
		r.Form.Add("url", url)
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/uploadgroupfacepic", addr), bodystr)
}

//设置群名片
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]toqq 指定群成员QQ [必须]card 指定群名片(存在特殊字符请使用URL编码)
func setgroupcard(fromqq int, togroup int, toqq int, card string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("card", card)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/setgroupcard", addr), bodystr)
}

//取昵称
//提交请求参数:[不使用缓存则必须,使用缓存则不须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [可选]fromcache 指定是否使用缓存(true,false)
func getnickname(fromqq int, toqq int, fromcache bool) string {
	var r http.Request
	_ = r.ParseForm()
	if fromcache {
		r.Form.Add("fromcache", "true")
	} else {
		r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
		r.Form.Add("fromcache", "false")
	}
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getnickname", addr), bodystr)
	if data == nil {
		logApi.Debug("取昵称出错")
		return ""
	}
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debugf("qq[%v]的昵称：[%v]", toqq, ret.Ret)
	return ret.Ret
}

//从缓存取群名称
//提交请求参数:[必须]group 指定群号
func getgroupnamefromcache(group int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getgroupnamefromcache", addr), bodystr)
	if data == nil {
		logApi.Debug("从缓存取群名称出错")
		return ""
	}
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debugf("qq群[%v]的群名称：[%v]", group, ret.Ret)
	return ret.Ret
}

//取好友列表
//提交请求参数:[必须]logonqq 指定框架QQ
func getfriendlist(logonqq int) *FriendListRet {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getfriendlist", addr), bodystr)
	if data == nil {
		logApi.Error("取好友列表出错")
		return nil
	}

	var ret = new(FriendListRet)
	_ = json.Unmarshal(data, &ret)
	for _, v := range ret.List {
		logApi.Debugf("\nqq:[%v]\nnickName:[%v]\nRemark:[%v] \nemail:[%v]", v.UIN, v.NickName, v.Remark, v.Email)
	}
	return ret
}

type FriendListRet struct {
	Ret  string `json:"ret"`
	List []struct {
		UIN      int    `json:"UIN"`
		NickName string `json:"NickName"`
		Remark   string `json:"Remark"`
		Email    string `json:"Email"`
	} `json:"List"`
}

//取群列表
//提交请求参数:[必须]logonqq 指定框架QQ
func getgrouplist(logonqq int) *GroupListRet {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getgrouplist", addr), bodystr)
	if data == nil {
		logApi.Error("取群列表出错")
		return nil
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(GroupListRet)
	_ = json.Unmarshal(data, &ret)
	for _, v := range ret.List {
		logApi.Debugf("\n群号:[%v]\n群名称:[%v]\n群所有者QQ:[%v] \n群成员数量:[%v]", v.GIN, v.StrGroupName, v.DwGroupOwnerUin, v.DwMemberNum)
	}
	return ret
}

type GroupListRet struct {
	Ret  string `json:"ret"`
	List []struct {
		GroupID                int    `json:"GroupID"`
		GIN                    int    `json:"GIN"`
		CFlag                  int    `json:"cFlag"`
		GroupInfoSeq           int    `json:"GroupInfoSeq"`
		DwGroupFlagExt         int    `json:"dwGroupFlagExt"`
		DwGroupRankSeq         int    `json:"dwGroupRankSeq"`
		DwCertificationType    int    `json:"dwCertificationType"`
		DwShutupTimestamp      int    `json:"dwShutupTimestamp"`
		DwMyShutupTimestamp    int    `json:"dwMyShutupTimestamp"`
		DwCmdUinUinFlag        int    `json:"dwCmdUinUinFlag"`
		DwAdditionalFlag       int    `json:"dwAdditionalFlag"`
		DwGroupTypeFlag        int    `json:"dwGroupTypeFlag"`
		DwGroupSecType         int    `json:"dwGroupSecType"`
		DwGroupSecTypeInfo     int    `json:"dwGroupSecTypeInfo"`
		DwGroupClassExt        int    `json:"dwGroupClassExt"`
		DwAppPrivilegeFlag     int    `json:"dwAppPrivilegeFlag"`
		DwSubscriptionUin      int    `json:"dwSubscriptionUin"`
		DwMemberNum            int    `json:"dwMemberNum"`
		DwMemberNumSeq         int    `json:"dwMemberNumSeq"`
		DwMemberCardSeq        int    `json:"dwMemberCardSeq"`
		DwGroupFlagExt3        int    `json:"dwGroupFlagExt3"`
		DwGroupOwnerUin        int    `json:"dwGroupOwnerUin"`
		CIsConfGroup           int    `json:"cIsConfGroup"`
		CIsModifyConfGroupFace int    `json:"cIsModifyConfGroupFace"`
		CIsModifyConfGroupName int    `json:"cIsModifyConfGroupName"`
		DwCmduinJoinTime       int    `json:"dwCmduinJoinTime"`
		StrGroupName           string `json:"strGroupName"`
		StrGroupMemo           string `json:"strGroupMemo"`
	} `json:"List"`
}

//取群成员列表
//提交请求参数:[必须]logonqq 指定框架QQ [必须]group 指定群号
func getgroupmemberlist(logonqq int, group int) *GroupMemberListRet {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getgroupmemberlist", addr), bodystr)
	if data == nil {
		logApi.Error("取群成员列表出错")
		return nil
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(GroupMemberListRet)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	for _, v := range ret.List {
		logApi.Debugf("\nQQ:[%v]\n昵称:[%v]\n性别:[%v] \n年龄:[%v]", v.UIN, v.NickName, v.Sex, v.Sex)
	}
	return ret
}

type GroupMemberListRet struct {
	Ret  string `json:"ret"`
	List []struct {
		UIN              int    `json:"UIN"`
		Age              int    `json:"Age"`
		Sex              int    `json:"Sex"`
		NickName         string `json:"NickName"`
		Email            string `json:"Email"`
		Card             string `json:"Card"`
		Remark           string `json:"Remark"`
		SpecTitle        string `json:"SpecTitle"`
		Phone            string `json:"Phone"`
		SpecTitleExpired int    `json:"SpecTitleExpired"`
		MuteTime         int    `json:"MuteTime"`
		AddGroupTime     int    `json:"AddGroupTime"`
		LastMsgTime      int    `json:"LastMsgTime"`
		GroupLevel       int    `json:"GroupLevel"`
	} `json:"List"`
}

//设置管理员
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]toqq 指定对方QQ [必须]bemgr 是否成为管理员(true,false)
func setgroupmgr(fromqq int, group int, toqq int, bemgr bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	if bemgr == true {
		r.Form.Add("bemgr", "true")
	} else {
		r.Form.Add("bemgr", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/setgroupmgr", addr), bodystr)
}

//取管理层列表
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号
func getgroupmgrlist(fromqq int, group int) []string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getgroupmgrlist", addr), bodystr)
	if data == nil {
		logApi.Error("取群管理层列表出错")
		return nil
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	split := strings.Split(ret.Ret, "\r\n")
	if len(split) > 0 {
		return split
	}
	return nil
}

//取群名片
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]toqq 指定对方QQ
func getgroupcard(fromqq int, group int, toqq int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getgroupcard", addr), bodystr)
	if data == nil {
		logApi.Error("取群名片出错")
		return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

//取个性签名
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ
func getsignat(fromqq int, toqq int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getsignat", addr), bodystr)
	if data == nil {
		logApi.Error("取个性签名出错")
		return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

//设置昵称
//提交请求参数:[必须]fromqq 指定框架QQ [必须]nickname 指定昵称
func setnickname(fromqq int, nickname string) bool {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("nickname", nickname)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setnickname", addr), bodystr)
	if data == nil {
		logApi.Error("设置昵称出错")
		//return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseBool, err := strconv.ParseBool(ret.Ret)
	if err != nil {
		logApi.Error()
	}
	return parseBool
}

//设置个性签名
//提交请求参数:[必须]fromqq 指定框架QQ [必须]signature 指定个性签名
//测试时会崩溃
func setsignat(fromqq int, signature string) bool {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("signature", signature)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setsignat", addr), bodystr)
	if data == nil {
		logApi.Error("设置个性签名出错")
		return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseBool, err := strconv.ParseBool(ret.Ret)
	if err != nil {
		logApi.Error()
	}
	return parseBool
}

//移出群成员
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]toqq 指定对方QQ [可选]ignoreaddgrequest 拒绝再加群申请(true,false)
func kickgroupmember(fromqq int, group int, toqq int, ignoreaddgrequest bool) bool {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	if ignoreaddgrequest {
		r.Form.Add("ignoreaddgrequest", "true")
	} else {
		r.Form.Add("ignoreaddgrequest", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/kickgroupmember", addr), bodystr)
	if data == nil {
		logApi.Error("移出群成员出错")
		return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseBool, err := strconv.ParseBool(ret.Ret)
	if err != nil {
		logApi.Error()
	}
	return parseBool
}

//禁言群成员
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]toqq 指定对方QQ [必须]time 指定禁言时长(以秒计)
func mutegroupmember(fromqq int, group int, toqq int, time int) bool {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("time", fmt.Sprintf("%v", time))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/mutegroupmember", addr), bodystr)
	if data == nil {
		logApi.Error("禁言群成员出错")
		return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseBool, err := strconv.ParseBool(ret.Ret)
	if err != nil {
		logApi.Error()
	}
	return parseBool
}

//退群
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号
func exitgroup(fromqq int, group int) bool {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/exitgroup", addr), bodystr)
	if data == nil {
		logApi.Error("退群出错")
		return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseBool, err := strconv.ParseBool(ret.Ret)
	if err != nil {
		logApi.Error()
	}
	return parseBool
}

//解散群
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号
func dispgroup(fromqq int, group int) bool {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/dispgroup", addr), bodystr)
	if data == nil {
		logApi.Error("解散群出错")
		return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseBool, err := strconv.ParseBool(ret.Ret)
	if err != nil {
		logApi.Error()
	}
	return parseBool
}

//全员禁言
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]ismute 指定是否禁言(true,false)
func setgroupwholemute(fromqq int, group int, ismute bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if ismute {
		r.Form.Add("ismute", "true")
	} else {
		r.Form.Add("ismute", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgroupwholemute", addr), bodystr)
	if data == nil {
		logApi.Error("全员禁言出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_发起新的群聊
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_newgroup(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_newgroup", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_发起新的群聊出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_发起临时会话
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_newtempsession(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_newtempsession", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_发起临时会话出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_上传文件
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_uploadfile(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_uploadfile", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_上传文件出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_上传相册
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_uploadphotoalbum(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_uploadphotoalbum", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_上传相册出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_邀请他人加群
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_invitein(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_invitein", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_邀请他人加群出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_匿名聊天
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_anonymous(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_anonymous", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_匿名聊天出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_坦白说
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_tanbaishuo(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_tanbaishuo", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_坦白说出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_新成员查看历史消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]allow 指定是否允许(true,false)
func setgrouppriv_newmembercanviewhistorymsg(fromqq int, group int, allow bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	if allow {
		r.Form.Add("allow", "true")
	} else {
		r.Form.Add("allow", "false")
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_newmembercanviewhistorymsg", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_新成员查看历史消息出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//置群员权限_邀请方式
//提交请求参数:[必须]fromqq 指定框架QQ [必须]togroup 指定群号 [必须]way 指定方式(1.无需审核;2.需要管理员审核;3.100人以内无需审核)
func setgrouppriv_inviteway(fromqq int, togroup int, way int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("way", fmt.Sprintf("%v", way))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgrouppriv_inviteway", addr), bodystr)
	if data == nil {
		logApi.Error("置群员权限_邀请方式出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//撤回群聊消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]random 发送消息返回(或事件给出)的random [必须]req 发送消息返回(或事件给出)的req
func deletegroupmsg(fromqq int, group int, random int, req int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("random", fmt.Sprintf("%v", random))
	r.Form.Add("req", fmt.Sprintf("%v", req))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/deletegroupmsg", addr), bodystr)
	if data == nil {
		logApi.Error("撤回群聊消息出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//撤回私聊消息
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ [必须]random 发送消息返回的random [必须]req 发送消息返回的req [必须]time 发送消息返回的
func deleteprivatemsg(fromqq int, toqq int, random int, req int, time int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("random", fmt.Sprintf("%v", random))
	r.Form.Add("req", fmt.Sprintf("%v", req))
	r.Form.Add("time", fmt.Sprintf("%v", time))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/deleteprivatemsg", addr), bodystr)
	if data == nil {
		logApi.Error("撤回私聊消息出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//设置位置共享
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]posx 指定经度 [必须]posy 指定纬度 [必须]enable 指定是否开启
func setsharepos(fromqq int, group int, posx float32, posy float64, enable bool) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("posx", fmt.Sprintf("%v", posx))
	r.Form.Add("posy", fmt.Sprintf("%v", posy))
	r.Form.Add("enable", strconv.FormatBool(enable))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setsharepos", addr), bodystr)
	if data == nil {
		logApi.Error("设置位置共享出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//上报当前位置
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]posx 指定经度 [必须]posy 指定纬度
func uploadpos(fromqq int, group int, posx float32, posy float64) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("posx", fmt.Sprintf("%v", posx))
	r.Form.Add("posy", fmt.Sprintf("%v", posy))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/uploadpos", addr), bodystr)
	if data == nil {
		logApi.Error("上报当前位置出错")
		//return false
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	//parseBool, err := strconv.ParseBool(ret.Ret)
	//if err != nil {
	//	logApi.Error()
	//}
	//return parseBool
}

//取禁言时间
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号
//return 秒
func getmutetime(fromqq int, group int) int64 {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getmutetime", addr), bodystr)
	if data == nil {
		logApi.Error("取禁言时间出错")
		return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	parseInt, err := strconv.ParseInt(ret.Ret, 10, 64)
	if err != nil {
		logApi.Error(err.Error())
		return 0
	}
	return parseInt
}

//处理群验证事件
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]group 指定群号
//[必须]qq 指定来源QQ
//[必须]seq 指定seq
//[必须]op 指定处理类型(11同意 12拒绝  14忽略)
//[必须]type 指定事件类型(群事件_某人申请加群:3 群事件_我被邀请加入群:1)
func setgroupaddrequest(fromqq int, group int, qq int, seq int, op int, eventType int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("qq", fmt.Sprintf("%v", qq))
	r.Form.Add("seq", fmt.Sprintf("%v", seq))
	r.Form.Add("op", fmt.Sprintf("%v", op))
	r.Form.Add("type", fmt.Sprintf("%v", eventType))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setgroupaddrequest", addr), bodystr)
	if data == nil {
		logApi.Error("处理群验证事件出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//处理好友验证事件
//提交请求参数:[必须]fromqq 指定框架QQ [必须]qq 指定来源QQ [必须]seq 指定seq [必须]op 指定处理类型(1同意 2拒绝)
func setfriendaddrequest(fromqq int, qq int, seq int, op int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("qq", fmt.Sprintf("%v", qq))
	r.Form.Add("seq", fmt.Sprintf("%v", seq))
	r.Form.Add("op", fmt.Sprintf("%v", op))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setfriendaddrequest", addr), bodystr)
	if data == nil {
		logApi.Error("处理好友验证事件出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//上传文件
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]path 指定文件名(存在特殊字符请使用URL编码)
func uploadfile(fromqq int, group int, path string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("path", path)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/uploadfile", addr), bodystr)
	if data == nil {
		logApi.Error("上传文件出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//创建群文件夹
//提交请求参数:[必须]fromqq 指定框架QQ [必须]group 指定群号 [必须]folder 指定文件夹名称(存在特殊字符请使用URL编码)
func newgroupfolder(fromqq int, group int, folder string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("folder", folder)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/newgroupfolder", addr), bodystr)
	if data == nil {
		logApi.Error("创建群文件夹出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//设置在线状态
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]state 指定在线主状态(11在线 31离开 41隐身 50忙碌 60Q我吧 70请勿打扰)
//[当state=11时可选]sun 指定在线子状态1(0普通在线 1000我的电量 1011信号弱 1024在线学习 1025在家旅游 1027TiMi中 1016睡觉中 1017游戏中 1018学习中 1019吃饭中 1021煲剧中 1022度假中 1032熬夜中)
//[当sun=1000时可选]power 自动电量(取值1到100)
func setonlinestate(fromqq int, state int, sun int, power int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("state", fmt.Sprintf("%v", state))
	if state == 11 && sun >= 0 {
		r.Form.Add("sun", fmt.Sprintf("%v", sun))
		if sun == 1000 && power >= 1 && power <= 100 {
			r.Form.Add("power", fmt.Sprintf("%v", power))
		}
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/setonlinestate", addr), bodystr)
	if data == nil {
		logApi.Error("设置在线状态出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//发送名片赞
//提交请求参数:[必须]fromqq 指定框架QQ [必须]toqq 指定对方QQ
func sendlike(fromqq int, toqq int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/sendlike", addr), bodystr)
	if data == nil {
		logApi.Error("发送名片赞出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//取图片下载地址
//提交请求参数:[必须]photo 指定图片代码(存在特殊字符请使用URL编码) [群聊图片必填，私聊图片不填]fromqq 指定框架QQ [群聊图片必填，私聊图片不填]group 指定群号
func getphotourl(photo string, fromqq int, group int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	if group > 0 {
		r.Form.Add("group", fmt.Sprintf("%v", group))
	}
	r.Form.Add("photo", photo)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getphotourl", addr), bodystr)
	if data == nil {
		logApi.Error("取图片下载地址出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

//群文件转发至群
//提交请求参数:[必须]fromqq 指定框架QQ [必须]fromgroup 指定来源群 [必须]togroup 指定目标群 [必须]fileid 指定文件ID(存在特殊字符请使用URL编码)
func forwardgroupfiletogroup(fromqq int, fromgroup int, togroup int, fileid string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("fromgroup", fmt.Sprintf("%v", fromgroup))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("fileid", fileid)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/forwardgroupfiletogroup", addr), bodystr)
	if data == nil {
		logApi.Error("群文件转发至群出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//群文件转发至好友
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]fromgroup 指定来源群
//[必须]toqq 指定目标QQ
//[必须]fileid 指定文件ID(存在特殊字符请使用URL编码)
//[必须]filename 指定文件名(存在特殊字符请使用URL编码)
func forwardgroupfiletofriend(fromqq int, fromgroup int, toqq int, fileid string, filename string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("fromgroup", fmt.Sprintf("%v", fromgroup))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("fileid", fileid)
	r.Form.Add("filename", filename)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/forwardgroupfiletofriend", addr), bodystr)
	if data == nil {
		logApi.Error("群文件转发至好友出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//好友文件转发至好友
//提交请求参数:
//[必须]logonqq 指定框架QQ
//[必须]fromqq 指定来源QQ
//[必须]toqq 指定目标QQ
//[必须]fileid 指定文件ID(存在特殊字符请使用URL编码)
//[必须]filename 指定文件名(存在特殊字符请使用URL编码)
func forwardfriendfiletofriend(loginqq int, fromqq int, toqq int, fileid string, filename string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("loginqq", fmt.Sprintf("%v", loginqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("fileid", fileid)
	r.Form.Add("filename", filename)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/forwardfriendfiletofriend", addr), bodystr)
	if data == nil {
		logApi.Error("好友文件转发至好友出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//查看转发聊天记录内容
//提交请求参数:[必须]logonqq 指定框架QQ [必须]resid 指定resid(xml消息中包含)
func getforwardedmsg(logonqq int, resid string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("resid", resid)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getforwardedmsg", addr), bodystr)
	if data == nil {
		logApi.Error("查看转发聊天记录内容出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//查询用户信息
//提交请求参数:[必须]logonqq 指定框架QQ [必须]qq 指定欲查询QQ
func queryuserinfo(logonqq int, qq int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("qq", fmt.Sprintf("%v", qq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/queryuserinfo", addr), bodystr)
	if data == nil {
		logApi.Error("查询用户信息出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//查询群信息
//提交请求参数:[必须]logonqq 指定框架QQ [必须]group 指定欲查群号
func querygroupinfo(logonqq int, group int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/querygroupinfo", addr), bodystr)
	if data == nil {
		logApi.Error("查询群信息出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//发送免费礼物
//提交请求参数:
//[必须]fromqq 指定框架QQ
//[必须]group 指定群号
//[必须]toqq 指定对方QQ
//[必须]pkgid 指定礼物类型(299卡布奇诺;302猫咪手表;280牵你的手;281可爱猫咪;284神秘面具;285甜wink;286我超忙的;289快乐肥宅水;290幸运手链;313坚强;307绒绒手套; 312爱心口罩;308彩虹糖果)
func sendfreepackage(fromqq int, group int, toqq int, pkgid int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("group", fmt.Sprintf("%v", group))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("pkgid", fmt.Sprintf("%v", pkgid))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/sendfreepackage", addr), bodystr)
	if data == nil {
		logApi.Error("发送免费礼物出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//取QQ在线状态
//提交请求参数:[必须]logonqq 指定框架QQ [必须]qq 指定欲查询QQ
func getqqonlinestate(logonqq int, qq int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("qq", fmt.Sprintf("%v", qq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getqqonlinestate", addr), bodystr)
	if data == nil {
		logApi.Error("取QQ在线状态出错")
		return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

//分享音乐
//提交请求参数:
//[必须]logonqq 指定框架QQ
//[可选]totype 指定分享对象类型(0私聊 1群聊  默认0)
//[必须]to 指定分享对象(分享的群或分享的好友QQ)
//[必须]musicname 指定歌曲名(存在特殊字符请使用URL编码)
//[必须]singername 指定歌手名(存在特殊字符请使用URL编码)
//[必须]jumpurl 指定跳转地址(点击音乐json后跳转的地址)(存在特殊字符请使用URL编码)
//[必须]wrapperurl 指定封面地址(音乐的封面图片地址)(存在特殊字符请使用URL编码)
//[必须]fileurl 指定文件地址(音乐源文件地址，如https://xxx.com/xxx.mp3)(存在特殊字符请使用URL编码)
//[可选]apptype 指定应用类型(0QQ音乐 1虾米音乐 2酷我音乐 3酷狗音乐 4网抑云音乐  默认0)
func sharemusic(logonqq int, totype int, to int, musicname string, singername string,
	jumpurl string, wrapperurl string, fileurl string, apptype int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("totype", fmt.Sprintf("%v", totype))
	r.Form.Add("to", fmt.Sprintf("%v", to))
	r.Form.Add("musicname", musicname)
	r.Form.Add("singername", singername)
	r.Form.Add("jumpurl", jumpurl)
	r.Form.Add("wrapperurl", wrapperurl)
	r.Form.Add("fileurl", fileurl)
	if apptype == 0 || apptype == 1 || apptype == 2 || apptype == 3 || apptype == 4 {
		r.Form.Add("apptype", fmt.Sprintf("%v", apptype))
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/sharemusic", addr), bodystr)
	if data == nil {
		logApi.Error("分享音乐出错")
		//return 0
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
}

//获取skey
//提交请求参数:
//[必须]logonqq 指定框架QQ
//[必须]domain 指定域(tenpay.com;openmobile.qq.com;docs.qq.com;connect.qq.com;qzone.qq.com;vip.qq.com;gamecenter.qq.com;qun.qq.com;game.qq.com;qqweb.qq.com;ti.qq.com;office.qq.com;mail.qq.com;mma.qq.com)
//测试接口未实现
func getpskey(logonqq int, domain string) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	r.Form.Add("domain", domain)
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getpskey", addr), bodystr)
	if data == nil {
		logApi.Error("获取skey出错")
		return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

//获取skey
//提交请求参数:[必须]logonqq 指定框架QQ
//测试接口未实现
func getskey(logonqq int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getskey", addr), bodystr)
	if data == nil {
		logApi.Error("获取skey出错")
		return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

//获取clientkey
//测试接口未实现
func getclientkey(logonqq int) string {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("logonqq", fmt.Sprintf("%v", logonqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	data := postFormData(fmt.Sprintf("http://%v/getclientkey", addr), bodystr)
	if data == nil {
		logApi.Error("获取clientkey出错")
		return ""
	}

	//fmt.Println(fmt.Sprintf("%s", data))
	var ret = new(Ret)
	_ = json.Unmarshal(data, &ret)
	logApi.Debug(ret.Ret)
	return ret.Ret
}

func postFormData(url string, bodystr string) []byte {
	logApi.Debugf("post form data: %v?%v", url, bodystr)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(bodystr))
	if err != nil {
		logApi.Debugf("POST请求出错了 %v", err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logApi.Debugf("ioutil.ReadAll出错了 %v", err.Error())
		return nil
	}
	_ = resp.Body.Close()
	logApi.Debugf("post返回值:%s", body)
	return body
}

func get(url string) []byte {
	logApi.Debugf("get url %v", url)
	resp, err := http.Get(url)
	if err != nil {
		logApi.Debugf("get error %v", err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logApi.Debugf("get read error %v", err.Error())
		return nil
	}
	_ = resp.Body.Close()
	logApi.Debugf("get返回值:%s", body)
	return body
}
