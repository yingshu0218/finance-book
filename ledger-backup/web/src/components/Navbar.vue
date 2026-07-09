<template>
  <nav class="bg-card border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
    <div class="max-w-6xl mx-auto px-4 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-lg">L</span>
          </div>
          <span class="text-xl font-semibold">Ledger</span>
        </div>

        <div class="flex items-center gap-4">
          <select
            :value="selectedBook"
            @change="handleBookChange"
            class="bg-gray-100 dark:bg-gray-700 border-none"
          >
            <option v-for="book in books" :key="book.name" :value="book.name">
              {{ book.name }}
            </option>
          </select>

          <div class="flex items-center gap-2">
            <button
              v-for="item in navItems"
              :key="item.path"
              @click="$emit('navigate', item.path)"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
                currentPath === item.path
                  ? 'bg-primary text-white'
                  : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700'
              ]"
            >
              {{ item.label }}
            </button>
          </div>

          <button
            @click="toggleTheme"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" clip-rule="evenodd" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  books: { type: Array, default: () => [] },
  selectedBook: { type: String, default: '' },
  currentPath: { type: String, default: '/' }
})

const emit = defineEmits(['navigate', 'bookChange', 'themeChange'])

const isDark = ref(false)

const navItems = [
  { path: '/', label: '首页' },
  { path: '/add', label: '记账' },
  { path: '/list', label: '记录' },
  { path: '/stats', label: '统计' }
]

function toggleTheme() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
  emit('themeChange', isDark.value)
}

function handleBookChange(event) {
  emit('bookChange', event.target.value)
}

onMounted(() => {
  isDark.value = document.documentElement.classList.contains('dark')
})
</script>
