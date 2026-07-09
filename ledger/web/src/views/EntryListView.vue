<template>
  <div class="animate-fade-in">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">记录列表</h1>
      <p class="text-gray-500">查看和管理您的记账记录</p>
    </div>

    <div class="card mb-6">
      <div class="flex flex-wrap gap-4 items-center">
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-600">月份</label>
          <input
            v-model="filters.month"
            type="month"
            class="bg-gray-100 dark:bg-gray-700 border-none"
            @change="loadEntries"
          />
        </div>
        <div class="flex items-center gap-2">
          <label class="text-sm font-medium text-gray-600">分类</label>
          <select
            v-model="filters.category"
            class="bg-gray-100 dark:bg-gray-700 border-none"
            @change="loadEntries"
          >
            <option value="">全部</option>
            <option v-for="cat in categories" :key="cat" :value="cat">
              {{ cat }}
            </option>
          </select>
        </div>
        <button
          @click="resetFilters"
          class="btn-secondary ml-auto"
        >
          重置筛选
        </button>
      </div>
    </div>

    <div class="card">
      <div v-if="entries.length === 0" class="text-center py-12 text-gray-500">
        暂无记录
      </div>

      <table v-else class="w-full">
        <thead>
          <tr class="border-b border-gray-200 dark:border-gray-700">
            <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">日期</th>
            <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">分类</th>
            <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">备注</th>
            <th class="text-right py-3 px-4 text-sm font-medium text-gray-500">金额</th>
            <th class="text-center py-3 px-4 text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="entry in entries"
            :key="entry.id"
            class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            <td class="py-4 px-4">{{ entry.date }}</td>
            <td class="py-4 px-4">{{ entry.category }}</td>
            <td class="py-4 px-4 text-gray-500">{{ entry.note || '-' }}</td>
            <td class="py-4 px-4 text-right font-mono font-semibold" :class="entry.amount > 0 ? 'text-green-600' : 'text-red-600'">
              {{ entry.amount > 0 ? '+' : '' }}¥{{ formatNumber(entry.amount) }}
            </td>
            <td class="py-4 px-4 text-center">
              <button
                @click="deleteEntry(entry.id)"
                class="text-red-500 hover:text-red-700 p-2 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { fetchEntries, fetchCategories, deleteEntry as apiDeleteEntry } from '../api/client'

const props = defineProps({
  book: { type: String, default: '' }
})

const entries = ref([])
const categories = ref([])
const filters = ref({
  month: '',
  category: ''
})

function formatNumber(num) {
  return Math.abs(num).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

async function loadEntries() {
  const params = { book: props.book }
  if (filters.value.month) {
    params.month = filters.value.month
  }
  if (filters.value.category) {
    params.category = filters.value.category
  }
  entries.value = await fetchEntries(params)
}

function resetFilters() {
  filters.value = { month: '', category: '' }
  loadEntries()
}

async function deleteEntry(id) {
  if (!confirm('确定删除这条记录？')) return

  await apiDeleteEntry(id, props.book)
  loadEntries()
}

onMounted(async () => {
  categories.value = await fetchCategories()
  loadEntries()
})

watch(() => props.book, loadEntries)
</script>
