<script setup>
import { ref } from 'vue';
import { messagesService } from './messagesService';

var importFile = null

const { uploadFile, error, loading } = messagesService()
const uploaded = ref(false)

const handleUploadFile = (e) => {
    if (e.target.files.length) importFile = e.target.files[0]
}

const dataUploaded = (data) => {
    if (data?.status == "Ok") {
        uploaded.value = true
    }
    else error.value = data
}

const sendFile = () => {
    uploadFile(importFile, dataUploaded)
}

</script>
<template>
    <p v-if="uploaded" class="text-green-300">File uploaded</p>
    <div v-if="loading" class="">
        <div class="loader"></div>
    </div>
    <div v-else class="text-center">
        <input v-on:change="handleUploadFile" class="text-slate-400 m-3 
              file:mr-4 file:py-2 file:px-4 file:cursor-pointer
              file:rounded-full file:border-0
              file:text-sm file:font-semibold
            file:bg-blue-50
             file:text-blue-700  hover:file:bg-blue-100" type="file">
        <button @click="sendFile" class="bg-blue-100 text-blue-700 hover:bg-blue-200 py-2 px-4 rounded-full">Import
            messages</button>

    </div>
    <div v-if="error" class="absolute shadow-md bg-red-500 right-10 bottom-10 px-4 py-2 rounded-sm text-slate-100">Error: {{
        error }}</div>
</template>