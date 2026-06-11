import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import SignUpView from '../views/SignUpView.vue'
import MovieDetailView from '../views/MovieDetailView.vue'
import SeatSelectionView from '../views/SeatSelectionView.vue'
import AdminDashboardView from '../views/AdminDashboardView.vue'
import TicketView from '../views/TicketView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: false }
    },
    {
      path: '/movie/:id',
      name: 'movie-detail',
      component: MovieDetailView,
      meta: { requiresAuth: false }
    },
    {
      path: '/booking/:showtimeId',
      name: 'seat-selection',
      component: SeatSelectionView,
      meta: { requiresAuth: true } // Must be logged in to pick seats
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/signup',
      name: 'signup',
      component: SignUpView
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminDashboardView,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/tickets',
      name: 'tickets',
      component: TicketView,
      meta: { requiresAuth: true }
    }
  ],
})


router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Wait for Firebase to finish checking the session
  await authStore.isReady

  const requiresAuth = to.meta.requiresAuth
  const requiresAdmin = to.meta.requiresAdmin

  if (requiresAuth && !authStore.user) {
    next('/login')
  } else if (requiresAdmin && !authStore.isAdmin) {
    // If user is logged in but not admin, kick back to home
    next('/')
  } else if (to.path === '/login' && authStore.user) {
    next('/')
  } else {
    next()
  }
})

export default router
