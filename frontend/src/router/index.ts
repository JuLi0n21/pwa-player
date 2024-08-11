import { createRouter, createWebHistory } from 'vue-router'
import MenuView from '../views/MenuView.vue'
import RecentView from '../views/RecentView.vue'
import FavouritView from '../views/FavouritView.vue'
import CollectionView from '../views/CollectionView.vue'
import NowPlayingView from '../views/NowPlayingView.vue'
import MeView from '../views/MeView.vue'
import SearchView from '../views/SearchView.vue'
import CollectionItem from '../components/CollectionItem.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '',
      component: MeView,
    },
    {
      path: '/menu',
      name: 'menu',
      component: MenuView,
      children: [
        {
          path: '',
          name: 'default',
          component: RecentView
        },
        {
          path: 'recent',
          name: 'recent',
          component: RecentView
        },
        {
          path: 'favourites',
          name: 'favourites',
          component: FavouritView
        },
        {
          path: 'collections',
          name: 'collections',
          component: CollectionView
        }
      ]
    },
    {
      path: '/nowplaying',
      name: 'nowplaying',
      component: NowPlayingView
    },
    {
      path: '/me',
      name: 'me',
      component: MeView
    },
    {
      path: '/search',
      name: 'search',
      component: SearchView
    },
    {
      path: '/collection/:id',
      name: 'collection',
      component: CollectionItem
    }
  ]
})

export default router
