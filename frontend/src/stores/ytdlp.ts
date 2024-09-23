import {defineStore} from "pinia";
import {ref, watch} from "vue";
import {
    appAll,
    appDbPersist,
    appDownloadYoutube,
    appDownloadYoutubeByKey,
    appGetVideoMeta, deleteProcess,
    UpdateYtDlpConfig
} from "@/bridge/ytdlp";
import {message} from "ant-design-vue";
import {Readfile, Writefile} from "@/bridge";
import {parse, stringify} from "yaml";
import {debounce, updateTrayMenus} from "@/utils";
import {type Menu, useAppStore} from "@/stores/app";

export type ProcessType = {
    id: string,
    url:string
    progress :{
        process_status:number
        percentage:string
        speed:number
        eta:number
    }
    info     :{
        url:string
        title :string
        thumbnail :string
        resolution :string
        size        :number
        vCodec      :string
        aCodec      :string
        extension   :string
        originalURL :string
        fileName    :string
        createdAt   :number
    }
    output   :{
        path          :string
        filename      :string
        savedFilePath :string
    }
    params   : string[]
}

export function formatSize(bytes: number): string {
    const threshold = 1024
    const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    let i = 0
    while (bytes >= threshold) {
        bytes /= threshold
        i = i + 1
    }

    return `${bytes.toFixed(i === 0 ? 0 : 2)} ${units[i]}`;
}
export const formatSpeedMiB = (val: number) =>
    `${(val / 1_048_576).toFixed(2)}MiB/s`

export const formatResolution=(val: string):string=>{
   if(val.includes('4320')){
       return '8K'
   }else if (val.includes('2160')||val.includes('3840')){
       return '4K'
   }else if (val.includes('1440')){
       return '2K'
   }else if (val.includes('1080')){
       return '1080P'
   }else if (val.includes('720')){
       return '720P'
   }else if (val.includes('480')){
       return '480P'
   }else {
       return '未知'
   }
}

export const useYtdlpStore = defineStore('ytdlp', () => {
    const process = ref<ProcessType[]>([])
    const resProcess = ref<ProcessType>({
        id: "",
        url:"",
        output: {filename: "", path: "", savedFilePath: ""},
        params: [],
        progress: {eta: 0, percentage: "", process_status: 0, speed: 0},
        info: {
            aCodec: "",
            createdAt: 0,
            extension: "",
            fileName: "",
            originalURL: "",
            resolution: "未知",
            size: 0,
            thumbnail: "",
            title: "未知",
            url: "",
            vCodec: ""
        }
    })
    const menuShow = ref(true)
    const videoUrl = ref<string>("")
    const downloadUrl = ref<string>("")
    const youtubeRegex = /^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$/;
    const bilibiliRegex = /^https:\/\/www\.bilibili\.com\/video\/BV/;
    const parseing = ref<boolean>(false)
    const determineUrl = (url:string) =>{
        if(youtubeRegex.test(url)){
            getBaseUrl(url)
            return
        }else if (bilibiliRegex.test(url)){
            downloadUrl.value=url;
            return;
        }
        downloadUrl.value=""
    }

    const getBaseUrl = (url:string) =>{
        const index = url.indexOf('?list');
        if(index !== -1){
            downloadUrl.value=url.substring(0, index);
        }else{
            downloadUrl.value=url;
        }
    }

    const getVideoMeta =async ()=>{
        parseing.value=true
        if(downloadUrl.value.length===0){
            resProcess.value.info.title = '未知'
            resProcess.value.info.resolution = '未知'
            parseing.value=false
            return
        }
        try {
            resProcess.value = await appGetVideoMeta(downloadUrl.value)
            console.log(resProcess.value)
            // videoTitle.value= data.info.title
            // videoBestFormats.value = formatResolution(data.info.resolution)
        }catch (error:any){
            resProcess.value.info.title = '未知'
            resProcess.value.info.resolution = '未知'
            message.error(error,1)
        }
        parseing.value=false
    }

    const downloadYoutube =async (isKey:boolean,retry:boolean)=>{
        try {
            if(isKey){
                await appDownloadYoutubeByKey(resProcess.value,retry)
            }else {
                await appDownloadYoutube(downloadUrl.value,[])
            }
            return await getAllVideoInfo()
        }catch (error:any){
            throw error
        }
    }

    const getAllVideoInfo = async () => {

        process.value = await appAll()
        // menuShow.value = false
        // menuShow.value = true

        console.log(process.value)
    }

    const dbPersist = async ()=>{
        try {
            await appDbPersist()
        }catch (error:any){
            throw error
        }
    }


    return {
        videoUrl,
        downloadUrl,
        parseing,
        process,
        resProcess,
        menuShow,
        determineUrl,
        getVideoMeta,
        downloadYoutube,
        dbPersist,
        getAllVideoInfo,
        formatSize,
        formatSpeedMiB,
        formatResolution
    }
})

type YtDlpCookie ={
    bilibili?:string
}
type YtdlpSetting = {
    downloadPath: string
    queueSize:string
    cookies:YtDlpCookie
}
export const useYtdlpSettingsStore = defineStore('ytdlp-settings', () =>{
    let latestYtdlpConfig = ''
    const ytdlpConfig = ref<YtdlpSetting>({cookies: {bilibili:""}, queueSize: "", downloadPath: ""})

    const setupYtdlpSettings = async () => {
        try {
            const b = await Readfile('data/ytdlp.yaml')
            ytdlpConfig.value = Object.assign(ytdlpConfig.value, parse(b))
        } catch (error) {
            console.log(error)
        }

    }
    const saveYtdlpSettings = debounce(async (config: string) => {
        console.log('save ytdlp settings')
        try {
            await Writefile('data/ytdlp.yaml', config)
            await UpdateYtDlpConfig()
        }catch (error:any){
            message.error(error)
        }

    }, 1500)

    watch(
        ytdlpConfig,
        (settings) => {

            const lastModifiedConfig = stringify(settings)
            if (latestYtdlpConfig !== lastModifiedConfig) {
                saveYtdlpSettings(lastModifiedConfig).then(() => {
                    latestYtdlpConfig = lastModifiedConfig
                })
            } else {
                saveYtdlpSettings.cancel()
            }
        },
        { deep: true }
    )
    return {
        ytdlpConfig,
        setupYtdlpSettings,
    }
})