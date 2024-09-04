<script setup lang="ts">
import StableBranch from './StableBranch.vue'
import AlphaBranch from './AlphaBranch.vue'
import SwitchBranch from './SwitchBranch.vue'
import Ffmpeg from "@/views/KernelView/Ffmpeg.vue";
</script>

<template>
  <div>
    <StableBranch />
    <AlphaBranch />
    <Ffmpeg />
    <SwitchBranch />
  </div>
</template>

<style lang="less" scoped></style>
