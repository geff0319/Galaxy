<script setup lang="ts">
import FlipItem from './FlipItem.vue'
import {onMounted,onBeforeUnmount, ref} from "vue";


const timeArr = ref<any[]>([])
const timer = ref(0)
// 更换数组类型
const toArr = (n:number) =>{
  return n >= 10 ? ('' + n).split('').map(item => Number(item)) : [0, n]
}

const getTimeArr = (now:Date = new Date())=> {
  const h = now.getHours()
  const m = now.getMinutes()
  const s = now.getSeconds()
  return  [
    ...toArr(h),
    ...toArr(m),
    ...toArr(s)
  ]
}

const startTimer = () => {
  stopTimer()
  timer.value =setTimeout(() => {
      timeArr.value = getTimeArr()
      startTimer()
    }, 1000)
}

const stopTimer = () =>{
  clearTimeout(timer.value)
}
timeArr.value = getTimeArr()
onMounted(() => startTimer())
onBeforeUnmount(() => stopTimer())
</script>

<template>
  <div class="clock-container">
    <FlipItem :total="2" :current="timeArr[0]"/>
    <FlipItem :total="9" :current="timeArr[1]"/>
    <div class="colon"></div>
    <FlipItem :total="5" :current="timeArr[2]"/>
    <FlipItem :total="9" :current="timeArr[3]"/>
    <div class="colon"></div>
    <FlipItem :total="5" :current="timeArr[4]"/>
    <FlipItem :total="9" :current="timeArr[5]"/>
  </div>
</template>

<style lang="less" scoped>
.clock-container {
  display: flex;
  align-items: center;

  .colon {
    height: 50px;
    padding: 0 2px;
    display: flex;
    justify-content: space-around;
    flex-direction: column;

    &::after,
    &::before {
      content: '';
      display: block;
      width: 10px;
      height: 10px;
      background: rgba(0, 0, 0, 0.7);
      border-radius: 50%;
    }
  }
}
</style>