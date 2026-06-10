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
        <span class="duration">
          <!-- <span class="icon">⏱️</span> -->
          {{ movie.duration }} min
        </span>
        <div class="genres-list">
          <span v-for="genre in movie.genre" :key="genre" class="genre-badge">{{ genre }}</span>
        </div>
      </div>
      <p class="description">{{ movie.description }}</p>
      <button class="book-btn">Book Tickets</button>
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
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

h3 {
  margin: 0 0 4px;
  font-size: 1rem;
  color: var(--text-main);
}

.meta {
  margin-bottom: 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.duration {
  font-size: 0.75rem;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon {
  font-size: 0.85rem;
}

.genres-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.genre-badge {
  display: inline-block;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2px 6px;
  border-radius: 12px;
  font-size: 0.65rem;
  font-weight: 500;
}

.description {
  font-size: 0.8rem;
  color: var(--text-muted);
  line-height: 1.4;
  margin-bottom: 12px;
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
  padding: 8px;
  border-radius: var(--radius-md);
  font-weight: 700;
  font-size: 0.9rem;
  width: 100%;
}

.book-btn:hover {
  background: var(--primary-hover);
}
</style>
