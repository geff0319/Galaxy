<script setup lang="ts">
import FlipClock from "./FlipClockView/FlipClock.vue";
import YtdlpWidget from "@/views/YtdlpView/YtdlpWidget.vue";
import {
  WindowSetSize,
  WindowSetPosition,
  WindowSetAlwaysOnTop,
} from "@wails/runtime";
import * as Stores from "@/stores";
import type {Menu} from "@/stores";
import {ShowToolWindow} from "@/bridge";

const appStore = Stores.useAppStore()
const closeWidgets =async () => {
  appStore.widgetsEnable = false
  appStore.widgetsType = ''
  await WindowSetSize(1000, 640)
  WindowSetPosition(650,250)
  WindowSetAlwaysOnTop(false)
  await ShowToolWindow()
}
const menus: Menu[] = [
  {
    label: '回到主界面',
    handler: closeWidgets
  }
]

</script>

<template>
    <div style="--wails-draggable: drag" class="container" v-if="appStore.widgetsType==='clock'"  v-menu="menus">
  <!--    后续可以试试<Component></Component>组件-->
      <FlipClock />
  <!--    <Input v-model="n"/>-->
    </div>
    <YtdlpWidget v-else-if="appStore.widgetsType==='ytdlp'" />
</template>

<style lang="less" scoped>
.container {
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  height: 100vh; /* 设置父容器高度为视口的高度 */
  background-color: rgba(255, 255, 255, 0.85); /* 白色背景，透明度为 0.5 */
  //background-color: transparent; /* 完全透明背景 */
}
</style>
