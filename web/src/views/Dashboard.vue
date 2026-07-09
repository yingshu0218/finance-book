<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">账本概览</h1>
      <p class="text-gray-500">查看当前账本的收支情况</p>
    </div>

    <div class="card mb-8">
      <div class="text-center py-8">
        <p class="text-gray-500 text-sm mb-2">当前余额</p>
        <p class="text-5xl font-bold font-mono" :class="balanceClass">
          ¥{{ formatNumber(balance) }}
        </p>
      </div>
    </div>

    <div class="grid grid-cols-3 gap-6 mb-8">
      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">本月收入</p>
            <p class="text-2xl font-bold font-mono text-green-600">
              ¥{{ formatNumber(income) }}
            </p>
          </div>
          <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">本月支出</p>
            <p class="text-2xl font-bold font-mono text-red-600">
              ¥{{ formatNumber(expense) }}
            </p>
          </div>
          <div class="w-12 h-12 bg-red-100 dark:bg-red-900/30 rounded-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">本月结余</p>
            <p class="text-2xl font-bold font-mono" :class="balanceClass">
              ¥{{ formatNumber(balance - prevBalance) }}
            </p>
          </div>
          <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <h2 class="text-lg font-semibold mb-4">近期记录</h2>
      <div v-if="recentEntries.length === 0" class="text-center py-8 text-gray-500">
        暂无记录
      </div>
      <div v-else class="space-y-3">
        <div
          v-for="entry in recentEntries"
          :key="entry.id"
          class="flex items-center justify-between p-4 rounded-lg bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        >
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium"
              :class="entry.amount > 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
            >
              {{ entry.amount > 0 ? '收' : '支' }}
            </div>
            <div>
              <p class="font-medium">{{ entry.category }}</p>
              <p class="text-sm text-gray-500">{{ entry.date }} {{ entry.note }}</p>
            </div>
          </div>
          <p class="font-mono font-semibold" :class="entry.amount > 0 ? 'text-green-600' : 'text-red-600'">
            {{ entry.amount > 0 ? '+' : '' }}¥{{ formatNumber(entry.amount) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { fetchBalance, fetchEntries } from '../api/client'

const props = defineProps({
  book: { type: String, default: '' }
})

const income = ref(0)
const expense = ref(0)
const balance = ref(0)
const prevBalance = ref(0)
const recentEntries = ref([])

const balanceClass = computed(() => {
  return balance.value >= 0 ? 'text-green-600' : 'text-red-600'
})

function formatNumber(num) {
  return Math.abs(num).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

async function loadData() {
  const now = new Date()
  const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`

  const balanceRes = await fetchBalance({ book: props.book, month: currentMonth })
  income.value = balanceRes.income
  expense.value = balanceRes.expense
  balance.value = balanceRes.balance

  const entriesRes = await fetchEntries({ book: props.book })
  recentEntries.value = entriesRes.slice(0, 10)
}

onMounted(loadData)

watch(() => props.book, loadData)
</script>
