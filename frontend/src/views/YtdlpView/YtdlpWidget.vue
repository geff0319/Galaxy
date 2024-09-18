<script setup lang="ts">
import {formatResolution, useAppSettingsStore, useYtdlpStore} from '@/stores'
import { LoadingOutlined, DownloadOutlined } from '@ant-design/icons-vue';
import { ref } from 'vue';
import {WindowHide, WindowSetAlwaysOnTop, WindowSetPosition, WindowSetSize} from "@wails/runtime";
import {ShowToolWindow} from "@/bridge";
import * as Stores from "@/stores";
import {message} from "ant-design-vue";
import {appGetVideoMeta} from "@/bridge/ytdlp";
import router from "@/router";

const ytdlpStore = useYtdlpStore()
const appStore = Stores.useAppStore()


const closeYtWidget = ()=>{
  appStore.widgetsEnable = false
  appStore.widgetsType = ''
  WindowSetSize(800, 540)
  WindowSetAlwaysOnTop(false)
  WindowHide()
}

const parseUrl = async ()=>{
  // ytdlpStore.videoTitle = ''
  // ytdlpStore.videoBestFormats = ''
  // const youtubeRegex = /^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$/;
  // if(!youtubeRegex.test(ytdlpStore.downloadUrl)){
  //   message.error("解析失败,请检查链接",500)
  //   return
  // }
  // const res = await appGetVideoMeta(ytdlpStore.downloadUrl)
  // ytdlpStore.videoTitle= res['title']
  // ytdlpStore.videoBestFormats = res['requested_formats'][0]['width'] + '*' + res['requested_formats'][0]['height']
}
const download = async () =>{
  try {
    router.push({path: '/ytdlp'})
    appStore.widgetsEnable = false
    appStore.widgetsType = ''
    await WindowSetSize(1000, 640)
    WindowSetPosition(650,250)
    WindowSetAlwaysOnTop(false)
    await ShowToolWindow()
    const res = await ytdlpStore.downloadYoutube(true,false)
    message.info("开始下载："+res)
  }catch (error:any){
    message.error(error)
  }
}
</script>

<template>
  <div class="ytcontainer">
    <div style="--wails-draggable: drag;height: 40px">
      <div class="action">
        <h3 style="text-align: center;">解析链接</h3>
        <Button
            @click.stop="closeYtWidget"
            :class="{ 'hover-red': true }"
            type="text"
        >
          <Icon icon="close" />
        </Button>
      </div>
    </div>
    <a-flex gap="small" class="spaced-div" vertical>
      <div style="width: 500px;padding-top: 2%" :title=ytdlpStore.resProcess.info.title>
        <a-typography-text
            style="font-size:18px; "
            strong
            ellipsis
            :level="2"
            :content=ytdlpStore.resProcess.info.title
            v-if="!ytdlpStore.parseing"
        />
        <a-skeleton-button :block="true" :active="true" size="small" shape="round" v-else />
      </div>
      <div>
        <a-flex gap="small">
          <a-input :pressEnter="parseUrl()" disabled v-model:value="ytdlpStore.videoUrl">
            <template #addonAfter>
              <LoadingOutlined v-if='ytdlpStore.parseing'/>
              <a-typography-text v-else style="max-width: 90px" :content="formatResolution(ytdlpStore.resProcess.info.resolution)" ellipsis :title="formatResolution(ytdlpStore.resProcess.info.resolution)"/>
            </template>
          </a-input>
          <a-button type="primary" size="middle" :disabled="ytdlpStore.parseing||ytdlpStore.videoUrl.length===0||formatResolution(ytdlpStore.resProcess.info.resolution)==='未知'" @click="download">
            <template #icon>
              <DownloadOutlined />
            </template>
          </a-button>
        </a-flex>
      </div>
    </a-flex>
  </div>

</template>

<style lang="less" scoped>
.ytcontainer {
  height: 100vh; /* 设置父容器高度为视口的高度 */
  background-color: rgba(255, 255, 255, 1); /* 白色背景，透明度为 0.5 */
  //background-color: transparent; /* 完全透明背景 */
}
.spaced-div {
  //width: 95%;
  margin-left: 13px;  /* 左间隔 */
  margin-right: 13px; /* 右间隔 */
}
.action {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  //justify-content: center;
  font-size: 14px;
  &-btn {
    width: 32px;
    height: 32px;
    line-height: 32px;
    text-align: center;
    border-radius: 4px;
    &:hover {
      background-color: var(--hover-bg-color);
    }
  }
}
.hover-red:hover {
  background: rgba(255, 0, 0, 0.6);
}
</style>