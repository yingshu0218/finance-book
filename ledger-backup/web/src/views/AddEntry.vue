<template>
  <div class="animate-fade-in max-w-lg mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-bold mb-2">记一笔</h1>
      <p class="text-gray-500">记录您的收支情况</p>
    </div>

    <div class="card">
      <form @submit.prevent="handleSubmit">
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">金额</label>
          <div class="relative">
            <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-lg">¥</span>
            <input
              v-model="form.amount"
              type="number"
              step="0.01"
              placeholder="0.00"
              class="pl-10 text-3xl font-mono font-bold"
            />
          </div>
          <div class="flex gap-4 mt-3">
            <button
              type="button"
              @click="form.type = 'expense'"
              :class="[
                'flex-1 py-3 rounded-lg font-medium transition-all',
                form.type === 'expense'
                  ? 'bg-red-100 text-red-700 ring-2 ring-red-500'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              支出
            </button>
            <button
              type="button"
              @click="form.type = 'income'"
              :class="[
                'flex-1 py-3 rounded-lg font-medium transition-all',
                form.type === 'income'
                  ? 'bg-green-100 text-green-700 ring-2 ring-green-500'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              收入
            </button>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">分类</label>
          <div class="grid grid-cols-3 gap-2">
            <button
              v-for="cat in categories"
              :key="cat"
              type="button"
              @click="form.category = cat"
              :class="[
                'py-2 px-3 rounded-lg text-sm font-medium transition-all',
                form.category === cat
                  ? 'bg-primary text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              {{ cat }}
            </button>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">日期</label>
            <input
              v-model="form.date"
              type="date"
              class="bg-gray-100 dark:bg-gray-700 border-none"
            />
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">备注</label>
          <textarea
            v-model="form.note"
            rows="3"
            placeholder="添加备注..."
            class="bg-gray-100 dark:bg-gray-700 border-none resize-none"
          ></textarea>
        </div>

        <button
          type="submit"
          class="w-full btn-primary text-lg py-4"
          :disabled="!form.amount || !form.category"
        >
          保存
        </button>
      </form>
    </div>

    <div v-if="success" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl p-8 text-center animate-fade-in">
        <div class="w-16 h-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold mb-2">记账成功</h3>
        <p class="text-gray-500 mb-4">记录已保存</p>
        <button @click="resetForm" class="btn-primary">继续记账</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { createEntry, fetchCategories } from '../api/client'

const props = defineProps({
  book: { type: String, default: '' }
})

const form = ref({
  amount: '',
  type: 'expense',
  category: '',
  date: '',
  note: ''
})

const categories = ref([])
const success = ref(false)

onMounted(async () => {
  categories.value = await fetchCategories()

  const now = new Date()
  form.value.date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
})

async function handleSubmit() {
  const amount = parseFloat(form.value.amount)
  if (!amount || !form.value.category) return

  const finalAmount = form.value.type === 'expense' ? -amount : amount

  await createEntry({
    book: props.book,
    amount: finalAmount,
    category: form.value.category,
    date: form.value.date,
    note: form.value.note
  })

  success.value = true
}

function resetForm() {
  success.value = false
  form.value = {
    amount: '',
    type: 'expense',
    category: '',
    date: '',
    note: ''
  }

  const now = new Date()
  form.value.date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
}
</script>
