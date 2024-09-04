<script setup lang="ts">
import {computed, ref} from "vue";
import {useI18n} from "vue-i18n";
import {getKernelFileName, getYtDlpFileName, KernelWorkDirectory, YtDlpWorkDirectory} from "@/constant";
import {Download, Exec, GetEnv, HttpGet, Makedir, Movefile, Removefile, UnzipGZFile, UnzipZIPFile} from "@/bridge";
import { useMessage } from '@/hooks'
import {useEnvStore} from "@/stores";
import {ignoredError} from "@/utils";
const { t } = useI18n()
const { message } = useMessage()

const releaseUrl = 'https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest'
const localVersion = ref('')
const remoteVersion = ref('')
const localVersionLoading = ref(false)
const remoteVersionLoading = ref(false)
const downloadLoading = ref(false)
const downloadSuccessful = ref(false)

const needUpdate = computed(() => remoteVersion.value && localVersion.value !== remoteVersion.value)
const updateLocalVersion = async (showTips = false) => {
  localVersion.value = await getLocalVersion(showTips)
}
const getLocalVersion = async (showTips = false) => {
  localVersionLoading.value = true
  try {
    const fileName = await getYtDlpFileName()
    const kernelFilePath = YtDlpWorkDirectory + '/' + fileName
    const res = await Exec(kernelFilePath, ['--version'])
    return res.trim()
  } catch (error: any) {
    console.log(error)
    showTips && message.error(error)
  } finally {
    localVersionLoading.value = false
  }
  return ''
}
const updateRemoteVersion = async () => {
  remoteVersion.value = await getRemoteVersion()
}
const getRemoteVersion = async () => {
  remoteVersionLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(releaseUrl)
    const { tag_name } = body
    return tag_name as string
  } catch (error: any) {
    console.log(error)
    message.error(error)
  } finally {
    remoteVersionLoading.value = false
  }
  return ''
}
const downloadCore = async () => {
  downloadLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(releaseUrl)
    const { os, arch } = await GetEnv()

    const { assets, message: msg } = body
    if (msg) throw msg

    const envStore = useEnvStore()
    // const amd64Compatible = arch === 'amd64' && envStore.env.x64Level < 3 ? '-compatible' : ''
    // const suffix = { windows: '.zip', linux: '.gz', darwin: '.gz' }[os]
    const assetName = `yt-dlp.exe`

    const asset = assets.find((v: any) => v.name === assetName)
    if (!asset) throw 'Asset Not Found:' + assetName

    const tmp = `data/yt-dlp.exe` // data/core-alpha.zip or data/core-alpha.gz

    await Makedir('data/yt-dlp')

    const { id } = message.info('Downloading...', 10 * 60 * 1_000)

    await Download(asset.browser_download_url, tmp, undefined, (progress, total) => {
      message.update(id, 'Downloading...' + ((progress / total) * 100).toFixed(2) + '%')
    }).catch((err) => {
      message.destroy(id)
      throw err
    })

    message.destroy(id)

    const fileName = await getYtDlpFileName() // yt-dlp.exe

    const ytdlpFilePath = YtDlpWorkDirectory + '/' + fileName

    await ignoredError(Movefile, ytdlpFilePath, ytdlpFilePath + '.bak')

    // if (suffix === '.zip') {
    //   await UnzipZIPFile(tmp, 'data')
    // } else {
    //   await UnzipGZFile(tmp, alphaKernelFilePath)
    // }

    await Movefile(tmp, ytdlpFilePath)

    await Removefile(tmp)

    downloadSuccessful.value = true

    message.success('Download Successful')
  } catch (error: any) {
    console.log(error)
    message.error(error)
    downloadSuccessful.value = false
  }

  downloadLoading.value = false

  updateLocalVersion()
}

const ffUrl = 'https://api.github.com/repos/yt-dlp/FFmpeg-Builds/releases/tags/latest'
const localFfVersion = ref('')
const localFfVersionLoading = ref(false)
const downloadFfLoading = ref(false)
const downloadFfSuccessful = ref(false)

// const needFfUpdate = computed(() => remoteFfVersion.value && localFfVersion.value !== remoteFfVersion.value)
const updateFfLocalVersion = async (showTips = false) => {
  localFfVersion.value = await getLocalFfVersion(showTips)
}
const getLocalFfVersion = async (showTips = false) => {
  localFfVersionLoading.value = true
  try {
    const fileName = "ffmpeg.exe"
    const kernelFilePath = YtDlpWorkDirectory + '/' + fileName
    const res = await Exec(kernelFilePath, ['-version'])
    return res.trim().match(/version\s([^\s]+)\sCopyright/)?.[1].trim() || ''
  } catch (error: any) {
    console.log(error)
    showTips && message.error(error)
  } finally {
    localFfVersionLoading.value = false
  }
  return ''
}
const downloadFfCore = async () => {
  downloadFfLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(ffUrl)
    const { os, arch } = await GetEnv()

    const { assets, message: msg } = body
    if (msg) throw msg

    const envStore = useEnvStore()
    const amd64Compatible = arch === 'amd64' && envStore.env.x64Level < 3 ? '-compatible' : ''
    const suffix = { windows: '.zip', linux: '.tar.xz', darwin: '.gz' }[os]
    const assetName = `ffmpeg-master-latest-win64-gpl${suffix}`
    const asset = assets.find((v: any) => v.name === assetName)
    if (!asset) throw 'Asset Not Found:' + assetName

    const tmp = `data/core-ffmpeg${suffix}` // data/core-alpha.zip or data/core-alpha.gz

    await Makedir('data/yt-dlp')

    const { id } = message.info('Downloading...', 10 * 60 * 1_000)

    await Download(asset.browser_download_url, tmp, undefined, (progress, total) => {
      message.update(id, 'Downloading...' + ((progress / total) * 100).toFixed(2) + '%')
    }).catch((err) => {
      message.destroy(id)
      throw err
    })

    message.destroy(id)

    // const fileName = await getKernelFileName() // mihomo-windows-amd64.exe
    // const alphaFileName = await getKernelFileName(true) // mihomo-windows-amd64-alpha.exe

    const ffmpegFilePath = YtDlpWorkDirectory + '/ffmpeg.exe'
    const ffplayFilePath = YtDlpWorkDirectory + '/ffplay.exe'
    const ffprobeFilePath = YtDlpWorkDirectory + '/ffprobe.exe'

    await ignoredError(Movefile, ffmpegFilePath, ffmpegFilePath + '.bak')
    await ignoredError(Movefile, ffplayFilePath, ffplayFilePath + '.bak')
    await ignoredError(Movefile, ffprobeFilePath, ffprobeFilePath + '.bak')

    if (suffix === '.zip') {
      await UnzipZIPFile(tmp, 'data')
    } else {
      await UnzipGZFile(tmp, 'data')
    }
    //
    await Movefile('data/ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe' , ffmpegFilePath)
    await Movefile('data/ffmpeg-master-latest-win64-gpl/bin/ffplay.exe' , ffplayFilePath)
    await Movefile('data/ffmpeg-master-latest-win64-gpl/bin/ffprobe.exe' , ffprobeFilePath)

    await Removefile(tmp)
    await Removefile('data/ffmpeg-master-latest-win64-gpl')

    downloadFfSuccessful.value = true

    message.success('Download Successful')
  } catch (error: any) {
    console.log(error)
    message.error(error)
    downloadFfSuccessful.value = false
  }

  downloadFfLoading.value = false

  updateFfLocalVersion(true)
}
const initVersion = async () => {
  getLocalVersion(false)
      .then((v) => {
        localVersion.value = v
      })
      .catch((error: any) => {
        console.log(error)
      })

  getRemoteVersion()
      .then((versions) => {
        remoteVersion.value = versions
      })
      .catch((error: any) => {
        console.log(error)
      })

  getLocalFfVersion(false)
      .then((v) => {
        localFfVersion.value = v
      })
      .catch((error: any) => {
        console.log(error)
      })
}

initVersion()
</script>

<template>
  <h3>YT-DLP</h3>
  <div style="display: flex;align-items: center;">
    <Tag @click="updateLocalVersion(true)" style="cursor: pointer">
      {{ t('kernel.local') }}
      :
      {{ localVersionLoading ? 'Loading' : localVersion || t('kernel.notFound') }}
    </Tag>
    <Tag @click="updateRemoteVersion()" style="cursor: pointer">
      {{ t('kernel.remote') }}
      :
      {{ remoteVersionLoading ? 'Loading' : remoteVersion }}
    </Tag>
    <Button
        v-show="!localVersionLoading && !remoteVersionLoading && needUpdate"
        @click="downloadCore"
        :loading="downloadLoading"
        size="small"
        type="primary"
    >
      {{ t('kernel.update') }} : {{ remoteVersion }}
    </Button>
  </div>
  <h3>FFmpeg(YT-DLP优化版)</h3>
  <div style="display: flex;align-items: center;">
    <Tag @click="updateFfLocalVersion(true)" style="cursor: pointer">
      {{ t('kernel.local') }}
      :
      {{ localFfVersionLoading ? 'Loading' : localFfVersion || t('kernel.notFound') }}
    </Tag>
    <Button
        v-show="!localFfVersionLoading"
        @click="downloadFfCore"
        :loading="downloadFfLoading"
        size="small"
        type="primary"
    >
      {{ t('kernel.update') }}
    </Button>
  </div>
</template>

<style lang="less" scoped>
</style>