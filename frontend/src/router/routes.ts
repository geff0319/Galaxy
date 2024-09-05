import { type RouteRecordRaw } from 'vue-router'

import { isDev } from '@/utils'

import HomeView from '@/views/HomeView/index.vue'
import SubscribesView from '@/views/SubscribesView/index.vue'
import SettingsView from '@/views/SettingsView/index.vue'
import ProfilesView from '@/views/ProfilesView/index.vue'
import RulesetsView from '@/views/RulesetsView/index.vue'
import PluginsView from '@/views/PluginsView/index.vue'
import ScheduledTasksView from '@/views/ScheduledTasksView/index.vue'
import PlaygroundView from '@/views/PlaygroundView/index.vue'
import TranslationView from '@/views/TranslationView/index.vue'
import YtdlpView from '@/views/YtdlpView/index.vue'


const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
    meta: {
      name: 'router.overview'
    }
  },
  {
    path: '/profiles',
    name: 'Profiles',
    component: ProfilesView,
    meta: {
      name: 'router.profiles',
      hidden: !isDev
    }
  },
  {
    path: '/subscriptions',
    name: 'Subscriptions',
    component: SubscribesView,
    meta: {
      name: 'router.subscriptions',
      hidden: !isDev
    }
  },
  {
    path: '/rulesets',
    name: 'Rulesets',
    component: RulesetsView,
    meta: {
      name: 'router.rulesets',
      hidden: !isDev
    }
  },
  {
    path: '/plugins',
    name: 'PluginsView',
    component: PluginsView,
    meta: {
      name: 'router.plugins'
    }
  },
  {
    path: '/translation',
    name: 'Translation',
    component: TranslationView,
    meta: {
      name: 'router.translation'
    }
  },
  {
    path: '/ytdlp',
    name: 'Ytdlp',
    component: YtdlpView,
    meta: {
      name: '下载',
      keepAlive:false
    }
  },
  {
    path: '/scheduledtasks',
    name: 'ScheduledTasksView',
    component: ScheduledTasksView,
    meta: {
      name: 'router.scheduledtasks'
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView,
    meta: {
      name: 'router.settings'
    }
  },
  {
    path: '/playground',
    name: 'Playground',
    component: PlaygroundView,
    meta: {
      name: 'Playground',
      hidden: !isDev
    }
  }
]

export default routes
