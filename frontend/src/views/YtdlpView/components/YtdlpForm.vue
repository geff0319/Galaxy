<script setup lang="ts">
import {useYtdlpStore} from "@/stores";
import {inject, ref} from "vue";
import {message} from "ant-design-vue";
const handleCancel = inject('cancel') as any


const ytdlpStore = useYtdlpStore()
const val = ref<string>("")

const handleSave = async () => {
  try {
    ytdlpStore.determineUrl(val.value)
    if(ytdlpStore.downloadUrl.length===0){
      message.error("未解析有效链接")
      handleCancel()
    }
    await ytdlpStore.downloadYoutube(false,false)
    handleCancel()
  } catch (error: any) {
    message.error("添加任务失败"+error)
  }
}

</script>

<template>
  <div class="header" style="--wails-draggable: drag">
    <div class="header-title">添加任务</div>
  </div>

  <div class="form">
    <div class="form-item">
      <div class="name">链接 *</div>
      <Input v-model="val" auto-size autofocus class="flex-1 ml-8" />
    </div>
  </div>

  <div class="form-action">
    <Button @click="handleCancel">取消</Button>
    <Button @click="handleSave" :disable="val.length===0" type="primary">确定</Button>
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
  padding-right: 8px;
  overflow-y: auto;
  max-height: 60vh;
}
</style>