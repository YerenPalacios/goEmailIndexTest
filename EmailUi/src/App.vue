<script setup>
import { onMounted, ref } from 'vue';
import MessageCard from './components/MessageCard.vue'
import { messagesService } from './components/messagesService';
import moment from 'moment'

const { getMessages, data, error } = messagesService()

const currentMessage = ref(null)

function increaseCount(id) {
  const filteredMessage = data.value.filter(message => message.id === id)
  if (filteredMessage.length === 1) {
    currentMessage.value = {
      from: {
        name: filteredMessage[0].from_name,
        email: filteredMessage[0].from_email
      },
      to: {
        name: filteredMessage[0].to_name,
        email: filteredMessage[0].to_email
      },
      date: filteredMessage[0].date,
      subject: filteredMessage[0].subject,
      content: filteredMessage[0].content
    }
  }
}

function searchMessages(e) {
  getMessages(e.target.value)
}

onMounted(() => {
  getMessages('')
})
</script>

<template>
  <section class="bg-[#fafafa] flex overflow-hidden">

    <div>
      <div class="max-h-screen overflow-y-auto overflow-x-hidden lg:max-h-[calc(100vh-50px)]">
        <div
          class="flex shadow-sm sticky top-0 bg-[#fafafa] hover:bg-[#fff] items-center px-6 py-4 border-b-[1px] border-slate-300">
          <img class="opacity-50 select-none" width="20" height="20" src="./assets/search.png" alt="">
          <input @keydown.enter="searchMessages" v-on:blur="searchMessages" class="outline-none bg-transparent text-slate-600 ml-3" type="text">
        </div>
        <MessageCard v-for="message in data" @increase-by="increaseCount" :id=message.id :date=message.date
          :fromName=message.from_name :subject=message.subject :content=message.content />

        <div v-if="!data" class="flex justify-center items-center h-56">
          <p class="text-slate-300 select-none">No messages found...</p>
        </div>
      </div>
    </div>

    <article class="text-slate-600 min-w-[500px] bg-white max-h-screen overflow-auto lg:max-h-[calc(100vh-50px)]">
      <div v-if="currentMessage">
        <div class="p-8">
          <div class="flex gap-10 mb-2">
            <div>
              <h1 class="font-bold">From: {{ currentMessage.from.name }}</h1>
              <p class="text-sm text-slate-400">{{ currentMessage.from.email }}</p>
            </div>
            <div>
              <h1 class="font-bold">To: {{ currentMessage.to.name }}</h1>
              <p class="text-sm text-slate-400">{{ currentMessage.to.email }}</p>
            </div>
          </div>

          <p class="text-sm">{{ moment(currentMessage.date).format('LLLL') }}</p>
        </div>
        <div class="px-8 pb-8">
          <h1 class="text-2xl font-bold mb-3">{{ currentMessage.subject }}</h1>
          <p class="whitespace-pre">{{ currentMessage.content }}</p>
        </div>
      </div>
      <div v-else class="flex justify-center items-center h-full">
        <p class="text-slate-300 select-none">Select a message to view its content...</p>
      </div>
    </article>

  </section>
</template>
