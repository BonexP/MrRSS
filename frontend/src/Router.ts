import { createRouter, createWebHistory } from 'vue-router'
import Home from './components/home/Home.vue'
import Settings from './components/settings/Settings.vue'
import RssSettings from './components/settings/RssSettings.vue'
import PreferenceSettings from './components/settings/PreferenceSettings.vue'
import AboutSettings from './components/settings/AboutSettings.vue'

const routes = [
  { 
    path: '/', 
    component: Home 
  },
  { 
    path: '/settings', 
    component: Settings,
    redirect: '/settings/rss',
    children: [
      {
        path: 'rss',
        component: RssSettings
      },
      {
        path: 'preference',
        component: PreferenceSettings
      },
      {
        path: 'about',
        component: AboutSettings
      }
    ]
  }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router