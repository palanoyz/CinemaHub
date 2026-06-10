<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MovieCard from '@/components/MovieCard.vue'
import api from '@/services/api'

interface Movie {
  id: string
  title: string
  description: string
  poster_url: string
  duration: number
  genre: string[]
  rating: number
}

const movies = ref<Movie[]>([])
const loading = ref(true)
const error = ref('')

async function fetchMovies() {
  try {
    const response = await api.get('/movies')
    movies.value = response.data
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Failed to fetch movies'
  } finally {
    loading.value = false
  }
}

onMounted(fetchMovies)
</script>

<template>
  <div class="home-view container">
    <section class="hero">
      <h1>Now Showing</h1>
      <p>Discover the latest blockbusters and book your seats instantly.</p>
    </section>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading amazing movies...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="fetchMovies">Try Again</button>
    </div>

    <div v-else class="movie-grid">
      <MovieCard v-for="movie in movies" :key="movie.id" :movie="movie" />
    </div>

    <div v-if="!loading && movies.length === 0" class="empty-state">
      <p>No movies available at the moment. Check back later!</p>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  padding-top: 40px;
  padding-bottom: 80px;
}

.hero {
  margin-bottom: 48px;
  text-align: center;
}

.hero h1 {
  font-size: 2.5rem;
  margin-bottom: 12px;
  font-weight: 800;
}

.hero p {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.movie-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 24px;
}

.loading-state, .error-state, .empty-state {
  text-align: center;
  padding: 100px 0;
}

.error-state p {
  color: var(--error-text);
  margin-bottom: 20px;
}

.error-state button {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: var(--radius-md);
}
</style>
