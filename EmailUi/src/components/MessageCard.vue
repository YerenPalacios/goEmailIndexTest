<script setup>
import moment from 'moment'
import { ref } from 'vue';
defineProps(['id', 'subject', 'content', 'fromName', 'date'])

const readed = ref(false)

function getReaded(id) {
    const localreaded = localStorage.getItem('readed') || ""
    return localreaded.includes(id)
}
</script>

<template>
    <div @click="readed = true;$emit('updateCurrentMessage', id);"
        class="w-[300px] text-slate-700 border-b-[1px] border-slate-300 px-6 py-4 cursor-pointer hover:bg-[#f0f0f0] duration-100">
        <p class="float-right text-slate-400 text-xs">{{ moment(date).fromNow() }} <span
                v-if="getReaded(id) || readed"
                class="rounded-full bg-green-500 opacity-50 text-white text-[10px] w-4 text-center inline-block">✓</span>
        </p>
        <h1 class="text-sm text-slate-500 font-bold">{{ fromName }}</h1>
        <h2 class="font-bold text-lg text-slate-700 mb-2 text-nowrap overflow-hidden overflow-ellipsis">{{ subject }}</h2>
        <p class="text-slate-500 text-sm text-nowrap overflow-hidden overflow-ellipsis">{{ content }}</p>
    </div>
</template>
