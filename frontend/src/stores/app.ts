import { ref } from 'vue'
import { defineStore } from 'pinia'

export type Menu = {
  label: string
  handler?: (...args: any) => void
  separator?: boolean
  children?: Menu[]
}

export const useAppStore = defineStore('app', () => {
  // const widgetsType = ref('')
  const widgetsEnable = ref(false)
  const widgetsType = ref("")
  // const changeWin = (status:boolean,type:string) => {
  //   widgetsEnable.value = status
  //   widgetsType.value = type
  //   console.log(widgetsType.value)
  //   console.log(widgetsEnable.value)
  // }
  /* Global Menu */
  const menuShow = ref(false)
  const menuList = ref<Menu[]>([])
  const menuPosition = ref({
    x: 0,
    y: 0
  })

  /* Global Tips */
  const tipsShow = ref(false)
  const tipsMessage = ref('')
  const tipsPosition = ref({
    x: 0,
    y: 0
  })

  const showAbout = ref(false)


  return {
    widgetsEnable,
    widgetsType,
    menuShow,
    menuPosition,
    menuList,
    tipsShow,
    tipsMessage,
    tipsPosition,
    showAbout
  }
})
