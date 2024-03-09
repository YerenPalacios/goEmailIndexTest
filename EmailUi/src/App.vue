<script setup>
import { onMounted, ref } from 'vue';
import MessageCard from './components/MessageCard.vue'
import { messagesService } from './components/messagesService';
import FileImporter from './components/FileImporter.vue';
import moment from 'moment'

const { getMessages, addMessages, currentMessages, error } = messagesService()

const currentMessage = ref(null)
var currentSearch = ""

function updateCurrrentMessage(id) {
  const readed = localStorage.getItem('readed') || ""
  localStorage.setItem('readed', readed + id)

  const filteredMessage = currentMessages.value.filter(message => message.id === id)
  if (filteredMessage.length === 1) {
    currentMessage.value = {
      id: filteredMessage[0].id,
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

function onGetMessages(data){
  if (typeof data === "string"){
    error.value = data
  }
  const hits = data?.hits?.hits || []

  const get_from_name = xfrom => {
    const splitedText = xfrom?.split('<')
    if (splitedText && splitedText.length > 1) return splitedText[0]
    else return xfrom
  }

  const messages = hits.map(message => {
    const item = {
      id: message._id,
      subject: message._source.Subject || "(No subject provided)",
      content: message._source.content,
      from_email: message._source.From,
      from_name: get_from_name(message._source['X-From']),
      to_email: message._source.To,
      to_name: get_from_name(message._source['X-To']),
      date: message._source.Date
    }
    return item
  })
  currentMessages.value = messages
  return messages
}


function searchMessages() {
  getMessages(currentSearch, onGetMessages)
}

function appendMessages() {
  const from = currentMessages.value.length
  addMessages(currentSearch, from, from + 10)
}

onMounted(() => {
  getMessages('', onGetMessages)
})
</script>

<template>
  <section class="bg-[#fafafa] flex overflow-hidden">

    <div>
      <div class="max-h-screen overflow-y-auto overflow-x-hidden lg:max-h-[calc(100vh-50px)]">
        <div
          class="flex shadow-sm sticky top-0 bg-[#fafafa] hover:bg-[#fff] items-center px-6 py-4 border-b-[1px] border-slate-300">
          <img class="opacity-50 select-none" width="20" height="20" src="./assets/search.png" alt="">
          <input v-model="currentSearch" @keydown.enter="searchMessages" v-on:blur="searchMessages"
            class="outline-none bg-transparent text-slate-600 ml-3 " type="text">
        </div>
        <MessageCard v-for="message in currentMessages" @updateCurrentMessage="updateCurrrentMessage" :id=message.id
          :date=message.date :fromName=message.from_name :subject=message.subject :content=message.content />

        <button v-if="currentMessages?.length" @click="appendMessages"
          class="rounded-full hover:bg-blue-200 w-10 h-10 bg-blue-100 mx-auto my-5 text-blue-700 flex items-center justify-center text-2xl leading-none">+</button>

        <div v-if="!currentMessages?.length" class="flex flex-col justify-center items-center h-56">
          <p class="text-slate-300 select-none">No messages found...</p>
          <FileImporter />
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
  <div v-if="error" class="absolute shadow-md bg-red-500 right-10 bottom-10 px-4 py-2 rounded-sm text-slate-100">Error {{
    error }}</div>
</template>
