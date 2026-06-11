<script setup lang="ts">
/* eslint-disable  @typescript-eslint/no-explicit-any */
import { ref, onMounted, watch } from 'vue'
import api from '@/services/api'
import axios from 'axios'

const bookings = ref<any[]>([])
const loading = ref(true)
const error = ref('')
const filterMovie = ref('')
const filterUser = ref('')

// Modal state
const showConfirmModal = ref(false)
const bookingToCancel = ref<string | null>(null)

async function fetchBookings() {
  loading.value = true
  error.value = ""
  try {
    const params = new URLSearchParams()
    if (filterMovie.value) params.append('movie', filterMovie.value.trim())
    if (filterUser.value) params.append('user_id', filterUser.value.trim())

    const res = await api.get(`/admin/bookings?${params.toString()}`)
    bookings.value = res.data
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || "Failed to load bookings"
    } else {
      error.value = "An unexpected error occurred"
    }
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  }).replace(',', '')
}

// Debounced filters
let timeout: any = null
watch([filterMovie, filterUser], () => {
  if (timeout) clearTimeout(timeout)
  timeout = setTimeout(fetchBookings, 500)
})

onMounted(fetchBookings)

function confirmCancel(id: string) {
  bookingToCancel.value = id
  showConfirmModal.value = true
}

async function cancelBooking() {
  if (!bookingToCancel.value) return

  try {
    await api.delete(`/admin/bookings/${bookingToCancel.value}`)
    showConfirmModal.value = false
    bookingToCancel.value = null
    fetchBookings() // Refresh list
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      alert(err.response?.data?.error || "Failed to cancel")
    }
  }
}
</script>

<template>
  <div class="admin-dashboard container">
    <!-- Custom Confirmation Modal -->
    <div v-if="showConfirmModal" class="modal-overlay">
      <div class="confirm-modal">
        <h3>Confirm Cancellation</h3>
        <p>Are you sure you want to cancel this booking? This action cannot be undone.</p>
        <div class="modal-actions">
          <button @click="showConfirmModal = false" class="cancel-btn">Keep</button>
          <button @click="cancelBooking" class="delete-btn">Yes, Cancel</button>
        </div>
      </div>
    </div>

    <div class="header-section">
      <h1>Admin Dashboard</h1>
      <p>Manage and monitor all ticket bookings</p>
    </div>

    <div class="filters-card">
      <div class="filter-group">
        <label>Filter by Movie</label>
        <input v-model="filterMovie" placeholder="Search movie title..." />
      </div>
      <div class="filter-group">
        <label>Filter by User ID</label>
        <input v-model="filterUser" placeholder="Paste User ID..." />
      </div>
      <button @click="fetchBookings" class="refresh-btn">Refresh</button>
    </div>

    <div v-if="error" class="error-msg">{{ error }}</div>

    <div class="table-container shadow-md">
      <table v-if="!loading && bookings.length > 0">
        <thead>
          <tr>
            <th>Date</th>
            <th>Movie</th>
            <th>Hall</th>
            <th>Seats</th>
            <th>User ID</th>
            <th>Total</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="booking in bookings" :key="booking.id">
            <td>{{ formatDate(booking.created_at) }}</td>
            <td class="bold">{{ booking.movie_title }}</td>
            <td>{{ booking.hall_name }}</td>
            <td>
              <span class="seat-tag" v-for="s in booking.seat_ids" :key="s">{{ s }}</span>
            </td>
            <td class="user-id" :title="booking.user_id">{{ booking.user_id }}</td>
            <td class="price">{{ booking.total_price }} THB</td>
            <td>
              <button @click="confirmCancel(booking.id)" class="action-cancel-btn">Cancel</button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="loading" class="table-loading">
        <div class="spinner"></div>
        <p>Loading database records...</p>
      </div>

      <div v-else-if="bookings.length === 0" class="table-empty-state">
        <p>No bookings found matching your filters.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.action-cancel-btn {
  background: #fff;
  color: #d32f2f;
  border: 1px solid #d32f2f;
  padding: 6px 12px;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.action-cancel-btn:hover {
  background: #d32f2f;
  color: white;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.confirm-modal {
  background: white;
  padding: 30px;
  border-radius: var(--radius-lg);
  max-width: 400px;
  width: 90%;
  text-align: center;
}

.confirm-modal h3 {
  margin-top: 0;
}

.modal-actions {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-top: 20px;
}

.cancel-btn {
  background: #f0f0f0;
  border: none;
  padding: 10px 20px;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
}

.delete-btn {
  background: #d32f2f;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: var(--radius-md);
  font-weight: 600;
  cursor: pointer;
}

.admin-dashboard {
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

.filters-card {
  background: white;
  padding: 24px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  display: flex;
  gap: 20px;
  align-items: flex-end;
  margin-bottom: 32px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.filter-group label {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
}

.filter-group input {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: var(--radius-md);
  font-size: 14px;
}

.refresh-btn {
  background: #f0f0f0;
  border: 1px solid #ddd;
  padding: 12px 24px;
  border-radius: var(--radius-md);
  font-weight: 600;
}

.table-container {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

th {
  background: #f0f0f0;
  padding: 16px 24px;
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  border-bottom: 2px solid #eee;
}

td {
  padding: 16px 24px;
  border-bottom: 1px solid #eee;
  font-size: 14px;
}

.bold {
  font-weight: 700;
}

.seat-tag {
  background: #eee;
  padding: 2px 6px;
  border-radius: 4px;
  margin-right: 4px;
  font-size: 12px;
  font-weight: 600;
}

.user-id {
  color: #888;
  font-family: monospace;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price {
  font-weight: 600;
  color: #666;
}

.empty {
  text-align: center;
  padding: 60px;
  color: var(--text-muted);
}

.table-loading {
  padding: 100px;
  text-align: center;
}

.table-empty-state {
  padding: 100px;
  text-align: center;
  color: var(--text-muted);
  font-size: 1.1rem;
}

.error-msg {
  background-color: var(--error-bg);
  color: var(--error-text);
  padding: 12px;
  border-radius: var(--radius-md);
  margin-bottom: 20px;
  border: 1px solid var(--error-border);
}
</style>
