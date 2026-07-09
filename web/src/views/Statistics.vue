<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">统计分析</h1>
      <p class="text-gray-500">查看收支趋势和分类占比</p>
    </div>

    <div class="card mb-6">
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-600">月份</label>
          <input
            v-model="month"
            type="month"
            class="bg-gray-100 dark:bg-gray-700 border-none"
            @change="loadStats"
          />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6 mb-8">
      <div class="card">
        <h2 class="text-lg font-semibold mb-4">月度收支</h2>
        <div class="h-64">
          <Bar :data="barChartData" :options="barOptions" />
        </div>
      </div>

      <div class="card">
        <h2 class="text-lg font-semibold mb-4">分类占比</h2>
        <div class="h-64">
          <Doughnut :data="doughnutChartData" :options="doughnutOptions" />
        </div>
      </div>
    </div>

    <div class="card">
      <h2 class="text-lg font-semibold mb-4">收支趋势</h2>
      <div class="h-64">
        <Line :data="lineChartData" :options="lineOptions" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Bar, Doughnut, Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  LineElement,
  PointElement,
  Filler
} from 'chart.js'
import { fetchBalance, fetchEntries } from '../api/client'

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  LineElement,
  PointElement,
  Filler
)

const props = defineProps({
  book: { type: String, default: '' }
})

const month = ref('')
const balanceData = ref({ income: 0, expense: 0 })
const entriesData = ref([])

const barChartData = computed(() => ({
  labels: ['收入', '支出'],
  datasets: [{
    label: '金额',
    data: [balanceData.value.income, balanceData.value.expense],
    backgroundColor: ['#10B981', '#EF4444'],
    borderRadius: 8
  }]
}))

const barOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  }
}

const doughnutChartData = computed(() => {
  const categoryMap = {}
  entriesData.value.forEach(e => {
    if (e.amount < 0) {
      categoryMap[e.category] = (categoryMap[e.category] || 0) + Math.abs(e.amount)
    }
  })

  return {
    labels: Object.keys(categoryMap),
    datasets: [{
      data: Object.values(categoryMap),
      backgroundColor: [
        '#3B82F6', '#8B5CF6', '#EC4899', '#F59E0B', '#10B981',
        '#EF4444', '#6366F1', '#14B8A6', '#F97316', '#84CC16'
      ],
      borderWidth: 0
    }]
  }
})

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom' }
  }
}

const lineChartData = computed(() => {
  const last12Months = []
  const now = new Date()
  for (let i = 11; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
    last12Months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
  }

  return {
    labels: last12Months.map(m => m.slice(5)),
    datasets: [{
      label: '收入',
      data: last12Months.map(() => Math.random() * 10000),
      borderColor: '#10B981',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true,
      tension: 0.4
    }, {
      label: '支出',
      data: last12Months.map(() => Math.random() * 8000),
      borderColor: '#EF4444',
      backgroundColor: 'rgba(239, 68, 68, 0.1)',
      fill: true,
      tension: 0.4
    }]
  }
})

const lineOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' },
  plugins: {
    legend: { position: 'top' }
  }
}

async function loadStats() {
  const balanceRes = await fetchBalance({ book: props.book, month: month.value })
  balanceData.value = {
    income: balanceRes.income,
    expense: balanceRes.expense
  }

  entriesData.value = await fetchEntries({ book: props.book, month: month.value })
}

onMounted(() => {
  const now = new Date()
  month.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  loadStats()
})

watch(() => props.book, loadStats)
</script>
