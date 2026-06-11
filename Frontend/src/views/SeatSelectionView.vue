<script setup lang="ts">
/* eslint-disable  @typescript-eslint/no-explicit-any */
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const showtime = ref<any>(null)
const selectedSeats = ref<string[]>([])
const loading = ref(true)
const error = ref('')
const locking = ref(false)
const bookingSuccess = ref(false)
let ws: WebSocket | null = null

async function fetchShowtime() {
  try {
    const id = route.params.showtimeId
    const res = await api.get(`/showtimes/${id}`)
    showtime.value = res.data

    syncLocks()
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || err.message
    } else {
      error.value = 'Failed to load seat map'
    }
  } finally {
    loading.value = false
  }
}

function syncLocks() {
  if (showtime.value && authStore.user) {
    // Source of truth: Server-side locked seats belonging to me
    selectedSeats.value = showtime.value.seats
      .filter((s: any) => s.status === 'LOCKED' && s.locked_by === authStore.user?.uid)
      .map((s: any) => s.id)
  }
}

function setupWebSocket() {
  const wsBase = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api'
  const wsUrl = `${wsBase}/showtimes/${route.params.showtimeId}/ws`

  ws = new WebSocket(wsUrl)

  ws.onmessage = (event) => {
    const update = JSON.parse(event.data)

    if (showtime.value && update.showtime_id === route.params.showtimeId) {
      showtime.value.seats = showtime.value.seats.map((s: any) => {
        if (update.seat_ids.includes(s.id)) {
          return { ...s, status: update.status, locked_by: update.locked_by }
        }
        return s
      })
      syncLocks()
    }
  }
}

const rows = computed(() => {
  if (!showtime.value) return []
  const rowMap: any = {}
  showtime.value.seats.forEach((seat: any) => {
    if (!rowMap[seat.row]) rowMap[seat.row] = []
    rowMap[seat.row].push(seat)
  })
  return Object.entries(rowMap).sort()
})

async function toggleSeat(seat: any) {
  const isLockedByMe = seat.status === 'LOCKED' && seat.locked_by === authStore.user?.uid
  if ((seat.status !== 'AVAILABLE' && !isLockedByMe) || locking.value) return

  locking.value = true
  error.value = ""

  try {
    if (isLockedByMe) {
      // UNLOCK
      await api.post(`/protected/showtimes/${route.params.showtimeId}/unlock`, {
        seat_ids: [seat.id]
      })
    } else {
      // LOCK
      await api.post(`/protected/showtimes/${route.params.showtimeId}/lock`, {
        seat_ids: [seat.id]
      })
    }
    // Refresh from server to get final truth and call syncLocks()
    await fetchShowtime()
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || "Action failed"
    }
    await fetchShowtime()
  } finally {
    locking.value = false
  }
}

const totalPrice = computed(() => {
  return selectedSeats.value.length * 250
})

async function handleBooking() {
  if (selectedSeats.value.length === 0) return

  locking.value = true
  error.value = ""

  try {
    await api.post(`/protected/showtimes/${route.params.showtimeId}/confirm`, {
      seat_ids: selectedSeats.value
    })
    bookingSuccess.value = true
    selectedSeats.value = []
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value = err.response?.data?.error || "Booking failed"
    }
  } finally {
    locking.value = false
  }
}

onMounted(() => {
  fetchShowtime()
  setupWebSocket()
})

onUnmounted(() => {
  if (ws) ws.close()
})
</script>

<template>
  <div class="seat-selection container">
    <!-- Success Modal -->
    <div v-if="bookingSuccess" class="modal-overlay">
      <div class="success-modal">
        <div class="check-icon">✓</div>
        <h2>Booking Successful!</h2>
        <p>Your tickets have been reserved. We are processing your confirmation via RabbitMQ.</p>
        <button @click="router.push('/')" class="home-btn">Back to Movies</button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
    </div>

    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="fetchShowtime">Refresh Map</button>
    </div>

    <div v-else-if="showtime" class="booking-layout">
      <div class="map-section">
        <div class="screen">SCREEN</div>

        <div class="seats-grid">
          <div v-for="[rowName, seats] in rows" :key="rowName" class="row">
            <div class="row-label">{{ rowName }}</div>
            <div v-for="seat in (seats as any[])" :key="seat.id" class="seat" :class="[
              seat.status.toLowerCase(),
              {
                selected: selectedSeats.includes(seat.id),
                'locked-by-others': seat.status === 'LOCKED' && seat.locked_by !== authStore.user?.uid
              }
            ]" @click="toggleSeat(seat)">
              {{ seat.number }}
            </div>
          </div>
        </div>

        <div class="legend">
          <div class="legend-item"><span class="box available"></span> Available</div>
          <div class="legend-item"><span class="box selected"></span> Selected</div>
          <div class="legend-item"><span class="box booked"></span> Booked</div>
          <div class="legend-item"><span class="box locked"></span> Locked</div>
        </div>
      </div>

      <div class="summary-section">
        <div class="summary-card">
          <h2>Booking Summary</h2>
          <div class="movie-info">
            <p class="hall">{{ showtime.hall_name }}</p>
            <p class="time">
              {{
                new Date(showtime.start_time).toLocaleString('en-TH', {
                  day: '2-digit',
                  month: '2-digit',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                  hour12: true
                })
              }}
            </p>
          </div>

          <div class="selection-details">
            <p>Selected Seats: <strong>{{ selectedSeats.length > 0 ? selectedSeats.join(', ') : 'None' }}</strong></p>
            <div class="price-row">
              <span>Total Price:</span>
              <span class="amount">{{ totalPrice }} THB</span>
            </div>
          </div>

          <button class="confirm-btn" :disabled="selectedSeats.length === 0 || locking" @click="handleBooking">
            {{ locking ? 'Processing...' : 'Confirm Booking' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.seat-selection {
  padding: 40px 20px;
}

/* Success Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.success-modal {
  background: white;
  padding: 50px;
  border-radius: var(--radius-lg);
  text-align: center;
  max-width: 450px;
  box-shadow: var(--shadow-md);
  animation: popIn 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes popIn {
  from {
    transform: scale(0.8);
    opacity: 0;
  }

  to {
    transform: scale(1);
    opacity: 1;
  }
}

.check-icon {
  width: 80px;
  height: 80px;
  background: #4caf50;
  color: white;
  font-size: 40px;
  line-height: 80px;
  border-radius: 50%;
  margin: 0 auto 24px;
}

.home-btn {
  margin-top: 30px;
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 14px 30px;
  border-radius: var(--radius-md);
  font-weight: 700;
  cursor: pointer;
}

.booking-layout {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 40px;
}

.map-section {
  background: white;
  padding: 60px 40px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.screen {
  width: 80%;
  height: 8px;
  background: #ddd;
  margin-bottom: 80px;
  border-radius: 50% 50% 0 0;
  box-shadow: 0 15px 20px rgba(0, 0, 0, 0.1);
  text-align: center;
  color: #999;
  font-size: 0.7rem;
  line-height: 40px;
}

.seats-grid {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 60px;
}

.row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.row-label {
  width: 30px;
  font-weight: 700;
  color: var(--text-muted);
}

.seat {
  width: 35px;
  height: 35px;
  border-radius: 6px;
  background: #eee;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.seat.available:hover {
  background: #ddd;
  transform: scale(1.1);
}

.seat.selected,
.seat.locked {
  background: var(--primary-color);
  color: white;
}

.seat.booked,
.seat.locked-by-others {
  background: #333;
  color: #666;
  cursor: not-allowed;
}

/* User's own locks should be Red (using .locked when owned by user) */
.seat.locked {
  background: var(--primary-color);
  color: white;
}

/* Others' locks should be Yellow */
.seat.locked-by-others {
  background: #ffc107;
  color: #333;
  cursor: not-allowed;
}

.legend {
  display: flex;
  gap: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
  color: var(--text-muted);
}

.box {
  width: 15px;
  height: 15px;
  border-radius: 3px;
}

.box.available {
  background: #eee;
}

.box.selected {
  background: var(--primary-color);
}

.box.booked {
  background: #333;
}

.box.locked {
  background: #ffc107;
}

.summary-card {
  background: white;
  padding: 30px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  position: sticky;
  top: 100px;
}

h2 {
  margin-top: 0;
}

.movie-info {
  border-bottom: 1px solid #eee;
  padding-bottom: 20px;
  margin-bottom: 20px;
}

.time {
  font-weight: 700;
  color: var(--primary-color);
}

.selection-details {
  margin-bottom: 30px;
}

.price-row {
  display: flex;
  justify-content: space-between;
  font-size: 1.25rem;
  font-weight: 800;
  margin-top: 20px;
}

.amount {
  color: var(--primary-color);
}

.confirm-btn {
  width: 100%;
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 16px;
  border-radius: var(--radius-md);
  font-weight: 700;
  font-size: 1.1rem;
}

.confirm-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

@media (max-width: 1000px) {
  .booking-layout {
    grid-template-columns: 1fr;
  }
}
</style>
