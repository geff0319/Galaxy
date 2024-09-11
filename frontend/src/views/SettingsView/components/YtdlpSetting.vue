<script setup lang="ts">
import { FolderOpenOutlined } from '@ant-design/icons-vue';
import {ref} from "vue";
import { OpenDirectoryDialog } from '@/bridge'
import {useYtdlpSettingsStore} from "@/stores";
import {message} from "ant-design-vue";
import YtDlp from "@/views/KernelView/YtDlp.vue";
import Ffmpeg from "@/views/KernelView/Ffmpeg.vue";


const ytdlpSettingsStore = useYtdlpSettingsStore()
// const numVal = ref<string>(ytdlpSettingsStore.ytdlpConfig.queueSize)
const handelOpenFileDialog =async ()=>{
  try {
    const folder = await  OpenDirectoryDialog()
    if (folder.length!==0) {
      ytdlpSettingsStore.ytdlpConfig.downloadPath = folder
    }
  }catch (error :any){
    message.error("选择文件夹失败")
  }
}
</script>

<template>
  <div class="settings">
    <div class="settings-item">
      <div class="title">配置项</div>
      <div class="input-wrapper">
        <h5 class="label">下载目录: </h5>
        <div class="input-container">
          <a-tooltip :title="ytdlpSettingsStore.ytdlpConfig.downloadPath">
            <a-input style="margin-right: 10px"  disabled v-model:value="ytdlpSettingsStore.ytdlpConfig.downloadPath" placeholder="请选择视频下载目录" />
          </a-tooltip>
          <a-button style="width: 40px" type="primary" size="middle" @click="handelOpenFileDialog">
            <template #icon>
              <FolderOpenOutlined />
            </template>
          </a-button>
        </div>
      </div>
      <div class="input-wrapper">
        <h5 class="label">并发数量(重启生效): </h5>
        <div class="input-container">
          <a-input-number v-model:value="ytdlpSettingsStore.ytdlpConfig.queueSize" :min="1" :max="10" />
        </div>
      </div>
      <div class="title">Cookie</div>
      <div class="input-wrapper">
        <h5 class="label">bilibili: </h5>
        <div class="input-container">
          <a-input-password :bordered="false" v-model:value="ytdlpSettingsStore.ytdlpConfig.cookies.bilibili" placeholder="input cookie" />
        </div>
      </div>
      <YtDlp />
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
