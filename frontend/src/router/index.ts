import { createRouter, createWebHashHistory } from 'vue-router'

import routes from './routes'

console.log(routes)
const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes
})

export default router