<script setup lang="ts">
/* eslint-disable  @typescript-eslint/no-explicit-any */
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/services/api'
import axios from 'axios'

const route = useRoute()
const showtime = ref<any>(null)
const selectedSeats = ref<string[]>([])
const loading = ref(true)
const error = ref('')

async function fetchShowtime() {
  try {
    const id = route.params.showtimeId
    const res = await api.get(`/showtimes/${id}`)
    showtime.value = res.data
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

const rows = computed(() => {
  if (!showtime.value) return []
  const rowMap: any = {}
  showtime.value.seats.forEach((seat: any) => {
    if (!rowMap[seat.row]) rowMap[seat.row] = []
    rowMap[seat.row].push(seat)
  })
  return Object.entries(rowMap).sort()
})

function toggleSeat(seat: any) {
  if (seat.status !== 'AVAILABLE') return

  const index = selectedSeats.value.indexOf(seat.id)
  if (index === -1) {
    selectedSeats.value.push(seat.id)
  } else {
    selectedSeats.value.splice(index, 1)
  }
}

const totalPrice = computed(() => {
  return selectedSeats.value.length * 250 // Hardcoded price for now
})

async function handleBooking() {
  if (selectedSeats.value.length === 0) return
  alert(`Booking seats: ${selectedSeats.value.join(', ')}. \nTotal: ${totalPrice.value} THB`)
}

onMounted(fetchShowtime)
</script>

<template>
  <div class="seat-selection container">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
    </div>

    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
    </div>

    <div v-else-if="showtime" class="booking-layout">
      <div class="map-section">
        <div class="screen">SCREEN</div>

        <div class="seats-grid">
          <div v-for="[rowName, seats] in rows" :key="rowName" class="row">
            <div class="row-label">{{ rowName }}</div>
            <div v-for="seat in (seats as any[])" :key="seat.id" class="seat" :class="[
              seat.status.toLowerCase(),
              { selected: selectedSeats.includes(seat.id) }
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
            <p class="time">{{ new Date(showtime.start_time).toLocaleString() }}</p>
          </div>

          <div class="selection-details">
            <p>Selected Seats: <strong>{{ selectedSeats.length > 0 ? selectedSeats.join(', ') : 'None' }}</strong></p>
            <div class="price-row">
              <span>Total Price:</span>
              <span class="amount">{{ totalPrice }} THB</span>
            </div>
          </div>

          <button class="confirm-btn" :disabled="selectedSeats.length === 0" @click="handleBooking">
            Confirm Booking
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

.seat.selected {
  background: var(--primary-color);
  color: white;
}

.seat.booked {
  background: #333;
  color: #666;
  cursor: not-allowed;
}

.seat.locked {
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
