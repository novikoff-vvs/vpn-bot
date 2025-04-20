import { createMemoryHistory, createRouter } from 'vue-router'

import HomeView from '@/components/MainPage.vue'
import Payments from '@/components/PaymentsPage.vue'
import Profile from '@/components/ProfilePage.vue'

const routes = [
    { path: '/', component: HomeView },
    { path: '/payments', component: Payments },
    { path: '/profile', component: Profile },
]

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

