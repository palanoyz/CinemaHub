<script setup lang="ts">
defineProps<{
  movie: {
    id: string
    title: string
    description: string
    poster_url: string
    duration: number
    genre: string[]
    rating: number
  }
}>()
</script>

<template>
  <div class="movie-card">
    <div class="poster-wrapper">
      <img :src="movie.poster_url" :alt="movie.title" />
      <div class="rating-badge">★ {{ movie.rating }}</div>
    </div>
    <div class="info">
      <h3>{{ movie.title }}</h3>
      <div class="meta">
        <span class="duration">{{ movie.duration }} min</span>
        <div class="genres-list">
          <span v-for="genre in movie.genre" :key="genre" class="genre-badge">{{ genre }}</span>
        </div>
      </div>
      <p class="description">{{ movie.description }}</p>
      <router-link :to="'/movie/' + movie.id" class="book-btn">
        Book Tickets
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.movie-card {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: flex;
  flex-direction: column;
}

.movie-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
}

.poster-wrapper {
  position: relative;
  aspect-ratio: 2/3;
}

.poster-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.rating-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(0, 0, 0, 0.7);
  color: #ffc107;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 14px;
  backdrop-filter: blur(4px);
}

.info {
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

h3 {
  margin: 0 0 8px;
  font-size: 1.25rem;
  color: var(--text-main);
}

.meta {
  margin-bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.duration {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.genres-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.genre-badge {
  background: #f0f0f0;
  color: var(--text-muted);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.description {
  font-size: 0.9rem;
  color: var(--text-muted);
  line-height: 1.5;
  margin-bottom: 20px;
  display: -webkit-box;
  line-clamp: 3;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.book-btn {
  margin-top: auto;
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 12px;
  border-radius: var(--radius-md);
  font-weight: 700;
  width: 100%;
  text-decoration: none;
  text-align: center;
  display: block;
  transition: background 0.2s;
}

.book-btn:hover {
  background: var(--primary-hover);
}
</style>
