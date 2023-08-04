<script setup lang="ts">
</script>

<template>
    <div>
        <h1 class="text-center my-8 text-4xl font-bold text-gray-800">Type a word</h1>
        <div class="flex flex-wrap justify-center">
            <input
            v-model="key"
                type="text"
                class="form-input border border-gray-400 px-4 py-2 rounded-lg shadow-sm focus:outline-none focus:shadow-outline"
                placeholder="Search by a word"
            />
            <div class="input-group-append ml-2">
                <button @click="getEmailByKey" class="btn btn-outline-secondary bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="button">
                Button
                </button>
            </div>
        </div>
        <div class="flex justify-around space-x-4 p-4">
            <div class="w-1/2">
                <table v-if="emails.length > 0" class="w-full divide-y divide-gray-200 border border-gray-300 border-collapse">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col" class="px-6 py-3 text-left text-sm text-gray-500 uppercase border border-gray-300">Date</th>
                            <th scope="col" class="px-6 py-3 text-left text-sm text-gray-500 uppercase border border-gray-300">From</th>
                            <th scope="col" class="px-6 py-3 text-left text-sm text-gray-500 uppercase border border-gray-300">To</th>
                            <th scope="col" class="px-6 py-3 text-left text-sm text-gray-500 uppercase border border-gray-300">Subject</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        <tr v-for="(email, Message_ID) in emails" :key="Message_ID" class="hover:bg-gray-300 cursor-pointer" @click="selectRow(email)">
                            <td class="px-6 py-4 whitespace-normal border border-gray-300">{{ email.Date }}</td>
                            <td class="px-6 py-4 whitespace-normal border border-gray-300">{{ email.From }}</td>
                            <td class="px-6 py-4 whitespace-normal border border-gray-300">{{ email.To }}</td>
                            <td class="px-6 py-4 whitespace-normal border border-gray-300">{{ email.Subject }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div v-if="Object.keys(emailSelected).length > 0" cols="30" rows="10" class="w-1/2 whitespace-pre-line p-2 border border-gray-300 rounded">
                <div v-html="highlightWord(emailSelected.Body, key)"></div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import EmailService from '../services/EmailService'
// import type
import type Email from '@/types/Email'

export default {
  name: 'emails-list',
  data() {
    return {
      emails: [] as Email[],
      emailSelected: {} as Email,
      key: '',
      message: 'No info to show'
    }
  },
  methods: {
    selectRow(selection : Email){
        this.emailSelected = selection
    },
    async getEmailByKey(){
      try {
        const response = await EmailService.get(this.key)
        console.log(response)
        this.emails = response.data
      } catch (e: any) {
        console.log('error: ', e)
        this.emails = []
        this.message = e.response.data
      }
    },
    highlightWord(text: string, word: string) {
      const regex = new RegExp(`(${word})`, "gi");
      return text.replace(regex, "<span class='bg-yellow-200'>\$1</span>");
    }
  }
}
</script>