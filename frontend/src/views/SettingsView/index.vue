<script setup lang="ts">
import {ref} from 'vue'
import {useI18n} from 'vue-i18n'

import {useAppStore, useYtdlpSettingsStore, useYtdlpStore} from '@/stores'

import KernelView from '@/views/KernelView/index.vue'
import GeneralSetting from './components/GeneralSetting.vue'
import TranslateSetting from './components/TranslateSetting.vue'
import YtdlpSetting from "./components/YtdlpSetting.vue";
import WsSetting from "@/views/SettingsView/components/WsSetting.vue";


const settings = [
  { key: 'general', tab: 'settings.general' },
  { key: 'translate', tab: 'settings.translate.translate' },
  { key: 'ytdlp', tab: 'YT-DLP' },
  { key: 'websocket', tab: 'WebSocket' },
  // { key: 'kernel', tab: 'router.kernel' }
]

const activeKey = ref(settings[0].key)

const { t } = useI18n()
const appStore = useAppStore()
useYtdlpSettingsStore().setupYtdlpSettings()
</script>

<template>
  <Tabs v-model:active-key="activeKey" :items="settings" height="100%">
    <template #general>
      <GeneralSetting />
    </template>

    <template #translate>
      <TranslateSetting />
    </template>

    <template #ytdlp>
      <YtdlpSetting />
    </template>

    <template #websocket>
      <WsSetting />
    </template>

<!--    <template #kernel>-->
<!--      <KernelView />-->
<!--    </template>-->

    <template #extra>
      <Button @click="appStore.showAbout = true" type="text">
        {{ t('router.about') }}
      </Button>
    </template>
  </Tabs>
</template>

<style lang="less" scoped></style>
