<template>
  <div class="min-h-screen">
    <Navbar
      :books="books"
      :selectedBook="selectedBook"
      :currentPath="currentPath"
      @navigate="currentPath = $event"
      @bookChange="selectedBook = $event"
      @themeChange="handleThemeChange"
    />

    <main class="max-w-6xl mx-auto px-4 py-8">
      <Dashboard v-if="currentPath === '/' || currentPath === ''" :book="selectedBook" />
      <AddEntry v-else-if="currentPath === '/add'" :book="selectedBook" />
      <EntryListView v-else-if="currentPath === '/list'" :book="selectedBook" />
      <Statistics v-else-if="currentPath === '/stats'" :book="selectedBook" />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from './components/Navbar.vue'
import Dashboard from './views/Dashboard.vue'
import AddEntry from './views/AddEntry.vue'
import EntryListView from './views/EntryListView.vue'
import Statistics from './views/Statistics.vue'
import { fetchBooks } from './api/client'

const books = ref([])
const selectedBook = ref('')
const currentPath = ref('/')

async function loadBooks() {
  books.value = await fetchBooks()
  if (books.value.length > 0) {
    selectedBook.value = books.value[0].name
  }
}

function handleThemeChange(isDark) {
  localStorage.setItem('theme', isDark ? 'dark' : 'light')
}

onMounted(() => {
  loadBooks()

  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
  }
})
</script>
