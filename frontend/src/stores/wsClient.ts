import {defineStore} from "pinia";
import {appConnectWs, appDisConnectWs, appPing} from "@/bridge/ws";
import {ref, watch} from "vue";
import {Readfile, Writefile} from "@/bridge";
import {parse, stringify} from "yaml";
import {debounce} from "@/utils";
import {message} from "ant-design-vue";

type WsSettings = {
    id:string
    autoConnect:boolean
    domain:string
}
function generateRandomId(length: number): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
        const randomIndex = Math.floor(Math.random() * chars.length);
        result += chars[randomIndex];
    }
    return result;
}

export const useWsClientStore = defineStore('wsClient', () => {
    let firstOpen = true
    let latestWsConfig = ''
    const isConnected = ref(false)
    const ws = ref<WsSettings>({
        domain: "",
        autoConnect:false,
        id: generateRandomId(8)
    })
    const connectWs =async ()=>{
        if(ws.value.id.length===0 || ws.value.domain.length===0){
            message.error("WS连接信息不全")
            return
        }
        try {
            await appConnectWs(ws.value.domain, ws.value.id)
            isConnected.value=true
            message.info("WS连接成功")
        }catch (error:any){
            message.error(error)
            throw error
        }
    }
    const disConnectWs =async ()=>{
        try {
            await appDisConnectWs()
            isConnected.value=false
            message.info("WS已断开连接")
        }catch (error:any){
            message.error(error)
            throw error
        }
    }
    const connectionStatus = ()=>{
        isConnected.value =false
        console.log(ws)
        if (ws.value.domain.length === 0){
            return
        }
        appPing(ws.value.domain, ws.value.id).then((data)=>{
            isConnected.value =true
        }).catch((error:any)=>{
            isConnected.value = false
            throw error
        })
    }

    const setupWsSettings = async () => {
        try {
            const b = await Readfile('data/ws.yaml')
            ws.value = Object.assign(ws.value, parse(b))
        } catch (error) {
            firstOpen = false
            console.log(error)
        }
    }
    const saveWsSettings = debounce((config: string) => {
        console.log('save ws settings')
        Writefile('data/ws.yaml', config)

    }, 1500)

    watch(
        ws,
        (settings) => {
            if (!firstOpen) {
                const lastModifiedConfig = stringify(settings)
                if (latestWsConfig !== lastModifiedConfig) {
                    saveWsSettings(lastModifiedConfig).then(() => {
                        latestWsConfig = lastModifiedConfig
                    })
                } else {
                    saveWsSettings.cancel()
                }
            }

            firstOpen = false
        },
        { deep: true }
    )
    return {
        connectWs,
        disConnectWs,
        connectionStatus,
        setupWsSettings,
        ws,
        isConnected
    }
})
