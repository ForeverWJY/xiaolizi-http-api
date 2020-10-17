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
}

func reply(reqMsg *ReceiveMessage, msg string) {
	switch reqMsg.Type {
	case "PrivateMsg":
		sendPrivateMsg(LoginQQ, reqMsg.FromQQ.UIN, msg)
		break
	case "GroupMsg":
		sendGroupMsg(LoginQQ, reqMsg.FromGroup.GIN, msg)
		break
	}
}

type Ret struct {
	Ret string `json:"ret"`
}

func getLoginQQ() int {
	qq := get(fmt.Sprintf("http://%v/getlogonqq", *addr))
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
func sendPrivateMsg(fromqq int, toqq int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	r.Form.Add("text", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendprivatemsg", *addr), bodystr)
}

//发送群消息
func sendGroupMsg(fromqq int, togroup int, text string) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("togroup", fmt.Sprintf("%v", togroup))
	r.Form.Add("text", text)
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/sendgroupmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendgrouptempmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/addfriend", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/addgroup", *addr), bodystr)
}

//删除好友
func deletefriend(fromqq int, toqq int) {
	var r http.Request
	_ = r.ParseForm()
	r.Form.Add("fromqq", fmt.Sprintf("%v", fromqq))
	r.Form.Add("toqq", fmt.Sprintf("%v", toqq))
	bodystr := strings.TrimSpace(r.Form.Encode())
	postFormData(fmt.Sprintf("http://%v/deletefriend", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/setfriendignmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/setfriendcare", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendprivatexmlmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendgroupxmlmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendprivatejsonmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendgroupjsonlmsg", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendprivatepic", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendgrouppic", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendprivateaudio", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/sendgroupaudio", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/uploadfacepic", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/uploadgroupfacepic", *addr), bodystr)
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
	postFormData(fmt.Sprintf("http://%v/setgroupcard", *addr), bodystr)
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
	data := postFormData(fmt.Sprintf("http://%v/getnickname", *addr), bodystr)
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
	data := postFormData(fmt.Sprintf("http://%v/getgroupnamefromcache", *addr), bodystr)
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
	data := postFormData(fmt.Sprintf("http://%v/getfriendlist", *addr), bodystr)
	if data == nil {
		logApi.Debug("取好友列表出错")
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
	data := postFormData(fmt.Sprintf("http://%v/getgrouplist", *addr), bodystr)
	if data == nil {
		logApi.Debug("取群列表出错")
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

func postFormData(url string, bodystr string) []byte {
	logApi.Debugf("send %v", bodystr)
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
