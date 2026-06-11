<script setup lang="ts">
/* eslint-disable  @typescript-eslint/no-explicit-any */
import { ref, onMounted } from 'vue'
import api from '@/services/api'
import axios from 'axios'

const tickets = ref<any[]>([])
const loading = ref(true)
const error = ref('')

async function fetchTickets() {
  loading.value = true
  try {
    const res = await api.get('/protected/bookings')
    tickets.value = res.data
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || "Failed to load tickets"
    } else {
      error.value = "An unexpected error occurred"
    }
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('en-TH', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  }).replace(',', '')
}

onMounted(fetchTickets)
</script>

<template>
  <div class="tickets-page container">
    <div class="header-section">
      <h1>My Tickets</h1>
      <p>View your upcoming and past bookings</p>
    </div>

    <div v-if="loading" class="loading">Loading tickets...</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else class="tickets-grid">
      <div v-for="ticket in tickets" :key="ticket.id" class="ticket-card">
        <div class="ticket-header">
          <h3>{{ ticket.movie_title }}</h3>
          <span class="hall-badge">{{ ticket.hall_name }}</span>
        </div>
        <div class="ticket-body">
          <p>{{ formatDate(ticket.start_time) }}</p>
          <p>Seats: <strong>{{ ticket.seat_ids.join(', ') }}</strong></p>
        </div>
        <div class="ticket-footer">
          <span class="total-label">Total</span>
          <span class="price">{{ ticket.total_price }} THB</span>
        </div>
      </div>

      <div v-if="tickets.length === 0" class="empty">You have no bookings yet.</div>
    </div>
  </div>
</template>

<style scoped>
.tickets-page {
  padding: 40px 20px;
}

.header-section {
  margin-bottom: 40px;
}

h1 {
  font-size: 2.5rem;
  margin: 0;
  font-weight: 800;
}

.header-section p {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.tickets-grid {
  display: grid;
  gap: 24px;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
}

.ticket-card { 
  background: white; 
  padding: 24px; 
  border-radius: var(--radius-lg); 
  box-shadow: var(--shadow-sm);
  border: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.ticket-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
  border-color: var(--primary-color);
}

.ticket-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 15px;
}

.ticket-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-main);
  flex: 1;
  min-width: 0;
  word-wrap: break-word;
}

.hall-badge {
  background: #f0f0f0;
  padding: 4px 10px;
  border-radius: var(--radius-md);
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  flex-shrink: 0;
}

.ticket-body {
  border-top: 1px dashed #ddd;
  border-bottom: 1px dashed #ddd;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ticket-body p {
  margin: 0;
  font-size: 0.95rem;
  color: var(--text-muted);
}

.ticket-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}

.total-label {
  font-size: 0.85rem;
  color: var(--text-muted);
  font-weight: 600;
}

.price {
  font-weight: 800;
  font-size: 1.2rem;
  color: var(--primary-color);
}

.empty {
  text-align: center;
  padding: 60px;
  color: var(--text-muted);
}
</style>
