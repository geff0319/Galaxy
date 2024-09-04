<script setup lang="ts">
import { CaretRightOutlined ,PauseOutlined, BugOutlined} from '@ant-design/icons-vue';
import {ref} from "vue";

import {message} from "ant-design-vue";
import {useWsClientStore} from "@/stores/wsClient";

const appWsClient = useWsClientStore()
const loading = ref(false)
const connectWs= async ()=>{
  loading.value=true
  try{
    await appWsClient.connectWs()
  }catch (error){
    loading.value=false
    return
  }
  loading.value=false
}
const disConnectWs= async ()=>{
  loading.value=true
  try {
    await appWsClient.disConnectWs()
  }catch (e) {
    loading.value=false
  }
  loading.value=false
}
const initWs = () => {
  loading.value=true
  try{
    appWsClient.connectionStatus()
  }catch (err) {
    loading.value=false
  }
  loading.value=false
}
initWs()
</script>

<template>
  <div class="settings">
    <div class="settings-item">
      <div class="title">
        WebSocket
      </div>
      <div class="input-wrapper">
        <h5 class="label">自动连接: </h5>
        <div class="input-container">
          <a-switch v-model:checked="appWsClient.ws.autoConnect" />
        </div>
      </div>
      <div class="input-wrapper">
        <h5 class="label">客户端ID: </h5>
        <div class="input-container">
          <a-tooltip :title="appWsClient.ws.id">
            <a-input style="margin-right: 10px" :disabled="appWsClient.isConnected||loading" v-model:value="appWsClient.ws.id" :bordered="false" placeholder="请输入id" />
          </a-tooltip>
        </div>
      </div>
      <div class="input-wrapper">
        <h5 class="label">服务器地址: </h5>
        <div class="input-container">
          <a-tooltip :title="appWsClient.ws.domain">
            <a-input style="margin-right: 10px" :disabled="appWsClient.isConnected||loading" v-model:value="appWsClient.ws.domain" placeholder="ws://" />
          </a-tooltip>
          <a-button v-if="!appWsClient.isConnected" :loading="loading" style="width: 60px" type="primary" size="middle" @click="connectWs">
            <template #icon>
              <CaretRightOutlined />
            </template>
          </a-button>
          <a-button v-else :loading="loading" style="width: 60px" type="primary" size="middle" @click="disConnectWs" danger>
            <template #icon>
              <PauseOutlined />
            </template>
          </a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.settings {
  &-item {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 8px 16px;
    .title {
      align-self: flex-start;
      font-size: 18px;
      font-weight: bold;
      padding: 8px 0 16px 0;
      .tips {
        font-weight: normal;
        font-size: 12px;
      }
    }
  }
}
.input-wrapper{
  display: flex;
  align-items: center; /* 垂直居中 */
  .label {
    flex: 0 0 auto; /* 不随内容变化而变化 */
    margin-right: 10px; /* 右侧留出一些空间 */
  }
  .input-container{
    display: flex;
    flex: 1/* 填充剩余空间 */
  }
}
</style>
