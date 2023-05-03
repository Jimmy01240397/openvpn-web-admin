import { createRouter, createWebHashHistory } from 'vue-router'
//import HomeView from '../views/HomeView.vue'
import LogIN from '../views/Login.vue'
import LayOut from '../views/LayOut.vue'
import DashBoard from '../views/DashBoard.vue'
import LOG from '../views/LOG.vue'
import USER from '../views/USER.vue'
import GROUP from '../views/GROUP.vue'
import VPNs from '../views/VPNs.vue'

const routes = [
  {
    path: '/',
    name: 'LogIN',
    component: LogIN
  },
  {
    path: '/layout',
    name: 'LayOut',
    component: LayOut,
    children: [
        {
            path: 'dashboard',
            name: 'DashBoard',
            component: DashBoard,
        },
        {
            path: 'log',
            name: 'LOG',
            component: LOG,
        },
        {
            path: 'user',
            name: 'USER',
            component: USER,
        },
        {
            path: 'group',
            name: 'GROUP',
            component: GROUP,
        },
        {
            path: 'vpns',
            name: 'VPNs',
            component: VPNs,
        },
    ]
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
