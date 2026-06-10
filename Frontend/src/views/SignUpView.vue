<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const authStore = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const loading = ref(false)

async function handleSignUp() {
  if (password.value !== confirmPassword.value) {
    error.value = "Passwords do not match"
    return
  }

  error.value = ""
  loading.value = true

  try {
    await authStore.signUp(email.value, password.value)
    router.push('/')
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || err.message
    } else if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'An unexpected error occurred'
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
        <h1>Create Account</h1>
        <p>Join CinemaHub to book your seats</p>
      </div>
      
      <div v-if="error" class="error-msg">{{ error }}</div>

      <form @submit.prevent="handleSignUp" class="email-form">
        <div class="input-group">
          <input v-model="email" type="email" placeholder="Email address" required />
        </div>
        <div class="input-group">
          <input v-model="password" type="password" placeholder="Password" required />
        </div>
        <div class="input-group">
          <input v-model="confirmPassword" type="password" placeholder="Confirm Password" required />
        </div>
        <button type="submit" class="email-btn" :disabled="loading">
          {{ loading ? 'Creating account...' : 'Sign Up' }}
        </button>
      </form>

      <div class="footer">
        Already have an account? <router-link to="/login">Login</router-link>
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

.email-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
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
