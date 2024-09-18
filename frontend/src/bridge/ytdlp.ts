import * as App from '@wails/go/bridge/App'
import {bridge} from "@wails/go/models";
import type {ProcessType} from "@/stores";


export const appGetVideoMeta = async (path:string) => {
    const { flag, msg,data } = await App.GetVideoMeta(path)
    if (!flag) {
        throw msg
    }
    return data
}

export const appDownloadYoutube = async (url:string, params:string[]) => {
    const { flag,data } = await App.DownloadYoutube(url,params)
    if (!flag) {
        throw "下载视频失败"
    }
    return data
}
export const appDownloadYoutubeByKey = async (p:ProcessType,retry:boolean) => {
    const { flag,data } = await App.DownloadYoutubeByKey(JSON.stringify(p),retry)
    if (!flag) {
        throw data
    }
    return data
}


export const appDbPersist = async () => {
    const { flag,data } = await App.Persist()
    if (!flag) {
        throw data
    }
    return data
}

export const appAll = async () => {
    const { flag,data } = await App.All()
    if (!flag) {
        throw data
    }
    return data
}
export const UpdateYtDlpConfig = async () => {
    const { flag,data } = await App.UpdateYtDlpConfig()
    if (!flag) {
        throw data
    }
    return data
}
export const deleteProcess = async (id:string) => {
    const { flag,data } = await App.Delete(id)
    if (!flag) {
        throw data
    }
    return data
}
export const AppCheckBiliLogin = async () => {
    const { flag,data } = await App.CheckBiliLogin()
    if (!flag) {
        throw data
    }
    return data
}
