// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {bridge} from '../models';
import {options} from '../models';

export function AbsolutePath(arg1:string):Promise<bridge.FlagResult>;

export function AddScheduledTask(arg1:string,arg2:string):Promise<bridge.FlagResult>;

export function All():Promise<bridge.FlagResultWithData>;

export function ConnectWs(arg1:string,arg2:string):Promise<bridge.FlagResult>;

export function Delete(arg1:string):Promise<bridge.FlagResult>;

export function DisConnectWs():Promise<bridge.FlagResult>;

export function Download(arg1:string,arg2:string,arg3:{[key: string]: string},arg4:string,arg5:string):Promise<bridge.HTTPResult>;

export function DownloadYoutube(arg1:string,arg2:Array<string>):Promise<bridge.FlagResult>;

export function DownloadYoutubeByKey(arg1:string,arg2:boolean):Promise<bridge.FlagResult>;

export function Exec(arg1:string,arg2:Array<string>,arg3:bridge.ExecOptions):Promise<bridge.FlagResult>;

export function ExecBackground(arg1:string,arg2:Array<string>,arg3:string,arg4:string,arg5:bridge.ExecOptions):Promise<bridge.FlagResult>;

export function ExitApp():Promise<void>;

export function ExitKey():Promise<void>;

export function FileExists(arg1:string):Promise<bridge.FlagResult>;

export function GetEnv():Promise<bridge.EnvResult>;

export function GetInterfaces():Promise<bridge.FlagResult>;

export function GetVideoMeta(arg1:string):Promise<bridge.FlagResultWithData>;

export function HideToolWindow():Promise<void>;

export function HttpDelete(arg1:string,arg2:{[key: string]: string},arg3:string):Promise<bridge.HTTPResult>;

export function HttpGet(arg1:string,arg2:{[key: string]: string},arg3:string):Promise<bridge.HTTPResult>;

export function HttpPost(arg1:string,arg2:{[key: string]: string},arg3:string,arg4:string):Promise<bridge.HTTPResult>;

export function HttpPut(arg1:string,arg2:{[key: string]: string},arg3:string,arg4:string):Promise<bridge.HTTPResult>;

export function KillProcess(arg1:number):Promise<bridge.FlagResult>;

export function Makedir(arg1:string):Promise<bridge.FlagResult>;

export function Movefile(arg1:string,arg2:string):Promise<bridge.FlagResult>;

export function Notify(arg1:string,arg2:string,arg3:string):Promise<bridge.FlagResult>;

export function OnSecondInstanceLaunch(arg1:options.SecondInstanceData):Promise<void>;

export function OpenDirectoryDialog():Promise<bridge.FlagResult>;

export function Persist():Promise<bridge.FlagResult>;

export function Ping(arg1:string,arg2:string):Promise<bridge.FlagResult>;

export function ProcessInfo(arg1:number):Promise<bridge.FlagResult>;

export function Readfile(arg1:string):Promise<bridge.FlagResult>;

export function RemoveScheduledTask(arg1:number):Promise<void>;

export function Removefile(arg1:string):Promise<bridge.FlagResult>;

export function RestartApp():Promise<bridge.FlagResult>;

export function ShowToolWindow():Promise<void>;

export function TencentTextTranslate(arg1:string,arg2:string,arg3:string):Promise<bridge.FlagResult>;

export function UnzipGZFile(arg1:string,arg2:string):Promise<bridge.FlagResult>;

export function UnzipZIPFile(arg1:string,arg2:string):Promise<bridge.FlagResult>;

export function UpdateTray(arg1:bridge.TrayContent):Promise<void>;

export function UpdateTrayMenus(arg1:Array<bridge.MenuItem>):Promise<void>;

export function UpdateYtDlpConfig():Promise<bridge.FlagResult>;

export function Upload(arg1:string,arg2:string,arg3:{[key: string]: string},arg4:string,arg5:string):Promise<bridge.HTTPResult>;

export function Writefile(arg1:string,arg2:string):Promise<bridge.FlagResult>;
