<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const authStore = useAuthStore()
const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleGoogleLogin() {
  try {
    await authStore.loginWithGoogle()
    router.push('/')
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || err.message
    } else if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'An unexpected error occurred'
    }
  }
}

async function handleEmailLogin() {
  error.value = ""
  loading.value = true
  try {
    await authStore.loginWithEmail(email.value, password.value)
    router.push('/')
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || err.message
    } else if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'Invalid email or password'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="header">
        <div class="logo-icon">🎬</div>
        <h1>Welcome Back</h1>
        <p>Login to book your favorite seats</p>
      </div>

      <div v-if="error" class="error-msg">{{ error }}</div>

      <div class="login-actions">
        <button @click="handleGoogleLogin" class="google-btn">
          <svg width="18" height="18" viewBox="0 0 18 18">
            <path
              d="M17.64 9.2c0-.63-.06-1.25-.16-1.84H9v3.49h4.84a4.14 4.14 0 0 1-1.8 2.71v2.26h2.91c1.71-1.58 2.69-3.91 2.69-6.62z"
              fill="#4285F4" />
            <path
              d="M9 18c2.43 0 4.47-.8 5.96-2.18l-2.91-2.26c-.8.54-1.83.85-3.05.85-2.35 0-4.33-1.58-5.04-3.71H.95v2.33A8.99 8.99 0 0 0 9 18z"
              fill="#34A853" />
            <path d="M3.96 10.7a5.41 5.41 0 0 1 0-3.4V4.97H.95a8.99 8.99 0 0 0 0 8.06l3.01-2.33z" fill="#FBBC05" />
            <path
              d="M9 3.58c1.32 0 2.5.45 3.44 1.35L15.02 2.3A8.99 8.99 0 0 0 .95 4.97l3.01 2.33c.71-2.13 2.69-3.71 5.04-3.71z"
              fill="#EA4335" />
          </svg>
          Continue with Google
        </button>

        <div class="divider">
          <span>OR</span>
        </div>

        <form @submit.prevent="handleEmailLogin" class="email-form">
          <div class="input-group">
            <input v-model="email" type="email" placeholder="Email address" required />
          </div>
          <div class="input-group">
            <input v-model="password" type="password" placeholder="Password" required />
          </div>
          <button type="submit" class="email-btn" :disabled="loading">
            {{ loading ? 'Logging in...' : 'Login with Email' }}
          </button>
        </form>
      </div>

      <div class="footer">
        Don't have an account? <router-link to="/signup">Sign up</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: calc(100vh - var(--nav-height));
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-card {
  background: var(--white);
  width: 100%;
  max-width: 420px;
  padding: 40px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.header {
  margin-bottom: 32px;
  text-align: center;
}

.logo-icon {
  font-size: 40px;
  margin-bottom: 16px;
}

h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: var(--text-main);
}

p {
  margin: 8px 0 0;
  color: var(--text-muted);
  font-size: 15px;
}

.login-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.google-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background-color: var(--white);
  color: #3c4043;
  border: 1px solid #dadce0;
  padding: 12px;
  border-radius: var(--radius-md);
  font-size: 15px;
  font-weight: 500;
}

.google-btn:hover {
  background-color: #f8f9fa;
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  margin: 8px 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid #eee;
}

.divider span {
  padding: 0 15px;
  color: var(--text-light);
  font-size: 13px;
  font-weight: 600;
}

.email-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.input-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e0e0e0;
  border-radius: var(--radius-md);
  font-size: 15px;
  box-sizing: border-box;
  transition: border-color 0.2s;
}

.input-group input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.email-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 14px;
  border-radius: var(--radius-md);
  font-size: 15px;
  font-weight: 600;
}

.email-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.footer {
  margin-top: 32px;
  text-align: center;
  font-size: 14px;
  color: #666;
}

.footer a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 600;
}

.error-msg {
  background-color: var(--error-bg);
  color: var(--error-text);
  padding: 12px;
  border-radius: var(--radius-md);
  margin-bottom: 20px;
  font-size: 14px;
  border: 1px solid var(--error-border);
}
</style>
