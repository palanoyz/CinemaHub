<script setup lang="ts">
/* eslint-disable  @typescript-eslint/no-explicit-any */
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/services/api'
import axios from 'axios'

const route = useRoute()
const movie = ref<any>(null)
const showtimes = ref<any[]>([])
const loading = ref(true)
const error = ref('')

async function fetchData() {
  try {
    const movieId = route.params.id

    // Fetch movie
    const moviesRes = await api.get('/movies')
    movie.value = moviesRes.data.find((m: any) => m.id === movieId)
    
    // Fetch showtimes
    const showtimesRes = await api.get(`/movies/${movieId}/showtimes`)
    showtimes.value = showtimesRes.data
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || err.message
    } else {
      error.value = 'Failed to load movie details'
    }
  } finally {
    loading.value = false
  }
}

function formatTime(dateStr: string) {
  return new Date(dateStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric' })
}

onMounted(fetchData)
</script>

<template>
  <div class="movie-detail container">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
    </div>

    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
    </div>

    <div v-else-if="movie" class="detail-grid">
      <div class="poster-section">
        <img :src="movie.poster_url" :alt="movie.title" class="main-poster" />
      </div>

      <div class="info-section">
        <h1>{{ movie.title }}</h1>
        <div class="badges">
          <span class="rating">★ {{ movie.rating }}</span>
          <span class="duration">{{ movie.duration }} min</span>
        </div>
        
        <div class="genres">
          <span v-for="g in movie.genre" :key="g" class="genre-tag">{{ g }}</span>
        </div>

        <h3>Synopsis</h3>
        <p class="description">{{ movie.description }}</p>

        <div class="showtimes-container">
          <h3>Select Showtime</h3>
          <div v-if="showtimes.length > 0" class="showtime-grid">
            <router-link 
              v-for="s in showtimes" 
              :key="s.id" 
              :to="'/booking/' + s.id"
              class="showtime-card"
            >
              <span class="hall">{{ s.hall_name }}</span>
              <span class="time">{{ formatTime(s.start_time) }}</span>
              <span class="date">{{ formatDate(s.start_time) }}</span>
            </router-link>
          </div>
          <p v-else class="no-showtimes">No showtimes available for this movie yet.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.movie-detail {
  padding: 60px 20px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 60px;
}

.main-poster {
  width: 100%;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

h1 {
  font-size: 3rem;
  margin: 0 0 16px;
  font-weight: 800;
}

.badges {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
  font-weight: 600;
}

.rating { color: #ffc107; }

.genres {
  display: flex;
  gap: 10px;
  margin-bottom: 32px;
}

.genre-tag {
  background: white;
  border: 1px solid #ddd;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.9rem;
}

h3 {
  margin: 0 0 12px;
  font-size: 1.25rem;
}

.description {
  color: var(--text-muted);
  line-height: 1.7;
  font-size: 1.1rem;
  margin-bottom: 40px;
}

.showtime-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.showtime-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: white;
  border: 2px solid transparent;
  border-radius: var(--radius-md);
  text-decoration: none;
  color: inherit;
  box-shadow: var(--shadow-sm);
  transition: all 0.2s;
}

.showtime-card:hover {
  border-color: var(--primary-color);
  transform: translateY(-3px);
}

.hall {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.time {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--primary-color);
}

.date {
  font-size: 0.8rem;
  margin-top: 4px;
}

.loading-state {
  display: flex;
  justify-content: center;
  padding: 100px;
}

@media (max-width: 900px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
  .poster-section {
    max-width: 300px;
    margin: 0 auto;
  }
}
</style>
