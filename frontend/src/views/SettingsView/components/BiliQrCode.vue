<script setup lang="ts">
import {inject, ref} from "vue";
import {HttpGet} from "@/bridge";
import {message} from "ant-design-vue";
import { sleep } from '@/utils'
import {useYtdlpSettingsStore} from "@/stores";

const ytdlpSettingsStore = useYtdlpSettingsStore()
const checkStatusOut = ref(false)
const qrUrl = ref("")
const statusCode = ref(1)
const handleCancel = inject('cancel') as any
const getQqrUrl = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
const pollUrl = "https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key="

const qrHandleCancel = ()=>{
  checkStatusOut.value = true
  console.log("二维码退出")
  handleCancel()
}

const generate =async ()=>{
  try{
    const { body:body1 } = await HttpGet<Record<string, any>>(getQqrUrl)
    if(body1.code!==0){
      message.error(body1.message)
      statusCode.value = 86038
      return
    }
    qrUrl.value = body1.data.url
    return body1.data.qrcode_key
  }catch (error: any) {
    message.error("获取二维码失败")
    statusCode.value = 86038
  }
}
const formattingStatus = () =>{
  switch(statusCode.value){
    case 86101:
      return undefined;
    case 86090:
      return "scanned"
    case 86038:
      return "expired"
    case 0:
      return "success"
    case 1:
      return "loading"
    case -1:
      return "out"
  }
}
const saveBilibiliCookie = (str:string) =>{
  const parsedUrl = new URL(str);
  const sessdata: string | null = parsedUrl.searchParams.get("SESSDATA")
  if (sessdata !== null) {
    ytdlpSettingsStore.ytdlpConfig.cookies.bilibili = sessdata; // 这里没有错误
  } else {
    ytdlpSettingsStore.ytdlpConfig.cookies.bilibili= undefined; // 或者不赋值
  }
}
const checkLoginStatus = async (key: string)=>{
  if(key === ""){
    message.error("获取二维码失败")
    statusCode.value = 86038
    return
  }
  while (true){
    if(checkStatusOut.value){
      return
    }
    const { body:body1 } = await HttpGet<Record<string, any>>(pollUrl + key)
    if(body1.code!==0){
      message.error(body1.message)
      statusCode.value = 86038
      return
    }
    statusCode.value = body1.data.code
    if(statusCode.value === 0){
      message.info("登录成功")
      saveBilibiliCookie(body1.data.url)
      handleCancel()
      return
    }else if(statusCode.value === 86038){
      c
    }
    await sleep(1000)
  }
}
const init = async ()=>{
  const key = await generate()
  await checkLoginStatus(key)
}

init()
</script>

<template >
  <div class="header" style="--wails-draggable: drag">
    <div class="header-title">请扫描登录二维码</div>
  </div>

  <div class="form">
    <a-qrcode :value="qrUrl" :status="formattingStatus()"/>
  </div>

  <div class="form-action">
    <Button @click="qrHandleCancel">取消</Button>
<!--    <Button @click="handleSave" :disable="val.length===0" type="primary">确定</Button>-->
  </div>
</template>


<style lang="less" scoped>
.header {
  display: flex;
  align-items: center;
  margin-top: 8px;
  &-title {
    font-size: 20px;
    font-weight: bold;
    margin: 8px 0 16px 0;
  }
}
.form {
  display: flex;
  align-items: center;
  justify-content:center;
}
</style>