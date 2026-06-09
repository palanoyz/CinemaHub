<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header v-if="!authStore.loading" class="main-header">
    <nav class="container">
      <div class="brand" @click="router.push('/')">
        <span class="logo">🎬</span>
        <span class="text">CinemaHub</span>
      </div>
      
      <div v-if="authStore.user" class="user-info">
        <div class="profile">
          <img v-if="authStore.user?.photoURL" :src="authStore.user.photoURL" alt="Profile" class="avatar" referrerpolicy="no-referrer" />
          <div v-else class="avatar-placeholder">{{ authStore.user?.email?.[0]?.toUpperCase() || '?' }}</div>
          <span class="email">{{ authStore.user?.email }}</span>
        </div>
        <button @click="handleLogout" class="logout-btn">Logout</button>
      </div>
      <div v-else-if="$route.name !== 'login'" class="nav-actions">
        <button @click="router.push('/login')" class="login-btn">Login</button>
      </div>
    </nav>
  </header>

  <main>
    <div v-if="authStore.loading" class="loading-screen">
      <div class="spinner"></div>
      <p>Initializing CinemaHub...</p>
    </div>
    <div v-else class="content-container">
      <RouterView />
    </div>
  </main>
</template>

<style>
:root {
  --primary-color: #e50914;
  --text-main: #1a1a1a;
  --text-muted: #666;
  --bg-body: #f8f9fa;
  --nav-height: 72px;
}

body {
  margin: 0;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  background-color: var(--bg-body);
  color: var(--text-main);
  -webkit-font-smoothing: antialiased;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.main-header {
  height: var(--nav-height);
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

nav {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.brand .logo {
  font-size: 28px;
}

.brand .text {
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: -0.5px;
  background: linear-gradient(45deg, #e50914, #ff4b5c);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.user-info {
  display: flex;
  gap: 24px;
  align-items: center;
}

.profile {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar, .avatar-placeholder {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  background: #eee;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  color: var(--text-muted);
}

.email {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
}

.logout-btn {
  background: transparent;
  border: 1px solid #ddd;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.logout-btn:hover {
  background: #f5f5f5;
  border-color: #ccc;
}

.loading-screen {
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(229, 9, 20, 0.1);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
