<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();
</script>

<template>
  <div class="content">
<!--     <RouterView />-->
    <RouterView v-slot="{ Component }">
      <KeepAlive>
        <component v-if="route.meta.keepAlive" :is="Component" />
      </KeepAlive>
      <component v-if="!route.meta.keepAlive" :is="Component" />
    </RouterView>
<!--    <RouterView v-slot="{ Component }">-->
<!--      <KeepAlive>-->
<!--        <component :is="Component" />-->
<!--      </KeepAlive>-->
<!--    </RouterView>-->
  </div>
</template>

<style scoped>
.content {
  overflow-y: auto;
  margin-top: 8px;
  padding: 0 8px;
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
