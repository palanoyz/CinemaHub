<script setup lang="ts">
defineProps<{
  movie: {
    id: string
    title: string
    poster_url: string
    rating: number
  }
}>()
</script>

<template>
  <router-link :to="'/movie/' + movie.id" class="movie-card">
    <div class="poster-wrapper">
      <img :src="movie.poster_url" :alt="movie.title" />
      <div class="rating-badge">★ {{ movie.rating }}</div>
    </div>
    <div class="info">
      <h3>{{ movie.title }}</h3>
      <div class="book-btn">
        View Details
      </div>
    </div>
  </router-link>
</template>

<style scoped>
.movie-card {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  text-decoration: none;
  color: inherit;
}

.movie-card:hover {
  transform: translateY(-8px);
  box-shadow: var(--shadow-md);
}

.poster-wrapper {
  position: relative;
  aspect-ratio: 2/3;
  overflow: hidden;
}

.poster-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.movie-card:hover img {
  transform: scale(1.05);
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
  font-size: 12px;
  backdrop-filter: blur(4px);
}

.info {
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 15px;
}

h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.3;
}

.book-btn {
  background: var(--primary-color);
  color: white;
  padding: 10px;
  border-radius: var(--radius-md);
  font-weight: 700;
  font-size: 0.9rem;
  text-align: center;
  transition: background 0.2s;
}

.movie-card:hover .book-btn {
  background: var(--primary-hover);
}
</style>
