<script setup lang="ts">
import { ref ,reactive} from 'vue'
import { useI18n } from 'vue-i18n'

import { useMessage } from '@/hooks'
import { useAppSettingsStore, useEnvStore } from '@/stores'


const { t } = useI18n()
const { message } = useMessage()
const appSettings = useAppSettingsStore()
const envStore = useEnvStore()
const hidestr = ref('***********************')
const showOptions =
    reactive<Record<string, boolean>>({
      "tencentTanslateSecretId": false,
      "tencentTanslateSecretKey": false
    })

const changeShow = (key:string) => {
  showOptions[key] = !showOptions[key];
  console.log( showOptions[key] )
}

</script>

<template>
  <div class="settings">
    <div class="settings-item">
      <div class="title">{{ t('settings.translate.tencentTanslate') }}</div>
      <div class="input-wrapper">
        <h5 class="label">{{ t('settings.translate.tencentTanslateSecretId')+'：'}}</h5>
        <Input
            v-if="showOptions.tencentTanslateSecretId"
            class="input-container"
            v-model.lazy="appSettings.app.translate.tencentTanslateSecretId"
            placeholder="请输入SecretId"
            type="text"
            :editable="true"
        >
          <template #extra>
            <Button @click="changeShow('tencentTanslateSecretId')" type="text">
              <Icon icon="show"/>
            </Button>
          </template>
        </Input>
        <Input
            v-else
            class="input-container"
            v-model="hidestr"
            type="text"
            :disabled="true"
            editable
        >
          <template #extra>
            <Button @click="changeShow('tencentTanslateSecretId')" type="text">
              <Icon icon="hide"/>
            </Button>
          </template>
        </Input>
      </div>
      <div class="input-wrapper">
        <h5 class="label">{{ t('settings.translate.tencentTanslateSecretKey')+'：' }}</h5>
        <Input
            v-if="showOptions.tencentTanslateSecretKey"
            class="input-container"
            v-model.lazy="appSettings.app.translate.tencentTanslateSecretKey"
            placeholder="请输入SecretKey"
            type="text"
            :editable="true"
        >
          <template #extra>
            <Button @click="changeShow('tencentTanslateSecretKey')" type="text">
              <Icon icon="show"/>
            </Button>
          </template>
        </Input>
        <Input
            v-else
            class="input-container"
            v-model.lazy="hidestr"
            type="text"
            :disabled="true"
            :editable="true"
        >
          <template #extra>
            <Button @click="changeShow('tencentTanslateSecretKey')" type="text">
              <Icon icon="hide"/>
            </Button>
          </template>
        </Input>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
.list-item {
  display: flex;
  padding: 0 0 0 8px;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  margin: 2px 0;
  border: #2e2e2e;
}

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
    flex: 1; /* 填充剩余空间 */
  }
}
</style>
