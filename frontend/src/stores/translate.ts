import {defineStore} from "pinia";
import {ref, watch} from "vue";
import icons from '@/components/Icon/icons'

interface LanguageOptionInterface {
    label: string;
    value: (typeof icons)[number]
}
type TranslateType = {
    originalText:string
    translationText:string
    // 翻译接口来源，后续拓展
    source:string
    sourceLanguage:string
    targetLanguage:string
    languageOptions:LanguageOptionInterface []
}



export const useTranslateStore = defineStore('translate', () => {
    const loading = ref(false)
    const trans = ref<TranslateType>({
        originalText : '',
        translationText :'',
        source :'',
        // sourceOptions : [{ label: '腾讯', value: 'tengxun' }],
        sourceLanguage :'auto',
        targetLanguage : 'zh',
        languageOptions : [
            { label: '简体中文', value: 'zh' },
            { label: '英文', value: 'en' },
            { label: '日文', value: 'ja' },
            { label: '韩文', value: 'ko' }
        ]
    })
    const sourceLanguageOptions =(): LanguageOptionInterface[]=> {
        return [{ label: '自动识别', value: 'auto' }, ...trans.value.languageOptions];
    };
    // const tencentTextTranslate = (newValue:string) => {
    //   console.log(newValue)
    // }
    // watch(() =>trans.value.originalText,tencentTextTranslate)

    return {trans,loading,sourceLanguageOptions}
})