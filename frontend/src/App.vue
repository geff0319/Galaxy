<script setup lang="ts">
import { ref } from 'vue'

import * as Stores from '@/stores'
import {
  EventsOn, HideToolWindow,
  ShowToolWindow,
  ClipboardGetText,
  WindowHide,
  WindowSetAlwaysOnTop,
  WindowSetPosition,
  WindowSetSize,
  WindowShow
} from '@/bridge'
import { exitApp, ignoredError, sampleID, sleep } from '@/utils'
import { useMessage, usePicker, useConfirm, usePrompt, useAlert } from '@/hooks'
import {message as antdMessage} from "ant-design-vue";

import AboutView from '@/views/AboutView.vue'
import SplashView from '@/views/SplashView.vue'
import CommandView from './views/CommandView.vue'
import { NavigationBar, MainPage, TitleBar } from '@/components'
import WidgetsView from "@/views/WidgetsView.vue";

const loading = ref(true)

const envStore = Stores.useEnvStore()
const appStore = Stores.useAppStore()
const pluginsStore = Stores.usePluginsStore()
const profilesStore = Stores.useProfilesStore()
const rulesetsStore = Stores.useRulesetsStore()
const appSettings = Stores.useAppSettingsStore()
const kernelApiStore = Stores.useKernelApiStore()
const subscribesStore = Stores.useSubscribesStore()
const scheduledTasksStore = Stores.useScheduledTasksStore()
const ytdlpStore = Stores.useYtdlpStore()
const wsClientStore = Stores.useWsClientStore()

const { message } = useMessage()
const { picker } = usePicker()
const { confirm } = useConfirm()
const { prompt } = usePrompt()
const { alert } = useAlert()

window.Plugins.message = message
window.Plugins.picker = picker
window.Plugins.confirm = confirm
window.Plugins.prompt = prompt
window.Plugins.alert = alert


EventsOn('launchArgs', async (args: string[]) => {
  console.log('launchArgs', args)
  const url = new URL(args[0])
  if (url.pathname === '//install-config/') {
    const _url = url.searchParams.get('url')
    const _name = url.searchParams.get('name') || sampleID()

    if (!_url) {
      message.error('URL missing')
      return
    }

    try {
      await subscribesStore.importSubscribe(_name, _url)
      message.success('common.success')
    } catch (error: any) {
      message.error(error)
    }
  }
})

EventsOn('beforeClose', async () => {
  if (appSettings.app.exitOnClose) {
    exitApp()
  } else {
    WindowHide()
  }
})

EventsOn('quitApp', () => exitApp())

window.addEventListener('beforeunload', scheduledTasksStore.removeScheduledTasks)

EventsOn('keyChangeWidgets', async ()=>{
  try {
    const text = await ClipboardGetText()
    ytdlpStore.determineUrl(text)
    console.log(ytdlpStore.downloadUrl)
    if(appStore.widgetsEnable){
      if(ytdlpStore.videoUrl === text){
        return
      }
      if(ytdlpStore.parseing){
        antdMessage.info('链接解析中..',1)
        return
      }
      ytdlpStore.videoUrl = text
      await ytdlpStore.getVideoMeta()
    }else if(ytdlpStore.downloadUrl!=='') {
      appStore.widgetsEnable = true
      appStore.widgetsType = 'ytdlp'
      // HideToolWindow()
      WindowSetAlwaysOnTop(true)
      await WindowSetSize(600, 150)
      WindowSetPosition(750,450)
      WindowShow()
      if(ytdlpStore.videoUrl === text){
        return
      }
      if(ytdlpStore.parseing){
        antdMessage.info('链接解析中..',1)
        return
      }
      ytdlpStore.videoUrl = text
      // ytdlpStore.determineUrl()
      await ytdlpStore.getVideoMeta()
    }
  }catch (error:any){
    appStore.widgetsEnable = true
    appStore.widgetsType = 'ytdlp'
    // HideToolWindow()
    WindowSetAlwaysOnTop(true)
    await WindowSetSize(600, 150)
    WindowSetPosition(550,350)
    WindowShow()
    if(error==='The operation completed successfully.'){
      antdMessage.info('粘贴板没有获取到视频链接',1)
    }else {
      antdMessage.error('解析链接出错：'+error,2)
    }
    return
  }
})

EventsOn('notify',(type:string,data:string)=>{
  switch (type) {
    case "info":
      antdMessage.info(data);
      break;
    case "error":
      antdMessage.error(data)
      break;
    case "warn":
      antdMessage.warn(data);
      break;
    default:
      antdMessage.info(data)
      break;
  }
})

appSettings.setupAppSettings().then(async () => {
  await Promise.all([
    ignoredError(envStore.setupEnv),
    ignoredError(profilesStore.setupProfiles),
    ignoredError(subscribesStore.setupSubscribes),
    ignoredError(rulesetsStore.setupRulesets),
    ignoredError(pluginsStore.setupPlugins),
    ignoredError(scheduledTasksStore.setupScheduledTasks),
    ignoredError(wsClientStore.setupWsSettings)
  ])
  await sleep(1000)

  //ws
  if(wsClientStore.ws.autoConnect){
    try {
      wsClientStore.connectWs()
    }catch (error){
      message.error("ws连接失败：" + error)
    }
  }


  loading.value = false
  kernelApiStore.updateKernelStatus().then(async (running) => {
    kernelApiStore.statusLoading = false
    try {
      if (running) {
        await kernelApiStore.refreshConfig()
        await kernelApiStore.refreshProviderProxies()
        await envStore.updateSystemProxyStatus()
      } else if (appSettings.app.autoStartKernel) {
        await kernelApiStore.startKernel()
      }
    } catch (error: any) {
      message.error(error)
    }
  })

  try {
    await pluginsStore.onStartupTrigger()
  } catch (error: any) {
    message.error(error)
  }
})
</script>

<template>
  <SplashView v-if="loading" />
  <template v-else-if="appStore.widgetsEnable"><widgetsView /></template>
<!--  <widgetsView v-else-if="appStore.widgetsEnable"/>-->
  <template v-else-if="!appStore.widgetsEnable">
    <TitleBar />
    <div class="main">
      <NavigationBar />
      <MainPage />
    </div>
  </template>

  <Modal
    v-model:open="appStore.showAbout"
    :cancel="false"
    :submit="false"
    mask-closable
    min-width="50"
  >
    <AboutView />
  </Modal>

  <Menu
    v-model="appStore.menuShow"
    :position="appStore.menuPosition"
    :menu-list="appStore.menuList"
  />

  <Tips
    v-model="appStore.tipsShow"
    :position="appStore.tipsPosition"
    :message="appStore.tipsMessage"
  />

  <CommandView />
</template>

<style scoped>
.main {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  padding: 8px;
}
</style>
